package dto

type Input struct {
	CampaignGroupURNs      []string `json:"campaignGroupURNs" jsonschema:"Filter by Campaign Group URNs (urn:li:sponsoredCampaignGroup:{id})"`
	AssociatedEntityValues []string `json:"associatedEntityValues" jsonschema:"Filter by associated entity"`
	CampaignURNs           []string `json:"campaignURNs" jsonschema:"Filter by Campaign URNs (urn:li:sponsoredCampaign:{id})"`
	Status                 []string `json:"status" jsonschema:"Filter by status: ACTIVE, PAUSED, ARCHIVED, COMPLETED, CANCELED, DRAFT, PENDING_DELETION, REMOVED"`
	Type                   []string `json:"type" jsonschema:"Filter by type: TEXT_AD, SPONSORED_UPDATES, SPONSORED_INMAILS, DYNAMIC"`
	Name                   []string `json:"name" jsonschema:"Filter by name (exact match)"`
	SortOrder              string   `json:"sortOrder" jsonschema:"Sort by campaign ID: ASCENDING or DESCENDING (default ASCENDING)"`
	PageSize               int      `json:"pageSize" jsonschema:"Results per page (1-1000). Default 100"`
	PageToken              string   `json:"pageToken" jsonschema:"Opaque cursor for pagination"`
}
