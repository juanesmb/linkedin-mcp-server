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
	// Always include the account from AccountID input, merged with any explicit accounts facet.
	accountURN := fmt.Sprintf("urn:li:sponsoredAccount:%s", input.AccountID)
	accountFacets := mergeAccountFacets(accountURN, input.Accounts)
	params = append(params, fmt.Sprintf("accounts=List(%s)", qb.buildListParam(accountFacets)))

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
	if len(input.Companies) > 0 {
		companiesList := qb.buildListParam(input.Companies)
		params = append(params, fmt.Sprintf("companies=List(%s)", companiesList))
	}

	// Campaign type - based on sample format
	if input.CampaignType != "" {
		params = append(params, fmt.Sprintf("campaignType=(value:%s)", strings.TrimSpace(input.CampaignType)))
	}

	// Sort parameters serialized as a single RestLi tuple: sortBy=(field:X,order:Y).
	// LinkedIn rejects split sortBy params; callers are expected to supply either both keys
	// or neither (enforced by the tool validator).
	sortByParts := make([]string, 0, 2)
	if field := strings.TrimSpace(input.SortBy.Field); field != "" {
		sortByParts = append(sortByParts, fmt.Sprintf("field:%s", field))
	}
	if order := strings.TrimSpace(input.SortBy.Order); order != "" {
		sortByParts = append(sortByParts, fmt.Sprintf("order:%s", order))
	}
	if len(sortByParts) > 0 {
		params = append(params, fmt.Sprintf("sortBy=(%s)", strings.Join(sortByParts, ",")))
	}

	// Fields parameter (required).
	//
	// LinkedIn treats `fields` as a strict projection: non-metric metadata such as
	// `dateRange` (time-bucketed responses) and `pivotValues` (pivoted responses)
	// is only returned when the caller lists it explicitly. We inject those
	// transparently so the tool always returns bucket/pivot labels when applicable.
	fields := augmentFieldsProjection(input.Fields, input.TimeGranularity, input.Pivot)

	if len(fields) > 0 {
		fieldsList := strings.Join(fields, ",")
		params = append(params, fmt.Sprintf("fields=%s", fieldsList))
	}

	return strings.Join(params, "&")
}

// augmentFieldsProjection returns the caller-supplied fields with metadata fields
// injected based on request shape: `dateRange` when timeGranularity is a bucketed
// value (anything other than ALL), `pivotValues` when a pivot is set. Already-present
// fields are preserved and not duplicated.
func augmentFieldsProjection(requested []string, timeGranularity, pivot string) []string {
	seen := make(map[string]struct{}, len(requested)+2)
	out := make([]string, 0, len(requested)+2)
	for _, f := range requested {
		normalized := strings.TrimSpace(f)
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}

	addMeta := func(name string) {
		if _, exists := seen[name]; exists {
			return
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}

	granularity := strings.ToUpper(strings.TrimSpace(timeGranularity))
	if granularity != "" && granularity != "ALL" {
		addMeta("dateRange")
	}
	if strings.TrimSpace(pivot) != "" {
		addMeta("pivotValues")
	}

	return out
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
		encoded = append(encoded, url.QueryEscape(normalized))
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

func mergeAccountFacets(defaultAccount string, explicitAccounts []string) []string {
	facets := make([]string, 0, len(explicitAccounts)+1)
	seen := map[string]struct{}{}

	add := func(value string) {
		normalized := strings.TrimSpace(value)
		if normalized == "" {
			return
		}
		if _, exists := seen[normalized]; exists {
			return
		}
		seen[normalized] = struct{}{}
		facets = append(facets, normalized)
	}

	add(defaultAccount)
	for _, account := range explicitAccounts {
		add(account)
	}

	return facets
}
