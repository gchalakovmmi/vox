package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"vox/internal/config"
)

type refreshReq struct {
	RefreshToken	string	`json:"refresh_token"`
}

type refreshRsp struct {
	AccessToken	string	`json:"access_token"`
	RefreshToken	string	`json:"refresh_token"`
}

// TokenRefreshHandler returns a new token pair and rotates the refresh token in Redis.
func TokenRefreshHandler(cfg *config.Config, rdb *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req refreshReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}

		sub, err := ValidateAndExtractSub(req.RefreshToken, cfg.GetRefreshTokenSecret())
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		// ---- Redis rotation ----
		ctx := context.Background()
		oldKey := fmt.Sprintf("rt:%d:%s", sub, req.RefreshToken)
		if err := rdb.Get(ctx, oldKey).Err(); err != nil {
			http.Error(w, "refresh token not found", http.StatusUnauthorized)
			return
		}

		newAT, _ := GenerateAccessToken(cfg, sub)
		newRT, _ := GenerateRefreshToken(cfg, sub)
		newKey := fmt.Sprintf("rt:%d:%s", sub, newRT)

		pipe := rdb.Pipeline()
		pipe.Set(ctx, newKey, "1", 30*24*time.Hour)
		pipe.Del(ctx, oldKey)
		if _, err := pipe.Exec(ctx); err != nil {
			http.Error(w, "redis error", http.StatusInternalServerError)
			return
		}
		// ---- end rotation ----

		resp := refreshRsp{AccessToken: newAT, RefreshToken: newRT}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}
