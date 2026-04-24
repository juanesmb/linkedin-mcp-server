package creatives

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

// BuildSearchCreativesByCampaignsQuery builds a GET URL for the creatives criteria finder
// (Microsoft Learn: Create and Manage Creatives — Search For Creatives).
func (qb *QueryBuilder) BuildSearchCreativesByCampaignsQuery(input SearchInput) string {
	accountID := strings.TrimSpace(input.AccountID)
	endpoint := fmt.Sprintf("%s/adAccounts/%s/creatives", strings.TrimRight(qb.baseURL, "/"), url.PathEscape(accountID))
	params := []string{"q=criteria"}

	if len(input.CampaignURNs) > 0 {
		list := buildRestLiEncodedList(input.CampaignURNs)
		params = append(params, fmt.Sprintf("campaigns=List(%s)", list))
	}

	sortOrder := strings.TrimSpace(input.SortOrder)
	if sortOrder == "" {
		sortOrder = "ASCENDING"
	}
	params = append(params, fmt.Sprintf("sortOrder=%s", url.QueryEscape(sortOrder)))

	if input.PageSize > 0 {
		params = append(params, fmt.Sprintf("pageSize=%d", input.PageSize))
	}

	if token := strings.TrimSpace(input.PageToken); token != "" {
		params = append(params, fmt.Sprintf("pageToken=%s", url.QueryEscape(token)))
	}

	return endpoint + "?" + strings.Join(params, "&")
}

func buildRestLiEncodedList(items []string) string {
	encoded := make([]string, 0, len(items))
	for _, item := range items {
		normalized := strings.TrimSpace(item)
		if normalized == "" {
			continue
		}
		if decoded, err := url.QueryUnescape(normalized); err == nil {
			normalized = decoded
		}
		encoded = append(encoded, url.QueryEscape(normalized))
	}
	return strings.Join(encoded, ",")
}
