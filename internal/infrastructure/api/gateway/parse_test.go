package gateway

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLinkedInRESTProxyTarget(t *testing.T) {
	path, query, err := ParseLinkedInRESTProxyTarget(
		"https://api.linkedin.com/rest/adAccounts?q=search&start=0",
	)
	require.NoError(t, err)
	require.Equal(t, "adAccounts", path)
	require.Equal(t, "search", query["q"])
	require.Equal(t, "0", query["start"])
}
