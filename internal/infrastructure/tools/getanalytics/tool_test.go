package getanalytics

import (
	"strings"
	"testing"

	"linkedin-mcp/internal/infrastructure/api/reporting"
	"linkedin-mcp/internal/infrastructure/tools/getanalytics/dto"
)

func TestNormalizeFieldName_MapsV8Aliases(t *testing.T) {
	tool := &Tool{}

	cases := map[string]string{
		"SPEND":                          "costInLocalCurrency",
		"COST":                           "costInLocalCurrency",
		"CONVERSIONS":                    "externalWebsiteConversions",
		"LEADS":                          "oneClickLeads",
		"UNIQUE_IMPRESSIONS":             "approximateMemberReach",
		"APPROXIMATE_UNIQUE_IMPRESSIONS": "approximateMemberReach",
		"ENGAGEMENTS":                    "totalEngagements",
	}

	for input, expected := range cases {
		got := tool.normalizeFieldName(input)
		if got != expected {
			t.Fatalf("normalizeFieldName(%q) = %q; want %q", input, got, expected)
		}
	}
}

func TestValidateAndNormalizeInput_AppliesFieldNormalization(t *testing.T) {
	tool := &Tool{}

	input := dto.Input{
		AccountID:       "512247261",
		DateRangeStart:  dto.Date{Year: 2025, Month: 1, Day: 1},
		TimeGranularity: "MONTHLY",
		Fields:          []string{"SPEND", "CONVERSIONS", "UNIQUE_IMPRESSIONS"},
	}

	normalized, derived, err := tool.validateAndNormalizeInput(input)
	if err != nil {
		t.Fatalf("expected no validation error, got %v", err)
	}

	expected := []string{"costInLocalCurrency", "externalWebsiteConversions", "approximateMemberReach"}
	if len(normalized.Fields) != len(expected) {
		t.Fatalf("unexpected normalized fields length: got %d want %d", len(normalized.Fields), len(expected))
	}

	for i := range expected {
		if normalized.Fields[i] != expected[i] {
			t.Fatalf("field[%d] = %q; want %q", i, normalized.Fields[i], expected[i])
		}
	}
	if len(derived) != 0 {
		t.Fatalf("expected no derived fields, got %v", derived)
	}
}

func TestValidateAndNormalizeInput_RejectsHalfSpecifiedSortBy(t *testing.T) {
	tool := &Tool{}

	cases := []dto.Input{
		{
			AccountID:       "512247261",
			DateRangeStart:  dto.Date{Year: 2025, Month: 1, Day: 1},
			TimeGranularity: "MONTHLY",
			Fields:          []string{"impressions"},
			SortByField:     "IMPRESSIONS",
		},
		{
			AccountID:       "512247261",
			DateRangeStart:  dto.Date{Year: 2025, Month: 1, Day: 1},
			TimeGranularity: "MONTHLY",
			Fields:          []string{"impressions"},
			SortByOrder:     "DESCENDING",
		},
	}

	for i, input := range cases {
		if _, _, err := tool.validateAndNormalizeInput(input); err == nil {
			t.Fatalf("case %d: expected validation error for half-specified sortBy", i)
		} else if !strings.Contains(err.Error(), "sortByField and sortByOrder") {
			t.Fatalf("case %d: expected actionable error about sortBy pair; got %v", i, err)
		}
	}
}

func TestValidateAndNormalizeInput_AcceptsFullyOrAbsentSortBy(t *testing.T) {
	tool := &Tool{}

	cases := []dto.Input{
		{
			AccountID:       "512247261",
			DateRangeStart:  dto.Date{Year: 2025, Month: 1, Day: 1},
			TimeGranularity: "MONTHLY",
			Fields:          []string{"impressions"},
		},
		{
			AccountID:       "512247261",
			DateRangeStart:  dto.Date{Year: 2025, Month: 1, Day: 1},
			TimeGranularity: "MONTHLY",
			Fields:          []string{"impressions"},
			SortByField:     "IMPRESSIONS",
			SortByOrder:     "DESCENDING",
		},
	}

	for i, input := range cases {
		if _, _, err := tool.validateAndNormalizeInput(input); err != nil {
			t.Fatalf("case %d: expected no validation error, got %v", i, err)
		}
	}
}

func TestValidateAndNormalizeInput_SplitsRawAndDerivedFields(t *testing.T) {
	tool := &Tool{}

	input := dto.Input{
		AccountID:       "512247261",
		DateRangeStart:  dto.Date{Year: 2025, Month: 1, Day: 1},
		TimeGranularity: "MONTHLY",
		// Mix of raw fields, a derived camelCase name, and aliases (CTR,
		// COST_PER_CLICK). The duplicate clicks requirement should be deduped.
		Fields: []string{"impressions", "clicks", "clickThroughRate", "COST_PER_CLICK", "CTR"},
	}

	normalized, derived, err := tool.validateAndNormalizeInput(input)
	if err != nil {
		t.Fatalf("expected no validation error, got %v", err)
	}

	// Raw fields: impressions + clicks (from user) + costInLocalCurrency (from
	// costPerClick dep). clicks must not be duplicated.
	expectedRaw := map[string]bool{
		"impressions":         true,
		"clicks":              true,
		"costInLocalCurrency": true,
	}
	if len(normalized.Fields) != len(expectedRaw) {
		t.Fatalf("unexpected raw fields: got %v", normalized.Fields)
	}
	for _, f := range normalized.Fields {
		if !expectedRaw[f] {
			t.Fatalf("unexpected raw field %q in %v", f, normalized.Fields)
		}
	}

	// Derived fields: clickThroughRate (dedup from CTR) and costPerClick.
	expectedDerived := map[string]bool{
		"clickThroughRate": true,
		"costPerClick":     true,
	}
	if len(derived) != len(expectedDerived) {
		t.Fatalf("unexpected derived fields: got %v", derived)
	}
	for _, f := range derived {
		if !expectedDerived[f] {
			t.Fatalf("unexpected derived field %q in %v", f, derived)
		}
	}
}

func TestInjectDerivedMetrics_ComputesRatios(t *testing.T) {
	result := &reporting.AnalyticsResult{
		Elements: []reporting.AnalyticsElement{
			{
				Metrics: map[string]interface{}{
					"impressions":         float64(1000),
					"clicks":              float64(50),
					"costInLocalCurrency": "25.00",
				},
			},
		},
	}

	injectDerivedMetrics(result, []string{"clickThroughRate", "costPerClick"})

	metrics := result.Elements[0].Metrics
	ctr, ok := metrics["clickThroughRate"].(float64)
	if !ok || ctr != 0.05 {
		t.Fatalf("expected clickThroughRate=0.05, got %v", metrics["clickThroughRate"])
	}
	cpc, ok := metrics["costPerClick"].(float64)
	if !ok || cpc != 0.5 {
		t.Fatalf("expected costPerClick=0.5, got %v", metrics["costPerClick"])
	}
}

func TestInjectDerivedMetrics_NilOnZeroDenominator(t *testing.T) {
	result := &reporting.AnalyticsResult{
		Elements: []reporting.AnalyticsElement{
			{
				Metrics: map[string]interface{}{
					"impressions": float64(0),
					"clicks":      float64(0),
				},
			},
		},
	}

	injectDerivedMetrics(result, []string{"clickThroughRate"})

	if _, present := result.Elements[0].Metrics["clickThroughRate"]; !present {
		t.Fatalf("expected clickThroughRate key to be present")
	}
	if result.Elements[0].Metrics["clickThroughRate"] != nil {
		t.Fatalf("expected clickThroughRate to be nil on zero denominator, got %v", result.Elements[0].Metrics["clickThroughRate"])
	}
}

func TestInjectDerivedMetrics_NilWhenDependencyMissing(t *testing.T) {
	result := &reporting.AnalyticsResult{
		Elements: []reporting.AnalyticsElement{
			{
				Metrics: map[string]interface{}{
					"impressions": float64(100),
					// clicks missing
				},
			},
		},
	}

	injectDerivedMetrics(result, []string{"clickThroughRate"})

	if result.Elements[0].Metrics["clickThroughRate"] != nil {
		t.Fatalf("expected clickThroughRate nil when clicks missing, got %v", result.Elements[0].Metrics["clickThroughRate"])
	}
}
