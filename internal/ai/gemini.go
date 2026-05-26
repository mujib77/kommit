package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/mujib77/kommit/config"
)

type GeminiProvider struct {
	cfg    config.Config
	client *http.Client
}

func NewGemini(cfg config.Config) *GeminiProvider {
	return &GeminiProvider{
		cfg:    cfg,
		client: &http.Client{},
	}
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (g *GeminiProvider) GenerateMessages(
	ctx context.Context,
	diff string,
	style string,
) ([]string, error) {
	prompt := BuildPrompt(diff, style)

	model := "gemini-3.5-flash"

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		model,
		g.cfg.APIKey,
	)

	req, err := http.NewRequestWithContext(
		ctx, "POST", url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("content-type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("gemini error %d: %s", resp.StatusCode, string(respBody))
	}

	var geminiResp geminiResponse
	err = json.Unmarshal(respBody, &geminiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from gemini")
	}

	content := strings.TrimSpace(
		geminiResp.Candidates[0].Content.Parts[0].Text,
	)
	return parseMessages(content)
}