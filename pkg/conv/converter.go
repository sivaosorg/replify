package conv

import (
	"time"
)

// ///////////////////////////
// Section:  Converter struct
// ///////////////////////////

// Converter implements type conversions with configurable options.
// It is safe for concurrent use by multiple goroutines.
type Converter struct {
	strictMode  bool     // If true, returns error for lossy conversions
	dateFormats []string // Custom date formats for time parsing
	locale      string   // Locale for parsing (future use)
	trimStrings bool     // If true, trims whitespace from strings before conversion
	nilAsZero   bool     // If true, nil values return zero value instead of error
	emptyAsZero bool     // If true, empty strings return zero value instead of error
}

// ///////////////////////////
// Section: Constructor and Builder methods
// ///////////////////////////

// NewConverter creates a new Converter instance with default settings.
//
// Returns:
//   - A pointer to a newly created Converter instance.
func NewConverter() *Converter {
	return &Converter{
		dateFormats: defaultDateFormats(),
		nilAsZero:   true,
		emptyAsZero: true,
		trimStrings: true,
	}
}

// WithStrictMode enables or disables strict mode for conversions.
// In strict mode, lossy conversions (e.g., float to int) return an error.
//
// Parameters:
//   - v: Boolean indicating whether strict mode should be enabled.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) WithStrictMode(v bool) *Converter {
	c.strictMode = v
	return c
}

// WithDateFormats sets custom date formats for time parsing.
// Formats are tried in order when parsing time strings.
//
// Parameters:
//   - formats:  Variadic list of date format strings.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) WithDateFormats(formats ...string) *Converter {
	c.dateFormats = formats
	return c
}

// WithLocale sets the locale for parsing operations.
//
// Parameters:
//   - locale:  Locale string (e.g., "en_US", "vi_VN").
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) WithLocale(locale string) *Converter {
	c.locale = locale
	return c
}

// WithTrimStrings enables or disables string trimming before conversion.
//
// Parameters:
//   - v: Boolean indicating whether strings should be trimmed.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) WithTrimStrings(v bool) *Converter {
	c.trimStrings = v
	return c
}

// EnableTrimStrings enables string trimming before conversion.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) EnableTrimStrings() *Converter {
	return c.WithTrimStrings(true)
}

// DisableTrimStrings disables string trimming before conversion.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) DisableTrimStrings() *Converter {
	return c.WithTrimStrings(false)
}

// WithNilAsZero enables or disables returning zero value for nil inputs.
//
// Parameters:
//   - v: Boolean indicating whether nil should return zero value.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) WithNilAsZero(v bool) *Converter {
	c.nilAsZero = v
	return c
}

// EnableNilAsZero enables returning zero value for nil inputs.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) EnableNilAsZero() *Converter {
	return c.WithNilAsZero(true)
}

// DisableNilAsZero disables returning zero value for nil inputs.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) DisableNilAsZero() *Converter {
	return c.WithNilAsZero(false)
}

// WithEmptyAsZero enables or disables returning zero value for empty string inputs.
//
// Parameters:
//   - v: Boolean indicating whether empty strings should return zero value.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) WithEmptyAsZero(v bool) *Converter {
	c.emptyAsZero = v
	return c
}

// EnableEmptyAsZero enables returning zero value for empty string inputs.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) EnableEmptyAsZero() *Converter {
	return c.WithEmptyAsZero(true)
}

// DisableEmptyAsZero disables returning zero value for empty string inputs.
//
// Returns:
//   - A pointer to the modified Converter instance (enabling method chaining).
func (c *Converter) DisableEmptyAsZero() *Converter {
	return c.WithEmptyAsZero(false)
}

// ///////////////////////////
// Section: Getter methods
// ///////////////////////////

// IsStrictMode returns whether strict mode is enabled.
//
// Returns:
//   - A boolean indicating if strict mode is enabled.
func (c *Converter) IsStrictMode() bool {
	return c.strictMode
}

// DateFormats returns the configured date formats.
//
// Returns:
//   - A slice of strings representing the date formats.
func (c *Converter) DateFormats() []string {
	return c.dateFormats
}

// Locale returns the configured locale.
//
// Returns:
//   - A string representing the locale.
func (c *Converter) Locale() string {
	return c.locale
}

// ///////////////////////////
// Section:  Clone and Reset
// ///////////////////////////

// Clone creates a deep copy of the Converter instance.
//
// Returns:
//   - A pointer to the cloned Converter instance.
func (c *Converter) Clone() *Converter {
	clone := &Converter{
		strictMode:  c.strictMode,
		locale:      c.locale,
		trimStrings: c.trimStrings,
		nilAsZero:   c.nilAsZero,
		emptyAsZero: c.emptyAsZero,
	}
	if c.dateFormats != nil {
		clone.dateFormats = make([]string, len(c.dateFormats))
		copy(clone.dateFormats, c.dateFormats)
	}
	return clone
}

// Reset resets the Converter to its default settings.
//
// Returns:
//   - A pointer to the reset Converter instance.
func (c *Converter) Reset() *Converter {
	c.strictMode = false
	c.dateFormats = defaultDateFormats()
	c.locale = ""
	c.trimStrings = true
	c.nilAsZero = true
	c.emptyAsZero = true
	return c
}

// ///////////////////////////
// Section:  Helper functions
// ///////////////////////////

// defaultDateFormats returns the default list of date formats for parsing.
//
// Returns:
//   - A slice of strings representing the default date formats.
func defaultDateFormats() []string {
	return []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02",
		"02-01-2006",
		"02/01/2006",
		"01-02-2006",
		"01/02/2006",
		"Jan 2, 2006",
		"January 2, 2006",
		"2 Jan 2006",
		"2 January 2006",
		"Mon, 02 Jan 2006 15:04:05",
		"Mon, 2 Jan 2006 15:04:05",
		"02 Jan 2006 15:04 MST",
		"2 Jan 2006 15:04:05",
		"2 Jan 2006 15:04:05 MST",
	}
}
