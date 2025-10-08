package dto

type Input struct {
	AccountID       string   `json:"accountID" jsonschema:"LinkedIn Ad Account ID (numeric value, e.g., 512345678) - used as accounts facet in query"`
	Pivot           string   `json:"pivot,omitempty" jsonschema:"Pivot of results: COMPANY, ACCOUNT, SHARE, CAMPAIGN, CREATIVE, CAMPAIGN_GROUP, CONVERSION, CONVERSATION_NODE, CONVERSATION_NODE_OPTION_INDEX, SERVING_LOCATION, CARD_INDEX, MEMBER_COMPANY_SIZE, MEMBER_INDUSTRY, MEMBER_SENIORITY, MEMBER_JOB_TITLE, MEMBER_JOB_FUNCTION, MEMBER_COUNTRY_V2, MEMBER_REGION_V2, MEMBER_COMPANY, PLACEMENT_NAME, IMPRESSION_DEVICE_TYPE, EVENT_STAGE"`
	DateRangeStart  Date     `json:"dateRangeStart" jsonschema:"Start date for analytics (required)"`
	DateRangeEnd    *Date    `json:"dateRangeEnd,omitempty" jsonschema:"End date for analytics (optional)"`
	TimeGranularity string   `json:"timeGranularity" jsonschema:"Time granularity: ALL, DAILY, MONTHLY, YEARLY (required)"`
	CampaignType    string   `json:"campaignType,omitempty" jsonschema:"Campaign type: TEXT_AD, SPONSORED_UPDATES, SPONSORED_INMAILS, DYNAMIC"`
	Shares          []string `json:"shares,omitempty" jsonschema:"Array of Share URNs"`
	Campaigns       []string `json:"campaigns,omitempty" jsonschema:"Array of Campaign URNs (urn:li:sponsoredCampaign:{id})"`
	CampaignGroups  []string `json:"campaignGroups,omitempty" jsonschema:"Array of Campaign Group URNs (urn:li:sponsoredCampaignGroup:{id})"`
	Accounts        []string `json:"accounts,omitempty" jsonschema:"Array of Account URNs (urn:li:sponsoredAccount:{id})"`
	Companies       []string `json:"companies,omitempty" jsonschema:"Array of Organization URNs (urn:li:organization:{id})"`
	SortByField     string   `json:"sortByField,omitempty" jsonschema:"Field to sort by: COST_IN_LOCAL_CURRENCY, IMPRESSIONS, CLICKS, ONE_CLICK_LEADS, OPENS, SENDS, EXTERNAL_WEBSITE_CONVERSIONS"`
	SortByOrder     string   `json:"sortByOrder,omitempty" jsonschema:"Sort order: ASCENDING, DESCENDING"`
	Fields          []string `json:"fields" jsonschema:"List of metric fields to fetch (required)."`
}

type Date struct {
	Year  int `json:"year" jsonschema:"Year (e.g., 2024)"`
	Month int `json:"month" jsonschema:"Month (1-12)"`
	Day   int `json:"day" jsonschema:"Day (1-31)"`
}
