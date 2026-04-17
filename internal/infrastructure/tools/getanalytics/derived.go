package getanalytics

import (
	"strconv"
	"strings"
)

// derivedMetric defines a metric the MCP server computes client-side from raw
// LinkedIn fields. These are arithmetic-only definitions so they do not need
// to be maintained when LinkedIn evolves the AdAnalytics schema.
type derivedMetric struct {
	// Canonical name emitted in the response under element.Metrics.
	Name string
	// Raw LinkedIn fields the computation depends on. These are unioned into
	// the outbound request whenever the caller asks for a derived metric.
	RequiredFields []string
	// Compute returns the derived value and a boolean indicating whether the
	// value is defined for the given element (false when a denominator is
	// zero or a required raw field is missing).
	Compute func(metrics map[string]interface{}) (float64, bool)
}

// Derived metric definitions. Keys must be canonical lowerCamelCase names.
// Aliases (UPPER_SNAKE, shorthand) are declared separately in derivedAliases
// so the canonical set stays small and self-contained.
var derivedMetrics = map[string]derivedMetric{
	"costPerClick": {
		Name:           "costPerClick",
		RequiredFields: []string{"costInLocalCurrency", "clicks"},
		Compute:        ratio("costInLocalCurrency", "clicks"),
	},
	"clickThroughRate": {
		Name:           "clickThroughRate",
		RequiredFields: []string{"clicks", "impressions"},
		Compute:        ratio("clicks", "impressions"),
	},
	"costPerLead": {
		Name:           "costPerLead",
		RequiredFields: []string{"costInLocalCurrency", "oneClickLeads"},
		Compute:        ratio("costInLocalCurrency", "oneClickLeads"),
	},
	"costPerMille": {
		Name:           "costPerMille",
		RequiredFields: []string{"costInLocalCurrency", "impressions"},
		Compute: func(metrics map[string]interface{}) (float64, bool) {
			cost, okCost := metricAsFloat(metrics["costInLocalCurrency"])
			impressions, okImpr := metricAsFloat(metrics["impressions"])
			if !okCost || !okImpr || impressions == 0 {
				return 0, false
			}
			return cost * 1000.0 / impressions, true
		},
	},
	"videoCompletionRate": {
		Name:           "videoCompletionRate",
		RequiredFields: []string{"videoCompletions", "videoViews"},
		Compute:        ratio("videoCompletions", "videoViews"),
	},
}

// derivedAliases maps user-supplied field names (any case/underscore variant)
// to the canonical derived-metric key. Unknown aliases fall through to the
// raw-field passthrough path so LinkedIn remains the source of truth for
// schema fields.
var derivedAliases = map[string]string{
	"COST_PER_CLICK":        "costPerClick",
	"COSTPERCLICK":          "costPerClick",
	"CPC":                   "costPerClick",
	"CLICK_THROUGH_RATE":    "clickThroughRate",
	"CLICKTHROUGHRATE":      "clickThroughRate",
	"CLICK_THRU_RATE":       "clickThroughRate",
	"CTR":                   "clickThroughRate",
	"COST_PER_LEAD":         "costPerLead",
	"COSTPERLEAD":           "costPerLead",
	"CPL":                   "costPerLead",
	"COST_PER_MILLE":        "costPerMille",
	"COSTPERMILLE":          "costPerMille",
	"CPM":                   "costPerMille",
	"CPA":                   "costPerLead", // common interchangeable shorthand for cost per action/lead
	"VIDEO_COMPLETION_RATE": "videoCompletionRate",
	"VIDEOCOMPLETIONRATE":   "videoCompletionRate",
}

// lookupDerivedMetric resolves a user-supplied field name to a derived metric
// definition. It accepts both canonical lowerCamelCase names and the aliases
// in derivedAliases. Returns (metric, true) on a hit.
func lookupDerivedMetric(field string) (derivedMetric, bool) {
	trimmed := strings.TrimSpace(field)
	if trimmed == "" {
		return derivedMetric{}, false
	}

	if metric, ok := derivedMetrics[trimmed]; ok {
		return metric, true
	}

	if canonical, ok := derivedAliases[strings.ToUpper(trimmed)]; ok {
		return derivedMetrics[canonical], true
	}

	return derivedMetric{}, false
}

// ratio builds a Compute function that divides numeratorKey by denominatorKey.
// Returns (0, false) if either field is missing or the denominator is zero.
func ratio(numeratorKey, denominatorKey string) func(metrics map[string]interface{}) (float64, bool) {
	return func(metrics map[string]interface{}) (float64, bool) {
		numerator, okNum := metricAsFloat(metrics[numeratorKey])
		denominator, okDen := metricAsFloat(metrics[denominatorKey])
		if !okNum || !okDen || denominator == 0 {
			return 0, false
		}
		return numerator / denominator, true
	}
}

// metricAsFloat coerces a LinkedIn metric value (number, numeric string) to a
// float64. LinkedIn returns monetary amounts as JSON strings (e.g. "17.16")
// and counts as numbers (e.g. 677), so both shapes must be accepted.
func metricAsFloat(value interface{}) (float64, bool) {
	switch typed := value.(type) {
	case nil:
		return 0, false
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case string:
		trimmed := strings.TrimSpace(typed)
		if trimmed == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}
