package reporting

import (
	"fmt"
	"net/url"
	"strings"
)

type QueryBuilder struct {
	baseURL string
}

func NewQueryBuilder(baseURL string) *QueryBuilder {
	return &QueryBuilder{baseURL: baseURL}
}

func (qb *QueryBuilder) BuildAnalyticsQuery(input AnalyticsInput) string {
	endpoint := fmt.Sprintf("%s/adAnalytics", strings.TrimRight(qb.baseURL, "/"))
	queryParams := qb.buildQueryParams(input)
	fullURL := endpoint + "?" + queryParams

	return fullURL
}

func (qb *QueryBuilder) buildQueryParams(input AnalyticsInput) string {
	var params []string

	// Required parameters
	params = append(params, "q=analytics")

	// Pivot parameter as plain enum symbol.
	if input.Pivot != "" {
		params = append(params, fmt.Sprintf("pivot=%s", url.QueryEscape(input.Pivot)))
	}

	// Date range as dotted nested query keys.
	params = append(params, qb.buildDateRangeParams(input.DateRange)...)

	// Time granularity as plain enum symbol.
	if input.TimeGranularity != "" {
		params = append(params, fmt.Sprintf("timeGranularity=%s", url.QueryEscape(input.TimeGranularity)))
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
	fields := make([]string, 0, len(input.Fields))
	fields = append(fields, input.Fields...)

	if len(fields) > 0 {
		fieldsList := strings.Join(fields, ",")
		params = append(params, fmt.Sprintf("fields=%s", fieldsList))
	}

	return strings.Join(params, "&")
}

func (qb *QueryBuilder) buildListParam(items []string) string {
	encoded := make([]string, len(items))
	for i, item := range items {
		encoded[i] = url.QueryEscape(item)
	}
	return strings.Join(encoded, ",")
}

func (qb *QueryBuilder) buildDateRangeParams(dateRange DateRange) []string {
	params := []string{
		fmt.Sprintf("dateRange.start.year=%d", dateRange.Start.Year),
		fmt.Sprintf("dateRange.start.month=%d", dateRange.Start.Month),
		fmt.Sprintf("dateRange.start.day=%d", dateRange.Start.Day),
	}

	if dateRange.End != nil {
		params = append(params,
			fmt.Sprintf("dateRange.end.year=%d", dateRange.End.Year),
			fmt.Sprintf("dateRange.end.month=%d", dateRange.End.Month),
			fmt.Sprintf("dateRange.end.day=%d", dateRange.End.Day),
		)
	}

	return params
}
