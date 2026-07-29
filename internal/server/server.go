package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
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
	"github.com/lazzerex/gitrpg/internal/equipment"
	"github.com/lazzerex/gitrpg/internal/leaderboards"
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
	equipment    *equipment.Service
	leaderboards *leaderboards.Service
	syncStart    sync.Map // userID int64 → time.Time
	httpServer   *http.Server
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, logger *slog.Logger, w *worker.Worker, charSvc *characters.Service, achSvc *achievements.Service, eqSvc *equipment.Service, lbSvc *leaderboards.Service, userStore *users.Store) *Server {
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
		equipment:    eqSvc,
		leaderboards: lbSvc,
	}
	s.registerMiddleware()
	s.registerRoutes()
	return s
}

var classIcons = map[string]string{
	"Guardian":   "shield",
	"Knight":     "sword",
	"Paladin":    "scroll-text",
	"Berserker":  "axe",
	"Warlord":    "flame",
	"Sage":       "book-open",
	"Battlemage": "wand-2",
	"Rogue":      "key",
	"Wanderer":   "compass",
}

var templateFuncs = template.FuncMap{
	"inc": func(n int) int { return n + 1 },
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
	"classIcon": func(class string) string {
		icon, ok := classIcons[class]
		if !ok {
			return classIcons["Wanderer"]
		}
		return icon
	},
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
	pages := []string{"index.html", "profile.html", "public.html", "cards.html", "leaderboard.html"}
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

func (s *Server) renderFragment(w http.ResponseWriter, tmplName, blockName string, data any) {
	tmpl, ok := s.templates[tmplName]
	if !ok {
		http.Error(w, "template not found: "+tmplName, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, blockName, data); err != nil {
		s.logger.Error("template render failed", "name", blockName, "error", err)
	}
}

func (s *Server) registerMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Recoverer)
	s.router.Use(s.securityHeaders)
	s.router.Use(s.requestLogger)
	s.router.Use(s.auth.LoadUser)
}

// importmapHash pins the sole remaining inline <script> in base.html (a
// type="importmap" block — the Import Maps spec needs it inline for browser
// support, so it can't move to an external file like the rest of the JS did).
// Recompute if that script's content changes even by one byte of whitespace.
const importmapHash = "'sha256-/LPWhjT5pvTvaGMtuaM1D1UezzGKyq0GNIvBMjhCUgY='"

const cspPolicy = "default-src 'self'; " +
	"script-src 'self' " + importmapHash + " https://unpkg.com https://cdn.jsdelivr.net; " +
	"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
	"font-src 'self' https://fonts.gstatic.com; " +
	"img-src 'self' data: blob: https:; " +
	"connect-src 'self' blob: https://unpkg.com; " +
	"frame-ancestors 'none'"

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", cspPolicy)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		if s.cfg.Server.Env == "production" {
			h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) registerRoutes() {
	s.router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/robots.txt", s.handleRobotsTxt)
	s.router.Get("/sitemap.xml", s.handleSitemap)
	s.router.Get("/", s.handleIndex)

	s.router.Get("/auth/github", s.auth.Login)
	s.router.Get("/auth/github/callback", s.auth.Callback)
	s.router.Get("/logout", s.auth.Logout)

	s.router.Get("/card/demo", s.handleCardDemo)
	s.router.Get("/card/{username}", s.handleCard)

	s.router.Get("/u/{username}", s.handlePublicProfile)
	s.router.Get("/cards", s.handleCards)
	s.router.Get("/leaderboard", s.handleLeaderboard)

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
	s.render(w, "index.html", indexData{User: u, CharClass: charClass, BaseURL: requestBaseURL(r)})
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

	loadout, err := s.equipment.GetForUser(r.Context(), user.ID)
	if err != nil {
		s.logger.Error("equipment load failed", "user_id", user.ID, "error", err)
	}

	s.render(w, "profile.html", profileData{
		User:         user,
		Character:    char,
		IsStale:      isStale,
		AccentColor:  accentColor,
		BaseURL:      requestBaseURL(r),
		Achievements: achs,
		Equipment:    loadout,
	})
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(users.ContextKey).(*users.User)

	rlKey := fmt.Sprintf("ratelimit:sync:%d", user.ID)
	acquired, err := s.redis.SetNX(r.Context(), rlKey, 1, 30*time.Second).Result()
	if err != nil {
		s.logger.Warn("sync rate limit check failed", "user_id", user.ID, "error", err)
	} else if !acquired {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, syncCooldownHTML)
		return
	}

	s.syncStart.Store(user.ID, time.Now())
	s.worker.SyncUser(user)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, syncingHTML)
}

const syncingHTML = `<div id="sync-status" hx-get="/sync/status" hx-trigger="every 2s" hx-swap="outerHTML"><p class="blink" style="color:var(--gold);font-size:8px;letter-spacing:1px;">SYNCING...</p></div>`

const syncCooldownHTML = `<div id="sync-status"><p style="color:var(--gold);font-size:8px;letter-spacing:1px;">SYNCED RECENTLY, TRY AGAIN SHORTLY.</p></div>`

func (s *Server) handleSyncStatus(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(users.ContextKey).(*users.User)

	startVal, ok := s.syncStart.Load(user.ID)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	startTime := startVal.(time.Time)

	char, err := s.characters.GetByUserID(r.Context(), user.ID)
	if err == nil && char != nil && char.UpdatedAt.After(startTime) {
		s.syncStart.Delete(user.ID)
		_ = s.redis.Del(r.Context(), "svg:card:"+user.Login).Err()
		w.Header().Set("HX-Redirect", "/profile")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if syncStatus, statusErr := s.worker.LatestSyncStatus(r.Context(), user.ID); statusErr == nil &&
		syncStatus != nil && syncStatus.Status == "failed" &&
		syncStatus.CompletedAt != nil && syncStatus.CompletedAt.After(startTime) {
		s.syncStart.Delete(user.ID)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, syncErrorHTML)
		return
	}

	if time.Since(startTime) > 3*time.Minute {
		s.syncStart.Delete(user.ID)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, syncBtnHTML)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, syncingHTML)
}

const syncBtnHTML = `<div id="sync-status"><button hx-post="/sync" hx-target="#sync-status" hx-swap="outerHTML" class="px-btn" style="font-size:8px;padding:8px 14px;display:inline-flex;align-items:center;gap:6px;"><i aria-hidden="true" data-lucide="refresh-cw" style="width:12px;height:12px;stroke:var(--gold);stroke-width:2;"></i>SYNC NOW</button></div>`

const syncErrorHTML = `<div id="sync-status"><p style="color:var(--red);font-size:8px;margin-bottom:8px;">SYNC FAILED. TRY AGAIN?</p><button hx-post="/sync" hx-target="#sync-status" hx-swap="outerHTML" class="px-btn px-btn-red" style="font-size:8px;padding:8px 14px;display:inline-flex;align-items:center;gap:6px;"><i aria-hidden="true" data-lucide="refresh-cw" style="width:12px;height:12px;stroke:var(--red);stroke-width:2;"></i>RETRY SYNC</button></div>`

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

func (s *Server) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintf(w, "User-agent: *\nAllow: /\nDisallow: /profile\nDisallow: /sync\n\nSitemap: %s/sitemap.xml\n", requestBaseURL(r))
}

func (s *Server) handleSitemap(w http.ResponseWriter, r *http.Request) {
	base := requestBaseURL(r)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>%s/</loc></url>
  <url><loc>%s/cards</loc></url>
  <url><loc>%s/leaderboard</loc></url>
</urlset>
`, base, base, base)
}

type baseData struct {
	User    *users.User
	BaseURL string
}

type indexData struct {
	User      *users.User
	CharClass string
	BaseURL   string
}

type profileData struct {
	User         *users.User
	Character    *stats.Character
	IsStale      bool
	AccentColor  string
	BaseURL      string
	Achievements []achievements.UserAchievement
	Equipment    equipment.Loadout
}

type publicProfileData struct {
	User         *users.User
	ProfileUser  *users.User
	Character    *stats.Character
	AccentColor  string
	BaseURL      string
	Achievements []achievements.UserAchievement
	Equipment    equipment.Loadout
}

func (s *Server) baseData(r *http.Request) baseData {
	u, _ := r.Context().Value(users.ContextKey).(*users.User)
	return baseData{User: u, BaseURL: requestBaseURL(r)}
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

// githubUsernameRe matches GitHub's username rules: 1-39 chars, alphanumeric
// or single hyphens, cannot start/end with a hyphen.
var githubUsernameRe = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$`)

func (s *Server) handlePublicProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if !githubUsernameRe.MatchString(username) {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
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

	loadout, err := s.equipment.GetForUser(r.Context(), profileUser.ID)
	if err != nil {
		s.logger.Error("equipment load failed", "user_id", profileUser.ID, "error", err)
	}

	s.render(w, "public.html", publicProfileData{
		User:         viewer,
		ProfileUser:  profileUser,
		Character:    char,
		AccentColor:  accentColor,
		BaseURL:      requestBaseURL(r),
		Achievements: achs,
		Equipment:    loadout,
	})
}

func (s *Server) handleCards(w http.ResponseWriter, r *http.Request) {
	s.render(w, "cards.html", s.baseData(r))
}

type leaderboardData struct {
	User        *users.User
	Leaderboard leaderboards.Page
	BaseURL     string
}

func (s *Server) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	lb, err := s.leaderboards.GetXPPage(r.Context(), page)
	if err != nil {
		s.logger.Error("leaderboard load failed", "error", err)
		http.Error(w, "leaderboard unavailable", http.StatusInternalServerError)
		return
	}

	data := leaderboardData{User: s.baseData(r).User, Leaderboard: lb, BaseURL: requestBaseURL(r)}

	if wantsLeaderboardFragment(r) {
		s.renderFragment(w, "leaderboard.html", "leaderboard-results", data)
		return
	}
	s.render(w, "leaderboard.html", data)
}

// Plain hx-boost nav also sends HX-Request, so check HX-Target instead -
// only the pagination links set hx-target="#leaderboard-results".
func wantsLeaderboardFragment(r *http.Request) bool {
	return r.Header.Get("HX-Target") == "leaderboard-results"
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
	s.serveCard(w, r)
}

func (s *Server) serveCard(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if len(username) > 4 && username[len(username)-4:] == ".svg" {
		username = username[:len(username)-4]
	}
	if !githubUsernameRe.MatchString(username) {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	cacheKey := "svg:card:" + username

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

	svgStr, err := svgpkg.Card(user.Login, char, "")
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
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}
