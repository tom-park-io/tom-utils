package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"soulcracker/app/core/helper/logger"
	"soulcracker/env"

	"github.com/sashabaranov/go-openai"
)

// AssistantService provides methods for working with the Assistants API.
type AssistantService struct {
	Client      *AssistantClient
	AssistantID string // The Assistant to use for message runs
}

// NewAssistantService creates a new AssistantService with the given Assistant ID.
func NewAssistantService(environment *env.Env, client *AssistantClient) *AssistantService {
	assistantId := environment.Infra.OpenAI.AssistantId
	if assistantId == "" {
		logger.Zap.Warn("OpenAI AssistantId key not set in environment config, GPT analysis will not work")
	}

	return &AssistantService{
		Client:      client,
		AssistantID: assistantId,
	}
}

// CreateThread creates a new empty thread.
func (s *AssistantService) CreateThread(ctx context.Context) (string, error) {
	resp, err := s.Client.do(ctx, "POST", "/threads", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.ID, nil
}

// SendMessage adds a user message to an existing thread.
func (s *AssistantService) SendMessage(ctx context.Context, threadID, content string) error {
	body := map[string]string{"role": openai.ChatMessageRoleUser, "content": content}
	_, err := s.Client.do(ctx, "POST", fmt.Sprintf("/threads/%s/messages", threadID), body)
	return err
}

// RunThread starts a run for the assistant on the specified thread.
func (s *AssistantService) RunThread(ctx context.Context, threadID string) (string, error) {
	body := map[string]string{"assistant_id": s.AssistantID}
	resp, err := s.Client.do(ctx, "POST", fmt.Sprintf("/threads/%s/runs", threadID), body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.ID, nil
}

// WaitForRunResult polls until the run completes and returns the final message.
func (s *AssistantService) WaitForRunResult(ctx context.Context, threadID, runID string) (string, error) {
	for {
		resp, err := s.Client.do(ctx, "GET", fmt.Sprintf("/threads/%s/runs/%s", threadID, runID), nil)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var run struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
			return "", err
		}

		if run.Status == "completed" {
			break
		}
	}

	// Get the latest message
	resp, err := s.Client.do(ctx, "GET", fmt.Sprintf("/threads/%s/messages", threadID), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Data []struct {
			Content []struct {
				Text struct {
					Value string `json:"value"`
				} `json:"text"`
			} `json:"content"`
			Role string `json:"role"`
		} `json:"data"`
	}
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &res)

	if len(res.Data) == 0 {
		return "", fmt.Errorf("no messages found")
	}

	// Find the most recent assistant message
	for i := len(res.Data) - 1; i >= 0; i-- {
		msg := res.Data[i]
		if msg.Role == "assistant" && len(msg.Content) > 0 {
			return msg.Content[0].Text.Value, nil
		}
	}

	return "", fmt.Errorf("no assistant message found")
}
