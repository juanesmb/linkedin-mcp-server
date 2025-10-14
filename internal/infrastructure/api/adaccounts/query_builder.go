package adaccounts

import (
	"fmt"
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

func (qb *QueryBuilder) BuildSearchAdAccountsQuery(input SearchInput) (string, map[string]string) {
	endpoint := fmt.Sprintf("%s/adAccounts", strings.TrimRight(qb.baseURL, "/"))
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
	params := []string{"q=search"}

	var searchParts []string

	addList := func(field string, values []string) {
		if len(values) == 0 {
			return
		}
		cleaned := make([]string, 0, len(values))
		for _, value := range values {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				continue
			}
			cleaned = append(cleaned, trimmed)
		}
		if len(cleaned) == 0 {
			return
		}
		searchParts = append(searchParts, fmt.Sprintf("%s:(values:List(%s))", field, strings.Join(cleaned, ",")))
	}

	addList("status", input.Status)
	addList("id", input.AccountIDs)
	addList("reference", input.References)
	addList("name", input.Names)

	if input.Test != nil {
		searchParts = append(searchParts, fmt.Sprintf("test:%t", *input.Test))
	}

	if len(searchParts) > 0 {
		params = append(params, fmt.Sprintf("search=(%s)", strings.Join(searchParts, ",")))
	}

	if input.Start > 0 {
		params = append(params, fmt.Sprintf("start=%d", input.Start))
	}

	if input.Count > 0 {
		params = append(params, fmt.Sprintf("count=%d", input.Count))
	}

	return strings.Join(params, "&")
}
