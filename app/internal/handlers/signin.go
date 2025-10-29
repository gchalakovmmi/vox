package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/signin"
	"vox/internal/status"
)

func Signin(page string, title string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(signin.Handler(page, title, status.StatusMessage{"TMP", "TMP"})).ServeHTTP(w, r)
	})
}

