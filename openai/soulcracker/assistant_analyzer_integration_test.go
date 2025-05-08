package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"soulcracker/app/common"
	"soulcracker/app/database"

	"github.com/stretchr/testify/assert"
)

// go test -v --run TestAssistantAnalyzer_AnalyzeDefiniteEvents_Integration
func TestAssistantAnalyzer_AnalyzeDefiniteEvents_Integration(t *testing.T) {
	// Load environment variables
	apiKey := "YOUR_API_KEY"
	assistantID := "YOUR_ASSISTANT_ID"

	if apiKey == "" || assistantID == "" {
		t.Skip("OPENAI_API_KEY and OPENAI_ASSISTANT_ID must be set for integration test")
	}

	// Step 1: Set up real HTTP client and componentsh
	client := &AssistantClient{
		APIKey:     apiKey,
		BaseURL:    "https://api.openai.com/v1",
		HTTPClient: &http.Client{Timeout: 20 * time.Second},
	}

	service := &AssistantService{
		Client:      client,
		AssistantID: assistantID,
	}

	now := time.Now().UnixMilli()
	analyzer := NewAssistantAnalyzer(service)

	// Format the message with channel name
	channel := "kronon_channel"
	message := "kronon forever!"
	formattedMessage := fmt.Sprintf("[%s] %s", channel, message)

	// Step 2: Prepare input data
	news := &common.NewsUpdated{
		Source:    common.Telegram,
		TgChannel: common.TgChannel(channel),
		News: common.News{
			Message:   formattedMessage,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	definiteEvents := []*database.DefiniteEvent{
		{ID: 1, Content: "Earthquake strikes Tokyo on April 5th."},
		{ID: 2, Content: "Local baseball game postponed due to weather."},
		{ID: 3, Content: "AWESOME KRONON."},
	}

	// Step 3: Execute real assistant-based analysis
	result, err := analyzer.AnalyzeDefiniteEvents(news, definiteEvents)

	// Step 4: Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// Optional: Pretty-print result
	pretty, _ := json.MarshalIndent(result, "", "  ")
	t.Logf("Assistant response:\n%s", string(pretty))
}
