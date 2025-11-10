package getanalytics

import (
	"context"
	"fmt"
	"strings"

	"linkedin-mcp/internal/infrastructure/api/reporting"
	"linkedin-mcp/internal/infrastructure/tools/getanalytics/dto"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool struct {
	repository *reporting.Repository
}

func NewTool(repository *reporting.Repository) *Tool {
	return &Tool{
		repository: repository,
	}
}

func (t *Tool) GetAnalytics(ctx context.Context, req *mcp.CallToolRequest, input dto.Input) (*mcp.CallToolResult, dto.Output, error) {
	result := &mcp.CallToolResult{}

	if err := t.validateInput(input); err != nil {
		return result, dto.Output{}, fmt.Errorf("input validation failed: %w", err)
	}

	analyticsInput := t.convertInput(input)

	analyticsResult, err := t.repository.GetAnalytics(ctx, analyticsInput)
	if err != nil {
		return result, dto.Output{}, fmt.Errorf("failed to get analytics: %w", err)
	}

	output := t.convertOutput(analyticsResult)

	return result, output, nil
}

func (t *Tool) validateInput(input dto.Input) error {
	// Validate AccountID
	if input.AccountID == "" {
		return fmt.Errorf("accountID is required")
	}
	input.AccountID = strings.TrimSpace(input.AccountID)
	if input.AccountID == "" {
		return fmt.Errorf("accountID cannot be empty or whitespace only")
	}

	// Validate required fields
	if input.DateRangeStart.Year == 0 {
		return fmt.Errorf("dateRangeStart is required")
	}
	if input.TimeGranularity == "" {
		return fmt.Errorf("timeGranularity is required")
	}
	if len(input.Fields) == 0 {
		return fmt.Errorf("fields is required and cannot be empty")
	}

	// Validate time granularity
	validGranularities := map[string]bool{
		"ALL":     true,
		"DAILY":   true,
		"MONTHLY": true,
		"YEARLY":  true,
	}
	if !validGranularities[input.TimeGranularity] {
		return fmt.Errorf("invalid timeGranularity: %s. Must be one of: ALL, DAILY, MONTHLY, YEARLY", input.TimeGranularity)
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
		return fmt.Errorf("invalid pivot: %s", input.Pivot)
	}

	// Validate campaign type
	validCampaignTypes := map[string]bool{
		"TEXT_AD":           true,
		"SPONSORED_UPDATES": true,
		"SPONSORED_INMAILS": true,
		"DYNAMIC":           true,
	}
	if input.CampaignType != "" && !validCampaignTypes[input.CampaignType] {
		return fmt.Errorf("invalid campaignType: %s. Must be one of: TEXT_AD, SPONSORED_UPDATES, SPONSORED_INMAILS, DYNAMIC", input.CampaignType)
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
		return fmt.Errorf("invalid sortByField: %s", input.SortByField)
	}

	// Validate sort order
	if input.SortByOrder != "" && input.SortByOrder != "ASCENDING" && input.SortByOrder != "DESCENDING" {
		return fmt.Errorf("sortByOrder must be either ASCENDING or DESCENDING")
	}

	// Validate date range
	if input.DateRangeEnd != nil {
		if input.DateRangeEnd.Year < input.DateRangeStart.Year ||
			(input.DateRangeEnd.Year == input.DateRangeStart.Year && input.DateRangeEnd.Month < input.DateRangeStart.Month) ||
			(input.DateRangeEnd.Year == input.DateRangeStart.Year && input.DateRangeEnd.Month == input.DateRangeStart.Month && input.DateRangeEnd.Day < input.DateRangeStart.Day) {
			return fmt.Errorf("dateRangeEnd must be after dateRangeStart")
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
	for i, field := range input.Fields {
		input.Fields[i] = strings.TrimSpace(field)
	}

	return nil
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
