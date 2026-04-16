package gateway

import (
	"encoding/json"
	"errors"
	"strings"

	"linkedin-mcp/internal/infrastructure/api"
)

var errLinkedInParamValidation = errors.New("linkedin request parameters are invalid")

type LinkedInInputError struct {
	Code        string
	Description string
	FieldPath   string
}

type LinkedInParamValidationError struct {
	Message        string
	ProviderStatus int
	InputErrors    []LinkedInInputError
}

func (e *LinkedInParamValidationError) Error() string {
	if strings.TrimSpace(e.Message) == "" {
		return "linkedin request parameters are invalid"
	}
	return e.Message
}

func (e *LinkedInParamValidationError) Unwrap() error {
	return errLinkedInParamValidation
}

func IsLinkedInParamValidation(err error) bool {
	return errors.Is(err, errLinkedInParamValidation)
}

func AsLinkedInParamValidation(err error) (*LinkedInParamValidationError, bool) {
	var target *LinkedInParamValidationError
	if !errors.As(err, &target) {
		return nil, false
	}
	return target, true
}

func ParseLinkedInParamValidationResponse(response *api.Response) (*LinkedInParamValidationError, bool) {
	if response == nil || response.StatusCode < 400 || response.StatusCode >= 500 || len(response.Body) == 0 {
		return nil, false
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(response.Body, &payload); err != nil {
		return nil, false
	}

	code, _ := payload["code"].(string)
	normalizedCode := strings.ToUpper(strings.TrimSpace(code))

	// Strict param validation payload (preferred).
	if normalizedCode != "LINKEDIN_PARAM_INVALID" {
		// Generic LinkedIn 400s sometimes come through as LINKEDIN_API_ERROR with a message like:
		// "Invalid query parameters passed to request".
		// Treat these as correctable parameter issues so tools can surface actionable context.
		if normalizedCode != "LINKEDIN_API_ERROR" {
			return nil, false
		}
	}

	message, _ := payload["message"].(string)
	trimmedMessage := strings.TrimSpace(message)
	if normalizedCode == "LINKEDIN_API_ERROR" {
		lower := strings.ToLower(trimmedMessage)
		if !strings.Contains(lower, "invalid query parameters") {
			return nil, false
		}
	}

	inputErrors := parseInputErrors(payload["inputErrors"])
	requestID := strings.TrimSpace(asString(payload["request_id"]))
	if requestID == "" {
		requestID = strings.TrimSpace(asString(payload["requestId"]))
	}
	if requestID != "" {
		trimmedMessage = strings.TrimSpace(trimmedMessage + " (request_id: " + requestID + ")")
	}

	return &LinkedInParamValidationError{
		Message:        trimmedMessage,
		ProviderStatus: response.StatusCode,
		InputErrors:    inputErrors,
	}, true
}

func parseInputErrors(raw any) []LinkedInInputError {
	list, ok := raw.([]interface{})
	if !ok {
		return nil
	}

	result := make([]LinkedInInputError, 0, len(list))
	for _, entry := range list {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		inputError := LinkedInInputError{
			Code:        strings.TrimSpace(asString(entryMap["code"])),
			Description: strings.TrimSpace(asString(entryMap["description"])),
			FieldPath:   strings.TrimSpace(asString(entryMap["fieldPath"])),
		}
		if inputError.Code == "" && inputError.Description == "" && inputError.FieldPath == "" {
			continue
		}
		result = append(result, inputError)
	}

	return result
}

func asString(value any) string {
	s, _ := value.(string)
	return s
}
