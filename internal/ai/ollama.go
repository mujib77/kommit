package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mujib77/kommit/config"
)

type OllamaProvider struct {
	cfg    config.Config
	client *http.Client
}

func NewOllama(cfg config.Config) *OllamaProvider {
	return &OllamaProvider{
		cfg:    cfg,
		client: &http.Client{},
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (o *OllamaProvider) GenerateMessages(
	ctx context.Context,
	diff string,
	style string,
) ([]string, error) {
	prompt := BuildPrompt(diff, style)

	model := o.cfg.Model
	if model == "gpt-4o" || model == "claude-sonnet-4-20250514" {
		model = "llama3.2"
	}

	reqBody := ollamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx, "POST",
		"http://localhost:11434/api/generate",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("content-type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama request failed — is ollama running? start with: ollama serve")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp ollamaResponse
	err = json.Unmarshal(respBody, &ollamaResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	content := strings.TrimSpace(ollamaResp.Response)
	return parseMessages(content)
}