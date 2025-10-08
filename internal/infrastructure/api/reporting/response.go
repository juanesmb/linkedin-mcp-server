package reporting

import "encoding/json"

type LinkedInAnalyticsResponse struct {
	Elements []json.RawMessage `json:"elements"`
	Paging   Paging            `json:"paging"`
}

type LinkedInAnalyticsElement struct {
	DateRange   *DateRange `json:"dateRange,omitempty"`
	PivotValues []string   `json:"pivotValues,omitempty"`
	// All other fields are dynamic based on the requested metrics
	// We'll use json.RawMessage to handle them dynamically
}
