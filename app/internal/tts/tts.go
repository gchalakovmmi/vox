package tts

import (
	"context"
	"fmt"
	"io"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func Synthesize(ctx context.Context, text, baseURL, model, voice, key string) ([]byte, error) {
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(key),
	)

	resp, err := client.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Model: openai.SpeechModel(model),
		Voice: openai.AudioSpeechNewParamsVoice(voice),
		Input: text,
	})
	if err != nil {
		return nil, fmt.Errorf("tts request failed: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
