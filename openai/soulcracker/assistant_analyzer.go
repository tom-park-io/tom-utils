package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"soulcracker/app/common"
	"soulcracker/app/core/helper/logger"
	"soulcracker/app/database"
)

type ArticleComparisonRequest struct {
	Article     string `json:"article"`
	RawCompared string `json:"comparisons"` // This is a stringified JSON array
}

// AssistantAnalyzer provides high-level analysis tasks built on top of the Assistants API.
type AssistantAnalyzer struct {
	// rootContext core.Context
	Service *AssistantService
}

// NewAssistantAnalyzer creates a new AssistantAnalyzer.
func NewAssistantAnalyzer(service *AssistantService) *AssistantAnalyzer {
	return &AssistantAnalyzer{Service: service}
}

// AnalyzeDefiniteEvents takes a NewsUpdated event and a list of known DefiniteEvents,
// and returns a list of semantic matches or structured analysis results using the Assistant.
func (a *AssistantAnalyzer) AnalyzeDefiniteEvents(event *common.NewsUpdated, definiteEvents []*database.DefiniteEvent) ([]map[string]interface{}, error) {

	ctx := context.Background()

	// Step 1: Format definite events into a slice of simple maps.
	eventsData := make([]map[string]interface{}, 0, len(definiteEvents))
	for _, defEvent := range definiteEvents {
		eventsData = append(eventsData, map[string]interface{}{
			"content_id": defEvent.ID,
			"content":    defEvent.Content,
		})
	}

	// Step 2: Marshal the events into JSON for assistant input.
	eventsJSON, err := json.Marshal(eventsData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal events data: %w", err)
	}

	// Step 3: Construct the prompt for the assistant.
	prompt := ArticleComparisonRequest{
		Article:     event.News.Message,
		RawCompared: string(eventsJSON), // assign the marshaled JSON string
	}
	promptBytes, err := json.Marshal(prompt)
	if err != nil {
		logger.Zap.Infof("failed to marshal request:", err)
	}

	startTime := time.Now()

	// Step 4: Use the assistant to get a structured response.
	threadID, err := a.Service.CreateThread(ctx)
	if err != nil {
		return nil, fmt.Errorf("create thread failed: %w", err)
	}

	if err := a.Service.SendMessage(ctx, threadID, string(promptBytes)); err != nil {
		return nil, fmt.Errorf("send message failed: %w", err)
	}

	runID, err := a.Service.RunThread(ctx, threadID)
	if err != nil {
		return nil, fmt.Errorf("run thread failed: %w", err)
	}

	response, err := a.Service.WaitForRunResult(ctx, threadID, runID)
	if err != nil {
		return nil, fmt.Errorf("wait for result failed: %w", err)
	}

	endTime := time.Now()
	fmt.Printf("Time taken for analysis: %v seconds\n", endTime.Sub(startTime).Seconds())

	// Step 5: Parse the assistant's JSON response into structured output.
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse assistant response: %w", err)
	}

	// Wrap the single map into a slice
	return []map[string]interface{}{result}, nil
}
