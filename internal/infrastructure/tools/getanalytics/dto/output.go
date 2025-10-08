package dto

type Output struct {
	Elements []AnalyticsElement `json:"elements" jsonschema:"Analytics results"`
	Paging   Paging             `json:"paging" jsonschema:"Pagination information"`
}

type AnalyticsElement struct {
	DateRange   *DateRange             `json:"dateRange,omitempty" jsonschema:"Date range for this data point"`
	PivotValues []string               `json:"pivotValues,omitempty" jsonschema:"Pivot values for this data point"`
	Metrics     map[string]interface{} `json:"metrics,omitempty" jsonschema:"Metric values (dynamic based on requested fields)"`
}

type DateRange struct {
	Start Date  `json:"start,omitempty"`
	End   *Date `json:"end,omitempty"`
}

type Paging struct {
	Count int                 `json:"count" jsonschema:"Number of elements returned"`
	Start int                 `json:"start" jsonschema:"Starting index"`
	Links []map[string]string `json:"links" jsonschema:"Pagination links"`
}
