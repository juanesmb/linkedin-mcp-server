package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeriveConnectURL(t *testing.T) {
	t.Run("explicit override", func(t *testing.T) {
		url := deriveConnectURL("https://jumon.example.com", "https://custom.example.com/link")
		require.Equal(t, "https://custom.example.com/link", url)
	})

	t.Run("derived from gateway base", func(t *testing.T) {
		url := deriveConnectURL("https://jumon.example.com", "")
		require.Equal(t, "https://jumon.example.com/connections", url)
	})

	t.Run("fallback path", func(t *testing.T) {
		url := deriveConnectURL("", "")
		require.Equal(t, "/connections", url)
	})
}
