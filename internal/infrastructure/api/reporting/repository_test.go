package reporting

import "testing"

func TestParseDateObject_DecodesFullDate(t *testing.T) {
	// encoding/json decodes JSON numbers into float64 by default.
	raw := map[string]interface{}{
		"year":  float64(2025),
		"month": float64(3),
		"day":   float64(15),
	}

	date, ok := parseDateObject(raw)
	if !ok {
		t.Fatal("expected parseDateObject to succeed for a full date payload")
	}
	if date.Year != 2025 || date.Month != 3 || date.Day != 15 {
		t.Fatalf("unexpected date: %+v", date)
	}
}

func TestParseDateObject_ReturnsFalseOnMissingField(t *testing.T) {
	raw := map[string]interface{}{
		"year":  float64(2025),
		"month": float64(3),
		// day missing
	}

	if _, ok := parseDateObject(raw); ok {
		t.Fatal("expected parseDateObject to fail when day is missing")
	}
}

func TestParseDateObject_ReturnsFalseOnNonObject(t *testing.T) {
	if _, ok := parseDateObject("2025-03-15"); ok {
		t.Fatal("expected parseDateObject to reject a non-object value")
	}
	if _, ok := parseDateObject(nil); ok {
		t.Fatal("expected parseDateObject to reject a nil value")
	}
}

func TestParseDateObject_ReturnsFalseOnNonNumericField(t *testing.T) {
	raw := map[string]interface{}{
		"year":  "2025",
		"month": float64(3),
		"day":   float64(15),
	}

	if _, ok := parseDateObject(raw); ok {
		t.Fatal("expected parseDateObject to fail when a field is not numeric")
	}
}
