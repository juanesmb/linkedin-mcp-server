package adaccounts

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

func (qb *QueryBuilder) BuildSearchAdAccountsQuery(input SearchInput) string {
	endpoint := fmt.Sprintf("%s/adAccounts", strings.TrimRight(qb.baseURL, "/"))
	queryParams := qb.buildQueryParams(input)
	fullURL := endpoint + "?" + queryParams

	return fullURL
}

func (qb *QueryBuilder) buildQueryParams(input SearchInput) string {
	params := []string{"q=search"}
	var searchParts []string

	addSearchList := func(field string, values []string, encode bool) {
		if len(values) == 0 {
			return
		}
		cleaned := make([]string, 0, len(values))
		for _, value := range values {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				continue
			}
			if encode {
				trimmed = url.QueryEscape(trimmed)
			}
			cleaned = append(cleaned, trimmed)
		}
		if len(cleaned) == 0 {
			return
		}
		searchParts = append(searchParts, fmt.Sprintf("%s:(values:List(%s))", field, strings.Join(cleaned, ",")))
	}

	addSearchList("status", input.Status, false)
	addSearchList("id", input.AccountIDs, false)
	addSearchList("reference", input.References, false)
	addSearchList("name", input.Names, true)

	if len(searchParts) > 0 {
		params = append(params, fmt.Sprintf("search=(%s)", strings.Join(searchParts, ",")))
	}

	if input.Test != nil {
		// LinkedIn expects this as a dedicated finder parameter, not inside search=(...).
		params = append(params, fmt.Sprintf("search.test=%t", *input.Test))
	}

	if input.Start > 0 {
		params = append(params, fmt.Sprintf("start=%d", input.Start))
	}

	if input.Count > 0 {
		params = append(params, fmt.Sprintf("count=%d", input.Count))
	}

	return strings.Join(params, "&")
}
