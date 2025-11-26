package ai

import (
	"context"
	"encoding/json"
	"fmt"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type ChatMessage struct {
	Role	string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletion contacts the LLM and returns the assistant’s reply.
func ChatCompletion(ctx context.Context,
	messages []ChatMessage,
	baseURL, apiKey, model string,
) (assistantMsg ChatMessage, err error) {

	// convert our DTO → SDK types
	chatMsgs := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, m := range messages {
		switch m.Role {
		case "user":
			chatMsgs[i] = openai.UserMessage(m.Content)
		case "assistant":
			chatMsgs[i] = openai.AssistantMessage(m.Content)
		default:
			chatMsgs[i] = openai.UserMessage(m.Content)
		}
	}

	cli := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)
	resp, err := cli.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: chatMsgs,
		Model:	openai.ChatModel(model),
	})
	if err != nil {
		return assistantMsg, fmt.Errorf("llm request: %w", err)
	}
	if len(resp.Choices) == 0 {
		return assistantMsg, fmt.Errorf("empty llm response")
	}
	assistantMsg.Content = resp.Choices[0].Message.Content
	assistantMsg.Role = "assistant"
	return assistantMsg, nil
}

// MessagesFromJSON unmarshals our own DTO slice.
func MessagesFromJSON(raw string) ([]ChatMessage, error) {
	if raw == "" {
		return nil, nil
	}
	var dto []ChatMessage
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		return nil, err
	}
	return dto, nil
}

// MessagesToJSON marshals our own DTO slice.
func MessagesToJSON(msgs []ChatMessage) (string, error) {
	b, err := json.Marshal(msgs)
	return string(b), err
}
