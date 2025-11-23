package server

import (
	"net/http"
	"vox/internal/config"
	"vox/internal/handlers"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"vox/internal/auth"
)

type Server struct {
	config	*config.Config
	db     *sql.DB
	rdb    *redis.Client
}

type Route struct {
	Endpoint	string
	Handler		http.Handler
}

func New(cfg *config.Config, db *sql.DB, rdb *redis.Client) *Server {
    return &Server{config: cfg, db: db, rdb: rdb}
}

func (srv *Server) CreateHandlers() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static/"))))

	mux.Handle("/health",		handlers.Health())
	mux.Handle("/signin",		handlers.Signin("signin", "Sign In"))
	mux.Handle("/process-signin",	handlers.ProcessSignin(srv.config))
	mux.Handle("/signup",		handlers.Signup("signup", "Sign Up"))
	mux.Handle("/process-signup",	handlers.ProcessSignup(srv.config))
	mux.Handle("/home",		auth.MiddlewareAccessToken(srv.rdb)(handlers.Home("home", "Home")))

	auth.RegisterAuthRoutes(mux, srv.db, srv.rdb)
	mux.Handle("/auth/", http.StripPrefix("/auth", mux))

	mux.Handle("/",			handlers.Home("home", "Home"))

	return mux
}
