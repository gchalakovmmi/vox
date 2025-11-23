package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"vox/internal/config"
)

// RequireAccessToken checks header || cookie; if the access token is expired
// or missing it rotates both tokens via the refresh cookie and continues.
func RequireAccessToken(cfg *config.Config, rdb *redis.Client, next func(http.ResponseWriter, *http.Request, uint)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := bearerOrCookie(r)
		// access token missing – try refresh
		if raw == "" { 
			c, err := r.Cookie("refresh_token")
			if err != nil {
				slog.Debug("internal/auth/jwt_middleware.go", "Message", "No access token and no refresh cookie")
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			slog.Debug("internal/auth/jwt_middleware.go", "Message", "About to call tryRefresh", "refresh_cookie", c.Value[:20]+"...")
			newAT, newRT, ok := tryRefresh(c.Value, cfg, rdb)
			if !ok {
				slog.Debug("internal/auth/jwt_middleware.go", "Message", "Refresh failed")
		http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			// set new cookies and continue with new access token
			setTokenCookies(w, newAT, newRT, cfg)
			raw = newAT
			// re-extract user id
			sub, _ := ValidateAndExtractSub(newAT, cfg.GetAccessTokenSecret())
			next(w, r, sub)
			return
		}

		// access token present – validate it
		sub, err := ValidateAndExtractSub(raw, cfg.GetAccessTokenSecret())
		// access token expired/invalid – same refresh path
		if err != nil {
			c, _ := r.Cookie("refresh_token")
			if c == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			newAT, newRT, ok := tryRefresh(c.Value, cfg, rdb)
			if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			setTokenCookies(w, newAT, newRT, cfg)
			raw = newAT
			sub, _ = ValidateAndExtractSub(newAT, cfg.GetAccessTokenSecret())
		}
		next(w, r, sub)
	}
}

func tryRefresh(rt string, cfg *config.Config, rdb *redis.Client) (string, string, bool) {
	sub, err := ValidateAndExtractSub(rt, cfg.GetRefreshTokenSecret())
	if err != nil {
		slog.Debug("internal/auth/jwt_middleware.go", "Message", "Refresh JWT invalid", "error", err)
		return "", "", false
	}
	ctx := context.Background()
	key := fmt.Sprintf("rt:%d", sub)

	stored, err := rdb.Get(ctx, key).Result()

	slog.Debug("internal/auth/jwt_middleware.go", "Message", "tryRefresh", "key", key, "stored", stored, "err", err)

	if err != nil || stored != rt {
		return "", "", false
	}

	newAT, _ := GenerateAccessToken(cfg, sub)
	newRT, _ := GenerateRefreshToken(cfg, sub)
	if err := rdb.Set(ctx, key, newRT, cfg.GetRefreshTokenTTL()).Err(); err != nil {
		return "", "", false
	}
	return newAT, newRT, true
}

func setTokenCookies(w http.ResponseWriter, at, rt string, cfg *config.Config) {
	http.SetCookie(w, &http.Cookie{
		Name:	 "access_token",
		Value:	at,
		Path:	 "/",
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: cfg.CookieSameSite,
		MaxAge:   int(cfg.GetAccessTokenTTL().Seconds()),
	})
	http.SetCookie(w, &http.Cookie{
		Name:	 "refresh_token",
		Value:	rt,
		Path:	 "/",
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: cfg.CookieSameSite,
		MaxAge:   int(cfg.GetRefreshTokenTTL().Seconds()),
	})
}

func bearerOrCookie(r *http.Request) string {
	for _, c := range r.Cookies() {
		slog.Debug("internal/middleware/statusMessage.go", "Message", "bearerOrCookie", "name", c.Name, "value", c.Value[:10]+"...")
	}

	if h := r.Header.Get("Authorization"); h != "" {
		p := strings.SplitN(h, " ", 2)
		if len(p) == 2 && p[0] == "Bearer" {
			return p[1]
		}
	}
	if c, err := r.Cookie("access_token"); err == nil {
		return c.Value
	}
	return ""
}
