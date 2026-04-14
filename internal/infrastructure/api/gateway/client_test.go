package gateway

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	customhttp "linkedin-mcp/internal/infrastructure/http"

	"github.com/stretchr/testify/require"
)

func TestProxyLinkedIn_SendsSecretAndUserContext(t *testing.T) {
	var gotSecret string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSecret = r.Header.Get("x-gateway-secret")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := NewClient(customhttp.NewClient(nil), server.URL, "super-secret")
	_, err := client.ProxyLinkedIn(context.Background(), "user_123", "adAccounts", map[string]string{
		"q": "search",
	})

	require.NoError(t, err)
	require.Equal(t, "super-secret", gotSecret)
	require.Equal(t, "user_123", gotBody["userId"])
}
