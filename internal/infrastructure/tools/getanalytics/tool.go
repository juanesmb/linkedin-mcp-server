package getanalytics

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"linkedin-mcp/internal/infrastructure/api/reporting"
	"linkedin-mcp/internal/infrastructure/tools/getanalytics/dto"
	"linkedin-mcp/internal/infrastructure/tools/toolerrors"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool struct {
	repository *reporting.Repository
	connectURL string
}

func NewTool(repository *reporting.Repository, connectURL string) *Tool {
	return &Tool{
		repository: repository,
		connectURL: connectURL,
	}
}

func (t *Tool) GetAnalytics(ctx context.Context, req *mcp.CallToolRequest, input dto.Input) (*mcp.CallToolResult, dto.Output, error) {
	result := &mcp.CallToolResult{}

	normalizedInput, derivedFields, err := t.validateAndNormalizeInput(input)
	if err != nil {
		return result, dto.Output{}, fmt.Errorf("input validation failed: %w", err)
	}

	analyticsInput := t.convertInput(normalizedInput)

	analyticsResult, err := t.repository.GetAnalytics(ctx, analyticsInput)
	if err != nil {
		return result, dto.Output{}, toolerrors.WrapToolExecutionError("get analytics", err, t.connectURL)
	}

	// LinkedIn does not expose ratio metrics (CPC, CTR, CPL, CPM, video
	// completion rate) in its AdAnalytics schema, so we compute them from the
	// raw fields we requested on the caller's behalf.
	injectDerivedMetrics(analyticsResult, derivedFields)

	output := t.convertOutput(analyticsResult)

	return result, output, nil
}

// validateAndNormalizeInput validates the raw tool input, splits the requested
// fields into raw LinkedIn fields and server-computed derived metrics, and
// unions the derived metrics' required raw fields into the outbound request.
// The returned slice of derived field names is used after the LinkedIn call
// to compute and inject the derived values into each response element.
func (t *Tool) validateAndNormalizeInput(input dto.Input) (dto.Input, []string, error) {
	// Validate AccountID
	if input.AccountID == "" {
		return dto.Input{}, nil, fmt.Errorf("accountID is required")
	}
	input.AccountID = strings.TrimSpace(input.AccountID)
	if input.AccountID == "" {
		return dto.Input{}, nil, fmt.Errorf("accountID cannot be empty or whitespace only")
	}

	// Validate required fields
	if input.DateRangeStart.Year == 0 {
		return dto.Input{}, nil, fmt.Errorf("dateRangeStart is required")
	}
	if input.TimeGranularity == "" {
		return dto.Input{}, nil, fmt.Errorf("timeGranularity is required")
	}
	if len(input.Fields) == 0 {
		return dto.Input{}, nil, fmt.Errorf("fields is required and cannot be empty")
	}

	// Validate time granularity
	validGranularities := map[string]bool{
		"ALL":     true,
		"DAILY":   true,
		"MONTHLY": true,
		"YEARLY":  true,
	}
	if !validGranularities[input.TimeGranularity] {
		return dto.Input{}, nil, fmt.Errorf("invalid timeGranularity: %s. Must be one of: ALL, DAILY, MONTHLY, YEARLY", input.TimeGranularity)
	}

	// Validate pivot values
	validPivots := map[string]bool{
		"COMPANY":                        true,
		"ACCOUNT":                        true,
		"SHARE":                          true,
		"CAMPAIGN":                       true,
		"CREATIVE":                       true,
		"CAMPAIGN_GROUP":                 true,
		"CONVERSION":                     true,
		"CONVERSATION_NODE":              true,
		"CONVERSATION_NODE_OPTION_INDEX": true,
		"SERVING_LOCATION":               true,
		"CARD_INDEX":                     true,
		"MEMBER_COMPANY_SIZE":            true,
		"MEMBER_INDUSTRY":                true,
		"MEMBER_SENIORITY":               true,
		"MEMBER_JOB_TITLE":               true,
		"MEMBER_JOB_FUNCTION":            true,
		"MEMBER_COUNTRY_V2":              true,
		"MEMBER_REGION_V2":               true,
		"MEMBER_COMPANY":                 true,
		"PLACEMENT_NAME":                 true,
		"IMPRESSION_DEVICE_TYPE":         true,
		"EVENT_STAGE":                    true,
	}
	if input.Pivot != "" && !validPivots[input.Pivot] {
		return dto.Input{}, nil, fmt.Errorf("invalid pivot: %s", input.Pivot)
	}

	// Validate campaign type
	validCampaignTypes := map[string]bool{
		"TEXT_AD":           true,
		"SPONSORED_UPDATES": true,
		"SPONSORED_INMAILS": true,
		"DYNAMIC":           true,
	}
	if input.CampaignType != "" && !validCampaignTypes[input.CampaignType] {
		return dto.Input{}, nil, fmt.Errorf("invalid campaignType: %s. Must be one of: TEXT_AD, SPONSORED_UPDATES, SPONSORED_INMAILS, DYNAMIC", input.CampaignType)
	}

	// Validate sort fields
	validSortFields := map[string]bool{
		"COST_IN_LOCAL_CURRENCY":       true,
		"IMPRESSIONS":                  true,
		"CLICKS":                       true,
		"ONE_CLICK_LEADS":              true,
		"OPENS":                        true,
		"SENDS":                        true,
		"EXTERNAL_WEBSITE_CONVERSIONS": true,
	}
	if input.SortByField != "" && !validSortFields[input.SortByField] {
		return dto.Input{}, nil, fmt.Errorf("invalid sortByField: %s", input.SortByField)
	}

	// Validate sort order
	if input.SortByOrder != "" && input.SortByOrder != "ASCENDING" && input.SortByOrder != "DESCENDING" {
		return dto.Input{}, nil, fmt.Errorf("sortByOrder must be either ASCENDING or DESCENDING")
	}

	// LinkedIn requires both field and order together as a RestLi tuple. Half-specified
	// sortBy is always rejected upstream, so fail fast with an actionable message instead
	// of forwarding a broken request.
	if (input.SortByField != "") != (input.SortByOrder != "") {
		return dto.Input{}, nil, fmt.Errorf("sortByField and sortByOrder must be provided together; supply both or omit both")
	}

	// Validate date range
	if input.DateRangeEnd != nil {
		if input.DateRangeEnd.Year < input.DateRangeStart.Year ||
			(input.DateRangeEnd.Year == input.DateRangeStart.Year && input.DateRangeEnd.Month < input.DateRangeStart.Month) ||
			(input.DateRangeEnd.Year == input.DateRangeStart.Year && input.DateRangeEnd.Month == input.DateRangeStart.Month && input.DateRangeEnd.Day < input.DateRangeStart.Day) {
			return dto.Input{}, nil, fmt.Errorf("dateRangeEnd must be after dateRangeStart")
		}
	}

	// Note: AccountID is always provided as the accounts facet, so no need to validate additional facets

	// Sanitize string inputs
	for i, share := range input.Shares {
		input.Shares[i] = strings.TrimSpace(share)
	}
	for i, campaign := range input.Campaigns {
		input.Campaigns[i] = strings.TrimSpace(campaign)
	}
	for i, group := range input.CampaignGroups {
		input.CampaignGroups[i] = strings.TrimSpace(group)
	}
	for i, account := range input.Accounts {
		input.Accounts[i] = strings.TrimSpace(account)
	}
	for i, company := range input.Companies {
		input.Companies[i] = strings.TrimSpace(company)
	}

	// Split requested fields into:
	//   - rawFields: LinkedIn schema fields (passed through untouched so LinkedIn
	//     remains the schema source of truth — no allowlist).
	//   - derivedFields: ratio/computed metrics (CPC, CTR, CPL, CPM, video
	//     completion rate) we calculate from raw fields after the response.
	// For each derived field we union its required raw dependencies into the
	// outbound request, deduplicated against the raw fields the caller already
	// asked for.
	rawFieldSeen := map[string]struct{}{}
	derivedSeen := map[string]struct{}{}
	rawFields := make([]string, 0, len(input.Fields))
	derivedFields := make([]string, 0)

	addRawField := func(field string) {
		if field == "" {
			return
		}
		if _, exists := rawFieldSeen[field]; exists {
			return
		}
		rawFieldSeen[field] = struct{}{}
		rawFields = append(rawFields, field)
	}

	for _, field := range input.Fields {
		trimmed := strings.TrimSpace(field)
		if trimmed == "" {
			continue
		}
		if metric, ok := lookupDerivedMetric(trimmed); ok {
			if _, exists := derivedSeen[metric.Name]; !exists {
				derivedSeen[metric.Name] = struct{}{}
				derivedFields = append(derivedFields, metric.Name)
			}
			for _, required := range metric.RequiredFields {
				addRawField(required)
			}
			continue
		}
		addRawField(t.normalizeFieldName(trimmed))
	}

	if len(rawFields) == 0 && len(derivedFields) == 0 {
		return dto.Input{}, nil, fmt.Errorf("no valid fields after normalization")
	}

	// Even when the caller only asked for derived metrics, LinkedIn still
	// requires at least one raw field in the request; the addRawField path
	// above guarantees that because every derived metric declares raw deps.
	input.Fields = rawFields
	return input, derivedFields, nil
}

func (t *Tool) convertInput(input dto.Input) reporting.AnalyticsInput {
	dateRange := reporting.DateRange{
		Start: reporting.Date{
			Year:  input.DateRangeStart.Year,
			Month: input.DateRangeStart.Month,
			Day:   input.DateRangeStart.Day,
		},
	}

	if input.DateRangeEnd != nil {
		dateRange.End = &reporting.Date{
			Year:  input.DateRangeEnd.Year,
			Month: input.DateRangeEnd.Month,
			Day:   input.DateRangeEnd.Day,
		}
	}

	return reporting.AnalyticsInput{
		AccountID:       input.AccountID,
		Pivot:           input.Pivot,
		DateRange:       dateRange,
		TimeGranularity: input.TimeGranularity,
		CampaignType:    input.CampaignType,
		Shares:          input.Shares,
		Campaigns:       input.Campaigns,
		CampaignGroups:  input.CampaignGroups,
		Accounts:        input.Accounts,
		Companies:       input.Companies,
		SortBy: reporting.SortBy{
			Field: input.SortByField,
			Order: input.SortByOrder,
		},
		Fields: input.Fields,
	}
}

func (t *Tool) convertOutput(result *reporting.AnalyticsResult) dto.Output {
	elements := make([]dto.AnalyticsElement, len(result.Elements))
	for i, element := range result.Elements {
		elements[i] = dto.AnalyticsElement{
			DateRange:   t.convertDateRange(element.DateRange),
			PivotValues: element.PivotValues,
			CreativeID:  element.CreativeID,
			Metrics:     element.Metrics,
		}
	}

	return dto.Output{
		Elements: elements,
		Paging: dto.Paging{
			Count: result.Paging.Count,
			Start: result.Paging.Start,
			Links: result.Paging.Links,
		},
	}
}

// normalizeFieldName converts field names from various formats (UPPER_CASE, UPPERCASE) to camelCase
// which is the format expected by LinkedIn API.
// Also handles common aliases and variations.
func (t *Tool) normalizeFieldName(field string) string {
	if field == "" {
		return ""
	}

	// Handle common aliases for raw LinkedIn schema fields. Derived metrics
	// (CPC, CTR, CPL, CPM, video completion rate) are handled separately via
	// lookupDerivedMetric in the caller and must not appear here.
	aliases := map[string]string{
		"CONVERSIONS":                    "externalWebsiteConversions",
		"SPEND_IN_LOCAL_CURRENCY":        "costInLocalCurrency",
		"SPEND":                          "costInLocalCurrency",
		"COST":                           "costInLocalCurrency",
		"LEADS":                          "oneClickLeads",
		"ONE_CLICK_LEAD_FORM_OPENS":      "oneClickLeadFormOpens",
		"UNIQUE_IMPRESSIONS":             "approximateMemberReach",
		"APPROXIMATE_UNIQUE_IMPRESSIONS": "approximateMemberReach",
		"TOTAL_ENGAGEMENTS":              "totalEngagements",
		"ENGAGEMENTS":                    "totalEngagements",
	}

	if alias, ok := aliases[strings.ToUpper(field)]; ok {
		return alias
	}

	// If already in camelCase (has lowercase), return as-is
	if len(field) > 0 && unicode.IsLower(rune(field[0])) {
		return field
	}

	// Convert UPPER_CASE or UPPERCASE to camelCase
	parts := strings.Split(field, "_")
	if len(parts) == 1 {
		// Single word - convert UPPERCASE to lowercase
		return strings.ToLower(field)
	}

	// Multiple parts separated by underscores
	var result strings.Builder
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if i == 0 {
			// First part: all lowercase
			result.WriteString(strings.ToLower(part))
		} else {
			// Subsequent parts: capitalize first letter, rest lowercase
			partLower := strings.ToLower(part)
			if len(partLower) > 0 {
				result.WriteString(strings.ToUpper(partLower[:1]))
				if len(partLower) > 1 {
					result.WriteString(partLower[1:])
				}
			}
		}
	}

	return result.String()
}

func (t *Tool) convertDateRange(dateRange *reporting.DateRange) *dto.DateRange {
	if dateRange == nil {
		return nil
	}

	dr := &dto.DateRange{
		Start: dto.Date{
			Year:  dateRange.Start.Year,
			Month: dateRange.Start.Month,
			Day:   dateRange.Start.Day,
		},
	}

	if dateRange.End != nil {
		dr.End = &dto.Date{
			Year:  dateRange.End.Year,
			Month: dateRange.End.Month,
			Day:   dateRange.End.Day,
		}
	}
	return dr
}

// injectDerivedMetrics computes the requested derived metrics for each
// element in the analytics result and writes them into element.Metrics under
// their canonical lowerCamelCase name. Values are emitted as nil when the
// computation is undefined (missing dependency or zero denominator), so the
// MCP client can distinguish "unavailable" from "zero".
func injectDerivedMetrics(result *reporting.AnalyticsResult, derivedFields []string) {
	if result == nil || len(derivedFields) == 0 {
		return
	}

	for i := range result.Elements {
		metrics := result.Elements[i].Metrics
		if metrics == nil {
			metrics = map[string]interface{}{}
		}

		for _, name := range derivedFields {
			metric, ok := derivedMetrics[name]
			if !ok {
				continue
			}
			value, defined := metric.Compute(metrics)
			if !defined {
				metrics[metric.Name] = nil
				continue
			}
			metrics[metric.Name] = value
		}

		result.Elements[i].Metrics = metrics
	}
}
