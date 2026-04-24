package creatives

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeCreative_TextAd(t *testing.T) {
	out := NormalizeCreative(map[string]any{
		"id":             "urn:li:sponsoredCreative:123",
		"campaign":       "urn:li:sponsoredCampaign:456",
		"intendedStatus": "ACTIVE",
		"isServing":      true,
		"format":         "TEXT_AD",
		"review":         map[string]any{"status": "APPROVED"},
		"content": map[string]any{
			"textAd": map[string]any{
				"headline":    "HL",
				"description": "DESC",
				"landingPage": "https://example.com/landing",
			},
		},
	})

	require.Equal(t, "123", out.CreativeID)
	require.Equal(t, "urn:li:sponsoredCreative:123", out.CreativeURN)
	require.Equal(t, "urn:li:sponsoredCampaign:456", out.CampaignURN)
	require.Equal(t, "ACTIVE", out.IntendedStatus)
	require.Equal(t, "APPROVED", out.ReviewStatus)
	require.NotNil(t, out.IsServing)
	require.True(t, *out.IsServing)
	require.Equal(t, "TEXT_AD", out.Format)
	require.Equal(t, "text_ad", out.ContentKind)
	require.Equal(t, "HL", out.Headline)
	require.Equal(t, "DESC", out.Description)
	require.Equal(t, "https://example.com/landing", out.LandingPageURL)
}

func TestNormalizeCreative_ReferenceOnly(t *testing.T) {
	out := NormalizeCreative(map[string]any{
		"id": float64(999),
		"content": map[string]any{
			"reference": "urn:li:ugcPost:abc",
		},
	})

	require.Equal(t, "999", out.CreativeID)
	require.Equal(t, "content_reference", out.ContentKind)
}

func TestNormalizeCreative_InlinePost(t *testing.T) {
	out := NormalizeCreative(map[string]any{
		"id": "urn:li:sponsoredCreative:1",
		"inlineContent": map[string]any{
			"post": map[string]any{
				"commentary":               "Main copy",
				"contentCallToActionLabel": "LEARN_MORE",
				"contentLandingPage":       "https://dest.example",
			},
		},
	})

	require.Equal(t, "inline_post", out.ContentKind)
	require.Equal(t, "Main copy", out.Headline)
	require.Equal(t, "LEARN_MORE", out.CTA)
	require.Equal(t, "https://dest.example", out.LandingPageURL)
}
