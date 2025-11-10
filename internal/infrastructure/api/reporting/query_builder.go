package reporting

import (
	"fmt"
	"net/url"
	"strings"
)

type QueryBuilder struct {
	baseURL     string
	version     string
	accessToken string
}

func NewQueryBuilder(baseURL, version, accessToken string) *QueryBuilder {
	return &QueryBuilder{
		baseURL:     baseURL,
		version:     version,
		accessToken: accessToken,
	}
}

func (qb *QueryBuilder) BuildAnalyticsQuery(input AnalyticsInput) (string, map[string]string) {
	endpoint := fmt.Sprintf("%s/adAnalytics", strings.TrimRight(qb.baseURL, "/"))
	queryParams := qb.buildQueryParams(input)
	fullURL := endpoint + "?" + queryParams

	headers := map[string]string{
		"Authorization":             qb.accessToken,
		"LinkedIn-Version":          qb.version,
		"X-Restli-Protocol-Version": "2.0.0",
		"Accept":                    "application/json",
	}

	return fullURL, headers
}

func (qb *QueryBuilder) buildQueryParams(input AnalyticsInput) string {
	var params []string

	// Required parameters
	params = append(params, "q=analytics")

	// Pivot parameter - supports both formats: pivot=MEMBER_COMPANY or pivot=(value:CAMPAIGN)
	if input.Pivot != "" {
		// Use simple format for certain pivot values
		simplePivots := map[string]bool{
			"MEMBER_COMPANY":      true,
			"MEMBER_INDUSTRY":     true,
			"MEMBER_SENIORITY":    true,
			"MEMBER_JOB_TITLE":    true,
			"MEMBER_JOB_FUNCTION": true,
			"MEMBER_COUNTRY_V2":   true,
			"MEMBER_REGION_V2":    true,
		}

		if simplePivots[input.Pivot] {
			params = append(params, fmt.Sprintf("pivot=%s", url.QueryEscape(input.Pivot)))
		} else {
			params = append(params, fmt.Sprintf("pivot=(value:%s)", url.QueryEscape(input.Pivot)))
		}
	}

	// Date range - based on sample: dateRange=(start:(year:2024,month:1,day:1))
	startDate := qb.formatDate(input.DateRange.Start)
	params = append(params, fmt.Sprintf("dateRange=(start:%s)", startDate))

	if input.DateRange.End != nil {
		endDate := qb.formatDate(*input.DateRange.End)
		// Replace the existing dateRange parameter
		for i, param := range params {
			if strings.HasPrefix(param, "dateRange=") {
				params[i] = fmt.Sprintf("dateRange=(start:%s,end:%s)", startDate, endDate)
				break
			}
		}
	}

	// Time granularity - supports both formats: timeGranularity=ALL or timeGranularity=(value:ALL)
	if input.TimeGranularity != "" {
		params = append(params, fmt.Sprintf("timeGranularity=(value:%s)", url.QueryEscape(input.TimeGranularity)))
	}

	// Facets - at least one is required (accounts come before pivot in the example)
	// Always include the account from AccountID input
	accountURN := fmt.Sprintf("urn:li:sponsoredAccount:%s", input.AccountID)
	params = append(params, fmt.Sprintf("accounts=List(%s)", url.QueryEscape(accountURN)))

	if len(input.Shares) > 0 {
		sharesList := qb.buildListParam(input.Shares)
		params = append(params, fmt.Sprintf("shares=List(%s)", sharesList))
	}
	if len(input.Campaigns) > 0 {
		campaignsList := qb.buildListParam(input.Campaigns)
		params = append(params, fmt.Sprintf("campaigns=List(%s)", campaignsList))
	}
	if len(input.CampaignGroups) > 0 {
		campaignGroupsList := qb.buildListParam(input.CampaignGroups)
		params = append(params, fmt.Sprintf("campaignGroups=List(%s)", campaignGroupsList))
	}
	if len(input.Accounts) > 0 {
		accountsList := qb.buildListParam(input.Accounts)
		params = append(params, fmt.Sprintf("accounts=List(%s)", accountsList))
	}
	if len(input.Companies) > 0 {
		companiesList := qb.buildListParam(input.Companies)
		params = append(params, fmt.Sprintf("companies=List(%s)", companiesList))
	}

	// Campaign type - based on sample format
	if input.CampaignType != "" {
		params = append(params, fmt.Sprintf("campaignType=(value:%s)", url.QueryEscape(input.CampaignType)))
	}

	// Sort parameters - based on sample format
	if input.SortBy.Field != "" {
		params = append(params, fmt.Sprintf("sortBy=(field:%s)", url.QueryEscape(input.SortBy.Field)))
	}
	if input.SortBy.Order != "" {
		params = append(params, fmt.Sprintf("sortBy=(order:%s)", url.QueryEscape(input.SortBy.Order)))
	}

	// Fields parameter (required)
	// Always include pivotValues to ensure pivot identifiers are returned
	fields := make([]string, 0, len(input.Fields)+1)
	fields = append(fields, input.Fields...)
	
	// Check if pivotValues is already in the fields list
	hasPivotValues := false
	for _, field := range fields {
		if field == "pivotValues" {
			hasPivotValues = true
			break
		}
	}
	
	// Add pivotValues if not already present
	if !hasPivotValues {
		fields = append(fields, "pivotValues")
	}
	
	if len(fields) > 0 {
		fieldsList := strings.Join(fields, ",")
		params = append(params, fmt.Sprintf("fields=%s", fieldsList))
	}

	return strings.Join(params, "&")
}

func (qb *QueryBuilder) formatDate(date Date) string {
	return fmt.Sprintf("(day:%d,month:%d,year:%d)", date.Day, date.Month, date.Year)
}

func (qb *QueryBuilder) buildListParam(items []string) string {
	encoded := make([]string, len(items))
	for i, item := range items {
		encoded[i] = url.QueryEscape(item)
	}
	return strings.Join(encoded, ",")
}
