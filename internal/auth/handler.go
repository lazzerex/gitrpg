package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/lazzerex/gitrpg/internal/config"
	"github.com/lazzerex/gitrpg/internal/users"
)

type Handler struct {
	cfg       *config.Config
	users     *users.Store
	logger    *slog.Logger
	postLogin func(*users.User)
}

func NewHandler(cfg *config.Config, store *users.Store, logger *slog.Logger) *Handler {
	return &Handler{cfg: cfg, users: store, logger: logger}
}

// SetPostLogin registers a callback invoked (in a goroutine) after successful login.
func (h *Handler) SetPostLogin(fn func(*users.User)) {
	h.postLogin = fn
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := randomHex(16)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	authURL := "https://github.com/login/oauth/authorize?" + url.Values{
		"client_id":    {h.cfg.GitHub.ClientID},
		"redirect_uri": {h.cfg.GitHub.CallbackURL},
		"scope":        {"read:user user:email repo"},
		"state":        {state},
	}.Encode()

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", MaxAge: -1, Path: "/"})

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	token, err := h.exchangeCode(r.Context(), code)
	if err != nil {
		h.logger.Error("token exchange failed", "error", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	ghUser, err := h.fetchGitHubUser(r.Context(), token)
	if err != nil {
		h.logger.Error("github user fetch failed", "error", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	user, err := h.users.Upsert(r.Context(), &users.User{
		GitHubID:    ghUser.ID,
		Login:       ghUser.Login,
		Name:        ghUser.Name,
		AvatarURL:   ghUser.AvatarURL,
		Email:       ghUser.Email,
		AccessToken: token,
	})
	if err != nil {
		h.logger.Error("user upsert failed", "error", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	secure := h.cfg.Server.Env == "production"
	if err := setSession(w, h.cfg.Session.Secret, user.ID, h.cfg.Session.MaxAge, secure); err != nil {
		h.logger.Error("session set failed", "error", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	if h.postLogin != nil {
		go h.postLogin(user)
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getSession(r, h.cfg.Session.Secret)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		user, err := h.users.GetByID(r.Context(), userID)
		if err != nil {
			clearSession(w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), users.ContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getSession(r, h.cfg.Session.Secret)
		if err == nil {
			if user, err := h.users.GetByID(r.Context(), userID); err == nil {
				ctx := context.WithValue(r.Context(), users.ContextKey, user)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

type githubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

func (h *Handler) exchangeCode(ctx context.Context, code string) (string, error) {
	body := url.Values{
		"client_id":     {h.cfg.GitHub.ClientID},
		"client_secret": {h.cfg.GitHub.ClientSecret},
		"code":          {code},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://github.com/login/oauth/access_token",
		strings.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", fmt.Errorf("github oauth error: %s", result.Error)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("empty access token")
	}
	return result.AccessToken, nil
}

func (h *Handler) fetchGitHubUser(ctx context.Context, token string) (*githubUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API %d: %s", resp.StatusCode, body)
	}

	var u githubUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
