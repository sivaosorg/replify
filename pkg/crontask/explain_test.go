package crontask

import (
	"errors"
	"testing"
)

// TestExplain_AtEvery exercises the @every interval syntax.
func TestExplain_AtEvery(t *testing.T) {
	t.Parallel()
	cases := []struct {
		expr string
		want string
	}{
		{"@every 1s", "Every second"},
		{"@every 30s", "Every 30 seconds"},
		{"@every 1m", "Every minute"},
		{"@every 5m", "Every 5 minutes"},
		{"@every 1h", "Every hour"},
		{"@every 6h", "Every 6 hours"},
		{"@every 24h", "Every day"},
		{"@every 48h", "Every 2 days"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.expr, func(t *testing.T) {
			t.Parallel()
			got, err := Explain(tc.expr)
			if err != nil {
				t.Fatalf("Explain(%q): unexpected error: %v", tc.expr, err)
			}
			if got != tc.want {
				t.Errorf("Explain(%q) = %q, want %q", tc.expr, got, tc.want)
			}
		})
	}
}

// TestExplain_BuiltinAliases verifies predefined alias descriptions.
func TestExplain_BuiltinAliases(t *testing.T) {
	t.Parallel()
	cases := []struct {
		alias string
		want  string
	}{
		{"@yearly", "Once a year, at midnight on January 1st"},
		{"@annually", "Once a year, at midnight on January 1st"},
		{"@monthly", "Once a month, at midnight on the 1st"},
		{"@weekly", "Once a week, at midnight on Sunday"},
		{"@daily", "Once a day, at midnight"},
		{"@midnight", "Once a day, at midnight"},
		{"@hourly", "Every hour"},
		{"@minutely", "Every minute"},
		{"@weekdays", "At midnight, Monday through Friday"},
		{"@weekends", "At midnight, on Saturday and Sunday"},
		{"@businessDaily", "At 09:00, Monday through Friday"},
		{"@businessHourly", "Every hour from 09:00 to 17:00, Monday through Friday"},
		{"@quarterly", "Once a quarter, at midnight on the 1st"},
		{"@semiMonthly", "Twice a month, at midnight on the 1st and 15th"},
		{"@workhours", "Every minute from 09:00 to 17:59, Monday through Friday"},
		{"@marketOpen", "At 09:30, Monday through Friday"},
		{"@marketClose", "At 16:00, Monday through Friday"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.alias, func(t *testing.T) {
			t.Parallel()
			got, err := Explain(tc.alias)
			if err != nil {
				t.Fatalf("Explain(%q): unexpected error: %v", tc.alias, err)
			}
			if got != tc.want {
				t.Errorf("Explain(%q) = %q, want %q", tc.alias, got, tc.want)
			}
		})
	}
}

// TestExplain_FiveField covers the most common five-field patterns.
func TestExplain_FiveField(t *testing.T) {
	t.Parallel()
	cases := []struct {
		expr string
		want string
	}{
		// All wildcards.
		{"* * * * *", "Every minute"},
		// Step minute.
		{"*/15 * * * *", "Every 15 minutes"},
		{"*/1 * * * *", "Every minute"},
		// Step hour.
		{"0 */6 * * *", "Every 6 hours"},
		{"0 */1 * * *", "Every hour"},
		// Specific time.
		{"0 9 * * *", "At 09:00"},
		{"30 9 * * *", "At 09:30"},
		{"0 0 * * *", "At 00:00"},
		// Hourly (minute=0, hour=wild).
		{"0 * * * *", "Every hour"},
		// Every hour at :30.
		{"30 * * * *", "Every hour at :30"},
		// DOW restriction.
		{"0 0 * * 1-5", "At 00:00, Monday through Friday"},
		{"0 9 * * 1-5", "At 09:00, Monday through Friday"},
		{"0 9 * * 1", "At 09:00, Monday"},
		{"0 9 * * 0", "At 09:00, Sunday"},
		// Month restriction.
		{"0 0 1 1 *", "At 00:00, on the 1st of each month, in January"},
		{"0 0 1 1,4,7,10 *", "At 00:00, on the 1st of each month, in January, April, July, and October"},
		// DOM restriction.
		{"0 0 15 * *", "At 00:00, on the 15th of each month"},
		{"0 0 1,15 * *", "At 00:00, on the 1st and 15th of each month"},
		// Hour range.
		{"0 9-17 * * *", "Every hour from 09:00 to 17:00"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.expr, func(t *testing.T) {
			t.Parallel()
			got, err := Explain(tc.expr)
			if err != nil {
				t.Fatalf("Explain(%q): unexpected error: %v", tc.expr, err)
			}
			if got != tc.want {
				t.Errorf("Explain(%q) = %q, want %q", tc.expr, got, tc.want)
			}
		})
	}
}

// TestExplain_SixField covers common six-field (seconds-first) patterns.
func TestExplain_SixField(t *testing.T) {
	t.Parallel()
	cases := []struct {
		expr string
		want string
	}{
		{"* * * * * *", "Every second"},
		{"*/30 * * * * *", "Every 30 seconds"},
		{"*/1 * * * * *", "Every second"},
		{"0 0 9 * * 1-5", "At 09:00, Monday through Friday"},
		{"30 0 9 * * 1-5", "At 09:00, Monday through Friday (at second 30)"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.expr, func(t *testing.T) {
			t.Parallel()
			got, err := Explain(tc.expr)
			if err != nil {
				t.Fatalf("Explain(%q): unexpected error: %v", tc.expr, err)
			}
			if got != tc.want {
				t.Errorf("Explain(%q) = %q, want %q", tc.expr, got, tc.want)
			}
		})
	}
}

// TestExplain_WithTZ verifies that the TZ= prefix is stripped gracefully.
func TestExplain_WithTZ(t *testing.T) {
	t.Parallel()
	got, err := Explain("TZ=America/New_York 0 9 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "At 09:00, Monday through Friday"
	if got != want {
		t.Errorf("Explain with TZ = %q, want %q", got, want)
	}
}

// TestExplain_InvalidExpr verifies that Explain returns a typed error for
// invalid input.
func TestExplain_InvalidExpr(t *testing.T) {
	t.Parallel()
	_, err := Explain("bad expression here")
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
	if !errors.Is(err, ErrInvalidExpression) {
		t.Errorf("errors.Is(err, ErrInvalidExpression) = false, want true")
	}
}

// TestExplain_CustomAlias verifies that a custom alias registered via
// RegisterAlias is described by the field-based explainer.
func TestExplain_CustomAlias(t *testing.T) {
	t.Parallel()
	if err := RegisterAlias("@nightly", "0 2 * * *"); err != nil {
		t.Fatalf("RegisterAlias: %v", err)
	}
	got, err := Explain("@nightly")
	if err != nil {
		t.Fatalf("Explain(@nightly): %v", err)
	}
	want := "At 02:00"
	if got != want {
		t.Errorf("Explain(@nightly) = %q, want %q", got, want)
	}

	// Cleanup so other tests are not affected.
	aliasMapMu.Lock()
	delete(aliasMap, "@nightly")
	aliasMapMu.Unlock()
}

// TestExplain_DOWVariants exercises various day-of-week field formats.
func TestExplain_DOWVariants(t *testing.T) {
	t.Parallel()
	cases := []struct {
		expr string
		want string
	}{
		{"0 9 * * mon-fri", "At 09:00, Monday through Friday"},
		{"0 9 * * 0,6", "At 09:00, Sunday and Saturday"},
		{"0 9 * * 1,3,5", "At 09:00, Monday, Wednesday, and Friday"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.expr, func(t *testing.T) {
			t.Parallel()
			got, err := Explain(tc.expr)
			if err != nil {
				t.Fatalf("Explain(%q): unexpected error: %v", tc.expr, err)
			}
			if got != tc.want {
				t.Errorf("Explain(%q) = %q, want %q", tc.expr, got, tc.want)
			}
		})
	}
}

// TestExplain_OrdinalFormats verifies ordinal suffixes for day-of-month.
func TestExplain_OrdinalFormats(t *testing.T) {
	t.Parallel()
	cases := []struct {
		dom  int
		want string
	}{
		{1, "1st"}, {2, "2nd"}, {3, "3rd"}, {4, "4th"},
		{11, "11th"}, {12, "12th"}, {13, "13th"},
		{21, "21st"}, {22, "22nd"}, {23, "23rd"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()
			got := ordinal(tc.dom)
			if got != tc.want {
				t.Errorf("ordinal(%d) = %q, want %q", tc.dom, got, tc.want)
			}
		})
	}
}
