package autofix

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/securego/gosec/v2/issue"
)

const (
	AIProviderFlagHelp = `AI API provider to generate auto fixes to issues. Valid options are:
- gemini-2.5-pro, gemini-2.5-flash, gemini-2.5-flash-lite, gemini-2.0-flash, gemini-2.0-flash-lite (gemini, default);
- claude-sonnet-4-0 (claude, default), claude-sonnet-4-5, claude-opus-4-0, claude-opus-4-1, claude-haiku-4-5, claude-sonnet-3-7
- gpt-4o (openai, default), gpt-4o-mini`

	AIPrompt = `Provide a brief explanation and a solution to fix this security issue
  in Go programming language: %q.
  Answer in markdown format and keep the response limited to 200 words.`

	timeout = 30 * time.Second
)

type GenAIClient interface {
	GenerateSolution(ctx context.Context, prompt string) (string, error)
}

// GenerateSolution generates a solution for the given issues using the specified AI provider
func GenerateSolution(model, aiAPIKey, baseURL string, skipSSL bool, issues []*issue.Issue) (err error) {
	var client GenAIClient

	switch {
	case strings.HasPrefix(model, "claude"):
		client, err = NewClaudeClient(model, aiAPIKey)
	case strings.HasPrefix(model, "gemini"):
		client, err = NewGeminiClient(model, aiAPIKey)
	case strings.HasPrefix(model, "gpt"):
		config := OpenAIConfig{
			Model:   model,
			APIKey:  aiAPIKey,
			BaseURL: baseURL,
			SkipSSL: skipSSL,
		}
		client, err = NewOpenAIClient(config)
	default:
		// Default to OpenAI-compatible API for custom models
		config := OpenAIConfig{
			Model:   model,
			APIKey:  aiAPIKey,
			BaseURL: baseURL,
			SkipSSL: skipSSL,
		}
		client, err = NewOpenAIClient(config)
	}

	if err != nil {
		return fmt.Errorf("initializing AI client: %w", err)
	}

	return generateSolution(client, issues)
}

func generateSolution(client GenAIClient, issues []*issue.Issue) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cachedAutofix := make(map[string]string)
	for _, issue := range issues {
		if val, ok := cachedAutofix[issue.What]; ok {
			issue.Autofix = val
			continue
		}

		prompt := fmt.Sprintf(AIPrompt, issue.What)
		resp, err := client.GenerateSolution(ctx, prompt)
		if err != nil {
			return fmt.Errorf("generating autofix with gemini: %w", err)
		}

		if resp == "" {
			return errors.New("no autofix returned by gemini")
		}

		issue.Autofix = resp
		cachedAutofix[issue.What] = issue.Autofix
	}
	return nil
}
