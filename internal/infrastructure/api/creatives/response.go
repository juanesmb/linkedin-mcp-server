package creatives

// LinkedInListResponse is the creatives search envelope (cursor pagination uses metadata.nextPageToken).
type LinkedInListResponse struct {
	Elements []map[string]any `json:"elements"`
	Metadata map[string]any   `json:"metadata"`
}

// SearchResult is returned by the repository after normalization.
type SearchResult struct {
	Elements []NormalizedCreative
	Paging   PagingSummary
}

// NormalizedCreative is the stable v1 MCP DTO (no raw LinkedIn payload).
type NormalizedCreative struct {
	CreativeID     string `json:"creativeId"`
	CreativeURN    string `json:"creativeUrn"`
	CampaignURN    string `json:"campaignUrn,omitempty"`
	IntendedStatus string `json:"intendedStatus,omitempty"`
	ReviewStatus   string `json:"reviewStatus,omitempty"`
	IsServing      *bool  `json:"isServing,omitempty"`
	Format         string `json:"format,omitempty"`
	Headline       string `json:"headline,omitempty"`
	Description    string `json:"description,omitempty"`
	CTA            string `json:"cta,omitempty"`
	LandingPageURL string `json:"landingPageUrl,omitempty"`
	ContentKind    string `json:"contentKind,omitempty"`
}

// PagingSummary is a small, agent-friendly view of LinkedIn pagination.
type PagingSummary struct {
	PageSize      int    `json:"pageSize,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}
