package ai

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type HistoryStore struct {
	rdb	*redis.Client
	ttl	time.Duration
	keyPref	string
}

func NewHistoryStore(rdb *redis.Client, ttl time.Duration) *HistoryStore {
	return &HistoryStore{rdb: rdb, ttl: ttl, keyPref: "chat"}
}

func (h *HistoryStore) Load(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("%s:%d", h.keyPref, userID)
	return h.rdb.Get(ctx, key).Result()
}

func (h *HistoryStore) Save(ctx context.Context, userID uint, raw string) error {
	key := fmt.Sprintf("%s:%d", h.keyPref, userID)
	return h.rdb.Set(ctx, key, raw, h.ttl).Err()
}
