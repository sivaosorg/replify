package replify

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"

	"github.com/sivaosorg/replify/pkg/encoding"
	"github.com/sivaosorg/replify/pkg/slogger"
)

// calculateSize calculates the size of the marshaled data.
// It uses encoding.Marshal to marshal the data and returns the length of the resulting byte slice.
// If an error occurs during marshaling, it returns 0.
func calculateSize(data any) int {
	_bytes, err := encoding.MarshalJSONb(data)
	if err != nil {
		return 0
	}
	return len(_bytes)
}

// compress compresses the given data using gzip and encodes it in base64.
// It first marshals the data using encoding.Marshal, then compresses the resulting byte slice
// using gzip. The compressed data is then encoded in base64 and returned as a string.
// If any error occurs during marshaling or compression, it returns an empty string.
func compress(data any) string {
	_bytes, err := encoding.MarshalJSONb(data)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err = gz.Write(_bytes)
	if err != nil {
		return ""
	}
	err = gz.Close()
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// decompress decompresses the given data using gzip and decodes it from base64.
// It first decodes the base64 encoded data using base64.StdEncoding.DecodeString,
// then decompresses the resulting byte slice using gzip. The decompressed data is
// then unmarshaled using encoding.Unmarshal and returned as an interface{}.
// If any error occurs during decoding or decompression, it returns nil.
func decompress(data string) any {
	_bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	gz, err := gzip.NewReader(bytes.NewReader(_bytes))
	if err != nil {
		return nil
	}
	defer gz.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(gz)
	if err != nil {
		return nil
	}
	var result any
	if err := encoding.UnmarshalBytes(buf.Bytes(), &result); err != nil {
		return nil
	}
	return result
}

// chunk takes a response represented as a map and returns a slice of byte slices,
// where each byte slice is a chunk of the JSON representation of the response.
// This is useful for streaming large responses in smaller segments.
// If the JSON encoding fails, it returns nil.
func chunk(data map[string]any) [][]byte {
	_bytes, err := encoding.MarshalJSONb(data)
	if err != nil {
		return nil
	}
	var chunks [][]byte
	for i := 0; i < len(_bytes); i += defaultChunkSize {
		end := i + defaultChunkSize
		if end > len(_bytes) {
			end = len(_bytes)
		}
		// Create a copy of the chunk to avoid referencing the underlying array.
		// This is important to ensure that each chunk is independent and can be
		// processed separately without affecting the others.
		chunk := make([]byte, end-i)
		copy(chunk, _bytes[i:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}

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
func jsonpass(data any) string {
	return encoding.JSON(data)
}

// jsonpretty converts a Go value to its pretty-printed JSON string representation or returns the value directly if it is already a string.
//
// This function checks if the input data is a string; if so, it returns it directly.
// Otherwise, it marshals the input value `data` into a formatted JSON string using
// the MarshalIndent function. If an error occurs during marshalling, it returns an empty string.
//
// Parameters:
//   - `data`: The Go value to be converted to pretty-printed JSON, or a string to be returned directly.
//
// Returns:
//   - A string containing the pretty-printed JSON representation of the input value, or an empty string if an error occurs.
//
// Example:
//
//	jsonPrettyStr := jsonpretty(myStruct)
func jsonpretty(data any) string {
	return encoding.JSONPretty(data)
}

// httpStatusLevel maps an HTTP status code to its corresponding [slogger.Level].
//
//   - 1xx → Debug  (informational)
//   - 2xx → Info   (success)
//   - 3xx → Warn   (redirection)
//   - 4xx → Error  (client error)
//   - 5xx → Error  (server error; Fatal is avoided — it calls os.Exit(1))
//   - other → Trace
func httpStatusLevel(code int) slogger.Level {
	switch {
	case code >= 400:
		return slogger.ErrorLevel
	case code >= 300:
		return slogger.WarnLevel
	case code >= 200:
		return slogger.InfoLevel
	case code >= 100:
		return slogger.DebugLevel
	default:
		return slogger.TraceLevel
	}
}

// logAtLevel dispatches a single log entry to l at the given level.
// It uses the appropriate method of the slogger.Logger based on the provided slogger.Level.
//
// Parameters:
//   - `l`: The slogger.Logger instance to which the log entry will be dispatched.
//   - `lvl`: The slogger.Level indicating the severity of the log entry (e.g., ErrorLevel, WarnLevel, InfoLevel, DebugLevel, TraceLevel).
//   - `msg`: The message string to be logged.
//   - `f`: A slogger.Field containing additional structured data to be included in the log entry.
//
// The function uses a switch statement to determine which logging method to call on the logger based on the provided level.
// If the level does not match any of the defined levels (ErrorLevel, WarnLevel, InfoLevel, DebugLevel), it defaults to using Trace.
func logAtLevel(l *slogger.Logger, lvl slogger.Level, msg string, f slogger.Field) {
	switch lvl {
	case slogger.ErrorLevel:
		l.Error(msg, f)
	case slogger.WarnLevel:
		l.Warn(msg, f)
	case slogger.InfoLevel:
		l.Info(msg, f)
	case slogger.DebugLevel:
		l.Debug(msg, f)
	default:
		l.Trace(msg, f)
	}
}

// safeBody checks if the provided value is a valid JSON string or byte slice and returns a safe representation.
//
// This function takes an input value and determines if it is a valid JSON string or byte slice.
// If the value is a valid JSON string, it returns a `json.RawMessage` containing the JSON data.
// If the value is a valid JSON byte slice, it also returns a `json.RawMessage` containing the JSON data.
// For any other type of value, it returns the original value as is.
//
// Parameters:
//   - value: The input value to be checked and processed.
//
// Returns:
//   - A `json.RawMessage` if the input is a valid JSON string or byte slice.
//   - The original value for any other type of input.
func safeBody(value any) any {
	var result any
	switch v := value.(type) {
	case string:
		if encoding.IsValidJSON(v) {
			result = json.RawMessage(encoding.Ugly([]byte(v)))
		} else {
			result = v
		}
	case []byte:
		if encoding.IsValidJSONBytes(v) {
			result = json.RawMessage(encoding.Ugly(v))
		} else {
			result = v
		}
	default:
		result = value
	}

	return result
}
