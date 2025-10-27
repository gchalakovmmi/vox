package handlers

import (
	"log/slog"
	"net/http"
)

func ProcessSignup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		passwordRepeat := r.PostFormValue("password-repeat")

		slog.Debug("Form Values", "Email", email, "Password", password, "Repeated Password", passwordRepeat)
	})
}

