package replify

import (
	"fmt"
	"time"
)

// WithApiVersion sets the API version for the `meta` instance.
//
// This function updates the `apiVersion` field of the `meta` instance with the specified value
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `v`: A string representing the API version to set.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithApiVersion(v string) *meta {
	m.apiVersion = v
	return m
}

// WithApiVersionf sets the API version for the `meta` instance using a formatted string.
//
// This function constructs a formatted string for the API version using the provided `format` string
// and arguments (`args`). It then assigns the formatted value to the `apiVersion` field of the `meta` instance.
// The method supports method chaining by returning a pointer to the modified `meta` instance.
//
// Parameters:
//   - format: A format string to construct the API version.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithApiVersionf(format string, args ...any) *meta {
	return m.WithApiVersion(fmt.Sprintf(format, args...))
}

// WithRequestID sets the request ID for the `meta` instance.
//
// This function updates the `requestID` field of the `meta` instance with the specified value
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `v`: A string representing the request ID to set.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithRequestID(v string) *meta {
	m.requestID = v
	return m
}

// WithRequestIDf sets the request ID for the `meta` instance using a formatted string.
//
// This function constructs a formatted string for the request ID using the provided `format` string
// and arguments (`args`). It then assigns the formatted value to the `requestID` field of the `meta` instance.
// The method supports method chaining by returning a pointer to the modified `meta` instance.
//
// Parameters:
//   - format: A format string to construct the request ID.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithRequestIDf(format string, args ...any) *meta {
	return m.WithRequestID(fmt.Sprintf(format, args...))
}

// WithLocale sets the locale for the `meta` instance.
//
// This function updates the `locale` field of the `meta` instance with the specified value
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `v`: A string representing the locale to set.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithLocale(v string) *meta {
	m.locale = v
	return m
}

// WithRequestedTime sets the requested time for the `meta` instance.
//
// This function updates the `requestedTime` field of the `meta` instance with the specified value
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `v`: A `time.Time` object representing the requested time to set.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithRequestedTime(v time.Time) *meta {
	m.requestedTime = v
	return m
}

// WithCustomFields sets the custom fields for the `meta` instance.
//
// This function updates the `customFields` map of the `meta` instance with the provided values
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `values`: A map of string keys to interface{} values representing custom fields.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithCustomFields(values map[string]any) *meta {
	m.customFields = values
	return m
}

// WithCustomFieldKV sets a specific custom field key-value pair for the `meta` instance.
//
// This function adds or updates a custom field in the `customFields` map of the `meta` instance.
// If the `customFields` map is empty, it is initialized first.
//
// Parameters:
//   - `key`: A string representing the custom field key.
//   - `value`: An interface{} representing the value to associate with the custom field key.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithCustomFieldKV(key string, value any) *meta {
	if !m.IsCustomFieldPresent() {
		m.customFields = make(map[string]interface{})
	}
	m.customFields[key] = value
	return m
}

// WithCustomFieldKVf sets a specific custom field key-value pair for the `meta` instance
// using a formatted value.
//
// This function creates a formatted string value using the provided `format` string and
// `args`. It then calls `WithCustomFieldKV` to add or update the custom field with the
// specified key and the formatted value. The modified `meta` instance is returned for
// method chaining.
//
// Parameters:
//   - key: A string representing the key for the custom field.
//   - format: A format string to construct the value.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithCustomFieldKVf(key string, format string, args ...any) *meta {
	return m.WithCustomFieldKV(key, fmt.Sprintf(format, args...))
}

// Respond generates a map representation of the `meta` instance.
//
// This method collects various fields of the `meta` instance (e.g., `apiVersion`, `requestID`, etc.)
// and organizes them into a key-value map. Only fields that are available and valid
// (e.g., non-empty or initialized) are included in the resulting map.
//
// Fields included in the response:
//   - `api_version`: The API version, if present.
//   - `request_id`: The unique request identifier, if present.
//   - `locale`: The locale information, if present.
//   - `requested_time`: The requested time, if it is initialized.
//   - `custom_fields`: A map of custom fields, if present.
//
// Returns:
//   - A `map[string]interface{}` containing the structured metadata.
func (m *meta) Respond() map[string]any {
	mk := make(map[string]any)
	if !m.Available() {
		return mk
	}
	if m.IsApiVersionPresent() {
		mk["api_version"] = m.apiVersion
	}
	if m.IsRequestIDPresent() {
		mk["request_id"] = m.requestID
	}
	if m.IsLocalePresent() {
		mk["locale"] = m.locale
	}
	if m.IsRequestedTimePresent() {
		mk["requested_time"] = m.requestedTime
	}
	if m.IsCustomFieldPresent() {
		mk["custom_fields"] = m.customFields
	}
	return mk
}

// Json serializes the `meta` instance into a compact JSON string.
//
// This function uses the `encoding.Json` utility to create a compact JSON representation
// of the `meta` instance. The resulting string is formatted without additional whitespace,
// suitable for efficient storage or transmission of metadata.
//
// Returns:
//   - A compact JSON string representation of the `meta` instance.
func (m *meta) Json() string {
	return jsonpass(m.Respond())
}

// JsonPretty serializes the `meta` instance into a prettified JSON string.
//
// This function calls the `encoding.JsonPretty` utility to produce a formatted, human-readable
// JSON string representation of the `meta` instance. The output is useful for debugging
// or inspecting metadata in a more structured format.
//
// Returns:
//   - A prettified JSON string representation of the `meta` instance.
func (m *meta) JsonPretty() string {
	return jsonpretty(m.Respond())
}

// RandRequestID generates and sets a random request ID for the `meta` instance.
//
// This function utilizes the `CryptoID` function to generate a unique request ID
// and assigns it to the `requestID` field of the `meta` instance. The modified
// `meta` instance is returned to allow for method chaining.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) RandRequestID() *meta {
	return m.WithRequestID(cryptoID())
}
