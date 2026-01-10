package replify

// WithCode sets the `code` field of the `header` instance.
//
// This function assigns the provided integer value to the `code` field of the `header`
// and returns the updated `header` instance, allowing for method chaining.
//
// Parameters:
//   - `v`: The integer value to set as the HTTP status code.
//
// Returns:
//   - The updated `header` instance with the `code` field set to the provided value.
func (h *header) WithCode(v int) *header {
	h.code = v
	return h
}

// WithText sets the `text` field of the `header` instance.
//
// This function assigns the provided string value to the `text` field of the `header`
// and returns the updated `header` instance, allowing for method chaining.
//
// Parameters:
//   - `v`: The string value to set as the text message.
//
// Returns:
//   - The updated `header` instance with the `text` field set to the provided value.
func (h *header) WithText(v string) *header {
	h.text = v
	return h
}

// WithType sets the `Type` field of the `header` instance.
//
// This function assigns the provided string value to the `Type` field of the `header`
// and returns the updated `header` instance, allowing for method chaining.
//
// Parameters:
//   - `v`: The string value to set as the type of the header.
//
// Returns:
//   - The updated `header` instance with the `Type` field set to the provided value.
func (h *header) WithType(v string) *header {
	h.typez = v
	return h
}

// WithDescription sets the `description` field of the `header` instance.
//
// This function assigns the provided string value to the `description` field of the `header`
// and returns the updated `header` instance, allowing for method chaining.
//
// Parameters:
//   - `v`: The string value to set as the description of the header.
//
// Returns:
//   - The updated `header` instance with the `description` field set to the provided value.
func (h *header) WithDescription(v string) *header {
	h.description = v
	return h
}

// Respond generates a map representation of the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) and includes the
// values of its fields in the returned map. Only the fields that are present (i.e., non-empty)
// are added to the map, ensuring a clean and concise response.
//
// Fields included in the response:
//   - `code`: The HTTP status code, if present and greater than 0.
//   - `text`: The associated text message, if present and not empty.
//   - `type`: The type of the header, if present and not empty.
//   - `description`: A description related to the header, if present and not empty.
//
// Returns:
//   - A `map[string]interface{}` containing the fields of the `header` instance that are present.
func (h *header) Respond() map[string]any {
	m := make(map[string]any)
	if !h.Available() {
		return m
	}
	if h.IsCodePresent() {
		m["code"] = h.code
	}
	if h.IsTextPresent() {
		m["text"] = h.text
	}
	if h.IsTypePresent() {
		m["type"] = h.typez
	}
	if h.IsDescriptionPresent() {
		m["description"] = h.description
	}
	return m
}

// Json serializes the `header` instance into a compact JSON string.
//
// This function uses the `encoding.Json` utility to create a compact JSON representation
// of the `header` instance. The resulting string contains only the key information, formatted
// with minimal whitespace, making it suitable for compact storage or transmission of header data.
//
// Returns:
//   - A compact JSON string representation of the `header` instance.
func (h *header) Json() string {
	return jsonpass(h.Respond())
}

// JsonPretty serializes the `header` instance into a prettified JSON string.
//
// This function uses the `encoding.JsonPretty` utility to produce a formatted, human-readable
// JSON string representation of the `header` instance. The output is structured with indentation
// and newlines, making it ideal for inspecting header data in a clear, easy-to-read format, especially
// during debugging or development.
//
// Returns:
//   - A prettified JSON string representation of the `header` instance, formatted for improved readability.
func (h *header) JsonPretty() string {
	return jsonpretty(h.Respond())
}
