package handlers

import (
	up "vox/internal/userProfile"
	"log/slog"
	"net/http"
	"net/url"
	"vox/internal/middleware"
	"errors"
	"database/sql"
	"github.com/lib/pq"
)

func ProcessSignup() http.Handler {
	return middleware.WithDBConn(func(conn *sql.DB) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userProfile := up.UserProfile{
				FirstName:		r.PostFormValue("first_name"),
				LastName:		r.PostFormValue("last_name"),
				Email:			r.PostFormValue("email"),
				Password:		r.PostFormValue("password"),
				RepeatedPassword:	r.PostFormValue("password-repeat"),
			}

			queryParams := url.Values{}
			queryParams.Add("first_name", url.QueryEscape(userProfile.FirstName))
			queryParams.Add("last_name", url.QueryEscape(userProfile.LastName))
			queryParams.Add("email", url.QueryEscape(userProfile.Email))

			if errorMessage := userProfile.ValidateSignUpInfo(); errorMessage != "" {
				queryParams.Add("error", errorMessage)
				http.Redirect(w, r, "/signup?"+queryParams.Encode(), http.StatusSeeOther)
				return
			}

			slog.Debug("internal/handlers/process-signup.go", "Message", "SignUp info Correct", "Email", userProfile.Email, "Password", userProfile.Password, "Repeated Password", userProfile.RepeatedPassword)

			if _, err := userProfile.HashPassword(); err != nil {
				queryParams.Add("error", "sighup_hashing")
				http.Redirect(w, r, "/signup?"+queryParams.Encode(), http.StatusSeeOther)
				return
			}
			slog.Debug("internal/handlers/process-signup.go", "Message", "Hashing Successful", "Hashed Password", userProfile.HashedPassword)

			query := `
				INSERT INTO USERS(FIRST_NAME, LAST_NAME, EMAIL, PASSWORD_HASH)
				VALUES ($1, $2, $3, $4)
				RETURNING ID
			`
			err := conn.QueryRow(query, userProfile.FirstName, userProfile.LastName, userProfile.Email, userProfile.HashedPassword).Scan(&userProfile.ID)
			if err != nil {
				var pqErr *pq.Error
				if errors.As(err, &pqErr) && pqErr.Code == middleware.PGErrUniqueViolation {
					slog.Debug("internal/handlers/process-signup.go", "Message", "User with email exists", "Error", err)
					queryParams.Add("error", "signup_duplicate_email")
					http.Redirect(w, r, "/signup?"+queryParams.Encode(), http.StatusSeeOther)
					return
				} else {
					slog.Debug("internal/handlers/process-signup.go", "Message", "Database Error", "Error", err)
					queryParams.Add("error", "signup_insert_user")
					http.Redirect(w, r, "/signup?"+queryParams.Encode(), http.StatusSeeOther)
					return
				}
			}
			slog.Debug("internal/handlers/process-signup.go", "Message", "User created succesfully", "ID", userProfile.ID)
			queryParams.Add("info", "signin_user_created")
			http.Redirect(w, r, "/signin?"+queryParams.Encode(), http.StatusSeeOther)
			return
		})
	})
}
