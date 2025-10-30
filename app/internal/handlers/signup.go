package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/signup"
	"log/slog"
	sm "vox/internal/statusMessage"
	up "vox/internal/userProfile"
)

func WithStatusMessage(next func(sm.StatusMessage) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		statusMessage := sm.StatusMessage{
			Info:  queryParams.Get("info"),
			Err: queryParams.Get("error"),
		}

		slog.Debug("WithStatusMessage middleware", "statusMessage", statusMessage)
		next(statusMessage).ServeHTTP(w, r)
	})
}

func WithFormAutoload(next func(up.UserProfile) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		userProfile := up.UserProfile{
			Email: queryParams.Get("Email"),
		}
		slog.Debug("WithFormAutoload middleware", "Email", userProfile.Email)
		next(userProfile).ServeHTTP(w, r)
	})
}

func Signup(page string, title string) http.Handler {
	return WithStatusMessage(func(statusMessage sm.StatusMessage) http.Handler {
		return WithFormAutoload(func(userProfile up.UserProfile) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				templ.Handler(signup.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
			})
		})
	})
}
