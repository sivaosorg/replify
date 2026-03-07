package slogger

import (
	"fmt"
	"strconv"
	"time"
)

// String constructs a Field carrying a string value.
//
// Parameters:
//   - `key`: the field name
//   - `val`: the string value
//
// Returns:
//
// a Field of type StringType.
func String(key, val string) Field {
	return Field{Key: key, Type: StringType, strVal: val}
}

// Int constructs a Field carrying an int value stored as int64.
//
// Parameters:
//   - `key`: the field name
//   - `val`: the integer value
//
// Returns:
//
// a Field of type Int64Type.
func Int(key string, val int) Field {
	return Field{Key: key, Type: Int64Type, intVal: int64(val)}
}

// Int64 constructs a Field carrying an int64 value.
//
// Parameters:
//   - `key`: the field name
//   - `val`: the int64 value
//
// Returns:
//
// a Field of type Int64Type.
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: Int64Type, intVal: val}
}

// Float64 constructs a Field carrying a float64 value.
//
// Parameters:
//   - `key`: the field name
//   - `val`: the float64 value
//
// Returns:
//
// a Field of type Float64Type.
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: Float64Type, floatVal: val}
}

// Bool constructs a Field carrying a bool value.
//
// Parameters:
//   - `key`: the field name
//   - `val`: the boolean value
//
// Returns:
//
// a Field of type BoolType.
func Bool(key string, val bool) Field {
	return Field{Key: key, Type: BoolType, boolVal: val}
}

// Err constructs a Field carrying an error under the "error" key.
//
// Parameters:
//   - `err`: the error to log (may be nil)
//
// Returns:
//
// a Field of type ErrorType with key "error".
func Err(err error) Field {
	return Field{Key: "error", Type: ErrorType, errVal: err}
}

// Time constructs a Field carrying a time.Time value.
//
// Parameters:
//   - `key`: the field name
//   - `val`: the time value
//
// Returns:
//
// a Field of type TimeType.
func Time(key string, val time.Time) Field {
	return Field{Key: key, Type: TimeType, timeVal: val}
}

// Duration constructs a Field carrying a time.Duration value.
//
// Parameters:
//   - `key`: the field name
//   - `val`: the duration value
//
// Returns:
//
// a Field of type DurationType.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: DurationType, durVal: val}
}

// Any constructs a Field carrying an arbitrary value.
//
// Parameters:
//   - `key`: the field name
//   - `val`: any value; rendered via fmt.Sprintf("%v", val)
//
// Returns:
//
// a Field of type AnyType.
func Any(key string, val interface{}) Field {
	return Field{Key: key, Type: AnyType, anyVal: val}
}

// Value returns the field's value formatted as a string for use by formatters.
//
// Returns:
//
// the string representation of the stored value; for ErrorType it returns
// err.Error() or "<nil>"; for TimeType it uses time.RFC3339; for DurationType
// it uses Duration.String(); for AnyType it uses fmt.Sprintf("%v", val).
func (f *Field) Value() string {
	switch f.Type {
	case StringType:
		return f.strVal
	case Int64Type:
		return strconv.FormatInt(f.intVal, 10)
	case Float64Type:
		return strconv.FormatFloat(f.floatVal, 'f', -1, 64)
	case BoolType:
		if f.boolVal {
			return "true"
		}
		return "false"
	case ErrorType:
		if f.errVal == nil {
			return "<nil>"
		}
		return f.errVal.Error()
	case TimeType:
		return f.timeVal.Format(time.RFC3339)
	case DurationType:
		return f.durVal.String()
	case AnyType:
		return fmt.Sprintf("%v", f.anyVal)
	default:
		return ""
	}
}
