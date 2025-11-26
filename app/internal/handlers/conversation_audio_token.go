package handlers

import (
	"net/http"
	"strconv"

	"vox/internal/ai"
)

// ConversationAudioTokenHandler  GET /conversation/audio?token=...
// Streams the MP3 bytes associated with the token.
func ConversationAudioTokenHandler(store *ai.AudioStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		audio, ok := store.Get(token)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(audio)))
		w.Header().Set("Cache-Control", "no-store")
		w.Write(audio)
	})
}
