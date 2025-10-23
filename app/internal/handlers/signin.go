package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	"vox/web/templates/pages/signin"
)

func Signin(page string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(signin.Handler(page)).ServeHTTP(w, r)
	})
}

