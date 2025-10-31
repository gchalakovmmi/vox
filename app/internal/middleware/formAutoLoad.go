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
			FirstName: queryParams.Get("first_name"),
			LastName: queryParams.Get("last_name"),
			Email: queryParams.Get("email"),
		}
		slog.Debug("WithFormAutoload middleware", "FirstName", userProfile.FirstName, "LastName", userProfile.LastName, "Email", userProfile.Email)
		next(userProfile).ServeHTTP(w, r)
	})
}
