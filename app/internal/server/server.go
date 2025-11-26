package server

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"vox/internal/auth"
	"vox/internal/config"
	"vox/internal/handlers"
	"vox/internal/ai"
	"context"
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
	mux.Handle("/home", auth.RequireAccessToken(s.cfg, s.rdb,
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			ctx := context.WithValue(r.Context(), auth.CtxKeyUID, uid)
			handlers.Home("/home", "Home", s.cfg).ServeHTTP(w, r.WithContext(ctx))
	}))
	mux.Handle("/conversation", auth.RequireAccessToken(s.cfg, s.rdb, 
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			ctx := context.WithValue(r.Context(), auth.CtxKeyUID, uid)
			handlers.Conversation("/conversation", "Conversation", s.cfg).ServeHTTP(w, r.WithContext(ctx))
	}))

	audioStore := ai.NewAudioStore(s.cfg.GetTTSAudioTimeToLive())
	mux.Handle("/conversation/greeting", auth.RequireAccessToken(s.cfg, s.rdb,
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			handlers.ConversationGreeting(s.cfg, audioStore).ServeHTTP(w, r)
	}))
	mux.Handle("/conversation/prompt", auth.RequireAccessToken(s.cfg, s.rdb,
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			handlers.ConversationPrompt(s.cfg, audioStore, uid).ServeHTTP(w, r)
	}))
	mux.Handle("/conversation/audio", auth.RequireAccessToken(s.cfg, s.rdb,
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			handlers.ConversationAudioTokenHandler(audioStore).ServeHTTP(w, r)
	}))

	mux.Handle("/conversation/analysis", auth.RequireAccessToken(s.cfg, s.rdb, 
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			ctx := context.WithValue(r.Context(), auth.CtxKeyUID, uid)
			handlers.ConversationAnalysis("/conversation/analysis", "Conversation Analysis", s.cfg, uid).ServeHTTP(w, r.WithContext(ctx))
	}))

	mux.Handle("/conversation/analysis/feedback", auth.RequireAccessToken(s.cfg, s.rdb,
		func(w http.ResponseWriter, r *http.Request, uid uint) {
			ctx := context.WithValue(r.Context(), auth.CtxKeyUID, uid)
			handlers.ConversationAnalysisFeedback(uid, s.cfg).ServeHTTP(w, r.WithContext(ctx))
	}))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}))


	return mux
}
