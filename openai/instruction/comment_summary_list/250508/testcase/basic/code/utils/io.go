package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Input struct {
	Article  string    `json:"reference_article"`
	Compared []Article `json:"articles_list"`
}

type Article struct {
	Id      int    `json:"content_id"`
	Content string `json:"content"`
}

func LoadInputFromFile(filename string) (*Input, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("파일 읽기 실패: %w", err)
	}

	var input Input
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("언마샬 실패: %w", err)
	}

	return &input, nil
}

func SaveJSONToFile(filename string, data []byte) error {
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return fmt.Errorf("JSON 인덴트 실패: %w", err)
	}

	if err := os.WriteFile(filename, out.Bytes(), 0644); err != nil {
		return fmt.Errorf("파일 저장 실패: %w", err)
	}

	return nil
}
