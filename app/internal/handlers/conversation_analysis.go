package handlers

import (
	"net/http"
	"github.com/a-h/templ"
	ca "vox/web/templates/pages/conversation_analysis"
	up "vox/internal/userProfile"
	sm "vox/internal/status_message"
	"vox/internal/auth"
	"vox/internal/postgres"
	"vox/internal/ai"
	"vox/internal/config"
	"strconv"
	"database/sql"
)

func ConversationAnalysis(page, title string, cfg *config.Config, uid uint) http.Handler {
    return sm.WithStatusMessage(func(statusMessage sm.StatusMessage) http.Handler {
        return auth.WithAuthentication(cfg, func(userProfile up.UserProfile) http.Handler {
            return postgres.WithDBConn(cfg, func(db *sql.DB) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    convID, err := strconv.Atoi(r.URL.Query().Get("id"))
                    if err != nil || convID <= 0 {
                        http.Error(w, "missing or invalid conversation id", http.StatusBadRequest)
                        return
                    }

                    hist := ai.NewChatHistory(db, uid, convID)
                    jsonHistory, err := hist.Load(r.Context())
                    if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                    }
                    messages, _ := ai.MessagesFromJSON(jsonHistory)

                    templ.Handler(ca.Handler(page, title, statusMessage, &userProfile, messages)).
                        ServeHTTP(w, r)
                })
            })
        })
    })
}
