package handlers

import (
	up "vox/internal/userProfile"
	"log/slog"
	"net/http"
	"net/url"
	"vox/internal/postgres"
	"database/sql"
)

func ProcessSignin() http.Handler {
	return postgres.WithDBConn(func(conn *sql.DB) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userProfile := up.UserProfile{
				Email:			r.PostFormValue("email"),
				Password:		r.PostFormValue("password"),
			}

			queryParams := url.Values{}
			queryParams.Add("email", url.QueryEscape(userProfile.Email))

			err := conn.QueryRow(`SELECT PASSWORD_HASH FROM USERS WHERE EMAIL = $1`, userProfile.Email).Scan(&userProfile.HashedPassword)
			if err != nil {
				slog.Debug("internal/handlers/process-signin.go", "Message", "Database Error", "Error", err)
				queryParams.Add("error", "signin_wrong_credentials")
				http.Redirect(w, r, "/signin?"+queryParams.Encode(), http.StatusSeeOther)
				return
			}
			slog.Debug("internal/handlers/process-signin.go", "Message", "Hashed password taken from DB", "Hashed Password", userProfile.HashedPassword)

			if userProfile.CompareHashAndPassword() != nil {
				queryParams.Add("error", "signin_wrong_credentials")
				http.Redirect(w, r, "/signin?"+queryParams.Encode(), http.StatusSeeOther)
				return
			}
			queryParams = url.Values{}
			http.Redirect(w, r, "/home"), http.StatusSeeOther)
			return
		})
	})
}
