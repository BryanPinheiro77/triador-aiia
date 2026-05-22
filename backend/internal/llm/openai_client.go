package llm

import (
	"context"
	"errors"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
}

func NewOpenAIClient() *OpenAIClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")

	if apiKey == "" {
		panic("OPENAI_API_KEY is required")
	}

	if model == "" {
		model = "gpt-4o-mini"
	}

	return &OpenAIClient{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (c *OpenAIClient) Analyze(
	ctx context.Context,
	prompt string,
) (string, error) {
	response, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2,
		},
	)

	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", errors.New("OpenAI returned no choices")
	}

	return response.Choices[0].Message.Content, nil
}
