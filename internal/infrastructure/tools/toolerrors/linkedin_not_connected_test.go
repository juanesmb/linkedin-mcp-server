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
