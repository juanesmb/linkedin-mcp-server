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
	if !strings.EqualFold(strings.TrimSpace(code), "LINKEDIN_PARAM_INVALID") {
		return nil, false
	}

	message, _ := payload["message"].(string)
	inputErrors := parseInputErrors(payload["inputErrors"])

	return &LinkedInParamValidationError{
		Message:        strings.TrimSpace(message),
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
