package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"code/utils"
)

type ThreadInput struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type ContentValue struct {
	Data []struct {
		Content []struct {
			Text struct {
				Value string `json:"value"`
			} `json:"text"`
		} `json:"content"`
	} `json:"data"`
}

func main() {

	// Load environment variables
	apiKey := "YOUR_API_KEY"
	assistantID := "YOUR_ASSISTANT_ID"
	baseUrl := "https://api.openai.com/v1"

	input, err := utils.LoadInputFromFile("../../../input_data.json")
	if err != nil {
		fmt.Println("마샬 실패:", err)
		return
	}
	// fmt.Printf("읽어온 입력: %+v\n", input)
	marshaledInput, err := json.Marshal(input)
	if err != nil {
		fmt.Println("마샬 실패:", err)
		return
	}
	threadInput := ThreadInput{
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: string(marshaledInput),
			},
		},
	}
	marshaledThreadInput, err := json.Marshal(threadInput)
	if err != nil {
		fmt.Println("마샬 실패:", err)
		return
	}

	repeat := 5
	for i := 0; i < repeat; i++ {

		client := &http.Client{Timeout: 20 * time.Second}

		// 1. create thread or reuse thread
		var threadId string
		// threadId = ""
		runURL := fmt.Sprintf("%s/threads", baseUrl)
		req, err := http.NewRequest(http.MethodPost, runURL, bytes.NewBuffer(marshaledThreadInput))
		if err != nil {
			fmt.Println("요청 생성 실패:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("OpenAI-Beta", "assistants=v2")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("요청 실패:", err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("응답 읽기 실패:", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		var threadResp struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(body, &threadResp); err != nil {
			fmt.Println("스레드 ID 파싱 실패:", err)
			continue
		}
		fmt.Println("생성된 스레드 ID:", threadResp.ID)
		threadId = threadResp.ID

		// 2. run thread
		var assistantId struct {
			ID string `json:"assistant_id"`
		}
		assistantId.ID = assistantID
		marshaledAssistantId, err := json.Marshal(assistantId)
		if err != nil {
			fmt.Println("JSON 마샬 실패:", err)
			continue
		}

		runURL = fmt.Sprintf("%s/threads/%s/runs", baseUrl, threadId)
		req, err = http.NewRequest(http.MethodPost, runURL, bytes.NewBuffer(marshaledAssistantId))
		if err != nil {
			fmt.Println("요청 생성 실패:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("OpenAI-Beta", "assistants=v2")

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("요청 실패:", err)
			continue
		}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("응답 읽기 실패:", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		var runResp struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(body, &runResp); err != nil {
			fmt.Println("Run ID 파싱 실패:", err)
			continue
		}
		fmt.Println("생성된 Run ID:", runResp.ID)

		// 3. polling queue
		var status string
		for {
			runURL = fmt.Sprintf("%s/threads/%s/runs/%s", baseUrl, threadId, runResp.ID)
			req, err = http.NewRequest(http.MethodGet, runURL, nil)
			if err != nil {
				fmt.Println("상태 요청 생성 실패:", err)
				break
			}
			req.Header.Set("Authorization", "Bearer "+apiKey)
			req.Header.Set("OpenAI-Beta", "assistants=v2")

			resp, err = client.Do(req)
			if err != nil {
				fmt.Println("상태 요청 실패:", err)
				break
			}
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("상태 응답 읽기 실패:", err)
				resp.Body.Close()
				break
			}
			resp.Body.Close()

			var statusResp struct {
				Status string `json:"status"`
			}
			if err := json.Unmarshal(body, &statusResp); err != nil {
				fmt.Println("상태 파싱 실패:", err)
				break
			}
			fmt.Println("현재 상태:", statusResp.Status)

			status = statusResp.Status
			if status == "completed" {
				fmt.Println("실행 완료!")
				break
			} else if status == "failed" || status == "cancelled" {
				fmt.Println("실행 실패/취소됨:", status)
				break
			}

			time.Sleep(2 * time.Second) // 잠시 대기 후 다시 확인
		}

		// 4. save the result
		if status == "completed" {
			runURL = fmt.Sprintf("%s/threads/%s/messages", baseUrl, threadId)
			req, err = http.NewRequest(http.MethodGet, runURL, nil)
			if err != nil {
				fmt.Println("메시지 요청 생성 실패:", err)
				continue
			}
			req.Header.Set("Authorization", "Bearer "+apiKey)
			req.Header.Set("OpenAI-Beta", "assistants=v2")

			resp, err = client.Do(req)
			if err != nil {
				fmt.Println("메시지 요청 실패:", err)
				continue
			}
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("메시지 응답 읽기 실패:", err)
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			var result ContentValue
			if err := json.Unmarshal(body, &result); err != nil {
				fmt.Println("파싱 실패:", err)
				return
			}
			if len(result.Data) > 0 && len(result.Data[0].Content) > 0 {
				contentValue := result.Data[0].Content[0].Text.Value

				filename := fmt.Sprintf("../result/test_result_%d.json", i+1)
				if err := os.WriteFile(filename, []byte(contentValue), 0644); err != nil {
					fmt.Println("파일 저장 실패:", err)
				} else {
					fmt.Println("✅ Content가", filename, "에 저장되었습니다.")
				}
			}
		}
	}
}
