package reporting

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"linkedin-mcp/internal/infrastructure/api"
)

const (
	logMessageFailedRequest        = "failed to make request"
	logMessageLinkedInAPIError     = "linkedin api responded with error"
	logMessageFailedDecodeResponse = "failed to decode response"
	logMessageFailedDecodeElement  = "failed to decode analytics element"

	logTagURL    = "url"
	logTagError  = "error"
	logTagStatus = "status"
	logTagBody   = "body"
	logTagMetric = "metric"

	errFmtFailedRequest         = "failed to make request: %w"
	errFmtLinkedInAPIErrorJSON  = "linkedin api error: status %d, body: %v"
	errFmtLinkedInAPIErrorPlain = "linkedin api error: status %d, body: %s"
	errFmtLinkedInAPIError      = "linkedin api error: status %d"
	errFmtDecodeResponse        = "failed to decode response: %w"
	errFmtDecodeElement         = "failed to decode analytics element: %w"
)

type Logger interface {
	Error(ctx context.Context, message string, tags map[string]string)
}

type Repository struct {
	client       api.Client
	queryBuilder *QueryBuilder
	logger       Logger
}

func NewRepository(client api.Client, queryBuilder *QueryBuilder, logger Logger) *Repository {
	return &Repository{
		client:       client,
		queryBuilder: queryBuilder,
		logger:       logger,
	}
}

func (r *Repository) GetAnalytics(ctx context.Context, input AnalyticsInput) (*AnalyticsResult, error) {
	requestURL, headers := r.queryBuilder.BuildAnalyticsQuery(input)

	response, err := r.client.Get(ctx, requestURL, headers)
	if err != nil {
		r.logError(ctx, logMessageFailedRequest, map[string]string{
			logTagURL:   requestURL,
			logTagError: err.Error(),
		})
		return nil, fmt.Errorf(errFmtFailedRequest, err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyString := strings.TrimSpace(string(response.Body))
		tags := map[string]string{
			logTagURL:    requestURL,
			logTagStatus: strconv.Itoa(response.StatusCode),
		}
		if bodyString != "" {
			tags[logTagBody] = bodyString
		}

		r.logError(ctx, logMessageLinkedInAPIError, tags)

		var errBody any
		if err := json.Unmarshal(response.Body, &errBody); err == nil {
			return nil, fmt.Errorf(errFmtLinkedInAPIErrorJSON, response.StatusCode, errBody)
		}

		trimmedBody := strings.TrimSpace(string(response.Body))
		if trimmedBody != "" {
			return nil, fmt.Errorf(errFmtLinkedInAPIErrorPlain, response.StatusCode, trimmedBody)
		}

		return nil, fmt.Errorf(errFmtLinkedInAPIError, response.StatusCode)
	}

	var liResp LinkedInAnalyticsResponse
	err = json.Unmarshal(response.Body, &liResp)
	if err != nil {
		r.logError(ctx, logMessageFailedDecodeResponse, map[string]string{
			logTagURL:   requestURL,
			logTagError: err.Error(),
		})
		return nil, fmt.Errorf(errFmtDecodeResponse, err)
	}

	// Convert LinkedIn response to our domain model
	elements := make([]AnalyticsElement, len(liResp.Elements))
	for i, rawElement := range liResp.Elements {
		// Parse the raw element into a map to extract dynamic fields
		var elementMap map[string]interface{}
		err := json.Unmarshal(rawElement, &elementMap)
		if err != nil {
			r.logError(ctx, logMessageFailedDecodeElement, map[string]string{
				logTagMetric: strconv.Itoa(i),
				logTagError:  err.Error(),
			})
			return nil, fmt.Errorf(errFmtDecodeElement, err)
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

func (r *Repository) logError(ctx context.Context, message string, tags map[string]string) {
	if r.logger == nil {
		return
	}

	r.logger.Error(ctx, message, tags)
}
