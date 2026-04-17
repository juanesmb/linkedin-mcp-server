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

func TestBuildAnalyticsQuery_InjectsPivotValuesWhenPivotSet(t *testing.T) {
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

	if !strings.Contains(query, "pivotValues") {
		t.Fatalf("expected pivotValues projection when pivot is set, got query: %s", query)
	}
	// Must stay in the fields projection, not appear as a standalone query parameter.
	if strings.Contains(query, "&pivotValues=") || strings.Contains(query, "?pivotValues=") {
		t.Fatalf("pivotValues must be injected into fields, not as its own param: %s", query)
	}
}

func TestBuildAnalyticsQuery_InjectsDateRangeWhenTimeGranularityBucketed(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	for _, granularity := range []string{"DAILY", "MONTHLY", "YEARLY"} {
		query := qb.BuildAnalyticsQuery(AnalyticsInput{
			AccountID: "512247261",
			DateRange: DateRange{
				Start: Date{Year: 2025, Month: 1, Day: 1},
				End:   &Date{Year: 2025, Month: 1, Day: 7},
			},
			TimeGranularity: granularity,
			Fields:          []string{"impressions", "clicks"},
		})

		if !strings.Contains(query, "fields=impressions,clicks,dateRange") {
			t.Fatalf("expected dateRange appended to fields projection for %s, got: %s", granularity, query)
		}
	}
}

func TestBuildAnalyticsQuery_DoesNotInjectDateRangeForAllGranularity(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		Fields:          []string{"impressions", "clicks"},
	})

	// `dateRange=` as a request filter must be present, but the fields projection
	// must not include it (ALL returns a single aggregate with no per-bucket metadata).
	fieldsPart := ""
	for _, p := range strings.Split(query, "&") {
		if strings.HasPrefix(p, "fields=") {
			fieldsPart = p
		}
	}
	if fieldsPart == "" {
		t.Fatalf("expected fields= param, got: %s", query)
	}
	if strings.Contains(fieldsPart, "dateRange") {
		t.Fatalf("expected no dateRange in fields projection for ALL granularity, got: %s", fieldsPart)
	}
}

func TestBuildAnalyticsQuery_DoesNotDuplicateMetadataWhenCallerSupplied(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		Pivot:     "CAMPAIGN",
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
			End:   &Date{Year: 2025, Month: 3, Day: 31},
		},
		TimeGranularity: "MONTHLY",
		Fields:          []string{"impressions", "dateRange", "pivotValues", "clicks"},
	})

	if strings.Count(query, "dateRange,") > 1 || strings.Count(query, ",dateRange") > 1 {
		t.Fatalf("expected dateRange deduplicated in fields projection, got: %s", query)
	}
	if strings.Count(query, "pivotValues") != 1 {
		t.Fatalf("expected pivotValues deduplicated, got: %s", query)
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

func TestBuildAnalyticsQuery_SortBySerializedAsSingleRestLiTuple(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		SortBy:          SortBy{Field: "IMPRESSIONS", Order: "DESCENDING"},
		Fields:          []string{"impressions"},
	})

	expected := "sortBy=(field:IMPRESSIONS,order:DESCENDING)"
	if !strings.Contains(query, expected) {
		t.Fatalf("expected single RestLi sortBy tuple %q, got: %s", expected, query)
	}
	if strings.Count(query, "sortBy=") != 1 {
		t.Fatalf("expected exactly one sortBy param, got: %s", query)
	}
}

func TestBuildAnalyticsQuery_SortByOmittedWhenEmpty(t *testing.T) {
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		Fields:          []string{"impressions"},
	})

	if strings.Contains(query, "sortBy=") {
		t.Fatalf("expected no sortBy param when both field and order are empty, got: %s", query)
	}
}

func TestBuildAnalyticsQuery_SortByWithOnlyFieldStillEmitsSingleParam(t *testing.T) {
	// Defense-in-depth: the tool validator rejects half-specified sortBy, but the builder
	// should still emit a single well-formed param rather than two separate ones.
	qb := NewQueryBuilder("https://api.linkedin.com/rest")

	query := qb.BuildAnalyticsQuery(AnalyticsInput{
		AccountID: "512247261",
		DateRange: DateRange{
			Start: Date{Year: 2025, Month: 1, Day: 1},
		},
		TimeGranularity: "ALL",
		SortBy:          SortBy{Field: "IMPRESSIONS"},
		Fields:          []string{"impressions"},
	})

	expected := "sortBy=(field:IMPRESSIONS)"
	if !strings.Contains(query, expected) {
		t.Fatalf("expected field-only sortBy tuple %q, got: %s", expected, query)
	}
	if strings.Count(query, "sortBy=") != 1 {
		t.Fatalf("expected exactly one sortBy param, got: %s", query)
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
