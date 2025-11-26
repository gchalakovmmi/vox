package auth

import (
    "context"
    "database/sql"
    "net/http"

    "vox/internal/postgres"
    "vox/internal/config"
    up "vox/internal/userProfile"
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

func WithAuthentication(
    cfg *config.Config,
    next func(up.UserProfile) http.Handler,
) http.Handler {
    return postgres.WithDBConn(cfg, func(db *sql.DB) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            uid, ok := UIDFromContext(r.Context())
            if !ok || uid == 0 {
                http.Error(w, "missing uid in context", http.StatusBadRequest)
                return
            }

            ctx := context.Background()
            u := up.UserProfile{}
            if err := db.QueryRowContext(ctx,
                `SELECT id, first_name, last_name, email FROM users WHERE id=$1`, uid).
                Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
                if err == sql.ErrNoRows {
                    http.Error(w, "user not found", http.StatusNotFound)
                    return
                }
                http.Error(w, "db error", http.StatusInternalServerError)
                return
            }
            next(u).ServeHTTP(w, r)
        })
    })
}
