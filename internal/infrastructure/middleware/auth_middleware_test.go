package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeVerifier struct {
	userID string
	err    error
}

func (f fakeVerifier) Verify(ctx context.Context, token string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.userID, nil
}

func TestRequireBearerAuth_WhenMissingToken_ReturnsChallenge(t *testing.T) {
	handler := RequireBearerAuth(
		fakeVerifier{userID: "user_123"},
		"https://mcp.example.com/.well-known/oauth-protected-resource",
		"mcp:tools:read",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("handler should not be called without token")
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnauthorized, res.Code)
	require.Contains(t, res.Header().Get("WWW-Authenticate"), `resource_metadata="https://mcp.example.com/.well-known/oauth-protected-resource"`)
	require.Contains(t, res.Header().Get("WWW-Authenticate"), `scope="mcp:tools:read"`)
}

func TestRequireBearerAuth_WhenTokenValid_InjectsUserIntoContext(t *testing.T) {
	handler := RequireBearerAuth(
		fakeVerifier{userID: "user_123"},
		"https://mcp.example.com/.well-known/oauth-protected-resource",
		"",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			require.True(t, ok)
			require.Equal(t, "user_123", userID)
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestRequireBearerAuth_WhenTokenInvalid_ReturnsUnauthorized(t *testing.T) {
	handler := RequireBearerAuth(
		fakeVerifier{err: errors.New("invalid token")},
		"https://mcp.example.com/.well-known/oauth-protected-resource",
		"",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("handler should not be called with invalid token")
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnauthorized, res.Code)
}
