package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const cookieName = "session"

type sessionPayload struct {
	UserID int64 `json:"u"`
	Exp    int64 `json:"e"`
}

func setSession(w http.ResponseWriter, secret string, userID int64, maxAge int, secure bool) error {
	p := sessionPayload{
		UserID: userID,
		Exp:    time.Now().Add(time.Duration(maxAge) * time.Second).Unix(),
	}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	encoded := base64.RawURLEncoding.EncodeToString(data)
	value := encoded + "." + sign(encoded, secret)

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    value,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	return nil
}

func getSession(r *http.Request, secret string) (int64, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return 0, err
	}

	parts := strings.SplitN(cookie.Value, ".", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("malformed session")
	}
	encoded, sig := parts[0], parts[1]

	if !hmac.Equal([]byte(sign(encoded, secret)), []byte(sig)) {
		return 0, fmt.Errorf("invalid session signature")
	}

	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return 0, err
	}

	var p sessionPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return 0, err
	}

	if time.Now().Unix() > p.Exp {
		return 0, fmt.Errorf("session expired")
	}

	return p.UserID, nil
}

func clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

func sign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
