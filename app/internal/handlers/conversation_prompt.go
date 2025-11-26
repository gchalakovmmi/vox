package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"vox/internal/ai"
	"vox/internal/config"
	"vox/internal/postgres"
	"strconv"
)

type promptRsp struct {
	User		string	`json:"user"`
	Vox		string	`json:"vox"`
	Correction	string	`json:"correction"`
	AudioURL	string	`json:"audioUrl"`
	ConversationID	int	`json:"conversationID"`
}

func ConversationPrompt(cfg *config.Config, store *ai.AudioStore, uid uint) http.Handler {
	return postgres.WithDBConn(cfg, func(db *sql.DB) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. size check
			r.Body = http.MaxBytesReader(w, r.Body, 30<<20)
			if err := r.ParseMultipartForm(30 << 20); err != nil {
				http.Error(w, "file too big or invalid", http.StatusBadRequest)
				return
			}
			file, _, err := r.FormFile("audio")
			if err != nil {
				http.Error(w, "missing audio", http.StatusBadRequest)
				return
			}
			defer file.Close()
			audioBytes, _ := io.ReadAll(file)

			ctx := r.Context()

			// 2. STT
			userText, err := ai.Transcribe(ctx, audioBytes,
				cfg.GetSTTOpenAIURL(), cfg.GetSTTOpenAIModelName(), cfg.GetSTTOpenAIAPIKey())
			if err != nil {
				http.Error(w, "transcription failed", http.StatusInternalServerError)
				return
			}

			// 3. load history
			conversation_id, err := strconv.Atoi(r.URL.Query().Get("conversationID"))
			if err != nil {
				http.Error(w, "Error converting requestCount to int", http.StatusInternalServerError)
				return
			}

			chat_history := ai.NewChatHistory(db, uid, conversation_id)

			if chat_history.ConversationID == -1 {
				chat_history.ConversationID, err = chat_history.CreateConversation(ctx, cfg)
				if err != nil {
					http.Error(w, "create conversation error", http.StatusInternalServerError)
					return
				}
			}
			rawHist, _ := chat_history.Load(ctx)
			messages, _ := ai.MessagesFromJSON(rawHist)
			messages = append(messages, ai.ChatMessage{Role: "user", Content: userText})

			// 4. Correction
			correction_system_prompt := ai.ChatMessage{Role: "system", Content: cfg.GetConversationCorrectorPrompt()}
			correctionMessages := append([]ai.ChatMessage{correction_system_prompt}, messages...)
			correctionMessages = append(correctionMessages, ai.ChatMessage{Role: "user", Content: userText})

			correction, err := ai.ChatCompletion(ctx, correctionMessages, cfg.GetLLMOpenAIURL(), cfg.GetLLMOpenAIAPIKey(), cfg.GetLLMOpenAIModelName())
			if err != nil {
				http.Error(w, "correction error", http.StatusInternalServerError)
				return
			}

			// 4. LLM
			system_prompt := ai.ChatMessage{Role: "system", Content: cfg.GetConversationMatePrompt()}
			llmMessages := append([]ai.ChatMessage{system_prompt}, messages...)
			llmMessages = append(llmMessages, ai.ChatMessage{Role: "user", Content: userText})

			assistant, err := ai.ChatCompletion(ctx, llmMessages, cfg.GetLLMOpenAIURL(), cfg.GetLLMOpenAIAPIKey(), cfg.GetLLMOpenAIModelName())
			if err != nil {
				http.Error(w, "llm error", http.StatusInternalServerError)
				return
			}
			messages = append(messages, assistant)

			// 5. save history
			newMsgs, _ := ai.MessagesToJSON([]ai.ChatMessage{
				{Role: "user", Content: userText},
				assistant,
			})
			_ = chat_history.Save(ctx, cfg, newMsgs)

			// 6. TTS
			audioMP3, err := ai.Synthesize(ctx, assistant.Content,
				cfg.GetTTSOpenAIURL(), cfg.GetTTSOpenAIModel(), cfg.GetTTSOpenAIVoice(), cfg.GetTTSOpenAIAPIKey())
			if err != nil {
				http.Error(w, "tts error", http.StatusInternalServerError)
				return
			}
			token := store.Put(audioMP3)

			resp := promptRsp{
				User:			userText,
				Vox:			assistant.Content,
				Correction:		correction.Content,
				AudioURL:		"/conversation/audio?token=" + token,
				ConversationID:		chat_history.ConversationID,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})
	})
}
