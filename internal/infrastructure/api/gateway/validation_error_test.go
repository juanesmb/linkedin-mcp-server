package gateway

import (
	"testing"

	"linkedin-mcp/internal/infrastructure/api"

	"github.com/stretchr/testify/require"
)

func TestParseLinkedInParamValidationResponse(t *testing.T) {
	response := &api.Response{
		StatusCode: 400,
		Body: []byte(`{
			"code":"LINKEDIN_PARAM_INVALID",
			"message":"Invalid param",
			"inputErrors":[
				{"code":"PARAM_INVALID","description":"wrong type","fieldPath":"search"}
			]
		}`),
	}

	validationErr, ok := ParseLinkedInParamValidationResponse(response)
	require.True(t, ok)
	require.Equal(t, "Invalid param", validationErr.Message)
	require.Equal(t, 400, validationErr.ProviderStatus)
	require.Len(t, validationErr.InputErrors, 1)
	require.Equal(t, "search", validationErr.InputErrors[0].FieldPath)
	require.Equal(t, "PARAM_INVALID", validationErr.InputErrors[0].Code)
}

func TestParseLinkedInParamValidationResponse_NonValidation(t *testing.T) {
	response := &api.Response{
		StatusCode: 500,
		Body:       []byte(`{"code":"INTERNAL_ERROR"}`),
	}

	validationErr, ok := ParseLinkedInParamValidationResponse(response)
	require.False(t, ok)
	require.Nil(t, validationErr)
}

func TestParseLinkedInParamValidationResponse_GenericInvalidQueryParameters(t *testing.T) {
	response := &api.Response{
		StatusCode: 400,
		Body:       []byte(`{"code":"LINKEDIN_API_ERROR","message":"Invalid query parameters passed to request","request_id":"req_123"}`),
	}

	validationErr, ok := ParseLinkedInParamValidationResponse(response)
	require.True(t, ok)
	require.Contains(t, validationErr.Message, "Invalid query parameters passed to request")
	require.Contains(t, validationErr.Message, "req_123")
	require.Equal(t, 400, validationErr.ProviderStatus)
}
