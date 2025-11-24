package handlers

import (
	"strconv"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"vox/internal/config"
	"vox/internal/tts"
)

type replyReq struct {
	Messages []message `json:"messages"`
}
type message struct {
	Role	string	`json:"role"`
	Content	string	`json:"content"`
}

type replyRsp struct {
	Text		string	`json:"text"`
	AudioURL	string	`json:"audioUrl"`
}

var (
	audioMu sync.RWMutex
	audioMap  = map[string][]byte{}
)

func ConversationReply(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req replyReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		text := cfg.GetConversationDefaultGreeting()
		if len(req.Messages) != 0 {
			// TODO: call LLM here instead of echo
			last := req.Messages[len(req.Messages)-1]
			text = "You said: " + last.Content
		}

		audio, err := tts.Synthesize(
			context.Background(),
			text,
			cfg.GetTTSOpenAIURL(),
			cfg.GetTTSOpenAIModel(),
			cfg.GetTTSOpenAIVoice(),
			cfg.GetTTSOpenAIKey(),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// store audio
		token := randomToken()
		audioMu.Lock()
		audioMap[token] = audio
		audioMu.Unlock()

		// forget after 60 s
		time.AfterFunc(60*time.Second, func() {
			audioMu.Lock()
			delete(audioMap, token)
			audioMu.Unlock()
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(replyRsp{
			Text:	 text,
			AudioURL: "/conversation/tts?token=" + token,
		})
	})
}

func ConversationTTS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		audioMu.RLock()
		audio, ok := audioMap[token]
		audioMu.RUnlock()
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(audio)))
		w.Write(audio)
	})
}

func randomToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
