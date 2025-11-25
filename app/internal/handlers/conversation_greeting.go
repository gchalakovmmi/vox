package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"vox/internal/ai"
	"vox/internal/config"
)

type greetingRsp struct {
	Text		string `json:"text"`
	AudioURL	string `json:"audioUrl"`
}

func ConversationGreeting(cfg *config.Config, store *ai.AudioStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		text := cfg.GetConversationDefaultGreeting()

		ctx := context.Background()
		audioMP3, err := ai.Synthesize(ctx, text,
			cfg.GetTTSOpenAIURL(), cfg.GetTTSOpenAIModel(), cfg.GetTTSOpenAIVoice(), cfg.GetTTSOpenAIAPIKey())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token := store.Put(audioMP3)

		resp := greetingRsp{
			Text:	 text,
			AudioURL: "/conversation/audio?token=" + token,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}
