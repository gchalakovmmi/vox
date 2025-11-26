package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/signup"
	up "vox/internal/userProfile"
	fa "vox/internal/form_autoload"
	sm "vox/internal/status_message"
)

func Signup(page string, title string) http.Handler {
	return sm.WithStatusMessage(func(statusMessage sm.StatusMessage) http.Handler {
		return fa.WithFormAutoload(func(userProfile up.UserProfile) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				templ.Handler(signup.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
			})
		})
	})
}
