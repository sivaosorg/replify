package crontask

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// aliasDescriptions maps well-known alias names (lower-case) to their
// human-readable English descriptions. Custom aliases registered via
// RegisterAlias are not present here; they fall through to the field-based
// explainer using the expanded expression.
var aliasDescriptions = map[string]string{
	"@yearly":         "Once a year, at midnight on January 1st",
	"@annually":       "Once a year, at midnight on January 1st",
	"@monthly":        "Once a month, at midnight on the 1st",
	"@weekly":         "Once a week, at midnight on Sunday",
	"@daily":          "Once a day, at midnight",
	"@midnight":       "Once a day, at midnight",
	"@hourly":         "Every hour",
	"@minutely":       "Every minute",
	"@weekdays":       "At midnight, Monday through Friday",
	"@weekends":       "At midnight, on Saturday and Sunday",
	"@businessdaily":  "At 09:00, Monday through Friday",
	"@businesshourly": "Every hour from 09:00 to 17:00, Monday through Friday",
	"@quarterly":      "Once a quarter, at midnight on the 1st",
	"@semimonthly":    "Twice a month, at midnight on the 1st and 15th",
	"@workhours":      "Every minute from 09:00 to 17:59, Monday through Friday",
	"@marketopen":     "At 09:30, Monday through Friday",
	"@marketclose":    "At 16:00, Monday through Friday",
}

// Explain converts a cron expression into a natural English description.
//
// Supported input forms:
//
//   - "@every 5m"               → "Every 5 minutes"
//   - "@daily", "@hourly", etc  → predefined descriptions for built-in aliases
//   - "*/30 * * * * *"          → "Every 30 seconds"
//   - "0 9 * * 1-5"             → "At 09:00, Monday through Friday"
//   - "TZ=..." prefixes         → described without the timezone qualifier
//
// Custom aliases registered via RegisterAlias are described by expanding them
// to their underlying expression and applying the field-based explainer.
//
// Explain returns ErrInvalidExpression (wrapped) for invalid input and never
// panics.
//
// Example:
//
//	desc, err := crontask.Explain("0 0 * * 1-5")
//	// desc == "At 00:00, Monday through Friday"
func Explain(expr string) (string, error) {
	trimmed := strings.TrimSpace(expr)
	// Validate first — reuse Parse to avoid duplicating validation logic.
	if _, err := Parse(trimmed); err != nil {
		return "", err
	}

	// Strip optional leading TZ= specifier; timezone does not affect the
	// "when" part of the description.
	clean := trimmed
	if strings.HasPrefix(clean, "TZ=") {
		idx := strings.Index(clean, " ")
		if idx >= 0 {
			clean = strings.TrimSpace(clean[idx+1:])
		}
	}

	// @every interval expressions.
	if strings.HasPrefix(clean, "@every ") {
		durStr := strings.TrimSpace(strings.TrimPrefix(clean, "@every "))
		d, _ := time.ParseDuration(durStr)
		return explainIntervalDuration(d), nil
	}

	// @alias expressions.
	if strings.HasPrefix(clean, "@") {
		lower := strings.ToLower(clean)
		// Check predefined descriptions first for the best output.
		if desc, ok := aliasDescriptions[lower]; ok {
			return desc, nil
		}
		// Custom alias — expand and fall through to the field-based explainer.
		if expanded, ok := lookupAlias(lower); ok {
			return explainFields(strings.Fields(expanded)), nil
		}
		return "", ErrInvalidExpression
	}

	return explainFields(strings.Fields(clean)), nil
}

// explainFields produces an English description from a pre-split field slice.
func explainFields(fields []string) string {
	switch len(fields) {
	case 5:
		return explainFiveFields(fields[0], fields[1], fields[2], fields[3], fields[4])
	case 6:
		return explainSixFields(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5])
	default:
		return "Custom schedule"
	}
}

// explainSixFields describes a six-field (seconds-first) expression.
func explainSixFields(sec, min, hr, dom, mon, dow string) string {
	// All wildcards — fires every second.
	if isWild(sec) && isWild(min) && isWild(hr) && isWild(dom) && isWild(mon) && isWild(dow) {
		return "Every second"
	}
	// Pure second-based patterns.
	if step := extractStep(sec); step > 0 &&
		isWild(min) && isWild(hr) && isWild(dom) && isWild(mon) && isWild(dow) {
		if step == 1 {
			return "Every second"
		}
		return fmt.Sprintf("Every %d seconds", step)
	}
	if n, ok := parseSingleInt(sec); ok &&
		isWild(min) && isWild(hr) && isWild(dom) && isWild(mon) && isWild(dow) {
		return fmt.Sprintf("At second :%02d of every minute", n)
	}

	// Fall through: describe the five remaining fields and annotate the second
	// field only when it is not zero or wildcard.
	base := explainFiveFields(min, hr, dom, mon, dow)
	if sec != "0" && !isWild(sec) {
		return base + fmt.Sprintf(" (at second %s)", sec)
	}
	return base
}

// explainFiveFields describes a standard five-field expression by combining a
// time-of-day description with optional day-of-week, day-of-month, and month
// qualifiers.
func explainFiveFields(min, hr, dom, mon, dow string) string {
	timeDesc := explainTimePart(min, hr)

	var qualifiers []string
	if !isWild(dom) {
		qualifiers = append(qualifiers, "on the "+describeDOMField(dom)+" of each month")
	}
	if !isWild(mon) {
		qualifiers = append(qualifiers, "in "+describeMonthField(mon))
	}
	if !isWild(dow) {
		qualifiers = append(qualifiers, describeDOWField(dow))
	}

	if len(qualifiers) == 0 {
		return timeDesc
	}
	return timeDesc + ", " + strings.Join(qualifiers, ", ")
}

// explainTimePart builds the time-of-day sentence fragment from the minute and
// hour fields. It covers the most common patterns and falls back to a generic
// description for unusual combinations.
func explainTimePart(min, hr string) string {
	minWild := isWild(min)
	hrWild := isWild(hr)
	minStep := extractStep(min)
	hrStep := extractStep(hr)

	// All wildcards: "Every minute".
	if minWild && hrWild {
		return "Every minute"
	}
	// "*/N * ..." — every N minutes.
	if minStep > 0 && hrWild {
		if minStep == 1 {
			return "Every minute"
		}
		return fmt.Sprintf("Every %d minutes", minStep)
	}
	// "0 */N ..." — every N hours (only when minute is exactly 0).
	if min == "0" && hrStep > 0 {
		if hrStep == 1 {
			return "Every hour"
		}
		return fmt.Sprintf("Every %d hours", hrStep)
	}
	// "0 * ..." — every hour.
	if min == "0" && hrWild {
		return "Every hour"
	}
	// Both minute and hour are single integers — "At HH:MM".
	if mInt, mOK := parseSingleInt(min); mOK {
		if hInt, hOK := parseSingleInt(hr); hOK {
			return fmt.Sprintf("At %02d:%02d", hInt, mInt)
		}
		// Single minute, wild hour — "Every hour at :MM".
		if hrWild {
			return fmt.Sprintf("Every hour at :%02d", mInt)
		}
	}
	// Wild minute, single hour — "Every minute of hour HH".
	if hInt, hOK := parseSingleInt(hr); hOK && minWild {
		return fmt.Sprintf("Every minute of hour %02d", hInt)
	}
	// Step minute, single hour — "Every N minutes during hour HH".
	if minStep > 0 {
		if hInt, hOK := parseSingleInt(hr); hOK {
			return fmt.Sprintf("Every %d minutes during hour %02d", minStep, hInt)
		}
	}
	// Hour range with minute 0 — "Every hour from HH:00 to HH:00".
	if hrRange, ok := parseRange(hr); ok && min == "0" {
		return fmt.Sprintf("Every hour from %02d:00 to %02d:00", hrRange[0], hrRange[1])
	}
	// Generic fallback.
	return fmt.Sprintf("At minute %s of hour %s", min, hr)
}

// explainIntervalDuration converts a Duration into an English phrase such as
// "Every 5 minutes" or "Every 30 seconds".
func explainIntervalDuration(d time.Duration) string {
	if d < time.Minute {
		secs := int(d.Seconds())
		if secs == 1 {
			return "Every second"
		}
		return fmt.Sprintf("Every %d seconds", secs)
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		if mins == 1 {
			return "Every minute"
		}
		return fmt.Sprintf("Every %d minutes", mins)
	}
	if d < 24*time.Hour {
		hrs := int(d.Hours())
		if hrs == 1 {
			return "Every hour"
		}
		return fmt.Sprintf("Every %d hours", hrs)
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "Every day"
	}
	return fmt.Sprintf("Every %d days", days)
}

// describeDOWField returns an English description of a day-of-week field value.
// Common patterns are handled with fast-path switch cases; the general path
// parses the field as a list of numeric or named values.
func describeDOWField(field string) string {
	lower := strings.ToLower(field)
	switch lower {
	case "1-5", "mon-fri":
		return "Monday through Friday"
	case "0", "7", "sun":
		return "Sunday"
	case "1", "mon":
		return "Monday"
	case "2", "tue":
		return "Tuesday"
	case "3", "wed":
		return "Wednesday"
	case "4", "thu":
		return "Thursday"
	case "5", "fri":
		return "Friday"
	case "6", "sat":
		return "Saturday"
	}
	// Numeric range "lo-hi".
	if r, ok := parseRange(field); ok && r[0] >= 0 && r[0] <= 6 && r[1] >= 0 && r[1] <= 6 {
		return dowLabel(r[0]) + " through " + dowLabel(r[1])
	}
	// Comma-separated list — process in input order so the caller controls ordering.
	parts := strings.Split(field, ",")
	if len(parts) > 1 {
		names := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if n, err := strconv.Atoi(p); err == nil {
				if n == 7 {
					n = 0
				}
				names = append(names, dowLabel(n))
			} else {
				// Named value: capitalise and use as-is.
				names = append(names, capitalizeFirst(p))
			}
		}
		return joinEnglish(names)
	}
	return field
}

// describeDOMField returns an English description of a day-of-month field.
func describeDOMField(field string) string {
	if n, ok := parseSingleInt(field); ok {
		return ordinal(n)
	}
	if r, ok := parseRange(field); ok {
		return ordinal(r[0]) + " through the " + ordinal(r[1])
	}
	parts := strings.Split(field, ",")
	if len(parts) > 1 {
		ords := make([]string, len(parts))
		for i, p := range parts {
			if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
				ords[i] = ordinal(n)
			} else {
				ords[i] = p
			}
		}
		return joinEnglish(ords)
	}
	return field
}

// describeMonthField returns an English description of a month field.
func describeMonthField(field string) string {
	monthLabels := [13]string{
		"", "January", "February", "March", "April",
		"May", "June", "July", "August", "September",
		"October", "November", "December",
	}
	lower := strings.ToLower(field)
	if n, ok := monthNames[lower]; ok {
		return monthLabels[n]
	}
	if n, ok := parseSingleInt(field); ok && n >= 1 && n <= 12 {
		return monthLabels[n]
	}
	if r, ok := parseRange(field); ok && r[0] >= 1 && r[0] <= 12 && r[1] >= 1 && r[1] <= 12 {
		return monthLabels[r[0]] + " through " + monthLabels[r[1]]
	}
	parts := strings.Split(field, ",")
	if len(parts) > 1 {
		names := make([]string, len(parts))
		for i, p := range parts {
			p = strings.TrimSpace(p)
			if n, ok2 := monthNames[strings.ToLower(p)]; ok2 {
				names[i] = monthLabels[n]
			} else if n, err := strconv.Atoi(p); err == nil && n >= 1 && n <= 12 {
				names[i] = monthLabels[n]
			} else {
				names[i] = p
			}
		}
		return joinEnglish(names)
	}
	return field
}

// isWild reports whether a cron field token represents "all values".
func isWild(f string) bool {
	return f == "*" || f == "?"
}

// extractStep extracts the step value from a "*/n" field token. It returns 0
// when the field is not a step expression or the step is invalid.
func extractStep(f string) int {
	if !strings.HasPrefix(f, "*/") {
		return 0
	}
	n, err := strconv.Atoi(f[2:])
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

// parseSingleInt parses a bare integer field token.
func parseSingleInt(f string) (int, bool) {
	n, err := strconv.Atoi(f)
	return n, err == nil
}

// parseRange parses a "lo-hi" field token into its two components.
func parseRange(f string) ([2]int, bool) {
	idx := strings.Index(f, "-")
	if idx < 0 {
		return [2]int{}, false
	}
	lo, err1 := strconv.Atoi(f[:idx])
	hi, err2 := strconv.Atoi(f[idx+1:])
	if err1 != nil || err2 != nil {
		return [2]int{}, false
	}
	return [2]int{lo, hi}, true
}

// dowLabel returns the English day name for a 0-based weekday index.
func dowLabel(n int) string {
	labels := [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	if n >= 0 && n < 7 {
		return labels[n]
	}
	return strconv.Itoa(n)
}

// ordinal converts a positive integer to its English ordinal string
// (1 → "1st", 2 → "2nd", 3 → "3rd", 4 → "4th", …).
func ordinal(n int) string {
	switch n % 100 {
	case 11, 12, 13:
		return fmt.Sprintf("%dth", n)
	}
	switch n % 10 {
	case 1:
		return fmt.Sprintf("%dst", n)
	case 2:
		return fmt.Sprintf("%dnd", n)
	case 3:
		return fmt.Sprintf("%drd", n)
	default:
		return fmt.Sprintf("%dth", n)
	}
}

// joinEnglish joins a slice of strings with commas and an Oxford "and".
//
//	[]string{"a"} → "a"
//	[]string{"a", "b"} → "a and b"
//	[]string{"a", "b", "c"} → "a, b, and c"
func joinEnglish(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " and " + items[1]
	default:
		return strings.Join(items[:len(items)-1], ", ") + ", and " + items[len(items)-1]
	}
}

// capitalizeFirst returns s with its first rune upper-cased and the rest
// lower-cased. It is used to normalise day and month names that arrive from
// user input.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
