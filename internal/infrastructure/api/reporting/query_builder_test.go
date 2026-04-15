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

func TestBuildAnalyticsQuery_InjectsPivotValuesWithPivot(t *testing.T) {
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
		t.Fatalf("expected pivotValues when pivot is set, got query: %s", query)
	}
}
