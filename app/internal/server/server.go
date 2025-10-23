package server

import (
	"net/http"
	"vox/internal/config"
	"vox/internal/handlers"
)

type Server struct {
	config	*config.Config
	// db	*sql.DB
	// routes	[]Route
}

func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (srv *Server) CreateHandlers() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static/"))))

	mux.Handle("/health", handlers.Health())
	mux.Handle("/signin", handlers.Signin("signin"))
	
	return mux
}
