package toolerrors

import (
	"errors"
	"testing"

	"linkedin-mcp/internal/infrastructure/api/gateway"

	"github.com/stretchr/testify/require"
)

func TestWrapToolExecutionError(t *testing.T) {
	err := WrapToolExecutionError("search ad accounts", gateway.ErrLinkedInNotConnected, "https://app.example.com/connections")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no LinkedIn account is connected")
	require.Contains(t, err.Error(), "https://app.example.com/connections")
	require.Contains(t, err.Error(), "retry this tool call")
}

func TestWrapToolExecutionError_GenericError(t *testing.T) {
	err := WrapToolExecutionError("search campaigns", errors.New("gateway timed out"), "https://app.example.com/connections")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to search campaigns")
	require.Contains(t, err.Error(), "gateway timed out")
}

func TestWrapToolExecutionError_ValidationError(t *testing.T) {
	err := WrapToolExecutionError(
		"search ad accounts",
		&gateway.LinkedInParamValidationError{
			Message: "Invalid param",
			InputErrors: []gateway.LinkedInInputError{
				{FieldPath: "search", Description: "Invalid value for param", Code: "PARAM_INVALID"},
			},
		},
		"https://app.example.com/connections",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "LinkedIn rejected the request parameters")
	require.Contains(t, err.Error(), "Field `search`")
	require.Contains(t, err.Error(), "retry this tool call")
}
