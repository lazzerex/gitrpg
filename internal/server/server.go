package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/lazzerex/gitrpg/internal/achievements"
	"github.com/lazzerex/gitrpg/internal/auth"
	"github.com/lazzerex/gitrpg/internal/characters"
	"github.com/lazzerex/gitrpg/internal/config"
	"github.com/lazzerex/gitrpg/internal/stats"
	svgpkg "github.com/lazzerex/gitrpg/internal/svg"
	"github.com/lazzerex/gitrpg/internal/users"
	"github.com/lazzerex/gitrpg/internal/worker"
)

type Server struct {
	cfg          *config.Config
	router       *chi.Mux
	db           *pgxpool.Pool
	redis        *redis.Client
	logger       *slog.Logger
	templates    map[string]*template.Template
	auth         *auth.Handler
	worker       *worker.Worker
	users        *users.Store
	characters   *characters.Service
	achievements *achievements.Service
	syncStart    sync.Map // userID int64 → time.Time
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, logger *slog.Logger, w *worker.Worker, charSvc *characters.Service, achSvc *achievements.Service, userStore *users.Store) *Server {
	authHandler := auth.NewHandler(cfg, userStore, logger)
	authHandler.SetPostLogin(w.SyncUser)

	s := &Server{
		cfg:          cfg,
		router:       chi.NewRouter(),
		db:           db,
		redis:        rdb,
		logger:       logger,
		auth:         authHandler,
		worker:       w,
		users:        userStore,
		characters:   charSvc,
		achievements: achSvc,
	}
	s.registerMiddleware()
	s.registerRoutes()
	return s
}

var templateFuncs = template.FuncMap{
	"inc": func(n int) int { return n + 1 },
	"xpPercent": func(into, for_ int) int {
		if for_ <= 0 {
			return 0
		}
		pct := into * 100 / for_
		if pct > 100 {
			return 100
		}
		return pct
	},
	"fmtAge": func(t time.Time) string {
		d := time.Since(t)
		switch {
		case d < time.Minute:
			return "just now"
		case d < time.Hour:
			m := int(d.Minutes())
			if m == 1 {
				return "1 minute ago"
			}
			return fmt.Sprintf("%d minutes ago", m)
		case d < 24*time.Hour:
			h := int(d.Hours())
			if h == 1 {
				return "1 hour ago"
			}
			return fmt.Sprintf("%d hours ago", h)
		default:
			days := int(d.Hours() / 24)
			if days == 1 {
				return "1 day ago"
			}
			return fmt.Sprintf("%d days ago", days)
		}
	},
}

// LoadTemplates builds a per-page template set (base.html + page) for each page
// in dir. Each page gets its own isolated set so {{define "content"}} blocks don't collide.
func (s *Server) LoadTemplates(dir string) error {
	pages := []string{"index.html", "profile.html", "public.html", "cards.html"}
	base := filepath.Join(dir, "base.html")
	partial := filepath.Join(dir, "partials", "char-panel.html")
	s.templates = make(map[string]*template.Template, len(pages)+1)

	for _, page := range pages {
		files := []string{base, filepath.Join(dir, page)}
		if page == "profile.html" {
			files = append(files, partial)
		}
		tmpl, err := template.New("").Funcs(templateFuncs).ParseFiles(files...)
		if err != nil {
			return err
		}
		s.templates[page] = tmpl
	}

	// Standalone partial for HTMX responses
	tmpl, err := template.New("").Funcs(templateFuncs).ParseFiles(partial)
	if err != nil {
		return err
	}
	s.templates["char-panel"] = tmpl
	return nil
}

func (s *Server) renderPartial(w http.ResponseWriter, name string, data any) {
	tmpl, ok := s.templates[name]
	if !ok {
		http.Error(w, "partial not found: "+name, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		s.logger.Error("partial render failed", "name", name, "error", err)
	}
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
	s.router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/", s.handleIndex)

	s.router.Get("/auth/github", s.auth.Login)
	s.router.Get("/auth/github/callback", s.auth.Callback)
	s.router.Get("/logout", s.auth.Logout)

	s.router.Get("/card/demo", s.handleCardDemo)
	s.router.Get("/card/{username}", s.handleCard)
	s.router.Get("/card/compact/{username}", s.handleCardCompact)

	s.router.Get("/u/{username}", s.handlePublicProfile)
	s.router.Get("/cards", s.handleCards)

	s.router.Group(func(r chi.Router) {
		r.Use(s.auth.RequireAuth)
		r.Get("/profile", s.handleProfile)
		r.Post("/sync", s.handleSync)
		r.Get("/sync/status", s.handleSyncStatus)
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	u, _ := r.Context().Value(users.ContextKey).(*users.User)
	var charClass string
	if u != nil {
		char, err := s.characters.GetByUserID(r.Context(), u.ID)
		if err == nil && char != nil {
			charClass = char.Class
		}
	}
	s.render(w, "index.html", indexData{User: u, CharClass: charClass})
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(users.ContextKey).(*users.User)

	char, err := s.characters.GetByUserID(r.Context(), user.ID)
	if err != nil && err != pgx.ErrNoRows {
		s.logger.Error("character load failed", "user_id", user.ID, "error", err)
	}

	var isStale bool
	if char != nil {
		isStale = time.Since(char.UpdatedAt) > 12*time.Hour
	}

	accentColor := svgpkg.ClassColor("")
	if char != nil {
		accentColor = svgpkg.ClassColor(char.Class)
	}

	achs, err := s.achievements.GetForUser(r.Context(), user.ID)
	if err != nil {
		s.logger.Error("achievements load failed", "user_id", user.ID, "error", err)
	}

	s.render(w, "profile.html", profileData{
		User:         user,
		Character:    char,
		IsStale:      isStale,
		AccentColor:  accentColor,
		BaseURL:      requestBaseURL(r),
		Achievements: achs,
	})
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(users.ContextKey).(*users.User)
	s.syncStart.Store(user.ID, time.Now())
	s.worker.SyncUser(user)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, syncingHTML)
}

const syncingHTML = `<div id="sync-status" hx-get="/sync/status" hx-trigger="every 2s" hx-swap="outerHTML"><p class="blink" style="color:var(--gold);font-size:8px;letter-spacing:1px;">SYNCING...</p></div>`

func (s *Server) handleSyncStatus(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(users.ContextKey).(*users.User)

	startVal, ok := s.syncStart.Load(user.ID)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	startTime := startVal.(time.Time)

	if time.Since(startTime) > 3*time.Minute {
		s.syncStart.Delete(user.ID)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, syncBtnHTML)
		return
	}

	char, err := s.characters.GetByUserID(r.Context(), user.ID)
	if err != nil || char == nil || !char.UpdatedAt.After(startTime) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, syncingHTML)
		return
	}

	s.syncStart.Delete(user.ID)

	accentColor := svgpkg.ClassColor(char.Class)
	achs, _ := s.achievements.GetForUser(r.Context(), user.ID)

	w.Header().Set("HX-Retarget", "#char-panel")
	w.Header().Set("HX-Reswap", "outerHTML")
	s.renderPartial(w, "char-panel", profileData{
		User:        user,
		Character:   char,
		IsStale:     time.Since(char.UpdatedAt) > 12*time.Hour,
		AccentColor: accentColor,
		Achievements: achs,
	})
}

const syncBtnHTML = `<div id="sync-status"><button hx-post="/sync" hx-target="#sync-status" hx-swap="outerHTML" class="px-btn" style="font-size:8px;padding:8px 14px;display:inline-flex;align-items:center;gap:6px;"><i data-lucide="refresh-cw" style="width:12px;height:12px;stroke:var(--gold);stroke-width:2;"></i>SYNC NOW</button></div>`

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

type baseData struct {
	User *users.User
}

type indexData struct {
	User      *users.User
	CharClass string
}

type profileData struct {
	User         *users.User
	Character    *stats.Character
	IsStale      bool
	AccentColor  string
	BaseURL      string
	Achievements []achievements.UserAchievement
}

type publicProfileData struct {
	User         *users.User
	ProfileUser  *users.User
	Character    *stats.Character
	AccentColor  string
	BaseURL      string
	Achievements []achievements.UserAchievement
}

func (s *Server) baseData(r *http.Request) baseData {
	u, _ := r.Context().Value(users.ContextKey).(*users.User)
	return baseData{User: u}
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func (s *Server) handlePublicProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	viewer, _ := r.Context().Value(users.ContextKey).(*users.User)

	profileUser, err := s.users.GetByLogin(r.Context(), username)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	char, _ := s.characters.GetByUserID(r.Context(), profileUser.ID)

	accentColor := svgpkg.ClassColor("")
	if char != nil {
		accentColor = svgpkg.ClassColor(char.Class)
	}

	achs, err := s.achievements.GetForUser(r.Context(), profileUser.ID)
	if err != nil {
		s.logger.Error("achievements load failed", "user_id", profileUser.ID, "error", err)
	}

	s.render(w, "public.html", publicProfileData{
		User:         viewer,
		ProfileUser:  profileUser,
		Character:    char,
		AccentColor:  accentColor,
		BaseURL:      requestBaseURL(r),
		Achievements: achs,
	})
}

func (s *Server) handleCards(w http.ResponseWriter, r *http.Request) {
	s.render(w, "cards.html", s.baseData(r))
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

func (s *Server) handleCardDemo(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	svg, err := svgpkg.Demo(class)
	if err != nil {
		http.Error(w, "demo generation failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=86400")
	s.svgResponse(w, svg)
}

func (s *Server) handleCard(w http.ResponseWriter, r *http.Request) {
	s.serveCard(w, r, false)
}

func (s *Server) handleCardCompact(w http.ResponseWriter, r *http.Request) {
	s.serveCard(w, r, true)
}

func (s *Server) serveCard(w http.ResponseWriter, r *http.Request, compact bool) {
	username := chi.URLParam(r, "username")
	// strip .svg suffix — URLs are /card/username.svg
	if len(username) > 4 && username[len(username)-4:] == ".svg" {
		username = username[:len(username)-4]
	}

	style := r.URL.Query().Get("style")
	switch style {
	case "chart", "stats":
	default:
		style = ""
	}

	prefix := "svg:card:"
	if compact {
		prefix = "svg:compact:"
	}
	cacheKey := prefix + username
	if style != "" {
		cacheKey = prefix + username + ":" + style
	}

	cached, err := s.redis.Get(r.Context(), cacheKey).Result()
	if err == nil {
		s.svgResponse(w, cached)
		return
	}
	if err != redis.Nil {
		s.logger.Warn("redis get failed", "key", cacheKey, "error", err)
	}

	user, err := s.users.GetByLogin(r.Context(), username)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	char, err := s.characters.GetByUserID(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "character not found — sync first", http.StatusNotFound)
		return
	}

	var svgStr string
	if compact {
		svgStr, err = svgpkg.Compact(user.Login, char)
	} else {
		svgStr, err = svgpkg.Card(user.Login, char, style)
	}
	if err != nil {
		s.logger.Error("svg generation failed", "user", username, "error", err)
		http.Error(w, "svg generation failed", http.StatusInternalServerError)
		return
	}

	if setErr := s.redis.Set(r.Context(), cacheKey, svgStr, time.Hour).Err(); setErr != nil {
		s.logger.Warn("redis set failed", "key", cacheKey, "error", setErr)
	}

	s.svgResponse(w, svgStr)
}

func (s *Server) svgResponse(w http.ResponseWriter, svg string) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = fmt.Fprint(w, svg)
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
