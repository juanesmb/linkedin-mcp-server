package creatives

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_BuildSearchCreativesByCampaignsQuery(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	u := qb.BuildSearchCreativesByCampaignsQuery(SearchInput{
		AccountID:    "512247261",
		CampaignURNs: []string{"urn:li:sponsoredCampaign:394073893"},
		PageSize:     50,
		PageToken:    "next-cursor-token",
		SortOrder:    "DESCENDING",
	})

	parsed, err := url.Parse(u)
	require.NoError(t, err)
	require.Equal(t, "/rest/adAccounts/512247261/creatives", parsed.Path)

	q := parsed.Query()
	require.Equal(t, "criteria", q.Get("q"))

	campaigns := q.Get("campaigns")
	require.True(t, strings.HasPrefix(campaigns, "List("))
	require.True(t, strings.Contains(campaigns, "urn"))
	require.Equal(t, "50", q.Get("pageSize"))
	require.Equal(t, "next-cursor-token", q.Get("pageToken"))
	require.Equal(t, "DESCENDING", q.Get("sortOrder"))
}

func TestQueryBuilder_BuildSearchCreativesByCampaignsQuery_EncodesURN(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")
	u := qb.BuildSearchCreativesByCampaignsQuery(SearchInput{
		AccountID:    "512247261",
		CampaignURNs: []string{"urn:li:sponsoredCampaign:394073893"},
	})

	parsed, err := url.Parse(u)
	require.NoError(t, err)
	rawList := parsed.Query().Get("campaigns")
	require.Contains(t, rawList, "394073893")
}
