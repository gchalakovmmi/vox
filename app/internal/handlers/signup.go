package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/signup"
	"log/slog"
	sm "vox/internal/statusMessage"
	up "vox/internal/userProfile"
)

func Signup(page string, title string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// IGNORE THIS IS FOR LATER!!! Take this out into a handler function named WithStatusMessage. Put the StatusMessage type there as well. What should the file be named? Leave it in the handlers package.
		queryParams := r.URL.Query()
		statusMessage := sm.StatusMessage{ queryParams.Get("info"), queryParams.Get("error") }
		slog.Debug("internal/handlers/signup.go", "StatusMessage", statusMessage)

		userProfile := up.UserProfile{
			Email:	queryParams.Get("Email"),
		}
		slog.Debug("Extracted form values from url", "Email", userProfile.Email)
		templ.Handler(signup.Handler(page, title, statusMessage, &userProfile)).ServeHTTP(w, r)
	})
}
