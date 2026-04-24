package reporting

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"linkedin-mcp/internal/infrastructure/api/gateway"
	"linkedin-mcp/internal/infrastructure/middleware"
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

	// creativeURNPrefix is the standard LinkedIn URN prefix for sponsored creatives.
	// Format: urn:li:sponsoredCreative:{id}
	// Reference: https://learn.microsoft.com/en-us/linkedin/shared/api-guide/concepts/urns
	creativeURNPrefix = "urn:li:sponsoredCreative:"
)

type Logger interface {
	Error(ctx context.Context, message string, tags map[string]string)
}

type Repository struct {
	gatewayClient *gateway.Client
	queryBuilder  *QueryBuilder
	logger        Logger
}

func NewRepository(gatewayClient *gateway.Client, queryBuilder *QueryBuilder, logger Logger) *Repository {
	return &Repository{
		gatewayClient: gatewayClient,
		queryBuilder:  queryBuilder,
		logger:        logger,
	}
}

func (r *Repository) GetAnalytics(ctx context.Context, input AnalyticsInput) (*AnalyticsResult, error) {
	requestURL := r.queryBuilder.BuildAnalyticsQuery(input)
	resourcePath, query, err := gateway.ParseLinkedInRESTProxyTarget(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build gateway proxy target: %w", err)
	}
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing authenticated user in request context")
	}

	connectionResponse, err := r.gatewayClient.GetLinkedInConnection(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LinkedIn connection state from gateway: %w", err)
	}
	if gateway.IsLinkedInNotConnectedResponse(connectionResponse) {
		return nil, gateway.ErrLinkedInNotConnected
	}
	if connectionResponse.StatusCode < 200 || connectionResponse.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to fetch LinkedIn connection state from gateway: status %d", connectionResponse.StatusCode)
	}

	response, err := r.gatewayClient.ProxyLinkedInOrRefresh(ctx, userID, resourcePath, query, nil)
	if err != nil {
		r.logError(ctx, logMessageFailedRequest, map[string]string{
			logTagURL:   requestURL,
			logTagError: err.Error(),
		})
		return nil, fmt.Errorf(errFmtFailedRequest, err)
	}
	if gateway.IsLinkedInNotConnectedResponse(response) {
		return nil, gateway.ErrLinkedInNotConnected
	}
	if validationErr, ok := gateway.ParseLinkedInParamValidationResponse(response); ok {
		return nil, validationErr
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
	if err := json.Unmarshal(response.Body, &liResp); err != nil {
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

		// Extract dateRange. LinkedIn returns each bucket's date range as
		// {start:{year,month,day}[, end:{...}]}. Use safe type assertions so a
		// malformed payload never panics; if we see a dateRange key but fail
		// to parse it, log a warning and surface whatever partial data we can.
		if dateRangeRaw, exists := elementMap["dateRange"]; exists {
			dateRange, ok := dateRangeRaw.(map[string]interface{})
			if !ok {
				r.logError(ctx, "linkedin returned dateRange with unexpected shape", map[string]string{
					logTagMetric: strconv.Itoa(i),
				})
			} else {
				if start, ok := parseDateObject(dateRange["start"]); ok {
					element.DateRange = &DateRange{Start: start}
				}
				if end, ok := parseDateObject(dateRange["end"]); ok {
					if element.DateRange == nil {
						element.DateRange = &DateRange{}
					}
					element.DateRange.End = &end
				}
				if element.DateRange == nil {
					r.logError(ctx, "linkedin dateRange present but no valid start/end decoded", map[string]string{
						logTagMetric: strconv.Itoa(i),
					})
				}
			}
		}

		if pivotValues, ok := elementMap["pivotValues"].([]interface{}); ok {
			element.PivotValues = make([]string, len(pivotValues))
			for j, pv := range pivotValues {
				element.PivotValues[j] = pv.(string)
			}
		}

		// Extract creative ID from pivotValues if it's a creative URN
		element.CreativeID = r.extractCreativeID(element.PivotValues)

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

// extractCreativeID extracts the numeric creative ID from a creative URN in pivotValues.
// URN format: urn:li:sponsoredCreative:{ID}
// According to LinkedIn API documentation, this is the standard format for creative URNs.
// Reference: https://learn.microsoft.com/en-us/linkedin/shared/api-guide/concepts/urns
// Returns empty string if no creative URN is found.
func (r *Repository) extractCreativeID(pivotValues []string) string {
	if len(pivotValues) == 0 {
		return ""
	}

	// Check the first pivot value for a creative URN
	urn := pivotValues[0]
	if !strings.HasPrefix(urn, creativeURNPrefix) {
		return ""
	}

	// Extract the ID (everything after the prefix)
	creativeID := strings.TrimPrefix(urn, creativeURNPrefix)
	if creativeID == "" {
		return ""
	}

	return creativeID
}

func (r *Repository) logError(ctx context.Context, message string, tags map[string]string) {
	if r.logger == nil {
		return
	}

	r.logger.Error(ctx, message, tags)
}

// parseDateObject safely decodes a LinkedIn date payload of shape
// {"year": N, "month": N, "day": N} into a Date. Returns (zero, false) if the
// value is missing, not an object, or does not contain all three keys as
// numbers. Using this helper instead of unchecked type assertions prevents the
// response decoder from panicking on unexpected LinkedIn payload shapes.
func parseDateObject(value interface{}) (Date, bool) {
	raw, ok := value.(map[string]interface{})
	if !ok {
		return Date{}, false
	}

	year, okY := numericFieldAsInt(raw["year"])
	month, okM := numericFieldAsInt(raw["month"])
	day, okD := numericFieldAsInt(raw["day"])
	if !okY || !okM || !okD {
		return Date{}, false
	}

	return Date{Year: year, Month: month, Day: day}, true
}

// numericFieldAsInt coerces a JSON-decoded numeric value to an int. In practice
// encoding/json only emits float64 into interface{}, but the other numeric
// shapes are accepted for defensive robustness against upstream format changes.
func numericFieldAsInt(value interface{}) (int, bool) {
	switch typed := value.(type) {
	case float64:
		return int(typed), true
	case int:
		return typed, true
	case int64:
		return int(typed), true
	default:
		return 0, false
	}
}
