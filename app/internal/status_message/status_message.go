package status_message

import (
	"net/http"
	"log/slog"
)

type StatusMessage struct {
	Info	string
	Err	string
}


func WithStatusMessage(next func(StatusMessage) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		statusMessage := StatusMessage{
			Info:	queryParams.Get("info"),
			Err:	queryParams.Get("error"),
		}

		slog.Debug("internal/middleware/statusMessage.go", "Message", "Extracted status message from url", "statusMessage", statusMessage)
		next(statusMessage).ServeHTTP(w, r)
	})
}
