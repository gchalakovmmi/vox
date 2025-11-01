package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/signin"
	m "vox/internal/middleware"
	up "vox/internal/userProfile"
)

func Signin(page string, title string) http.Handler {
	return m.WithStatusMessage(func(statusMessage m.StatusMessage) http.Handler {
		return m.WithFormAutoload(func(userProfile up.UserProfile) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				templ.Handler(signin.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
			})
		})
	})
}
