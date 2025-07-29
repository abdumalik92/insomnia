package controller

import (
	"encoding/json"
	"fmt"
	"strings"
)

func extractUserId(path string) (string, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("expected payload.ResourcePath contains 2 parts")
	}
	return parts[1], nil
}

func unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
