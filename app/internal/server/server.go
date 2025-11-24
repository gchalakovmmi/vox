package server

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"vox/internal/auth"
	"vox/internal/config"
	"vox/internal/handlers"
)

type Server struct {
	cfg	*config.Config
	rdb	*redis.Client
}

func New(cfg *config.Config, rdb *redis.Client) *Server {
	return &Server{cfg: cfg, rdb: rdb}
}

func (s *Server) CreateHandlers() http.Handler {
	mux := http.NewServeMux()

	// static
	mux.Handle("/web/static/",	http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static"))))

	// public
	mux.Handle("/health",		handlers.Health())
	mux.Handle("/signin",		auth.WithoutAuthentication(s.cfg, handlers.Signin("signin", "Sign In")))
	mux.Handle("/process-signin",	handlers.ProcessSignin(s.cfg, s.rdb))
	mux.Handle("/signup",		auth.WithoutAuthentication(s.cfg, handlers.Signup("signup", "Sign Up")))

	mux.Handle("/process-signup",	handlers.ProcessSignup(s.cfg))

	// auth
	mux.Handle("/auth/refresh",	auth.TokenRefreshHandler(s.cfg, s.rdb))

	// protected
	mux.Handle("/home", auth.RequireAccessToken(s.cfg, s.rdb, func(w http.ResponseWriter, r *http.Request, uid uint) {
		handlers.Home("home", "Home").ServeHTTP(w, r)
	}))
	mux.Handle("/conversation", auth.RequireAccessToken(s.cfg, s.rdb, func(w http.ResponseWriter, r *http.Request, uid uint) {
		handlers.Conversation("conversation", "Conversation").ServeHTTP(w, r)
	}))
	mux.Handle("/conversation/reply", auth.RequireAccessToken(s.cfg, s.rdb, func(w http.ResponseWriter, r *http.Request, uid uint) {
		handlers.ConversationReply(s.cfg).ServeHTTP(w, r)
	}))
	mux.Handle("/conversation/tts", auth.RequireAccessToken(s.cfg, s.rdb,
	func(w http.ResponseWriter, r *http.Request, uid uint) {
		handlers.ConversationTTS().ServeHTTP(w, r)
	}))

	mux.Handle("/",			handlers.Home("home", "Home"))

	return mux
}
