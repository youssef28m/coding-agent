package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type LLM interface {
	GenerateResponse(ctx context.Context, history []openai.ChatCompletionMessage) (string ,error)
}

type GroqLLM struct {
	client *openai.Client
	model  string
}

func NewGroqLLM(apiKey string, model string) *GroqLLM {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"

	return &GroqLLM{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

func (g *GroqLLM) GenerateResponse(ctx context.Context, history []openai.ChatCompletionMessage) (string, error) {
	resp, err := g.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    g.model,
		Messages: history,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from LLM")
	}

	return resp.Choices[0].Message.Content, nil
}
