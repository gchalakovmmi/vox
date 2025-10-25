package server

import (
	"net/http"
	"vox/internal/config"
)

type Server struct {
	config	*config.Config
	// db	*sql.DB
	routes	[]Route
}

func New(cfg *config.Config, /*db *sqlSomething,*/ routes []Route) *Server {
	return &Server{
		config: cfg,
		routes:	routes,
	}
}

func (srv *Server) CreateHandlers() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static/"))))

	for _, route := range srv.routes {
		mux.Handle(route.Endpoint, route.Handler)
	}
	
	return mux
}
