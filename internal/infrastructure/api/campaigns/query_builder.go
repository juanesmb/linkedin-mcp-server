package campaigns

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

func (qb *QueryBuilder) BuildSearchCampaignsQuery(input SearchInput) (string, map[string]string) {
	endpoint := fmt.Sprintf("%s/adAccounts/%s/adCampaigns", strings.TrimRight(qb.baseURL, "/"), url.PathEscape(input.AccountID))
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

func (qb *QueryBuilder) buildQueryParams(input SearchInput) string {
	var params []string

	// Always include q=search
	params = append(params, "q=search")

	// Build single Rest.li-style composite search parameter
	var searchParts []string

	addList := func(field string, values []string) {
		cleaned := make([]string, 0, len(values))
		for _, v := range values {
			if v == "" {
				continue
			}
			cleaned = append(cleaned, v)
		}
		if len(cleaned) == 0 {
			return
		}
		searchParts = append(searchParts, fmt.Sprintf("%s:(values:List(%s))", field, strings.Join(cleaned, ",")))
	}

	addList("campaignGroup", input.CampaignGroupURNs)
	addList("associatedEntity", input.AssociatedEntityValues)
	addList("id", input.CampaignURNs)
	addList("status", input.Status)
	addList("type", input.Type)
	addList("name", input.Name)

	if input.Test != nil {
		searchParts = append(searchParts, fmt.Sprintf("test:%t", *input.Test))
	}

	if len(searchParts) > 0 {
		// Build the search parameter without URL encoding the Rest.li syntax
		searchParam := fmt.Sprintf("search=(%s)", strings.Join(searchParts, ","))
		params = append(params, searchParam)
	}

	if input.SortOrder != "" {
		params = append(params, fmt.Sprintf("sortOrder=%s", url.QueryEscape(input.SortOrder)))
	}
	if input.PageSize > 0 {
		params = append(params, fmt.Sprintf("pageSize=%d", input.PageSize))
	}
	if input.PageToken != "" {
		params = append(params, fmt.Sprintf("pageToken=%s", url.QueryEscape(input.PageToken)))
	}

	return strings.Join(params, "&")
}
