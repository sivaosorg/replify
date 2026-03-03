package crontask

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse converts a cron expression string into a Schedule that can be used
// to compute successive activation times.
//
// Supported formats
//
//   - Five fields:  "minute hour day-of-month month day-of-week"
//   - Six fields:   "second minute hour day-of-month month day-of-week"
//   - @alias:       see the Alias constants for the full list
//   - @every <d>:   interval expression, e.g. "@every 5m"
//
// Field syntax (per field)
//
//   - *                — every value in the valid range
//   - n                — exact value
//   - n-m              — inclusive range
//   - n-m/step         — range with step
//   - */step           — every step values across the full range
//   - a,b,c            — comma-separated list (each element may itself use
//     any of the above forms)
//
// Month and day-of-week fields additionally accept three-letter English
// abbreviations (jan-dec and sun-sat respectively), case-insensitive.
//
// Timezone
//
// An optional IANA timezone specifier may appear at the front of the
// expression, separated from the fields by a space:
//
//	"TZ=America/New_York 0 9 * * 1-5"
//
// When a timezone is provided, the returned Schedule activates at the
// specified local time. When omitted, UTC is used.
func Parse(expr string) (Schedule, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, &ExpressionError{Expression: expr, Field: -1, Reason: "expression is empty"}
	}

	// Extract optional leading TZ= specifier.
	loc := time.UTC
	if strings.HasPrefix(expr, "TZ=") {
		idx := strings.Index(expr, " ")
		if idx < 0 {
			return nil, &ExpressionError{Expression: expr, Field: -1, Reason: "missing fields after TZ specifier"}
		}
		tzName := expr[3:idx]
		var err error
		loc, err = time.LoadLocation(tzName)
		if err != nil {
			return nil, &ExpressionError{Expression: expr, Field: -1, Reason: fmt.Sprintf("unknown timezone %q: %v", tzName, err)}
		}
		expr = strings.TrimSpace(expr[idx+1:])
	}

	// Handle @every interval expressions.
	if strings.HasPrefix(expr, "@every ") {
		durationStr := strings.TrimPrefix(expr, "@every ")
		d, err := time.ParseDuration(strings.TrimSpace(durationStr))
		if err != nil {
			return nil, &ExpressionError{Expression: expr, Field: -1, Reason: fmt.Sprintf("invalid duration %q: %v", durationStr, err)}
		}
		if d <= 0 {
			return nil, &ExpressionError{Expression: expr, Field: -1, Reason: "interval duration must be positive"}
		}
		return &intervalSchedule{interval: d}, nil
	}

	// Handle @alias expressions.
	if strings.HasPrefix(expr, "@") {
		expanded, ok := lookupAlias(strings.ToLower(expr))
		if !ok {
			return nil, &ExpressionError{Expression: expr, Field: -1, Reason: fmt.Sprintf("unknown alias %q", expr)}
		}
		expr = expanded
	}

	// Split into fields.
	fields := strings.Fields(expr)
	switch len(fields) {
	case 5:
		return parseFiveField(expr, fields, loc)
	case 6:
		return parseSixField(expr, fields, loc)
	default:
		return nil, &ExpressionError{
			Expression: expr, Field: -1,
			Reason: fmt.Sprintf("expected 5 or 6 fields, got %d", len(fields)),
		}
	}
}

// Validate reports whether the given expression is syntactically and
// semantically valid. It returns nil on success or a typed *ExpressionError
// (which also satisfies errors.Is(err, ErrInvalidExpression)).
func Validate(expr string) error {
	_, err := Parse(expr)
	return err
}

// parseFiveField parses a standard five-field cron expression.
func parseFiveField(raw string, fields []string, loc *time.Location) (*cronSchedule, error) {
	s := &cronSchedule{
		minute:     make([]bool, 60),
		hour:       make([]bool, 24),
		dayOfMonth: make([]bool, 32), // index 0 unused; days 1..31
		month:      make([]bool, 13), // index 0 unused; months 1..12
		dayOfWeek:  make([]bool, 8),  // indices 0..6; 7 is normalised to 0
		loc:        loc,
	}
	specs := cronFields
	targets := [5]*[]bool{&s.minute, &s.hour, &s.dayOfMonth, &s.month, &s.dayOfWeek}
	for i, f := range fields {
		if err := parseField(raw, f, i, specs[i], targets[i]); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// parseSixField parses a six-field (seconds-first) cron expression.
func parseSixField(raw string, fields []string, loc *time.Location) (*cronSchedule, error) {
	s := &cronSchedule{
		second:     make([]bool, 60),
		minute:     make([]bool, 60),
		hour:       make([]bool, 24),
		dayOfMonth: make([]bool, 32),
		month:      make([]bool, 13),
		dayOfWeek:  make([]bool, 8), // indices 0..6; 7 is normalised to 0
		loc:        loc,
	}
	specs := cronFieldsWithSeconds
	targets := [6]*[]bool{&s.second, &s.minute, &s.hour, &s.dayOfMonth, &s.month, &s.dayOfWeek}
	for i, f := range fields {
		if err := parseField(raw, f, i, specs[i], targets[i]); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// parseField populates the bits slice for a single cron field expression.
// The fieldIndex parameter is used only for error reporting.
func parseField(expr, field string, fieldIndex int, spec fieldSpec, bits *[]bool) error {
	// Comma-separated list.
	parts := strings.Split(field, ",")
	for _, part := range parts {
		if err := parsePart(expr, part, fieldIndex, spec, bits); err != nil {
			return err
		}
	}
	return nil
}

// parsePart handles a single, comma-free element of a cron field: *, n, n-m,
// n-m/step, */step.
func parsePart(expr, part string, fieldIndex int, spec fieldSpec, bits *[]bool) error {
	// Step suffix.
	step := 1
	if idx := strings.Index(part, "/"); idx >= 0 {
		stepStr := part[idx+1:]
		var err error
		step, err = strconv.Atoi(stepStr)
		if err != nil || step <= 0 {
			return newExpressionError(expr, fieldIndex, "invalid step %q in %q", stepStr, part)
		}
		part = part[:idx]
	}

	isMonth := spec.name == "month"
	isDow := spec.name == "day-of-week"

	// Wildcard.
	if part == "*" || part == "?" {
		for v := spec.min; v <= spec.max; v += step {
			setBit(bits, v, spec)
		}
		return nil
	}

	// Range or single value.
	rangeMin, rangeMax := spec.min, spec.max
	if idx := strings.Index(part, "-"); idx >= 0 {
		loStr, hiStr := part[:idx], part[idx+1:]
		lo, err := parseNamedOrInt(loStr, isMonth, isDow)
		if err != nil {
			return newExpressionError(expr, fieldIndex, "invalid range start %q: %v", loStr, err)
		}
		hi, err := parseNamedOrInt(hiStr, isMonth, isDow)
		if err != nil {
			return newExpressionError(expr, fieldIndex, "invalid range end %q: %v", hiStr, err)
		}
		if lo < spec.min || lo > spec.max {
			return newExpressionError(expr, fieldIndex, "range start %d out of [%d,%d]", lo, spec.min, spec.max)
		}
		if hi < spec.min || hi > spec.max {
			return newExpressionError(expr, fieldIndex, "range end %d out of [%d,%d]", hi, spec.min, spec.max)
		}
		if lo > hi {
			return newExpressionError(expr, fieldIndex, "range start %d is greater than end %d", lo, hi)
		}
		rangeMin, rangeMax = lo, hi
	} else {
		// Single value.
		v, err := parseNamedOrInt(part, isMonth, isDow)
		if err != nil {
			return newExpressionError(expr, fieldIndex, "invalid value %q: %v", part, err)
		}
		if v < spec.min || v > spec.max {
			return newExpressionError(expr, fieldIndex, "value %d out of [%d,%d]", v, spec.min, spec.max)
		}
		rangeMin, rangeMax = v, v
	}

	for v := rangeMin; v <= rangeMax; v += step {
		setBit(bits, v, spec)
	}
	return nil
}

// setBit marks index v as active in the bits slice. For fields with a
// minimum of 1 (day-of-month, month) the index is stored as-is; for
// day-of-week, value 7 is normalised to 0 (both represent Sunday).
func setBit(bits *[]bool, v int, spec fieldSpec) {
	if spec.name == "day-of-week" && v == 7 {
		v = 0
	}
	if v >= 0 && v < len(*bits) {
		(*bits)[v] = true
	}
}

// parseNamedOrInt attempts to interpret s as a named month/day-of-week
// abbreviation first, then falls back to integer parsing.
func parseNamedOrInt(s string, isMonth, isDow bool) (int, error) {
	lower := strings.ToLower(s)
	if isMonth {
		if n, ok := monthNames[lower]; ok {
			return n, nil
		}
	}
	if isDow {
		if n, ok := dowNames[lower]; ok {
			return n, nil
		}
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("not a number or recognised name: %q", s)
	}
	return n, nil
}
