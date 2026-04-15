package getanalytics

import (
	"testing"

	"linkedin-mcp/internal/infrastructure/tools/getanalytics/dto"
)

func TestNormalizeFieldName_MapsV8Aliases(t *testing.T) {
	tool := &Tool{}

	cases := map[string]string{
		"SPEND":                         "costInLocalCurrency",
		"COST":                          "costInLocalCurrency",
		"CONVERSIONS":                   "externalWebsiteConversions",
		"LEADS":                         "oneClickLeads",
		"UNIQUE_IMPRESSIONS":            "approximateMemberReach",
		"APPROXIMATE_UNIQUE_IMPRESSIONS": "approximateMemberReach",
		"ENGAGEMENTS":                   "totalEngagements",
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

	normalized, err := tool.validateAndNormalizeInput(input)
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
}
