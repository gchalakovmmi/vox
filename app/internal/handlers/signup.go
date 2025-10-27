package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/signup"
)

func Signup(page string, title string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(signup.Handler(page, title)).ServeHTTP(w, r)
	})
}

