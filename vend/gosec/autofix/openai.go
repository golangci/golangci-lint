package autofix

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const (
	ModelGPT4o           = openai.ChatModelGPT4o
	ModelGPT4oMini       = openai.ChatModelGPT4oMini
	DefaultOpenAIBaseURL = "https://api.openai.com/v1"
)

var _ GenAIClient = (*openaiWrapper)(nil)

type OpenAIConfig struct {
	Model       string
	APIKey      string `json:"-"`
	BaseURL     string
	MaxTokens   int
	Temperature float64
	SkipSSL     bool
}

type openaiWrapper struct {
	client      openai.Client
	model       openai.ChatModel
	maxTokens   int
	temperature float64
}

func NewOpenAIClient(config OpenAIConfig) (GenAIClient, error) {
	var options []option.RequestOption

	if config.APIKey != "" {
		options = append(options, option.WithAPIKey(config.APIKey))
	}

	// Support custom base URL (for OpenAI-compatible APIs)
	if config.BaseURL != "" {
		options = append(options, option.WithBaseURL(config.BaseURL))
	}

	// Support skip SSL verification
	if config.SkipSSL {
		// Create custom HTTP client with InsecureSkipVerify
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // #nosec G402
				},
			},
		}
		options = append(options, option.WithHTTPClient(httpClient))
	}

	openaiModel := parseOpenAIModel(config.Model)

	// Set default values
	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1024
	}

	temperature := config.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	return &openaiWrapper{
		client:      openai.NewClient(options...),
		model:       openaiModel,
		maxTokens:   maxTokens,
		temperature: temperature,
	}, nil
}

func (o *openaiWrapper) GenerateSolution(ctx context.Context, prompt string) (string, error) {
	params := openai.ChatCompletionNewParams{
		Model: o.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	}

	// Set optional parameters if available
	// Using WithMaxTokens and WithTemperature methods if they exist in v3
	resp, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("generating autofix: %w", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		return "", errors.New("no autofix returned by openai")
	}

	content := resp.Choices[0].Message.Content
	if content == "" {
		return "", errors.New("nothing found in the first autofix returned by openai")
	}

	return content, nil
}

func parseOpenAIModel(model string) openai.ChatModel {
	switch model {
	case "gpt-4o":
		return openai.ChatModelGPT4o
	case "gpt-4o-mini":
		return openai.ChatModelGPT4oMini
	default:
		return model
	}
}
