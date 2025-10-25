package server

import (
	"net/http"
	"vox/internal/handlers"
)

type Route struct {
	Endpoint	string
	Handler		http.Handler
}

var Routes = []Route {
	{"/health", handlers.Health()},
	{"/signin", handlers.Signin("Sign In")},
	{"/signup", handlers.Signin("Sign Up")},
}
