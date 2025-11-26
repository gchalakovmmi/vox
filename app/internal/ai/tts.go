package ai

import (
	"context"
	"fmt"
	"io"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Synthesize returns MP3 bytes for the given text.
func Synthesize(ctx context.Context, text, baseURL, model, voice, key string) ([]byte, error) {
	cli := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(key),
	)
	resp, err := cli.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Model: openai.SpeechModel(model),
		Voice: openai.AudioSpeechNewParamsVoice(voice),
		Input: text,
	})
	if err != nil {
		return nil, fmt.Errorf("tts request: %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
