package toolerrors

import (
	"fmt"
	"strings"

	"linkedin-mcp/internal/infrastructure/api/gateway"
)

func WrapToolExecutionError(operation string, err error, connectURL string) error {
	if gateway.IsLinkedInNotConnected(err) {
		return fmt.Errorf(
			"cannot %s because no LinkedIn account is connected for this user. Ask the user to open %s, complete LinkedIn connection, and then retry this tool call",
			operation,
			connectURL,
		)
	}
	if validationErr, ok := gateway.AsLinkedInParamValidation(err); ok {
		details := strings.TrimSpace(validationFields(validationErr))
		if details != "" {
			details = " " + details
		}
		return fmt.Errorf(
			"cannot %s because LinkedIn rejected the request parameters: %s.%s Please adjust the parameters and retry this tool call",
			operation,
			validationMessage(validationErr),
			details,
		)
	}

	return fmt.Errorf("failed to %s: %w", operation, err)
}

func validationMessage(err *gateway.LinkedInParamValidationError) string {
	message := strings.TrimSpace(err.Message)
	if message == "" {
		return "invalid request parameters"
	}
	return message
}

func validationFields(err *gateway.LinkedInParamValidationError) string {
	if len(err.InputErrors) == 0 {
		return ""
	}

	parts := make([]string, 0, len(err.InputErrors))
	for _, inputErr := range err.InputErrors {
		description := strings.TrimSpace(inputErr.Description)
		fieldPath := strings.TrimSpace(inputErr.FieldPath)
		if fieldPath != "" && description != "" {
			parts = append(parts, fmt.Sprintf("Field `%s`: %s.", fieldPath, description))
			continue
		}
		if fieldPath != "" {
			parts = append(parts, fmt.Sprintf("Field `%s` is invalid.", fieldPath))
			continue
		}
		if description != "" {
			parts = append(parts, description+".")
		}
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}
