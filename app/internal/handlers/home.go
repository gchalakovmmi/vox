package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/home"
	up "vox/internal/userProfile"
	sm "vox/internal/status_message"
	"vox/internal/auth"
)

func Home(page string, title string) http.Handler {
	return sm.WithStatusMessage(func(statusMessage sm.StatusMessage) http.Handler {
		return auth.WithAuthentication(func(userProfile up.UserProfile) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				templ.Handler(home.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
			})
		})
	})
}
