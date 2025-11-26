package ai

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// AudioStore keeps MP3 blobs for a short time and returns a token to retrieve them.
type AudioStore struct {
	mu	sync.RWMutex
	blobs	map[string][]byte
	ttl	time.Duration
}

func NewAudioStore(ttl time.Duration) *AudioStore {
	return &AudioStore{blobs: make(map[string][]byte), ttl: ttl}
}

func (a *AudioStore) Put(blob []byte) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	token := randomToken()
	a.blobs[token] = blob
	time.AfterFunc(a.ttl, func() {
		a.mu.Lock()
		delete(a.blobs, token)
		a.mu.Unlock()
	})
	return token
}

func (a *AudioStore) Get(token string) ([]byte, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	b, ok := a.blobs[token]
	return b, ok
}

func randomToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
