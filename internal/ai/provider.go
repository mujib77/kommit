package ai

import (
	"context"
	"fmt"

	"github.com/mujib77/kommit/config"
)

type Provider interface {
	GenerateMessages(ctx context.Context, diff string, style string) ([]string, error)
}

func New(cfg config.Config) (Provider, error) {
	switch cfg.Provider {
	case "openai":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("no API key — set OPENAI_API_KEY")
		}
		return NewOpenAI(cfg), nil
	case "anthropic":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("no API key — set ANTHROPIC_API_KEY")
		}
		return NewAnthropic(cfg), nil
	case "gemini":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("no API key — set GEMINI_API_KEY or add to config")
		}
		return NewGemini(cfg), nil
	case "ollama":
		return NewOllama(cfg), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s (use openai/anthropic/gemini/ollama)", cfg.Provider)
	}
}


func BuildPrompt(diff string, style string) string {
	styleGuide := ""
	if style == "conventional" {
		styleGuide = `
Use Conventional Commits format:
  - feat: new feature
  - fix: bug fix  
  - docs: documentation
  - style: formatting
  - refactor: code restructure
  - test: adding tests
  - chore: maintenance

  Format: <type>(<optional scope>): <description>
  Example: feat(auth): add JWT token validation`
}

	return fmt.Sprintf(`You are an expert developer writing git commit messages.

Analyze this git diff and generate exactly 3 different commit message options.

%s

Rules:
- First line max 72 characters
- Be specific and descriptive
- Use present tense ("add" not "added")
- No period at end
- Each option should have a different angle/focus

Return ONLY a JSON array with exactly 3 strings, nothing else:
["message 1", "message 2", "message 3"]

Git diff:
%s`, styleGuide, diff)
}