package conv

import (
	"math"
	"math/cmplx"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sivaosorg/replify/pkg/encoding"
	"github.com/sivaosorg/replify/pkg/strutil"
)

// jsonpass converts a Go value to its JSON string representation or returns the value directly if it is already a string.
//
// This function checks if the input data is a string; if so, it returns it directly.
// Otherwise, it marshals the input value `data` into a JSON string using the
// MarshalToString function. If an error occurs during marshalling, it returns an empty string.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON, or a string to be returned directly.
//
// Returns:
//   - A string containing the JSON representation of the input value, or an empty string if an error occurs.
//
// Example:
//
//	jsonStr := jsonpass(myStruct)
func Jsonpass(data any) string {
	return encoding.JsonSafe(data)
}

// ///////////////////////////
// Section:  Time constants
// ///////////////////////////

var (
	emptyTime      = time.Time{}
	typeOfTime     = reflect.TypeOf(emptyTime)
	typeOfDuration = reflect.TypeOf(time.Duration(0))
)

// ///////////////////////////
// Section: Time conversion interfaces
// ///////////////////////////

// timeConverter is an interface for types that can convert themselves to time.Time.
type timeConverter interface {
	Time() (time.Time, error)
}

// durationConverter is an interface for types that can convert themselves to time.Duration.
type durationConverter interface {
	Duration() (time.Duration, error)
}

// ///////////////////////////
// Section: Duration conversion
// ///////////////////////////

// Duration attempts to convert the given value to time.Duration, returns the
// zero value and an error on failure.
//
// String parsing supports Go duration format (e.g., "1h30m") and numeric strings.
// Numeric values are treated as nanoseconds.
// Float values are treated as seconds (e.g., 1.5 = 1.5 seconds).
//
// Parameters:
//   - from: The value to convert.
//
// Returns:
//   - The converted time.Duration value.
//   - An error if conversion fails.
func (c *Converter) Duration(from any) (time.Duration, error) {
	if from == nil {
		if c.nilAsZero {
			return 0, nil
		}
		return 0, newConvError(from, "time.Duration")
	}

	// Fast path for common types
	switch v := from.(type) {
	case time.Duration:
		return v, nil
	case *time.Duration:
		if v == nil {
			return 0, nil
		}
		return *v, nil
	case string:
		return c.stringToDuration(v)
	case *string:
		if v == nil {
			return 0, nil
		}
		return c.stringToDuration(*v)
	case int:
		return time.Duration(v), nil
	case int64:
		return time.Duration(v), nil
	case int32:
		return time.Duration(v), nil
	case uint:
		return time.Duration(v), nil
	case uint64:
		if v > math.MaxInt64 {
			v = math.MaxInt64
		}
		return time.Duration(v), nil
	case float64:
		return c.float64ToDuration(v), nil
	case float32:
		return c.float64ToDuration(float64(v)), nil
	}

	// Check for custom converter interface
	if conv, ok := from.(durationConverter); ok {
		return conv.Duration()
	}

	// Use reflection for other types
	return c.durationFromReflect(from)
}

// ///////////////////////////
// Section: Time conversion
// ///////////////////////////

// Time attempts to convert the given value to time.Time, returns the zero value
// of time.Time and an error on failure.
//
// String parsing supports multiple date formats.  Use WithDateFormats() to customize.
//
// Parameters:
//   - from: The value to convert.
//
// Returns:
//   - The converted time.Time value.
//   - An error if conversion fails.
func (c *Converter) Time(from any) (time.Time, error) {
	if from == nil {
		if c.nilAsZero {
			return emptyTime, nil
		}
		return emptyTime, newConvError(from, "time.Time")
	}

	// Fast path for common types
	switch v := from.(type) {
	case time.Time:
		return v, nil
	case *time.Time:
		if v == nil {
			return emptyTime, nil
		}
		return *v, nil
	case string:
		return c.stringToTime(v)
	case *string:
		if v == nil {
			return emptyTime, nil
		}
		return c.stringToTime(*v)
	case int64:
		// Treat as Unix timestamp (seconds)
		return time.Unix(v, 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	case uint64:
		return time.Unix(int64(v), 0), nil
	case float64:
		// Treat as Unix timestamp with fractional seconds
		sec := int64(v)
		nanoSec := int64((v - float64(sec)) * 1e9)
		return time.Unix(sec, nanoSec), nil
	}

	// Check for custom converter interface
	if conv, ok := from.(timeConverter); ok {
		return conv.Time()
	}

	// Use reflection for other types
	return c.timeFromReflect(from)
}

// ///////////////////////////
// Section:  Package-level Time functions
// ///////////////////////////

// IsZeroTime checks if the given value converts to a zero time.Time.
//
// Parameters:
//   - from: The value to check.
//
// Returns:
//   - true if the converted time.Time is zero, false otherwise.
func IsZeroTime(from any) bool {
	if t, err := defaultConverter.Time(from); err == nil {
		return t.IsZero()
	}
	return true
}

// ///////////////////////////
// Section: Duration helpers
// ///////////////////////////

// Seconds creates a Duration representing the given number of seconds.
//
// Parameters:
//   - n: The number of seconds.
//
// Returns:
//   - A time.Duration representing the specified number of seconds.
func Seconds(n float64) time.Duration {
	return time.Duration(n * float64(time.Second))
}

// Minutes creates a Duration representing the given number of minutes.
//
// Parameters:
//   - n: The number of minutes.
//
// Returns:
//   - A time.Duration representing the specified number of minutes.
func Minutes(n float64) time.Duration {
	return time.Duration(n * float64(time.Minute))
}

// Hours creates a Duration representing the given number of hours.
//
// Parameters:
//   - n: The number of hours.
//
// Returns:
//   - A time.Duration representing the specified number of hours.
func Hours(n float64) time.Duration {
	return time.Duration(n * float64(time.Hour))
}

// Days creates a Duration representing the given number of days.
//
// Parameters:
//   - n: The number of days.
//
// Returns:
//   - A time.Duration representing the specified number of days.
func Days(n float64) time.Duration {
	return time.Duration(n * 24 * float64(time.Hour))
}

// stringToDuration converts a string to time.Duration.
//
// Parameters:
//   - v: The string value to convert.
//
// Returns:
//   - The converted time.Duration value.
//   - An error if the conversion fails.
func (c *Converter) stringToDuration(v string) (time.Duration, error) {
	if strutil.IsEmpty(v) {
		if c.emptyAsZero {
			return 0, nil
		}
		return 0, newConvErrorMsg("cannot convert empty string to time.Duration")
	}

	if c.trimStrings {
		v = strings.TrimSpace(v)
	}

	// Try Go duration format first (e.g., "1h30m")
	if parsed, err := time.ParseDuration(v); err == nil {
		return parsed, nil
	}

	// Try as integer nanoseconds
	if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
		return time.Duration(parsed), nil
	}

	// Try as float seconds
	if parsed, err := strconv.ParseFloat(v, 64); err == nil {
		return c.float64ToDuration(parsed), nil
	}

	return 0, newConvErrorf("cannot parse %q as time.Duration", v)
}

// float64ToDuration converts a float64 value representing seconds to time.Duration.
//
// Parameters:
//   - v: The float64 value to convert.
//
// Returns:
//   - The converted time.Duration value.
func (c *Converter) float64ToDuration(v float64) time.Duration {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return time.Duration(v * float64(time.Second))
}

// durationFromReflect converts a reflect.Value to time.Duration.
//
// Parameters:
//   - from: The input value to convert.
//
// Returns:
//   - The converted time.Duration value.
//   - An error if the conversion fails.
func (c *Converter) durationFromReflect(from any) (time.Duration, error) {
	value := indirectValue(reflect.ValueOf(from))
	if !value.IsValid() {
		if c.nilAsZero {
			return 0, nil
		}
		return 0, newConvError(from, "time.Duration")
	}

	kind := value.Kind()
	switch {
	case kind == reflect.String:
		return c.stringToDuration(value.String())
	case isKindInt(kind):
		return time.Duration(value.Int()), nil
	case isKindUint(kind):
		v := value.Uint()
		if v > math.MaxInt64 {
			v = math.MaxInt64
		}
		return time.Duration(v), nil
	case isKindFloat(kind):
		return c.float64ToDuration(value.Float()), nil
	case isKindComplex(kind):
		v := value.Complex()
		if cmplx.IsNaN(v) || cmplx.IsInf(v) {
			return 0, nil
		}
		return c.float64ToDuration(real(v)), nil
	}

	return 0, newConvError(from, "time.Duration")
}

// stringToTime converts a string to time.Time.
//
// Parameters:
//   - v: The string value to convert.
//
// Returns:
//   - The converted time.Time value.
//   - An error if the conversion fails.
func (c *Converter) stringToTime(v string) (time.Time, error) {
	if strutil.IsEmpty(v) {
		if c.emptyAsZero {
			return emptyTime, nil
		}
		return emptyTime, newConvErrorMsg("cannot convert empty string to time.Time")
	}

	if c.trimStrings {
		v = strings.TrimSpace(v)
	}

	// Try each configured format
	for _, format := range c.dateFormats {
		if t, err := time.Parse(format, v); err == nil {
			return t, nil
		}
	}

	// Try parsing as Unix timestamp
	if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
		return time.Unix(ts, 0), nil
	}

	// Try parsing as float timestamp
	if ts, err := strconv.ParseFloat(v, 64); err == nil {
		sec := int64(ts)
		nanoSec := int64((ts - float64(sec)) * 1e9)
		return time.Unix(sec, nanoSec), nil
	}

	return emptyTime, newConvErrorf("cannot parse %q as time.Time", v)
}

// timeFromReflect converts a reflect.Value to time.Time.
//
// Parameters:
//   - from: The input value to convert.
//
// Returns:
//   - The converted time.Time value.
//   - An error if the conversion fails.
func (c *Converter) timeFromReflect(from any) (time.Time, error) {
	value := indirectValue(reflect.ValueOf(from))
	if !value.IsValid() {
		if c.nilAsZero {
			return emptyTime, nil
		}
		return emptyTime, newConvError(from, "time.Time")
	}

	kind := value.Kind()
	switch {
	case kind == reflect.String:
		return c.stringToTime(value.String())
	case kind == reflect.Struct:
		if value.Type() == typeOfTime && value.CanInterface() {
			return value.Interface().(time.Time), nil
		}
	case isKindInt(kind):
		return time.Unix(value.Int(), 0), nil
	case isKindUint(kind):
		return time.Unix(int64(value.Uint()), 0), nil
	case isKindFloat(kind):
		v := value.Float()
		sec := int64(v)
		nanoSec := int64((v - float64(sec)) * 1e9)
		return time.Unix(sec, nanoSec), nil
	}

	return emptyTime, newConvError(from, "time.Time")
}
