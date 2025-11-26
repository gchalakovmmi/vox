package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/conversation"
	up "vox/internal/userProfile"
	sm "vox/internal/status_message"
	"vox/internal/auth"
	"vox/internal/config"
)

func Conversation(page string, title string, cfg *config.Config) http.Handler {
	return sm.WithStatusMessage(func(statusMessage sm.StatusMessage) http.Handler {
		return auth.WithAuthentication(cfg, func(userProfile up.UserProfile) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				templ.Handler(conversation.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
			})
		})
	})
}
