package gateway

import (
	"encoding/json"
	"errors"
	"strings"

	"linkedin-mcp/internal/infrastructure/api"
)

var ErrLinkedInNotConnected = errors.New("linkedin account is not connected")

func IsLinkedInNotConnected(err error) bool {
	return errors.Is(err, ErrLinkedInNotConnected)
}

func IsLinkedInNotConnectedResponse(response *api.Response) bool {
	if response == nil {
		return false
	}
	if response.StatusCode == 404 {
		return true
	}
	if len(response.Body) == 0 {
		return false
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(response.Body, &payload); err != nil {
		return false
	}

	if connected, ok := extractConnected(payload); ok {
		return !connected
	}
	if status, ok := payload["status"].(string); ok && strings.EqualFold(strings.TrimSpace(status), "not_connected") {
		return true
	}
	if code, ok := payload["code"].(string); ok {
		normalized := strings.ToUpper(strings.TrimSpace(code))
		if normalized == "LINKEDIN_NOT_CONNECTED" || normalized == "NO_LINKEDIN_CONNECTION" {
			return true
		}
	}
	if message, ok := payload["message"].(string); ok {
		lower := strings.ToLower(strings.TrimSpace(message))
		if strings.Contains(lower, "not connected") || strings.Contains(lower, "no linkedin connection") {
			return true
		}
	}

	return false
}

func extractConnected(payload map[string]interface{}) (bool, bool) {
	if connected, ok := payload["connected"].(bool); ok {
		return connected, true
	}

	connectionRaw, ok := payload["connection"].(map[string]interface{})
	if !ok {
		return false, false
	}
	connected, ok := connectionRaw["connected"].(bool)
	return connected, ok
}
