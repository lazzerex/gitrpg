package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/lazzerex/gitrpg/internal/auth"
	"github.com/lazzerex/gitrpg/internal/config"
	"github.com/lazzerex/gitrpg/internal/users"
	"github.com/lazzerex/gitrpg/internal/worker"
)

type Server struct {
	cfg       *config.Config
	router    *chi.Mux
	db        *pgxpool.Pool
	redis     *redis.Client
	logger    *slog.Logger
	templates map[string]*template.Template
	auth      *auth.Handler
	worker    *worker.Worker
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, logger *slog.Logger, w *worker.Worker) *Server {
	userStore := users.NewStore(db)
	authHandler := auth.NewHandler(cfg, userStore, logger)
	authHandler.SetPostLogin(w.SyncUser)

	s := &Server{
		cfg:    cfg,
		router: chi.NewRouter(),
		db:     db,
		redis:  rdb,
		logger: logger,
		auth:   authHandler,
		worker: w,
	}
	s.registerMiddleware()
	s.registerRoutes()
	return s
}

// LoadTemplates builds a per-page template set (base.html + page) for each page
// in dir. Each page gets its own isolated set so {{define "content"}} blocks don't collide.
func (s *Server) LoadTemplates(dir string) error {
	pages := []string{"index.html", "profile.html"}
	s.templates = make(map[string]*template.Template, len(pages))
	base := filepath.Join(dir, "base.html")
	for _, page := range pages {
		tmpl, err := template.ParseFiles(base, filepath.Join(dir, page))
		if err != nil {
			return err
		}
		s.templates[page] = tmpl
	}
	return nil
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	tmpl, ok := s.templates[name]
	if !ok {
		http.Error(w, "template not found: "+name, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		s.logger.Error("template render failed", "name", name, "error", err)
	}
}

func (s *Server) registerMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Recoverer)
	s.router.Use(s.requestLogger)
	s.router.Use(s.auth.LoadUser)
}

func (s *Server) registerRoutes() {
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/", s.handleIndex)

	s.router.Get("/auth/github", s.auth.Login)
	s.router.Get("/auth/github/callback", s.auth.Callback)
	s.router.Get("/logout", s.auth.Logout)

	s.router.Group(func(r chi.Router) {
		r.Use(s.auth.RequireAuth)
		r.Get("/profile", s.handleProfile)
		r.Post("/sync", s.handleSync)
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.render(w, "index.html", s.templateData(r))
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	s.render(w, "profile.html", s.templateData(r))
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(users.ContextKey).(*users.User)
	s.worker.SyncUser(user)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := s.db.Ping(ctx); err != nil {
		s.logger.Error("db ping failed", "error", err)
		http.Error(w, "db unavailable", http.StatusServiceUnavailable)
		return
	}

	if err := s.redis.Ping(ctx).Err(); err != nil {
		s.logger.Error("redis ping failed", "error", err)
		http.Error(w, "redis unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"status":"ok"}`)
}

type templateData struct {
	User *users.User
}

func (s *Server) templateData(r *http.Request) templateData {
	u, _ := r.Context().Value(users.ContextKey).(*users.User)
	return templateData{User: u}
}

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		s.logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration", time.Since(start),
			"request_id", middleware.GetReqID(r.Context()),
		)
	})
}

func (s *Server) Start() error {
	addr := ":" + s.cfg.Server.Port
	s.logger.Info("server starting", "addr", addr, "env", s.cfg.Server.Env)
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}
