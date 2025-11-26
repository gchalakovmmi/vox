package ai

import (
	"bytes"
	"context"
	"fmt"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Transcribe returns the text contained in the audio bytes.
func Transcribe(ctx context.Context, audio []byte, baseURL, model, key string) (string, error) {
	cli := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(key),
	)
	reader := bytes.NewReader(audio)
	resp, err := cli.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		Model: openai.AudioModel(model),
		File:  openai.File(reader, "audio.webm", ""),
	})
	if err != nil {
		return "", fmt.Errorf("stt request: %w", err)
	}
	return resp.Text, nil
}
