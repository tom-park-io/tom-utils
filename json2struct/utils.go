package json2struct

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"unicode"
)

// toTitle converts the first character of a string to uppercase
func toTitle(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// generateRandomFilename generates a random filename with timestamp
func generateRandomFilename(prefix string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s_%s_%s.txt", prefix, timestamp, string(b))
}

// interfaceToBytes converts an interface{} to []byte using JSON encoding
func interfaceToBytes(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert interface to bytes: %v", err)
	}
	return bytes, nil
}
