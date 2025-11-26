package form_autoload

import (
	"net/http"
	"log/slog"
	up "vox/internal/userProfile"
	"net/url"
)

func WithFormAutoload(next func(up.UserProfile) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		userProfile := up.UserProfile{
			FirstName:	queryParams.Get("first_name"),
			LastName:	queryParams.Get("last_name"),
			Email:		queryParams.Get("email"),
		}

		unescapedEmail, err := url.QueryUnescape(userProfile.Email)
		if err != nil {
			unescapedEmail = userProfile.Email
		}
		slog.Debug("internal/form_autoload/form_autoload.go", "Message", "Extracted form values from url", "FirstName", userProfile.FirstName, "LastName", userProfile.LastName, "Email", unescapedEmail)
		next(userProfile).ServeHTTP(w, r)
	})
}
