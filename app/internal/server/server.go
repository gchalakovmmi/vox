package server

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"vox/internal/auth"
	"vox/internal/config"
	"vox/internal/handlers"
	"vox/internal/ai"
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
	mux.Handle("/conversation", auth.RequireAccessToken(s.cfg, s.rdb, 
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			handlers.Conversation("conversation", "Conversation").ServeHTTP(w, r)
	}))
	audioStore := ai.NewAudioStore(60 * time.Second)

	mux.Handle("/conversation/greeting", auth.RequireAccessToken(s.cfg, s.rdb,
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			handlers.PostConversationGreeting(s.cfg, audioStore).ServeHTTP(w, r)
	}))
	mux.Handle("/conversation/prompt", auth.RequireAccessToken(s.cfg, s.rdb,
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			handlers.PostConversationPrompt(s.cfg, s.rdb, audioStore).ServeHTTP(w, r)
	}))
	mux.Handle("/conversation/audio", auth.RequireAccessToken(s.cfg, s.rdb,
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			handlers.ConversationAudioTokenHandler(audioStore).ServeHTTP(w, r)
	}))


	mux.Handle("/",			handlers.Home("home", "Home"))

	return mux
}
