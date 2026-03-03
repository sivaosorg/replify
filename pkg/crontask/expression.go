package crontask

import (
	"sync"
	"time"
)

// Schedule is the interface implemented by any type that can compute the next
// activation time for a scheduled job. The single method Next receives the
// current (or reference) time and returns the earliest future time at which
// the job should run next.
//
// Implementing Schedule allows users to inject custom scheduling logic — for
// example, a schedule driven by an external calendar API — without forking the
// package.
type Schedule interface {
	// Next returns the next activation time after the given reference time t.
	// If no further activation exists (e.g. a once-only schedule in the past),
	// Next returns the zero time.Time.
	Next(t time.Time) time.Time
}

// Alias is a short-hand name for a commonly used cron schedule, inspired by
// both Vixie-cron and gronx. Aliases are resolved by the parser before
// field-level parsing occurs.
//
// The following aliases are recognised out of the box:
//
//	@yearly   (or @annually)  — "0 0 1 1 *"
//	@monthly                  — "0 0 1 * *"
//	@weekly                   — "0 0 * * 0"
//	@daily    (or @midnight)  — "0 0 * * *"
//	@hourly                   — "0 * * * *"
//	@minutely                 — "* * * * *"
//	@weekdays                 — "0 0 * * 1-5"
//	@weekends                 — "0 0 * * 0,6"
//
// Business-oriented aliases:
//
//	@businessDaily   — "0 9 * * 1-5"       (09:00 on weekdays)
//	@businessHourly  — "0 9-17 * * 1-5"    (top of each hour, 09–17, weekdays)
//	@quarterly       — "0 0 1 1,4,7,10 *"  (midnight, first day of each quarter)
//	@semiMonthly     — "0 0 1,15 * *"      (midnight, 1st and 15th of each month)
//	@workhours       — "* 9-17 * * 1-5"    (every minute 09:00–17:59, weekdays)
//	@marketOpen      — "30 9 * * 1-5"      (09:30, weekdays)
//	@marketClose     — "0 16 * * 1-5"      (16:00, weekdays)
//
// Additional aliases can be registered at runtime with RegisterAlias.
//
// When using the six-field (seconds-first) format, the above expansions gain
// a leading "0" second field automatically.
type Alias = string

const (
	// AliasYearly fires once a year, at midnight on 1 January.
	AliasYearly Alias = "@yearly"

	// AliasAnnually is a synonym for AliasYearly.
	AliasAnnually Alias = "@annually"

	// AliasMonthly fires once a month, at midnight on the 1st.
	AliasMonthly Alias = "@monthly"

	// AliasWeekly fires once a week, at midnight on Sunday.
	AliasWeekly Alias = "@weekly"

	// AliasDaily fires once a day, at midnight.
	AliasDaily Alias = "@daily"

	// AliasMidnight is a synonym for AliasDaily.
	AliasMidnight Alias = "@midnight"

	// AliasHourly fires once an hour, at the top of the hour.
	AliasHourly Alias = "@hourly"

	// AliasMinutely fires once a minute, at the top of each minute.
	AliasMinutely Alias = "@minutely"

	// AliasWeekdays fires at midnight on every weekday (Monday–Friday).
	AliasWeekdays Alias = "@weekdays"

	// AliasWeekends fires at midnight on Saturday and Sunday.
	AliasWeekends Alias = "@weekends"

	// AliasBusinessDaily fires at 09:00 on every weekday (Monday–Friday).
	// It is suitable for once-per-business-day jobs that run at the start of
	// the working day.
	AliasBusinessDaily Alias = "@businessDaily"

	// AliasBusinessHourly fires at the top of each hour from 09:00 to 17:00
	// on every weekday. It is suitable for tasks that should run once per
	// business hour during core working hours.
	AliasBusinessHourly Alias = "@businessHourly"

	// AliasQuarterly fires at midnight on the first day of each calendar
	// quarter (January, April, July, and October).
	AliasQuarterly Alias = "@quarterly"

	// AliasSemiMonthly fires at midnight on the 1st and 15th of every month,
	// giving two activations per month.
	AliasSemiMonthly Alias = "@semiMonthly"

	// AliasWorkHours fires every minute during business hours: 09:00–17:59
	// on every weekday. Suitable for polling tasks that should only run
	// during office hours.
	AliasWorkHours Alias = "@workhours"

	// AliasMarketOpen fires at 09:30 on every weekday, aligned with the
	// standard US equity market open time.
	AliasMarketOpen Alias = "@marketOpen"

	// AliasMarketClose fires at 16:00 on every weekday, aligned with the
	// standard US equity market close time.
	AliasMarketClose Alias = "@marketClose"
)

// aliasMapMu guards aliasMap for concurrent reads and writes. The initial
// entries are populated at package-init time (single-goroutine), after which
// all reads go through lookupAlias and all writes through RegisterAlias.
var aliasMapMu sync.RWMutex

// aliasMap maps each recognised alias to its canonical five-field (minute-first)
// cron expression. The parser consults this map before attempting field-level
// parsing.
var aliasMap = map[string]string{
	// Standard aliases — keys stored in lower case to match lookupAlias.
	"@yearly":   "0 0 1 1 *",
	"@annually": "0 0 1 1 *",
	"@monthly":  "0 0 1 * *",
	"@weekly":   "0 0 * * 0",
	"@daily":    "0 0 * * *",
	"@midnight": "0 0 * * *",
	"@hourly":   "0 * * * *",
	"@minutely": "* * * * *",
	"@weekdays": "0 0 * * 1-5",
	"@weekends": "0 0 * * 0,6",
	// Business-oriented aliases.
	"@businessdaily":  "0 9 * * 1-5",
	"@businesshourly": "0 9-17 * * 1-5",
	"@quarterly":      "0 0 1 1,4,7,10 *",
	"@semimonthly":    "0 0 1,15 * *",
	"@workhours":      "* 9-17 * * 1-5",
	"@marketopen":     "30 9 * * 1-5",
	"@marketclose":    "0 16 * * 1-5",
}

// lookupAlias performs a concurrent-safe lookup in aliasMap using the
// lower-cased alias name (including the "@" prefix).
func lookupAlias(name string) (string, bool) {
	aliasMapMu.RLock()
	v, ok := aliasMap[name]
	aliasMapMu.RUnlock()
	return v, ok
}

// fieldSpec describes the valid range for a single cron field.
type fieldSpec struct {
	min  int
	max  int
	name string
}

// cronFields defines the valid ranges for each position in a five-field
// (minute, hour, day-of-month, month, day-of-week) expression.
// The day-of-week field accepts 0–7 where both 0 and 7 represent Sunday.
var cronFields = [5]fieldSpec{
	{0, 59, "minute"},
	{0, 23, "hour"},
	{1, 31, "day-of-month"},
	{1, 12, "month"},
	{0, 7, "day-of-week"},
}

// cronFieldsWithSeconds defines the valid ranges for each position in a
// six-field (second, minute, hour, day-of-month, month, day-of-week)
// expression.
// The day-of-week field accepts 0–7 where both 0 and 7 represent Sunday.
var cronFieldsWithSeconds = [6]fieldSpec{
	{0, 59, "second"},
	{0, 59, "minute"},
	{0, 23, "hour"},
	{1, 31, "day-of-month"},
	{1, 12, "month"},
	{0, 7, "day-of-week"},
}

// monthNames maps three-letter month abbreviations to their numeric
// equivalents. The parser accepts both numeric and abbreviated month names.
var monthNames = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4,
	"may": 5, "jun": 6, "jul": 7, "aug": 8,
	"sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

// dowNames maps three-letter day-of-week abbreviations (and the full name
// "sunday"…"saturday") to their numeric equivalents (0 = Sunday).
var dowNames = map[string]int{
	"sun": 0, "mon": 1, "tue": 2, "wed": 3,
	"thu": 4, "fri": 5, "sat": 6,
}

// cronSchedule is the internal, parsed representation of a standard cron
// expression. It implements the Schedule interface.
type cronSchedule struct {
	second     []bool // [0..59] — only populated in six-field mode
	minute     []bool // [0..59]
	hour       []bool // [0..23]
	dayOfMonth []bool // [1..31]
	month      []bool // [1..12]
	dayOfWeek  []bool // [0..6] (0 = Sunday)
	loc        *time.Location
}

// Next returns the earliest time after t at which the schedule would activate.
// The receiver's location is applied before field matching; if no activation
// exists within the next four years, the zero time is returned.
func (s *cronSchedule) Next(t time.Time) time.Time {
	// Normalise to the schedule's timezone.
	t = t.In(s.loc)

	// Advance by one second (or one minute in five-field mode) to ensure we
	// return a time strictly after t.
	if len(s.second) > 0 {
		t = t.Add(time.Second)
	} else {
		t = t.Add(time.Minute).Truncate(time.Minute)
	}

	// Search forward up to ~4 years to find the next matching instant.
	deadline := t.Add(4 * 365 * 24 * time.Hour)

WRAP:
	for t.Before(deadline) {
		// Month check.
		if !s.month[t.Month()] {
			// Advance to the first day of the next valid month.
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, s.loc)
			continue WRAP
		}
		// Day-of-month and day-of-week check.
		if !s.dayOfMonth[t.Day()] || !s.dayOfWeek[t.Weekday()] {
			t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, s.loc)
			continue WRAP
		}
		// Hour check.
		if !s.hour[t.Hour()] {
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, s.loc)
			continue WRAP
		}
		// Minute check.
		if !s.minute[t.Minute()] {
			t = t.Add(time.Minute).Truncate(time.Minute)
			continue WRAP
		}
		// Second check (six-field mode only).
		if len(s.second) > 0 && !s.second[t.Second()] {
			t = t.Add(time.Second).Truncate(time.Second)
			continue WRAP
		}
		return t
	}
	return time.Time{}
}

// intervalSchedule fires every fixed Duration starting from the first tick
// after the reference time. It implements Schedule.
type intervalSchedule struct {
	interval time.Duration
}

// Next returns the earliest multiple of the interval that is strictly after t.
func (s *intervalSchedule) Next(t time.Time) time.Time {
	return t.Add(s.interval - time.Duration(t.UnixNano())%s.interval)
}
