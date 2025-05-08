package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"soulcracker/app/core/helper/logger"
	"soulcracker/env"
)

// AssistantClient is a low-level HTTP client for the OpenAI Assistants API.
type AssistantClient struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewAssistantClient initializes the AssistantClient with default values.
func NewAssistantClient(environment *env.Env) *AssistantClient {
	apiKey := environment.Infra.OpenAI.ApiKey
	if apiKey == "" {
		logger.Zap.Warn("OpenAI API key not set in environment config, GPT analysis will not work")
	}

	return &AssistantClient{
		APIKey:     apiKey,
		BaseURL:    "https://api.openai.com/v1",
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// do sends an HTTP request with necessary headers for the Assistants API.
func (c *AssistantClient) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}
