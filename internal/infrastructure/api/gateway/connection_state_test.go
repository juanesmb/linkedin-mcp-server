package gateway

import (
	"testing"

	"linkedin-mcp/internal/infrastructure/api"

	"github.com/stretchr/testify/require"
)

func TestIsLinkedInNotConnectedResponse(t *testing.T) {
	t.Run("status 404", func(t *testing.T) {
		require.True(t, IsLinkedInNotConnectedResponse(&api.Response{StatusCode: 404}))
	})

	t.Run("body connected false", func(t *testing.T) {
		require.True(t, IsLinkedInNotConnectedResponse(&api.Response{
			StatusCode: 200,
			Body:       []byte(`{"connected": false}`),
		}))
	})

	t.Run("body nested connection connected false", func(t *testing.T) {
		require.True(t, IsLinkedInNotConnectedResponse(&api.Response{
			StatusCode: 200,
			Body:       []byte(`{"connection": {"connected": false}}`),
		}))
	})

	t.Run("body status not_connected", func(t *testing.T) {
		require.True(t, IsLinkedInNotConnectedResponse(&api.Response{
			StatusCode: 409,
			Body:       []byte(`{"status":"not_connected"}`),
		}))
	})

	t.Run("connected true", func(t *testing.T) {
		require.False(t, IsLinkedInNotConnectedResponse(&api.Response{
			StatusCode: 200,
			Body:       []byte(`{"connected": true}`),
		}))
	})
}
