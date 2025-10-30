package middleware

import (
	"net/http"
	"log/slog"
	up "vox/internal/userProfile"
)

func WithFormAutoload(next func(up.UserProfile) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		userProfile := up.UserProfile{
			Email: queryParams.Get("email"),
		}
		slog.Debug("WithFormAutoload middleware", "Email", userProfile.Email)
		next(userProfile).ServeHTTP(w, r)
	})
}
