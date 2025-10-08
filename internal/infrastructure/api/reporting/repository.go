package reporting

import (
	"context"
	"encoding/json"
	"fmt"

	"linkedin-mcp/internal/infrastructure/api"
)

type Repository struct {
	client       api.Client
	queryBuilder *QueryBuilder
}

func NewRepository(client api.Client, queryBuilder *QueryBuilder) *Repository {
	return &Repository{
		client:       client,
		queryBuilder: queryBuilder,
	}
}

func (r *Repository) GetAnalytics(ctx context.Context, input AnalyticsInput) (*AnalyticsResult, error) {
	requestURL, headers := r.queryBuilder.BuildAnalyticsQuery(input)

	response, err := r.client.Get(ctx, requestURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.Unmarshal(response.Body, &errBody)

		return nil, fmt.Errorf("linkedin api error: status %d, body: %v", response.StatusCode, errBody)
	}

	var liResp LinkedInAnalyticsResponse
	err = json.Unmarshal(response.Body, &liResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert LinkedIn response to our domain model
	elements := make([]AnalyticsElement, len(liResp.Elements))
	for i, rawElement := range liResp.Elements {
		// Parse the raw element into a map to extract dynamic fields
		var elementMap map[string]interface{}
		err := json.Unmarshal(rawElement, &elementMap)
		if err != nil {
			return nil, fmt.Errorf("failed to decode analytics element: %w", err)
		}

		element := AnalyticsElement{}

		// Extract known fields
		if dateRange, ok := elementMap["dateRange"].(map[string]interface{}); ok {
			if start, ok := dateRange["start"].(map[string]interface{}); ok {
				element.DateRange = &DateRange{
					Start: Date{
						Year:  int(start["year"].(float64)),
						Month: int(start["month"].(float64)),
						Day:   int(start["day"].(float64)),
					},
				}
			}
			if end, ok := dateRange["end"].(map[string]interface{}); ok {
				if element.DateRange == nil {
					element.DateRange = &DateRange{}
				}
				element.DateRange.End = &Date{
					Year:  int(end["year"].(float64)),
					Month: int(end["month"].(float64)),
					Day:   int(end["day"].(float64)),
				}
			}
		}

		if pivotValues, ok := elementMap["pivotValues"].([]interface{}); ok {
			element.PivotValues = make([]string, len(pivotValues))
			for j, pv := range pivotValues {
				element.PivotValues[j] = pv.(string)
			}
		}

		// Extract all other fields as metrics (excluding known fields)
		metrics := make(map[string]interface{})
		for key, value := range elementMap {
			if key != "dateRange" && key != "pivotValues" {
				metrics[key] = value
			}
		}
		element.Metrics = metrics

		elements[i] = element
	}

	result := &AnalyticsResult{
		Elements: elements,
		Paging:   liResp.Paging,
	}

	return result, nil
}
