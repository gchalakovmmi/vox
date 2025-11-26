package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"vox/internal/ai"
	"vox/internal/config"
	"vox/internal/postgres"
)

// ConversationAnalysisFeedback returns teacher-style feedback for the whole conversation.
// GET /conversation/analysis/feedback?id=<conversation-id>
func ConversationAnalysisFeedback(uid uint, cfg *config.Config) http.Handler {
	return postgres.WithDBConn(cfg, func(db *sql.DB) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			convID, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil || convID <= 0 {
				http.Error(w, "invalid conversation id", http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			// 1. load whole history (user + assistant only)
			hist := ai.NewChatHistory(db, uid, convID)
			jsonHist, err := hist.Load(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			messages, _ := ai.MessagesFromJSON(jsonHist)

			// 2. build the prompt
			var prompt string
			for _, m := range messages {
				prompt += m.Role + ": " + m.Content + "\n"
			}
			fullPrompt := cfg.GetConversationAnalysisTeacherPrompt() + "\n\n" + prompt

			// 3. ask the LLM
			reply, err := ai.ChatCompletion(ctx,
				[]ai.ChatMessage{{Role: "user", Content: fullPrompt}},
				cfg.GetLLMOpenAIURL(), cfg.GetLLMOpenAIAPIKey(), cfg.GetLLMOpenAIModelName())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(reply.Content))
		})
	})
}
