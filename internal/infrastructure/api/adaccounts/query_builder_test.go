package adaccounts

import (
	"strings"
	"testing"
)

func TestBuildSearchAdAccountsQuery_SetsTestAsDedicatedParam(t *testing.T) {
	testFlag := false
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildSearchAdAccountsQuery(SearchInput{
		Status: []string{"ACTIVE"},
		Test:   &testFlag,
		Start:  0,
		Count:  100,
	})

	if strings.Contains(query, "test:false") {
		t.Fatalf("expected test filter outside search composite, got: %s", query)
	}
	if !strings.Contains(query, "search.test=false") {
		t.Fatalf("expected dedicated search.test param, got: %s", query)
	}
	if !strings.Contains(query, "search=(status:(values:List(ACTIVE)))") {
		t.Fatalf("expected status search composite, got: %s", query)
	}
}
