package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOAuthProtectedResourceHandler_ReturnsMetadata(t *testing.T) {
	configs := Configs{
		AuthConfig: AuthConfig{
			AuthorizationURL: "https://clerk.example.com",
			RequiredScope:    "mcp:tools:read",
		},
		ServerConfig: ServerConfig{
			Path:      "/mcp",
			PublicURL: "https://linkedin-mcp.example.com",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/.well-known/oauth-protected-resource", nil)
	res := httptest.NewRecorder()

	oauthProtectedResourceHandler(configs).ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	require.Equal(t, "https://linkedin-mcp.example.com/mcp", payload["resource"])
	require.Equal(t, []interface{}{"https://clerk.example.com"}, payload["authorization_servers"])
	require.Equal(t, []interface{}{"mcp:tools:read"}, payload["scopes_supported"])
}
