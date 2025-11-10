package reporting

type AnalyticsInput struct {
	AccountID       string    `json:"accountID"`
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
