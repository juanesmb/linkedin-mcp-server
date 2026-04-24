package dto

type Input struct {
	AccountID   string  `json:"accountID" jsonschema:"LinkedIn Ad Account ID (numeric, e.g. 512247261)"`
	CampaignID  string  `json:"campaignID,omitempty" jsonschema:"Numeric campaign id when campaignURN is omitted (e.g. 394073893)"`
	CampaignURN string  `json:"campaignURN,omitempty" jsonschema:"Full campaign URN (urn:li:sponsoredCampaign:{id}); overrides campaignID when set"`
	PageSize    *int    `json:"pageSize,omitempty" jsonschema:"Page size for cursor pagination (1-100). Default 100"`
	PageToken   *string `json:"pageToken,omitempty" jsonschema:"Opaque cursor from prior response paging.nextPageToken"`
	SortOrder   string  `json:"sortOrder,omitempty" jsonschema:"ASCENDING or DESCENDING (default ASCENDING)"`
}
