package auth

import (
	"net/http"
	up "vox/internal/userProfile"
	"vox/internal/config"
)

func WithoutAuthentication(cfg *config.Config, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if c, _ := r.Cookie("access_token"); c != nil {
            if _, err := ValidateAndExtractSub(c.Value, cfg.GetAccessTokenSecret()); err == nil {
                http.Redirect(w, r, "/", http.StatusSeeOther)
                return
            }
        }
        next.ServeHTTP(w, r)
    })
}

func WithAuthentication(next func(up.UserProfile) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userProfile := up.UserProfile{}
		next(userProfile).ServeHTTP(w, r)
	})
}
