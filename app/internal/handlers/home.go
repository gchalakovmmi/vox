package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/home"
	m "vox/internal/middleware"
	up "vox/internal/userProfile"
)


// Temporairly to let me compile
func WithAuthentication(next func(up.UserProfile) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userProfile := up.UserProfile{}
		next(userProfile).ServeHTTP(w, r)
	})
}

func Home(page string, title string) http.Handler {
	return m.WithStatusMessage(func(statusMessage m.StatusMessage) http.Handler {
		return WithAuthentication(func(userProfile up.UserProfile) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				templ.Handler(home.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
			})
		})
	})
}
