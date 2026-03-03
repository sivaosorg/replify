package crontask

import (
	"errors"
	"sync"
)

// Sentinel errors returned by the crontask package. Callers can test for
// these values with errors.Is when they need to handle a specific condition.
var (
	// ErrInvalidExpression is returned by Parse and Validate when the provided
	// cron expression is syntactically or semantically invalid.
	ErrInvalidExpression = errors.New("crontask: invalid cron expression")

	// ErrJobNotFound is returned by Remove and similar methods when the given
	// job ID does not exist in the registry.
	ErrJobNotFound = errors.New("crontask: job not found")

	// ErrSchedulerRunning is returned by Start when the scheduler is already
	// in a running state.
	ErrSchedulerRunning = errors.New("crontask: scheduler is already running")

	// ErrSchedulerStopped is returned by Register and similar methods when the
	// caller attempts to mutate a scheduler that has been permanently shut down.
	ErrSchedulerStopped = errors.New("crontask: scheduler has been stopped")

	// ErrJobTimeout is wrapped into the error returned by the executor when a
	// job's execution deadline is exceeded.
	ErrJobTimeout = errors.New("crontask: job execution timed out")

	// ErrMaxRetriesExceeded is returned when a job exhausts its configured
	// retry budget without succeeding.
	ErrMaxRetriesExceeded = errors.New("crontask: maximum retries exceeded")
)

// Cron expression aliases are short-hand names for commonly used cron schedules.
// They are resolved by the parser before field-level parsing occurs.
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
	// TODO: Add more business-oriented aliases.
	"@businessdaily":  "0 9 * * 1-5",
	"@businesshourly": "0 9-17 * * 1-5",
	"@quarterly":      "0 0 1 1,4,7,10 *",
	"@semimonthly":    "0 0 1,15 * *",
	"@workhours":      "* 9-17 * * 1-5",
	"@marketopen":     "30 9 * * 1-5",
	"@marketclose":    "0 16 * * 1-5",
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

// aliasMapMu guards aliasMap for concurrent reads and writes. The initial
// entries are populated at package-init time (single-goroutine), after which
// all reads go through lookupAlias and all writes through RegisterAlias.
var aliasMapMu sync.RWMutex

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

// monthLabels maps month numbers to their English names.
// The zeroth element is unused.
var monthLabels = [13]string{
	"", "January", "February", "March", "April",
	"May", "June", "July", "August", "September",
	"October", "November", "December",
}
