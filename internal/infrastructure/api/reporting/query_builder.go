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
		params = append(params, fmt.Sprintf("pivot=%s", strings.TrimSpace(input.Pivot)))
	}

	// Date range as RestLi tuple syntax required by LinkedIn analytics finder.
	params = append(params, fmt.Sprintf("dateRange=%s", qb.buildDateRangeParam(input.DateRange)))

	// Time granularity as plain enum symbol.
	if input.TimeGranularity != "" {
		params = append(params, fmt.Sprintf("timeGranularity=%s", strings.TrimSpace(input.TimeGranularity)))
	}

	// Facets - at least one is required (accounts come before pivot in the example)
	// Always include the account from AccountID input
	accountURN := fmt.Sprintf("urn:li:sponsoredAccount:%s", input.AccountID)
	params = append(params, fmt.Sprintf("accounts=List(%s)", accountURN))

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
		params = append(params, fmt.Sprintf("campaignType=(value:%s)", strings.TrimSpace(input.CampaignType)))
	}

	// Sort parameters - based on sample format
	if input.SortBy.Field != "" {
		params = append(params, fmt.Sprintf("sortBy=(field:%s)", strings.TrimSpace(input.SortBy.Field)))
	}
	if input.SortBy.Order != "" {
		params = append(params, fmt.Sprintf("sortBy=(order:%s)", strings.TrimSpace(input.SortBy.Order)))
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
	encoded := make([]string, 0, len(items))
	for _, item := range items {
		normalized := strings.TrimSpace(item)
		if normalized == "" {
			continue
		}

		// If upstream already passed URL-encoded URNs, decode them so we emit canonical RestLi List() values.
		decoded, err := url.QueryUnescape(normalized)
		if err == nil {
			normalized = decoded
		}
		encoded = append(encoded, normalized)
	}
	return strings.Join(encoded, ",")
}

func (qb *QueryBuilder) buildDateRangeParam(dateRange DateRange) string {
	start := fmt.Sprintf("(year:%d,month:%d,day:%d)", dateRange.Start.Year, dateRange.Start.Month, dateRange.Start.Day)
	if dateRange.End == nil {
		return fmt.Sprintf("(start:%s)", start)
	}

	end := fmt.Sprintf("(year:%d,month:%d,day:%d)", dateRange.End.Year, dateRange.End.Month, dateRange.End.Day)
	return fmt.Sprintf("(start:%s,end:%s)", start, end)
}
