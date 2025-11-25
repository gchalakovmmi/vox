package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	ca "vox/web/templates/pages/conversation_analysis"
	up "vox/internal/userProfile"
	sm "vox/internal/status_message"
	"vox/internal/auth"
)

func ConversationAnalysis(page string, title string) http.Handler {
	return sm.WithStatusMessage(func(statusMessage sm.StatusMessage) http.Handler {
		return auth.WithAuthentication(func(userProfile up.UserProfile) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				templ.Handler(ca.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
			})
		})
	})
}
