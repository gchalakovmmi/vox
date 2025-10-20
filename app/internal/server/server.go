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
	mux.Handle("/health", handlers.Health())
	
	return mux
}
