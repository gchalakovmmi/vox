package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"fmt"

	"log/slog"
	"vox/internal/auth"
	"vox/internal/config"
	"vox/internal/postgres"
	up "vox/internal/userProfile"
	"github.com/redis/go-redis/v9"
)

func ProcessSignin(cfg *config.Config, rdb *redis.Client) http.Handler {
	return postgres.WithDBConn(cfg, func(conn *sql.DB) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userProfile := up.UserProfile{
				Email:		r.PostFormValue("email"),
				Password:	r.PostFormValue("password"),
			}

			queryParams := url.Values{}
			queryParams.Add("email", url.QueryEscape(userProfile.Email))

			if err := conn.QueryRow(`SELECT id, password_hash FROM users WHERE email=$1`, userProfile.Email).
				Scan(&userProfile.ID, &userProfile.HashedPassword); err != nil {
				slog.Debug("internal/handlers/process-signin.go" "Message", "Could not execute DB query", "error", err)
				queryParams.Add("error", "signin_wrong_credentials")
				http.Redirect(w, r, "/signin?"+queryParams.Encode(), http.StatusSeeOther)
				return
			}
			if err := userProfile.CompareHashAndPassword(); err != nil {
				queryParams.Add("error", "signin_wrong_credentials")
				http.Redirect(w, r, "/signin?"+queryParams.Encode(), http.StatusSeeOther)
				return
			}

			accessToken, _ := auth.GenerateAccessToken(cfg, uint(userProfile.ID))
			refreshToken, _ := auth.GenerateRefreshToken(cfg, uint(userProfile.ID))

			key := fmt.Sprintf("rt:%d", userProfile.ID)
			rdb.Set(context.Background(), key, refreshToken, cfg.GetRefreshTokenTTL())

			http.SetCookie(w, &http.Cookie{
				Name:		"access_token",
				Value:		accessToken,
				Path:		"/",
				HttpOnly:	true,
				Secure:		false,
				SameSite:	http.SameSiteLaxMode,
				MaxAge:		int(cfg.GetAccessTokenTTL().Seconds()),
			})
			http.SetCookie(w, &http.Cookie{
				Name:		"refresh_token",
				Value:		refreshToken,
				Path:		"/",
				HttpOnly:	true,
				Secure:		false,
				SameSite:	http.SameSiteLaxMode,
				MaxAge:		int(cfg.GetRefreshTokenTTL().Seconds()),
			})

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		})
	})
}
