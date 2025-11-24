package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"log/slog"

	"vox/internal/ai"
	"vox/internal/auth"
	"vox/internal/config"
	"github.com/redis/go-redis/v9"
)

type promptRsp struct {
	User		string	`json:"user"`		// what the user said
	Vox		string	`json:"vox"`		// what the assistant answered
	AudioURL	string	`json:"audioUrl"`	// url to stream mp3
}

// PostConversationPrompt  POST /conversation/prompt  (multipart)
func PostConversationPrompt(cfg *config.Config, rdb *redis.Client, store *ai.AudioStore) http.Handler {
	return auth.RequireAccessToken(cfg, rdb, func(w http.ResponseWriter, r *http.Request, uid uint) {
		slog.Debug("PostConversationPrompt", "uid", uid)

		// 1. enforce size limit
		r.Body = http.MaxBytesReader(w, r.Body, 30<<20)
		if err := r.ParseMultipartForm(30 << 20); err != nil {
			slog.Debug("PostConversationPrompt", "parse error", err)
			http.Error(w, "file too big or invalid", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("audio")
		if err != nil {
			slog.Debug("PostConversationPrompt", "missing audio", err)
			http.Error(w, "missing audio", http.StatusBadRequest)
			return
		}
		defer file.Close()
		audioBytes, _ := io.ReadAll(file)

		ctx := context.Background()

		// 2. STT
		userText, err := ai.Transcribe(ctx, audioBytes,
			cfg.GetSTTOpenAIURL(), cfg.GetSTTOpenAIModelName(), cfg.GetSTTOpenAIAPIKey())
		if err != nil {
			slog.Debug("PostConversationPrompt", "stt error", err)
			http.Error(w, "transcription failed", http.StatusInternalServerError)
			return
		}

		// 3. load history
		histStore := ai.NewHistoryStore(rdb, 24*time.Hour)
		rawHist, _ := histStore.Load(ctx, uid)
		messages, _ := ai.MessagesFromJSON(rawHist)
		messages = append(messages, ai.ChatMessage{Role: "user", Content: userText})

		// 4. LLM
		assistant, err := ai.ChatCompletion(ctx, messages,
			cfg.GetLLMOpenAIURL(), cfg.GetLLMOpenAIAPIKey(), cfg.GetLLMOpenAIModelName())
		if err != nil {
			slog.Debug("PostConversationPrompt", "llm error", err)
			http.Error(w, "llm error", http.StatusInternalServerError)
			return
		}
		messages = append(messages, assistant)

		// 5. save history
		raw, _ := ai.MessagesToJSON(messages)
		_ = histStore.Save(ctx, uid, raw)

		// 6. TTS
		audioMP3, err := ai.Synthesize(ctx, assistant.Content,
			cfg.GetTTSOpenAIURL(), cfg.GetTTSOpenAIModel(), cfg.GetTTSOpenAIVoice(), cfg.GetTTSOpenAIAPIKey())
		if err != nil {
			slog.Debug("PostConversationPrompt", "tts error", err)
			http.Error(w, "tts error", http.StatusInternalServerError)
			return
		}
		token := store.Put(audioMP3)

		resp := promptRsp{
			User:		userText,
			Vox:		assistant.Content,
			AudioURL:	"/conversation/audio?token=" + token,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}
