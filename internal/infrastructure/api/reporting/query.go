package reporting

type AnalyticsInput struct {
	Pivot           string    `json:"pivot,omitempty"`
	DateRange       DateRange `json:"dateRange"`
	TimeGranularity string    `json:"timeGranularity"`
	CampaignType    string    `json:"campaignType,omitempty"`
	Shares          []string  `json:"shares,omitempty"`
	Campaigns       []string  `json:"campaigns,omitempty"`
	CampaignGroups  []string  `json:"campaignGroups,omitempty"`
	Accounts        []string  `json:"accounts,omitempty"`
	Companies       []string  `json:"companies,omitempty"`
	SortBy          SortBy    `json:"sortBy,omitempty"`
	Fields          []string  `json:"fields"`
}

type DateRange struct {
	Start Date  `json:"start"`
	End   *Date `json:"end,omitempty"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type SortBy struct {
	Field string `json:"field,omitempty"`
	Order string `json:"order,omitempty"`
}

type AnalyticsResult struct {
	Elements []AnalyticsElement `json:"elements"`
	Paging   Paging             `json:"paging"`
}

type AnalyticsElement struct {
	DateRange   *DateRange `json:"dateRange,omitempty"`
	PivotValues []string   `json:"pivotValues,omitempty"`
	// Dynamic fields based on requested metrics
	Metrics map[string]interface{} `json:"metrics,omitempty"`
}

type Paging struct {
	Count int                 `json:"count"`
	Start int                 `json:"start"`
	Links []map[string]string `json:"links"`
}
