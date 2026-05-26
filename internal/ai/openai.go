package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/mujib77/kommit/config"
)

type OpenAIProvider struct {
	client *openai.Client
	cfg config.Config
}

func NewOpenAI(cfg config.Config) *OpenAIProvider {
	return &OpenAIProvider{
		client: openai.NewClient(cfg.APIKey),
		cfg: cfg,
	}
}

func (o *OpenAIProvider) GenerateMessages(
	ctx context.Context,
	diff string,
	style string,
) ([]string, error) {
	prompt := BuildPrompt(diff, style)

	resp, err := o.client.CreateChatCompletion(ctx,
	openai.ChatCompletionRequest{
		Model: o.cfg.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens: 500,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("openai error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from openai")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	return parseMessages(content)
		}

func parseMessages(content string) ([]string, error) {
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var messages []string
	err := json.Unmarshal([]byte(content), &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages returned")
	}
	return messages, nil
}