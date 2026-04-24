package creatives

// SearchInput lists creatives for one or more campaigns via the LinkedIn
// GET /rest/adAccounts/{adAccountId}/creatives finder (q=criteria).
// Pagination is cursor-based (pageSize / pageToken); max pageSize is 100 per LinkedIn.
type SearchInput struct {
	AccountID string
	// CampaignURNs are full URNs, e.g. urn:li:sponsoredCampaign:394073893.
	CampaignURNs []string
	PageSize     int
	PageToken    string
	SortOrder    string
}
