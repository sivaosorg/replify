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

	// AliasEndOfDay fires at 17:00 on every weekday — useful for daily
	// close-of-business summaries, digest emails, or EOD reports.
	AliasEndOfDay Alias = "@endOfDay"

	// AliasStartOfDay fires at 08:00 on every weekday — a minute before most
	// staff arrive, ideal for pre-loading caches, warming services, or sending
	// morning briefings.
	AliasStartOfDay Alias = "@startOfDay"

	// AliasLunchtime fires at 12:00 on every weekday — suitable for mid-day
	// digest notifications or low-priority batch jobs run during off-peak hours.
	AliasLunchtime Alias = "@lunchtime"

	// AliasEndOfWeek fires at 17:00 every Friday — perfect for weekly summary
	// emails, cleanup jobs, or end-of-week reporting pipelines.
	AliasEndOfWeek Alias = "@endOfWeek"

	// AliasStartOfWeek fires at 08:00 every Monday — ideal for weekly planning
	// notifications, metric resets, or Monday morning digest generation.
	AliasStartOfWeek Alias = "@startOfWeek"

	// AliasEndOfMonth fires at 23:59 on the last day of each month (28th used
	// as safe cross-month anchor). For true last-day logic, use a job-level
	// calendar check. Suitable for monthly billing runs and invoicing triggers.
	AliasEndOfMonth Alias = "@endOfMonth"

	// AliasPayroll fires at 08:00 on the 1st and 15th of every month —
	// matching the two most common semi-monthly payroll schedules.
	AliasPayroll Alias = "@payroll"

	// AliasNightlyMaintenance fires at 02:00 every day — a low-traffic window
	// suited for database vacuums, index rebuilds, and nightly backup jobs.
	AliasNightlyMaintenance Alias = "@nightlyMaintenance"

	// AliasPreMarket fires at 08:00 on every weekday — one hour before the US
	// equity market opens, useful for pre-market data ingestion or alert checks.
	AliasPreMarket Alias = "@preMarket"

	// AliasAfterMarket fires at 17:00 on every weekday — one hour after the US
	// equity market closes, suitable for after-hours reconciliation or reporting.
	AliasAfterMarket Alias = "@afterMarket"

	// AliasMidMarket fires at 12:30 on every weekday — the midpoint of the US
	// trading day, useful for intraday snapshot jobs or mid-session risk checks.
	AliasMidMarket Alias = "@midMarket"

	// AliasQuarterEnd fires at 23:59 on the last day of each fiscal quarter
	// (March, June, September, December). Useful for quarter-close accounting
	// jobs, regulatory filings, or board report generation.
	AliasQuarterEnd Alias = "@quarterEnd"

	// AliasTaxDeadline fires at 08:00 on April 15th — the standard US federal
	// tax filing deadline. Useful for annual compliance reminder pipelines.
	AliasTaxDeadline Alias = "@taxDeadline"

	// AliasRegulatoryOpen fires at 07:00 on every weekday — before market open,
	// aligned with common regulatory reporting windows (e.g. FINRA, SEC).
	AliasRegulatoryOpen Alias = "@regulatoryOpen"

	// AliasOffPeakHourly fires once an hour between 20:00 and 06:00 every day —
	// useful for infrastructure jobs, bulk imports, or ML training runs that
	// should avoid peak business hours.
	AliasOffPeakHourly Alias = "@offPeakHourly"

	// AliasDeploymentWindow fires at 22:00 on Tuesday and Thursday — a
	// conventional low-risk deployment window outside business hours but not
	// on a weekend, giving a full working day before and after for monitoring.
	AliasDeploymentWindow Alias = "@deploymentWindow"

	// AliasDatabaseBackup fires at 01:00 every day — a quiet early-morning
	// window well-suited for full or incremental database backup jobs.
	AliasDatabaseBackup Alias = "@databaseBackup"

	// AliasWeeklyReport fires at 08:00 every Monday — delivers weekly KPI
	// summaries, analytics digests, or newsletter generation at the start of
	// the work week, before most users are active.
	AliasWeeklyReport Alias = "@weeklyReport"

	// AliasMonthlyReport fires at 08:00 on the 1st of every month — aligns
	// with standard monthly reporting cycles for finance, product, or ops teams.
	AliasMonthlyReport Alias = "@monthlyReport"

	// AliasCustomerDigest fires at 09:00 on every weekday — sends customer
	// activity digests, CRM summaries, or support queue snapshots at the start
	// of each business day.
	AliasCustomerDigest Alias = "@customerDigest"

	// AliasSLACheck fires every 15 minutes during business hours on weekdays —
	// useful for SLA breach detection, ticket-age monitors, or uptime heartbeat
	// checks that need sub-hourly granularity without running 24/7.
	AliasSLACheck Alias = "@slaCheck"
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
	"@businessdaily":      "0 9 * * 1-5",
	"@businesshourly":     "0 9-17 * * 1-5",
	"@quarterly":          "0 0 1 1,4,7,10 *",
	"@semimonthly":        "0 0 1,15 * *",
	"@workhours":          "* 9-17 * * 1-5",
	"@marketopen":         "30 9 * * 1-5",
	"@marketclose":        "0 16 * * 1-5",
	"@endofday":           "0 17 * * 1-5",
	"@startofday":         "0 8 * * 1-5",
	"@lunchtime":          "0 12 * * 1-5",
	"@endofweek":          "0 17 * * 5",
	"@startofweek":        "0 8 * * 1",
	"@endofmonth":         "59 23 28 * *",
	"@payroll":            "0 8 1,15 * *",
	"@nightlymaintenance": "0 2 * * *",
	"@premarket":          "0 8 * * 1-5",
	"@aftermarket":        "0 17 * * 1-5",
	"@midmarket":          "30 12 * * 1-5",
	"@quarterend":         "59 23 31 3,6,9,12 *",
	"@taxdeadline":        "0 8 15 4 *",
	"@regulatoryopen":     "0 7 * * 1-5",
	"@offpeakhourly":      "0 20-23,0-6 * * *",
	"@deploymentwindow":   "0 22 * * 2,4",
	"@databasebackup":     "0 1 * * *",
	"@weeklyreport":       "0 8 * * 1",
	"@monthlyreport":      "0 8 1 * *",
	"@customerdigest":     "0 9 * * 1-5",
	"@slacheck":           "*/15 9-17 * * 1-5",
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
	"@yearly":             "Once a year, at midnight on January 1st",
	"@annually":           "Once a year, at midnight on January 1st",
	"@monthly":            "Once a month, at midnight on the 1st",
	"@weekly":             "Once a week, at midnight on Sunday",
	"@daily":              "Once a day, at midnight",
	"@midnight":           "Once a day, at midnight",
	"@hourly":             "Every hour",
	"@minutely":           "Every minute",
	"@weekdays":           "At midnight, Monday through Friday",
	"@weekends":           "At midnight, on Saturday and Sunday",
	"@businessdaily":      "At 09:00, Monday through Friday",
	"@businesshourly":     "Every hour from 09:00 to 17:00, Monday through Friday",
	"@quarterly":          "Once a quarter, at midnight on the 1st",
	"@semimonthly":        "Twice a month, at midnight on the 1st and 15th",
	"@workhours":          "Every minute from 09:00 to 17:59, Monday through Friday",
	"@marketopen":         "At 09:30, Monday through Friday",
	"@marketclose":        "At 16:00, Monday through Friday",
	"@endofday":           "At 17:00, Monday through Friday",
	"@startofday":         "At 08:00, Monday through Friday",
	"@lunchtime":          "At 12:00, Monday through Friday",
	"@endofweek":          "At 17:00 on Friday",
	"@startofweek":        "At 08:00 on Monday",
	"@endofmonth":         "At 23:59 on the 28th of every month",
	"@payroll":            "At 08:00 on the 1st and 15th of every month",
	"@nightlymaintenance": "At 02:00, every day",
	"@premarket":          "At 08:00, Monday through Friday (1 hour before US market open)",
	"@aftermarket":        "At 17:00, Monday through Friday (1 hour after US market close)",
	"@midmarket":          "At 12:30, Monday through Friday (US market midpoint)",
	"@quarterend":         "At 23:59 on the last day of each fiscal quarter",
	"@taxdeadline":        "At 08:00 on April 15th (US federal tax deadline)",
	"@regulatoryopen":     "At 07:00, Monday through Friday",
	"@offpeakhourly":      "Every hour from 20:00 to 06:00, every day",
	"@deploymentwindow":   "At 22:00 on Tuesday and Thursday",
	"@databasebackup":     "At 01:00, every day",
	"@weeklyreport":       "At 08:00 on Monday",
	"@monthlyreport":      "At 08:00 on the 1st of every month",
	"@customerdigest":     "At 09:00, Monday through Friday",
	"@slacheck":           "Every 15 minutes from 09:00 to 17:00, Monday through Friday",
}

// monthLabels maps month numbers to their English names.
// The zeroth element is unused.
var monthLabels = [13]string{
	"", "January", "February", "March", "April",
	"May", "June", "July", "August", "September",
	"October", "November", "December",
}
