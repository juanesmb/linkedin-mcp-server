package reporting

import (
	"strings"
	"testing"
)

func TestBuildAnalyticsQuery_DoesNotInjectPivotValuesWithoutPivot(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		DateRange: DateRange{
			Start: Date{Year: 2026, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		Fields:          []string{"impressions", "clicks", "costInLocalCurrency"},
	})

	if strings.Contains(query, "pivotValues") {
		t.Fatalf("expected no pivotValues when pivot is empty, got query: %s", query)
	}
}

func TestBuildAnalyticsQuery_DoesNotInjectPivotValuesWithPivot(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		Pivot:     "CAMPAIGN",
		DateRange: DateRange{
			Start: Date{Year: 2026, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		Fields:          []string{"impressions", "clicks"},
	})

	if strings.Contains(query, "pivotValues") {
		t.Fatalf("expected no implicit pivotValues injection, got query: %s", query)
	}
}

func TestBuildAnalyticsQuery_UsesPlainEnumSerialization(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		Pivot:     "CAMPAIGN",
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
		},
		TimeGranularity: "MONTHLY",
		Fields:          []string{"impressions"},
	})

	if !strings.Contains(query, "pivot=CAMPAIGN") {
		t.Fatalf("expected plain pivot enum, got query: %s", query)
	}
	if !strings.Contains(query, "timeGranularity=MONTHLY") {
		t.Fatalf("expected plain timeGranularity enum, got query: %s", query)
	}
	if strings.Contains(query, "pivot=(value:") || strings.Contains(query, "timeGranularity=(value:") {
		t.Fatalf("expected no wrapped enum syntax, got query: %s", query)
	}
}

func TestBuildAnalyticsQuery_UsesRestLiDateRangeSerialization(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
			End:   &Date{Year: 2026, Month: 4, Day: 15},
		},
		TimeGranularity: "ALL",
		Fields:          []string{"impressions"},
	})

	expected := "dateRange=(start:(year:2025,month:1,day:1),end:(year:2026,month:4,day:15))"
	if !strings.Contains(query, expected) {
		t.Fatalf("expected RestLi dateRange %q in query, got: %s", expected, query)
	}
}

func TestBuildAnalyticsQuery_UsesRestLiDateRangeStartOnlySerialization(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		Fields:          []string{"impressions"},
	})

	expected := "dateRange=(start:(year:2025,month:1,day:1))"
	if !strings.Contains(query, expected) {
		t.Fatalf("expected RestLi dateRange %q in query, got: %s", expected, query)
	}
}

func TestBuildAnalyticsQuery_UsesRestLiListFormattingForFacets(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		Campaigns: []string{"urn:li:sponsoredCampaign:474763193"},
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		Fields:          []string{"impressions", "clicks", "costInLocalCurrency"},
	})

	if !strings.Contains(query, "accounts=List(urn%3Ali%3AsponsoredAccount%3A512247261)") {
		t.Fatalf("expected RestLi accounts facet format, got query: %s", query)
	}
	if !strings.Contains(query, "campaigns=List(urn%3Ali%3AsponsoredCampaign%3A474763193)") {
		t.Fatalf("expected RestLi campaigns facet format, got query: %s", query)
	}
	if strings.Contains(query, "accounts=List(urn:li:sponsoredAccount:") || strings.Contains(query, "campaigns=List(urn:li:sponsoredCampaign:") {
		t.Fatalf("expected encoded URNs inside RestLi lists, got query: %s", query)
	}
}

func TestBuildAnalyticsQuery_MergesDefaultAndExplicitAccountsIntoSingleFacet(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		Accounts: []string{
			"urn:li:sponsoredAccount:512247261",
			"urn:li:sponsoredAccount:999999999",
		},
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		Fields:          []string{"impressions", "clicks"},
	})

	expected := "accounts=List(urn%3Ali%3AsponsoredAccount%3A512247261,urn%3Ali%3AsponsoredAccount%3A999999999)"
	if !strings.Contains(query, expected) {
		t.Fatalf("expected merged accounts facet %q, got query: %s", expected, query)
	}
	if strings.Count(query, "accounts=List(") != 1 {
		t.Fatalf("expected a single accounts facet, got query: %s", query)
	}
}
