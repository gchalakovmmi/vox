package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/signup"
	m "vox/internal/middleware"
	up "vox/internal/userProfile"
	fa "vox/internal/form_autoload"
)

func Signup(page string, title string) http.Handler {
	return m.WithStatusMessage(func(statusMessage m.StatusMessage) http.Handler {
		return fa.WithFormAutoload(func(userProfile up.UserProfile) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				templ.Handler(signup.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
			})
		})
	})
}
