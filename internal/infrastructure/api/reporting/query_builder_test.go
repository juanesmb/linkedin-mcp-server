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

func TestBuildAnalyticsQuery_UsesDottedDateRangeSerialization(t *testing.T) {
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

	expected := []string{
		"dateRange.start.year=2025",
		"dateRange.start.month=1",
		"dateRange.start.day=1",
		"dateRange.end.year=2026",
		"dateRange.end.month=4",
		"dateRange.end.day=15",
	}
	for _, part := range expected {
		if !strings.Contains(query, part) {
			t.Fatalf("expected %q in query, got: %s", part, query)
		}
	}

	if strings.Contains(query, "dateRange=(") {
		t.Fatalf("expected no tuple-style dateRange serialization, got query: %s", query)
	}
}
