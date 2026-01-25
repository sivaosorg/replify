package replify

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sivaosorg/replify/pkg/coll"
	"github.com/sivaosorg/replify/pkg/hashy"
	"github.com/sivaosorg/replify/pkg/randn"
	"github.com/sivaosorg/replify/pkg/strutil"
)

// Available checks whether the `wrapper` instance is non-nil.
//
// This function ensures that the `wrapper` object exists and is not nil.
// It serves as a safety check to avoid null pointer dereferences when accessing the instance's fields or methods.
//
// Returns:
//   - A boolean value indicating whether the `wrapper` instance is non-nil:
//   - `true` if the `wrapper` instance is non-nil.
//   - `false` if the `wrapper` instance is nil.
func (w *wrapper) Available() bool {
	return w != nil
}

// Error retrieves the error associated with the `wrapper` instance.
//
// This function returns the `errors` field of the `wrapper`, which contains
// any errors encountered during the operation of the `wrapper`.
//
// Returns:
//   - An error object, or `nil` if no errors are present.
func (w *wrapper) Error() string {
	if !w.Available() {
		return ""
	}
	cause := w.Cause()
	if cause == nil {
		return ""
	}
	return cause.Error()
}

// Cause traverses the error chain and returns the underlying cause of the error
// associated with the `wrapper` instance.
//
// This function checks if the error stored in the `wrapper` is itself another
// `wrapper` instance. If so, it recursively calls `Cause` on the inner error
// to find the ultimate cause. Otherwise, it returns the current error.
//
// Returns:
//   - The underlying cause of the error, which can be another error or the original error.
func (w *wrapper) Cause() error {
	if !w.Available() || w.errors == nil {
		// prevent nil pointer dereference
		// then just leave it as empty string
		return errors.New("")
	}
	// Traverse through wrapped errors.
	// We will use Unwrap() method to unwrap errors instead of checking for *wrapper explicitly.
	// This way, we can traverse to the innermost cause regardless of error type.
	visited := make(map[error]bool)
	cause := w.errors
	for cause != nil && !visited[cause] {
		visited[cause] = true
		if err, ok := cause.(interface{ Cause() error }); ok {
			next := err.Cause()
			if next == cause { // Prevent self-reference
				break
			}
			cause = next
		} else {
			break
		}
	}
	return cause
}

// StatusCode retrieves the HTTP status code associated with the `wrapper` instance.
//
// This function returns the `statusCode` field of the `wrapper`, which represents
// the HTTP status code for the response, indicating the outcome of the request.
//
// Returns:
//   - An integer representing the HTTP status code.
func (w *wrapper) StatusCode() int {
	if !w.Available() {
		return 0
	}
	if w.IsHeaderPresent() && w.statusCode <= 0 {
		w.statusCode = w.header.Code()
	}
	return w.statusCode
}

// StatusText returns a human-readable string representation of the HTTP status.
//
// This function combines the status code with its associated status text, which
// is retrieved using the `http.StatusText` function from the `net/http` package.
// The returned string follows the format "statusCode (statusText)".
//
// For example, if the status code is 200, the function will return "200 (OK)".
// If the status code is 404, it will return "404 (Not Found)".
//
// Returns:
//   - A string formatted as "statusCode (statusText)", where `statusCode` is the
//     numeric HTTP status code and `statusText` is the corresponding textual description.
func (w *wrapper) StatusText() string {
	return fmt.Sprintf("%d (%s)", w.StatusCode(), http.StatusText(w.StatusCode()))
}

// Message retrieves the message associated with the `wrapper` instance.
//
// This function returns the `message` field of the `wrapper`, which typically
// provides additional context or a description of the operation's outcome.
//
// Returns:
//   - A string representing the message.
func (w *wrapper) Message() string {
	if !w.Available() {
		return ""
	}
	return w.message
}

// Total retrieves the total number of items associated with the `wrapper` instance.
//
// This function returns the `total` field of the `wrapper`, which indicates
// the total number of items available, often used in paginated responses.
//
// Returns:
//   - An integer representing the total number of items.
func (w *wrapper) Total() int {
	if !w.Available() {
		return 0
	}
	return w.total
}

// Body retrieves the body data associated with the `wrapper` instance.
//
// This function returns the `data` field of the `wrapper`, which contains
// the primary data payload of the response.
//
// Returns:
//   - The body data (of any type), or `nil` if no body data is present.
func (w *wrapper) Body() any {
	if !w.Available() {
		return nil
	}
	return w.data
}

// OptimizeSafe compresses the body data if it exceeds a specified threshold.
//
// This function checks if the `wrapper` instance is available and if the body data
// exceeds the specified threshold for compression. If the body data is larger than
// the threshold, it compresses the data using gzip and updates the body with the
// compressed data. It also adds debugging information about the compression process,
// including the original and compressed sizes.
// If the threshold is not specified or is less than or equal to zero, it defaults to 1024 bytes (1KB).
// It also removes any empty debugging fields to clean up the response.
// Parameters:
//   - `threshold`: An integer representing the size threshold for compression.
//     If the body data size exceeds this threshold, it will be compressed.
//
// Returns:
//   - A pointer to the `wrapper` instance, allowing for method chaining.
//
// If the `wrapper` is not available, it returns the original instance without modifications.
func (w *wrapper) OptimizeSafe(threshold int) *wrapper {
	if !w.Available() {
		return w
	}
	if threshold <= 0 {
		threshold = 1024 // 1KB threshold for compression
	}
	// Check if the body data is present and exceeds the threshold
	// If the body is a string, we will check its length.
	if s, ok := w.data.(string); ok {
		if len(s) > threshold {
			compressed := compress(s)
			w.
				WithBody(compressed).
				WithDebuggingKV("compression", "gzip").
				WithDebuggingKV("original_size", len(s)).
				WithDebuggingKV("compressed_size", len(compressed))
			return w
		}
	}
	// If the body is not a string, we will check its size using CalculateSize.
	// Calculate the size of the body data.
	if data := w.Body(); data != nil {
		if size := calculateSize(data); size > threshold {
			compressed := compress(data)
			w.
				WithBody(compressed).
				WithDebuggingKV("compression", "gzip").
				WithDebuggingKV("original_size", size).
				WithDebuggingKV("compressed_size", calculateSize(compressed))
		}
	}
	return w
}

// Stream retrieves a channel that streams the body data of the `wrapper` instance.
//
// This function checks if the body data is present and, if so, streams the data
// in chunks. It creates a buffered channel to hold the streamed data, allowing
// for asynchronous processing of the response body.
// If the body is not present, it returns an empty channel.
// The streaming is done in a separate goroutine to avoid blocking the main execution flow.
// The body data is chunked into smaller parts using the `Chunk` function, which
// splits the response data into manageable segments for efficient streaming.
//
// Returns:
//   - A channel of byte slices that streams the body data.
//   - An empty channel if the body data is not present.
//
// This is useful for handling large responses in a memory-efficient manner,
// allowing the consumer to process each chunk as it becomes available.
// Note: The channel is closed automatically when the streaming is complete.
// If the body is not present, it returns an empty channel.
func (w *wrapper) Stream() <-chan []byte {
	ch := make(chan []byte, 1)
	if !w.IsBodyPresent() {
		return ch
	}
	go func() {
		defer close(ch)
		// Chunk the response data into smaller parts.
		// This is useful for streaming large responses in smaller segments.
		// We will use the Chunk function to split the response data into manageable chunks.
		chunks := chunk(w.Respond())
		for _, chunk := range chunks {
			ch <- chunk
		}
	}()
	return ch
}

// Debugging retrieves the debugging information from the `wrapper` instance.
//
// This function checks if the `wrapper` instance is available (non-nil) before returning
// the value of the `debug` field. If the `wrapper` is not available, it returns an
// empty map to ensure safe usage.
//
// Returns:
//   - A `map[string]interface{}` containing the debugging information.
//   - An empty map if the `wrapper` instance is not available.
func (w *wrapper) Debugging() map[string]any {
	if !w.Available() {
		return nil
	}
	return w.debug
}

// OnKeyDebugging retrieves the value of a specific debugging key from the `wrapper` instance.
//
// This function checks if the `wrapper` is available (non-nil) and if the specified debugging key
// is present in the `debug` map. If both conditions are met, it returns the value associated with
// the specified key. Otherwise, it returns `nil` to indicate the key is not available.
//
// Parameters:
//   - `key`: A string representing the debugging key to retrieve.
//
// Returns:
//   - The value associated with the specified debugging key if it exists.
//   - `nil` if the `wrapper` is unavailable or the key is not present in the `debug` map.
func (w *wrapper) OnKeyDebugging(key string) any {
	if !w.Available() || !w.IsDebuggingKeyPresent(key) {
		return nil
	}
	return w.debug[key]
}

// Pagination retrieves the `pagination` instance associated with the `wrapper`.
//
// This function returns the `pagination` field of the `wrapper`, allowing access to
// pagination details such as the current page, total pages, and total items. If no
// pagination information is available, it returns `nil`.
//
// Returns:
//   - A pointer to the `pagination` instance if available.
//   - `nil` if the `pagination` field is not set.
func (w *wrapper) Pagination() *pagination {
	return w.pagination
}

// Meta retrieves the `meta` information from the `wrapper` instance.
//
// This function returns the `meta` field, which contains metadata related to the response or data
// in the `wrapper` instance. If no `meta` information is set, it returns `nil`.
//
// Returns:
//   - A pointer to the `meta` instance associated with the `wrapper`.
//   - `nil` if no `meta` information is available.
func (w *wrapper) Meta() *meta {
	return w.meta
}

// Header retrieves the `header` associated with the `wrapper` instance.
//
// This function returns the `header` field from the `wrapper` instance, which contains
// information about the HTTP response or any other relevant metadata. If the `wrapper`
// instance is correctly initialized, it will return the `header`; otherwise, it may
// return `nil` if the `header` has not been set.
//
// Returns:
//   - A pointer to the `header` instance associated with the `wrapper`.
//   - `nil` if the `header` is not set or the `wrapper` is uninitialized.
func (w *wrapper) Header() *header {
	return w.header
}

// IsDebuggingPresent checks whether debugging information is present in the `wrapper` instance.
//
// This function verifies if the `debug` field of the `wrapper` is not nil and contains at least one entry.
// It returns `true` if debugging information is available; otherwise, it returns `false`.
//
// Returns:
//   - A boolean value indicating whether debugging information is present:
//   - `true` if `debug` is not nil and contains data.
//   - `false` if `debug` is nil or empty.
func (w *wrapper) IsDebuggingPresent() bool {
	return w.Available() && w.debug != nil && len(w.debug) > 0
}

// IsDebuggingKeyPresent checks whether a specific key exists in the `debug` information.
//
// This function first checks if debugging information is present using `IsDebuggingPresent()`.
// Then it uses `coll.MapContainsKey` to verify if the given key is present within the `debug` map.
//
// Parameters:
//   - `key`: The key to search for within the `debug` field.
//
// Returns:
//   - A boolean value indicating whether the specified key is present in the `debug` map:
//   - `true` if the `debug` field is present and contains the specified key.
//   - `false` if `debug` is nil or does not contain the key.
func (w *wrapper) IsDebuggingKeyPresent(key string) bool {
	return w.IsDebuggingPresent() && coll.MapContainsKey(w.debug, key)
}

// IsBodyPresent checks whether the body data is present in the `wrapper` instance.
//
// This function checks if the `data` field of the `wrapper` is not nil, indicating that the body contains data.
//
// Returns:
//   - A boolean value indicating whether the body data is present:
//   - `true` if `data` is not nil.
//   - `false` if `data` is nil.
func (w *wrapper) IsBodyPresent() bool {
	s, ok := w.data.(string)
	if ok {
		return w.Available() &&
			strutil.IsNotEmpty(s) &&
			strutil.IsNotBlank(s)
	}
	return w.Available() && w.data != nil
}

// IsHeaderPresent checks whether header information is present in the `wrapper` instance.
//
// This function checks if the `header` field of the `wrapper` is not nil, indicating that header information is included.
//
// Returns:
//   - A boolean value indicating whether header information is present:
//   - `true` if `header` is not nil.
//   - `false` if `header` is nil.
func (w *wrapper) IsHeaderPresent() bool {
	return w.Available() && w.header != nil
}

// IsMetaPresent checks whether metadata information is present in the `wrapper` instance.
//
// This function checks if the `meta` field of the `wrapper` is not nil, indicating that metadata is available.
//
// Returns:
//   - A boolean value indicating whether metadata is present:
//   - `true` if `meta` is not nil.
//   - `false` if `meta` is nil.
func (w *wrapper) IsMetaPresent() bool {
	return w.Available() && w.meta != nil
}

// IsPagingPresent checks whether pagination information is present in the `wrapper` instance.
//
// This function checks if the `pagination` field of the `wrapper` is not nil, indicating that pagination details are included.
//
// Returns:
//   - A boolean value indicating whether pagination information is present:
//   - `true` if `pagination` is not nil.
//   - `false` if `pagination` is nil.
func (w *wrapper) IsPagingPresent() bool {
	return w.Available() && w.pagination != nil
}

// IsErrorPresent checks whether an error is present in the `wrapper` instance.
//
// This function checks if the `errors` field of the `wrapper` is not nil, indicating that an error has occurred.
//
// Returns:
//   - A boolean value indicating whether an error is present:
//   - `true` if `errors` is not nil.
//   - `false` if `errors` is nil.
func (w *wrapper) IsErrorPresent() bool {
	return w.Available() && w.errors != nil
}

// IsTotalPresent checks whether the total number of items is present in the `wrapper` instance.
//
// This function checks if the `total` field of the `wrapper` is greater than or equal to 0,
// indicating that a valid total number of items has been set.
//
// Returns:
//   - A boolean value indicating whether the total is present:
//   - `true` if `total` is greater than or equal to 0.
//   - `false` if `total` is negative (indicating no total value).
func (w *wrapper) IsTotalPresent() bool {
	return w.Available() && w.total >= 0
}

// IsStatusCodePresent checks whether a valid status code is present in the `wrapper` instance.
//
// This function checks if the `statusCode` field of the `wrapper` is greater than 0,
// indicating that a valid HTTP status code has been set.
//
// Returns:
//   - A boolean value indicating whether the status code is present:
//   - `true` if `statusCode` is greater than 0.
//   - `false` if `statusCode` is less than or equal to 0.
func (w *wrapper) IsStatusCodePresent() bool {
	return w.Available() && w.statusCode > 0
}

// IsError checks whether there is an error present in the `wrapper` instance.
//
// This function returns `true` if the `wrapper` contains an error, which can be any of the following:
//   - An error present in the `errors` field.
//   - A client error (4xx status code) or a server error (5xx status code).
//
// Returns:
//   - A boolean value indicating whether there is an error:
//   - `true` if there is an error present, either in the `errors` field or as an HTTP client/server error.
//   - `false` if no error is found.
func (w *wrapper) IsError() bool {
	return w.IsErrorPresent() || w.IsClientError() || w.IsServerError()
}

// IsSuccess checks whether the HTTP status code indicates a successful response.
//
// This function checks if the `statusCode` is between 200 and 299, inclusive, which indicates a successful HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response was successful:
//   - `true` if the status code is between 200 and 299 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsSuccess() bool {
	return w.Available() && (200 <= w.statusCode) && (w.statusCode <= 299)
}

// IsInformational checks whether the HTTP status code indicates an informational response.
//
// This function checks if the `statusCode` is between 100 and 199, inclusive, which indicates an informational HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response is informational:
//   - `true` if the status code is between 100 and 199 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsInformational() bool {
	return w.Available() && (100 <= w.statusCode) && (w.statusCode <= 199)
}

// IsRedirection checks whether the HTTP status code indicates a redirection response.
//
// This function checks if the `statusCode` is between 300 and 399, inclusive, which indicates a redirection HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response is a redirection:
//   - `true` if the status code is between 300 and 399 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsRedirection() bool {
	return w.Available() && (300 <= w.statusCode) && (w.statusCode <= 399)
}

// IsClientError checks whether the HTTP status code indicates a client error.
//
// This function checks if the `statusCode` is between 400 and 499, inclusive, which indicates a client error HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response is a client error:
//   - `true` if the status code is between 400 and 499 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsClientError() bool {
	return w.Available() && (400 <= w.statusCode) && (w.statusCode <= 499)
}

// IsServerError checks whether the HTTP status code indicates a server error.
//
// This function checks if the `statusCode` is between 500 and 599, inclusive, which indicates a server error HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response is a server error:
//   - `true` if the status code is between 500 and 599 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsServerError() bool {
	return w.Available() && (500 <= w.statusCode) && (w.statusCode <= 599)
}

// IsLastPage checks whether the current page is the last page of results.
//
// This function verifies that pagination information is present and then checks if the current page is the last page.
// It combines the checks of `IsPagingPresent()` and `IsLast()` to ensure that the pagination structure exists
// and that it represents the last page.
//
// Returns:
//   - A boolean value indicating whether the current page is the last page:
//   - `true` if pagination is present and the current page is the last one.
//   - `false` if pagination is not present or the current page is not the last.
func (w *wrapper) IsLastPage() bool {
	return w.Available() && w.IsPagingPresent() && w.pagination.IsLast()
}

// Clone creates a deep copy of the `wrapper` instance.
//
// This function creates a new `wrapper` instance with the same fields as the original instance.
// It creates a new `header`, `meta`, and `pagination` instances and copies the values from the original instance.
// It also creates a new `debug` map and copies the values from the original instance.
//
// Returns:
//   - A pointer to the cloned `wrapper` instance.
//   - `nil` if the `wrapper` instance is not available.
func (w *wrapper) Clone() *wrapper {
	if !w.Available() {
		return New()
	}

	clone := &wrapper{
		statusCode: w.statusCode,
		total:      w.total,
		message:    w.message,
		data:       w.data,
		path:       w.path,
		errors:     w.errors,
	}

	// Clone header
	if w.header != nil {
		clone.header = Header().
			WithCode(w.header.code).
			WithText(w.header.text).
			WithType(w.header.typez).
			WithDescription(w.header.description)
	}

	// Clone meta
	if w.meta != nil {
		clone.meta = Meta().
			WithApiVersion(w.meta.apiVersion).
			WithRequestID(w.meta.requestID).
			WithLocale(w.meta.locale).
			WithRequestedTime(w.meta.requestedTime)

		if w.meta.customFields != nil {
			customFieldsCopy := make(map[string]any)
			for k, v := range w.meta.customFields {
				customFieldsCopy[k] = v
			}
			clone.meta.WithCustomFields(customFieldsCopy)
		}
	}

	// Clone pagination
	if w.pagination != nil {
		clone.pagination = Pages().
			WithPage(w.pagination.page).
			WithPerPage(w.pagination.perPage).
			WithTotalPages(w.pagination.totalPages).
			WithTotalItems(w.pagination.totalItems).
			WithIsLast(w.pagination.isLast)
	}

	// Clone debug
	if w.debug != nil {
		clone.debug = make(map[string]any)
		for k, v := range w.debug {
			clone.debug[k] = v
		}
	}

	return clone
}

// Reset resets the `wrapper` instance to its initial state.
//
// This function sets the `wrapper` instance to its initial state by resetting
// the `statusCode`, `total`, `message`, `path`, `cacheHash`, `data`, `debug`,
// `header`, `errors`, `pagination`, and `cachedWrap` fields to their default values.
// It also resets the `meta` instance to its initial state.
//
// Returns:
//   - A pointer to the reset `wrapper` instance.
//   - `nil` if the `wrapper` instance is not available.
func (w *wrapper) Reset() *wrapper {
	if !w.Available() {
		return New()
	}

	// Reset status code and total
	w.total = 0
	w.statusCode = 0

	// Reset message, path, and cache hash
	w.path = ""
	w.message = ""
	w.cacheHash = ""

	// Reset data, debug, header, errors, pagination, and cached wrap
	w.data = nil
	w.debug = nil
	w.header = nil
	w.errors = nil
	w.pagination = nil
	w.cachedWrap = nil

	// Reset meta
	w.meta = defaultMetaValues()

	return w
}

// DeltaValue retrieves the delta value from the `meta` instance.
//
// This function checks if the `meta` instance is present and returns the `deltaValue` field.
// If the `meta` instance is not present, it returns a default value of `0`.
//
// Returns:
//   - A float64 representing the delta value.
func (w *wrapper) DeltaValue() float64 {
	if !w.IsMetaPresent() {
		return 0
	}
	return w.meta.deltaValue
}

// DeltaCnt retrieves the delta count from the `meta` instance.
//
// This function checks if the `meta` instance is present and returns the `deltaCnt` field.
// If the `meta` instance is not present, it returns a default value of `0`.
//
// Returns:
//   - An integer representing the delta count.
func (w *wrapper) DeltaCnt() int {
	if !w.IsMetaPresent() {
		return 0
	}
	return w.meta.deltaCnt
}

// Available checks whether the `pagination` instance is non-nil.
//
// This function ensures that the `pagination` object exists and is not nil.
// It serves as a safety check to avoid null pointer dereferences when accessing the instance's fields or methods.
//
// Returns:
//   - A boolean value indicating whether the `pagination` instance is non-nil:
//   - `true` if the `pagination` instance is non-nil.
//   - `false` if the `pagination` instance is nil.
func (p *pagination) Available() bool {
	return p != nil
}

// Page retrieves the current page number from the `pagination` instance.
//
// This function checks if the `pagination` instance is available (non-nil) before
// returning the value of the `page` field. If the instance is not available, it
// returns a default value of `0`.
//
// Returns:
//   - An integer representing the current page number.
//   - `0` if the `pagination` instance is not available.
func (p *pagination) Page() int {
	if !p.Available() {
		return 0
	}
	return p.page
}

// PerPage retrieves the number of items per page from the `pagination` instance.
//
// This function checks if the `pagination` instance is available (non-nil) before
// returning the value of the `perPage` field. If the instance is not available, it
// returns a default value of `0`.
//
// Returns:
//   - An integer representing the number of items per page.
//   - `0` if the `pagination` instance is not available.
func (p *pagination) PerPage() int {
	if !p.Available() {
		return 0
	}
	return p.perPage
}

// TotalPages retrieves the total number of pages from the `pagination` instance.
//
// This function checks if the `pagination` instance is available (non-nil) before
// returning the value of the `totalPages` field. If the instance is not available, it
// returns a default value of `0`.
//
// Returns:
//   - An integer representing the total number of pages.
//   - `0` if the `pagination` instance is not available.
func (p *pagination) TotalPages() int {
	if !p.Available() {
		return 0
	}
	return p.totalPages
}

// TotalItems retrieves the total number of items from the `pagination` instance.
//
// This function checks if the `pagination` instance is available (non-nil) before
// returning the value of the `totalItems` field. If the instance is not available, it
// returns a default value of `0`.
//
// Returns:
//   - An integer representing the total number of items.
//   - `0` if the `pagination` instance is not available.
func (p *pagination) TotalItems() int {
	if !p.Available() {
		return 0
	}
	return p.totalItems
}

// IsLast checks whether the current pagination represents the last page.
//
// This function checks the `isLast` field of the `pagination` instance to determine if the current page is the last one.
// The `isLast` field is typically set to `true` when there are no more pages of data available.
//
// Returns:
//   - A boolean value indicating whether the current page is the last:
//   - `true` if `isLast` is `true`, indicating this is the last page of results.
//   - `false` if `isLast` is `false`, indicating more pages are available.
func (p *pagination) IsLast() bool {
	if !p.Available() {
		return true
	}
	return p.isLast
}

// Available checks whether the `meta` instance is available (non-nil).
//
// This function ensures that the `meta` instance is valid and not nil
// to safely access its fields and methods.
//
// Returns:
//   - `true` if the `meta` instance is not nil.
//   - `false` if the `meta` instance is nil.
func (m *meta) Available() bool {
	return m != nil
}

// IsApiVersionPresent checks whether the API version is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `apiVersion` field contains a non-empty value.
//
// Returns:
//   - `true` if `apiVersion` is non-empty.
//   - `false` if `meta` is unavailable or `apiVersion` is empty.
func (m *meta) IsApiVersionPresent() bool {
	return m.Available() && strutil.IsNotEmpty(m.apiVersion)
}

// IsRequestIDPresent checks whether the request ID is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `requestID` field contains a non-empty value.
//
// Returns:
//   - `true` if `requestID` is non-empty.
//   - `false` if `meta` is unavailable or `requestID` is empty.
func (m *meta) IsRequestIDPresent() bool {
	return m.Available() && strutil.IsNotEmpty(m.requestID)
}

// IsLocalePresent checks whether the locale information is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `locale` field contains a non-empty value.
//
// Returns:
//   - `true` if `locale` is non-empty.
//   - `false` if `meta` is unavailable or `locale` is empty.
func (m *meta) IsLocalePresent() bool {
	return m.Available() && strutil.IsNotEmpty(m.locale)
}

// IsRequestedTimePresent checks whether the requested time is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `requestedTime` field is not the zero value of `time.Time`.
//
// Returns:
//   - `true` if `requestedTime` is not the zero value of `time.Time`.
//   - `false` if `meta` is unavailable or `requestedTime` is uninitialized.
func (m *meta) IsRequestedTimePresent() bool {
	return m.Available() && m.requestedTime != time.Time{} && !m.requestedTime.IsZero()
}

// IsCustomFieldPresent checks whether custom fields are present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `customFields` field is non-nil and contains at least one entry.
//
// Returns:
//   - `true` if `customFields` is non-nil and has a non-zero length.
//   - `false` if `meta` is unavailable, `customFields` is nil, or it is empty.
func (m *meta) IsCustomFieldPresent() bool {
	return m.Available() && m.customFields != nil && len(m.customFields) > 0
}

// IsDeltaValuePresent checks whether the delta value is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `deltaValue` field is greater than zero.
//
// Returns:
//   - `true` if `deltaValue` is greater than zero.
//   - `false` if `meta` is unavailable or `deltaValue` is zero or negative.
func (m *meta) IsDeltaValuePresent() bool {
	return m.Available() && m.deltaValue > 0
}

// IsDeltaCntPresent checks whether the delta count is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `deltaCnt` field is greater than zero.
//
// Returns:
//   - `true` if `deltaCnt` is greater than zero.
//   - `false` if `meta` is unavailable or `deltaCnt` is zero or negative.
func (m *meta) IsDeltaCntPresent() bool {
	return m.Available() && m.deltaCnt > 0
}

// IsCustomFieldKeyPresent checks whether a specific key is present in the custom fields of the `meta` instance.
//
// This function first verifies that the `customFields` field is available and contains data using
// `IsCustomFieldPresent`. If so, it checks if the specified key exists within the `customFields` map.
//
// Parameters:
//   - `key`: A string representing the key to search for in the `customFields` map.
//
// Returns:
//   - `true` if the `customFields` map is available and contains the specified key.
//   - `false` if `customFields` is nil, empty, or does not contain the specified key.
func (m *meta) IsCustomFieldKeyPresent(key string) bool {
	return m.IsCustomFieldPresent() && coll.MapContainsKey(m.customFields, key)
}

// OnKeyCustomField retrieves the value associated with a specific key in the custom fields of the `meta` instance.
//
// This function checks whether the `meta` instance is available and whether the specified key exists
// in the `customFields` map. If both conditions are met, it returns the corresponding value for the key.
// If the `meta` instance is unavailable or the key is not present, it returns `nil`.
//
// Parameters:
//   - `key`: A string representing the key to search for in the `customFields` map.
//
// Returns:
//   - The value associated with the specified key in the `customFields` map if the key is present.
//   - `nil` if the `meta` instance is unavailable or the key does not exist in the `customFields` map.
func (m *meta) OnKeyCustomField(key string) any {
	if !m.Available() || !m.IsCustomFieldKeyPresent(key) {
		return nil
	}
	return m.customFields[key]
}

// ApiVersion retrieves the API version from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before attempting
// to retrieve the `apiVersion`. If the `meta` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The API version as a string if available.
//   - An empty string if the `meta` instance is unavailable.
func (m *meta) ApiVersion() string {
	if !m.Available() {
		return ""
	}
	return m.apiVersion
}

// RequestID retrieves the request ID from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before retrieving
// the `requestID`. If the `meta` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The request ID as a string if available.
//   - An empty string if the `meta` instance is unavailable.
func (m *meta) RequestID() string {
	if !m.Available() {
		return ""
	}
	if strutil.IsEmpty(m.requestID) || strutil.IsBlank(m.requestID) {
		m.WithRequestID(randn.CryptoID())
	}
	return m.requestID
}

// Locale retrieves the locale from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before retrieving
// the `locale`. If the `meta` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The locale as a string if available.
//   - An empty string if the `meta` instance is unavailable.
func (m *meta) Locale() string {
	if !m.Available() {
		return ""
	}
	return m.locale
}

// RequestedTime retrieves the requested time from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before retrieving
// the `requestedTime`. If the `meta` instance is unavailable, it returns the zero value
// of `time.Time` (January 1, year 1, 00:00:00 UTC).
//
// Returns:
//   - The requested time as a `time.Time` object if available.
//   - The zero value of `time.Time` if the `meta` instance is unavailable.
func (m *meta) RequestedTime() time.Time {
	if !m.Available() {
		return time.Time{}
	}
	return m.requestedTime
}

// CustomFields retrieves the custom fields from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before retrieving
// the `customFields`. If the `meta` instance is unavailable, it returns `nil`.
//
// Returns:
//   - A map of custom fields if available.
//   - `nil` if the `meta` instance is unavailable.
func (m *meta) CustomFields() map[string]any {
	if !m.Available() {
		return nil
	}
	return m.customFields
}

// Available checks if the `header` instance is non-nil.
//
// This function ensures that the `header` instance is not nil before performing any operations.
// It returns `true` if the `header` is non-nil, and `false` if the `header` is nil.
//
// Returns:
//   - `true` if the `header` instance is not nil.
//   - `false` if the `header` instance is nil.
func (h *header) Available() bool {
	return h != nil
}

// IsCodePresent checks if the `code` field in the `header` instance is present and greater than zero.
//
// This function first checks if the `header` is available (non-nil), and then checks if the `code`
// field is greater than zero, indicating that it is present and valid.
//
// Returns:
//   - `true` if the `code` field is greater than zero.
//   - `false` if the `code` field is either not present (nil) or zero.
func (h *header) IsCodePresent() bool {
	return h.Available() && h.code > 0
}

// IsTextPresent checks if the `text` field in the `header` instance is present and not empty.
//
// This function verifies if the `header` is available and if the `text` field is not empty, using
// the `strutil.IsNotEmpty` utility to ensure the presence of the `text` field.
//
// Returns:
//   - `true` if the `text` field is non-empty.
//   - `false` if the `text` field is either not present (nil) or empty.
func (h *header) IsTextPresent() bool {
	return h.Available() && strutil.IsNotEmpty(h.text)
}

// IsTypePresent checks if the `Type` field in the `header` instance is present and not empty.
//
// This function checks if the `header` instance is available and if the `Type` field is not empty,
// utilizing the `strutil.IsNotEmpty` utility to determine whether the `Type` field contains a value.
//
// Returns:
//   - `true` if the `Type` field is non-empty.
//   - `false` if the `Type` field is either not present (nil) or empty.
func (h *header) IsTypePresent() bool {
	return h.Available() && strutil.IsNotEmpty(h.typez)
}

// IsDescriptionPresent checks if the `description` field in the `header` instance is present and not empty.
//
// This function ensures that the `header` is available and that the `description` field is not empty,
// using `strutil.IsNotEmpty` to check for non-emptiness.
//
// Returns:
//   - `true` if the `description` field is non-empty.
//   - `false` if the `description` field is either not present (nil) or empty.
func (h *header) IsDescriptionPresent() bool {
	return h.Available() && strutil.IsNotEmpty(h.description)
}

// Code retrieves the code value from the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) before retrieving
// the `code` field. If the `header` instance is unavailable, it returns 0.
//
// Returns:
//   - The `code` as an integer if available.
//   - 0 if the `header` instance is unavailable.
func (h *header) Code() int {
	if !h.Available() {
		return 0
	}
	return h.code
}

// Text retrieves the text value from the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) before retrieving
// the `text` field. If the `header` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The `text` as a string if available.
//   - An empty string if the `header` instance is unavailable.
func (h *header) Text() string {
	if !h.Available() {
		return ""
	}
	return h.text
}

// Type retrieves the type value from the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) before retrieving
// the `Type` field. If the `header` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The `Type` as a string if available.
//   - An empty string if the `header` instance is unavailable.
func (h *header) Type() string {
	if !h.Available() {
		return ""
	}
	return h.typez
}

// Description retrieves the description value from the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) before retrieving
// the `description` field. If the `header` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The `description` as a string if available.
//   - An empty string if the `header` instance is unavailable.
func (h *header) Description() string {
	if !h.Available() {
		return ""
	}
	return h.description
}

// WithStatusCode sets the HTTP status code for the `wrapper` instance.
// Ensure that code is between 100 and 599, defaults to 500 if invalid value.
//
// This function updates the `statusCode` field of the `wrapper` and
// returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `code`: An integer representing the HTTP status code to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithStatusCode(code int) *wrapper {
	if code < 100 || code > 599 {
		code = http.StatusInternalServerError
	}
	w.statusCode = code
	w.header = Header().WithCode(code).WithText(http.StatusText(code))
	return w
}

// WithTotal sets the total number of items for the `wrapper` instance.
//
// This function updates the `total` field of the `wrapper` and
// returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `total`: An integer representing the total number of items to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithTotal(total int) *wrapper {
	w.total = total
	return w
}

// WithMessage sets a message for the `wrapper` instance.
//
// This function updates the `message` field of the `wrapper` with the provided string
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `message`: A string message to be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithMessage(message string) *wrapper {
	w.message = message
	return w
}

// WithMessagef sets a formatted message for the `wrapper` instance.
//
// This function constructs a formatted string using the provided format string and arguments,
// assigns it to the `message` field of the `wrapper`, and returns the modified instance.
//
// Parameters:
//   - message: A format string for constructing the message.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance, enabling method chaining.
func (w *wrapper) WithMessagef(message string, args ...any) *wrapper {
	w.message = fmt.Sprintf(message, args...)
	return w
}

// WithBody sets the body data for the `wrapper` instance.
//
// This function updates the `data` field of the `wrapper` with the provided value
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: The value to be set as the body data, which can be any type.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithBody(v any) *wrapper {
	w.data = v
	return w
}

// WithPath sets the request path for the `wrapper` instance.
//
// This function updates the `path` field of the `wrapper` with the provided string
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A string representing the request path.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithPath(v string) *wrapper {
	w.path = v
	return w
}

// WithPathf sets a formatted request path for the `wrapper` instance.
//
// This function constructs a formatted string using the provided format string `v` and arguments `args`,
// assigns the resulting string to the `path` field of the `wrapper`, and returns the modified instance.
//
// Parameters:
//   - v: A format string for constructing the request path.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance, enabling method chaining.
func (w *wrapper) WithPathf(v string, args ...any) *wrapper {
	w.path = fmt.Sprintf(v, args...)
	return w
}

// WithHeader sets the header for the `wrapper` instance.
//
// This function updates the `header` field of the `wrapper` with the provided `header`
// instance and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A pointer to a `header` struct that will be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithHeader(v *header) *wrapper {
	w.header = v
	w.WithStatusCode(w.Header().Code())
	return w
}

// WithMeta sets the metadata for the `wrapper` instance.
//
// This function updates the `meta` field of the `wrapper` with the provided `meta`
// instance and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A pointer to a `meta` struct that will be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithMeta(v *meta) *wrapper {
	w.meta = v
	return w
}

// WithPagination sets the pagination information for the `wrapper` instance.
//
// This function updates the `pagination` field of the `wrapper` with the provided `pagination`
// instance and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A pointer to a `pagination` struct that will be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithPagination(v *pagination) *wrapper {
	w.pagination = v
	return w
}

// WithDebugging sets the debugging information for the `wrapper` instance.
//
// This function updates the `debug` field of the `wrapper` with the provided map of debugging data
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A map containing debugging information to be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithDebugging(v map[string]any) *wrapper {
	w.debug = v
	return w
}

// WithError sets an error for the `wrapper` instance.
//
// This function updates the `errors` field of the `wrapper` with the provided error
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `err`: An error object to be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
// func (w *wrapper) WithError(err error) *wrapper {
// 	w.errors = err
// 	return w
// }

// WithError sets an error for the `wrapper` instance using a plain error message.
//
// This function creates an error object from the provided message, assigns it to
// the `errors` field of the `wrapper`, and returns the modified instance.
//
// Parameters:
//   - message: A string containing the error message to be wrapped as an error object.
//
// Returns:
//   - A pointer to the modified `wrapper` instance to support method chaining.
func (w *wrapper) WithError(message string) *wrapper {
	w.errors = NewError(message)
	return w
}

// WithErrorf sets a formatted error for the `wrapper` instance.
//
// This function uses a formatted string and arguments to construct an error object,
// assigns it to the `errors` field of the `wrapper`, and returns the modified instance.
//
// Parameters:
//   - format: A format string for constructing the error message.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance to support method chaining.
func (w *wrapper) WithErrorf(format string, args ...any) *wrapper {
	w.errors = NewErrorf(format, args...)
	return w
}

// WithErrSck sets an error with a stack trace for the `wrapper` instance.
//
// This function wraps the provided error with stack trace information, assigns it
// to the `errors` field of the `wrapper`, and returns the modified instance.
//
// Parameters:
//   - err: The error object to be wrapped with stack trace information.
//
// Returns:
//   - A pointer to the modified `wrapper` instance to support method chaining.
func (w *wrapper) WithErrSck(err error) *wrapper {
	w.errors = NewErrorAck(err)
	return w
}

// WithErrWrap wraps an existing error with an additional message and sets it for the `wrapper` instance.
//
// This function adds context to the provided error by wrapping it with an additional message.
// The resulting error is assigned to the `errors` field of the `wrapper`.
//
// Parameters:
//   - err: The original error to be wrapped.
//   - message: A string message to add context to the error.
//
// Returns:
//   - A pointer to the modified `wrapper` instance to support method chaining.
func (w *wrapper) WithErrWrap(err error, message string) *wrapper {
	w.errors = WithErrorAck(err, message)
	return w
}

// WithErrWrapf wraps an existing error with a formatted message and sets it for the `wrapper` instance.
//
// This function adds context to the provided error by wrapping it with a formatted message.
// The resulting error is assigned to the `errors` field of the `wrapper`.
//
// Parameters:
//   - err: The original error to be wrapped.
//   - format: A format string for constructing the contextual error message.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance to support method chaining.
func (w *wrapper) WithErrWrapf(err error, format string, args ...any) *wrapper {
	w.errors = NewErrorAckf(err, format, args...)
	return w
}

// WithErrMessage adds a plain contextual message to an existing error and sets it for the `wrapper` instance.
//
// This function wraps the provided error with an additional plain message and assigns it
// to the `errors` field of the `wrapper`.
//
// Parameters:
//   - err: The original error to be wrapped.
//   - message: A plain string message to add context to the error.
//
// Returns:
//   - A pointer to the modified `wrapper` instance to support method chaining.
func (w *wrapper) WithErrMessage(err error, message string) *wrapper {
	w.errors = WithMessage(err, message)
	return w
}

// WithErrMessagef adds a formatted contextual message to an existing error and sets it for the `wrapper` instance.
//
// This function wraps the provided error with an additional formatted message and assigns it
// to the `errors` field of the `wrapper`.
//
// Parameters:
//   - err: The original error to be wrapped.
//   - format: A format string for constructing the contextual error message.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance to support method chaining.
func (w *wrapper) WithErrMessagef(err error, format string, args ...any) *wrapper {
	w.errors = WithMessagef(err, format, args...)
	return w
}

// BindCause sets the error for the `wrapper` instance using its current message.
//
// This function creates an error object from the `message` field of the `wrapper`,
// assigns it to the `errors` field, and returns the modified instance.
// It allows for method chaining.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) BindCause() *wrapper {
	if strutil.IsNotEmpty(w.message) {
		w.errors = NewError(w.message)
	}
	return w
}

// WithDebuggingKV adds a key-value pair to the debugging information in the `wrapper` instance.
//
// This function checks if debugging information is already present. If it is not, it initializes
// an empty map. Then it adds the given key-value pair to the `debug` map and returns the modified
// `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `key`: The key for the debugging information to be added.
//   - `value`: The value associated with the key to be added to the `debug` map.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithDebuggingKV(key string, value any) *wrapper {
	if !w.IsDebuggingPresent() {
		w.debug = make(map[string]any)
	}
	w.debug[key] = value
	return w
}

// WithDebuggingKVf adds a formatted key-value pair to the debugging information in the `wrapper` instance.
//
// This function creates a formatted string value using the provided `format` string and `args`,
// then delegates to `WithDebuggingKV` to add the resulting key-value pair to the `debug` map.
// It returns the modified `wrapper` instance for method chaining.
//
// Parameters:
//   - key: A string representing the key for the debugging information.
//   - format: A format string for constructing the value.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance, enabling method chaining.
func (w *wrapper) WithDebuggingKVf(key string, format string, args ...any) *wrapper {
	return w.WithDebuggingKV(key, fmt.Sprintf(format, args...))
}

// WithApiVersion sets the API version in the `meta` field of the `wrapper` instance.
//
// This function checks if the `meta` information is present in the `wrapper`. If it is not,
// a new `meta` instance is created. Then, it calls the `WithApiVersion` method on the `meta`
// instance to set the API version.
//
// Parameters:
//   - `v`: A string representing the API version to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithApiVersion(v string) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithApiVersion(v)
	return w
}

// WithApiVersionf sets the API version in the `meta` field of the `wrapper` instance using a formatted string.
//
// This function ensures that the `meta` field in the `wrapper` is initialized. If the `meta`
// field is not present, a new `meta` instance is created using the `NewMeta` function.
// Once the `meta` instance is ready, it updates the API version using the `WithApiVersionf` method
// on the `meta` instance. The API version is constructed by interpolating the provided `format`
// string with the variadic arguments (`args`).
//
// Parameters:
//   - format: A format string used to construct the API version.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance, enabling method chaining.
func (w *wrapper) WithApiVersionf(format string, args ...any) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithApiVersionf(format, args...)
	return w
}

// WithRequestID sets the request ID in the `meta` field of the `wrapper` instance.
//
// This function ensures that if `meta` information is not already set in the `wrapper`, a new
// `meta` instance is created. Then, it calls the `WithRequestID` method on the `meta` instance
// to set the request ID.
//
// Parameters:
//   - `v`: A string representing the request ID to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithRequestID(v string) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithRequestID(v)
	return w
}

// WithRequestIDf sets the request ID in the `meta` field of the `wrapper` instance using a formatted string.
//
// This function ensures that the `meta` field in the `wrapper` is initialized. If the `meta` field
// is not already present, a new `meta` instance is created using the `NewMeta` function.
// Once the `meta` instance is ready, it updates the request ID by calling the `WithRequestIDf`
// method on the `meta` instance. The request ID is constructed using the provided `format` string
// and the variadic `args`.
//
// Parameters:
//   - format: A format string used to construct the request ID.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance, allowing for method chaining.
func (w *wrapper) WithRequestIDf(format string, args ...any) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithRequestIDf(format, args...)
	return w
}

// RandRequestID generates and sets a random request ID in the `meta` field of the `wrapper` instance.
//
// This function checks if the `meta` field is present in the `wrapper`. If it is not,
// a new `meta` instance is created. Then, it calls the `RandRequestID` method on the `meta`
// instance to generate and set a random request ID.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) RandRequestID() *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.RandRequestID()
	return w
}

// RandDeltaValue generates and sets a random delta value in the `meta` field of the `wrapper` instance.
//
// This function checks if the `meta` field is present in the `wrapper`. If it is not,
// a new `meta` instance is created. Then, it calls the `RandDeltaValue` method on the `meta`
// instance to generate and set a random delta value.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) RandDeltaValue() *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.RandDeltaValue()
	w.meta.IncreaseDeltaCnt()
	return w
}

// IncreaseDeltaCnt increments the delta count in the `meta` field of the `wrapper` instance.
//
// This function ensures the `meta` field is present, creating a new instance if needed, and
// increments the delta count in the `meta` using the `IncreaseDeltaCnt` method.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) IncreaseDeltaCnt() *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.IncreaseDeltaCnt()
	return w
}

// DecreaseDeltaCnt decrements the delta count in the `meta` field of the `wrapper` instance.
//
// This function ensures the `meta` field is present, creating a new instance if needed, and
// decrements the delta count in the `meta` using the `DecreaseDeltaCnt` method.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) DecreaseDeltaCnt() *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.DecreaseDeltaCnt()
	return w
}

// WithLocale sets the locale in the `meta` field of the `wrapper` instance.
//
// This function ensures the `meta` field is present, creating a new instance if needed, and
// sets the locale in the `meta` using the `WithLocale` method.
//
// Parameters:
//   - `v`: A string representing the locale to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithLocale(v string) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithLocale(v)
	return w
}

// WithRequestedTime sets the requested time in the `meta` field of the `wrapper` instance.
//
// This function ensures that the `meta` field exists, and if not, creates a new one. It then
// sets the requested time in the `meta` using the `WithRequestedTime` method.
//
// Parameters:
//   - `v`: A `time.Time` value representing the requested time.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithRequestedTime(v time.Time) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithRequestedTime(v)
	return w
}

// WithCustomFields sets the custom fields in the `meta` field of the `wrapper` instance.
//
// This function checks if the `meta` field is present. If not, it creates a new `meta` instance
// and sets the provided custom fields using the `WithCustomFields` method.
//
// Parameters:
//   - `values`: A map representing the custom fields to set in the `meta`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithCustomFields(values map[string]any) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithCustomFields(values)
	return w
}

// WithCustomFieldKV sets a specific custom field key-value pair in the `meta` field of the `wrapper` instance.
//
// This function ensures that if the `meta` field is not already set, a new `meta` instance is created.
// It then adds the provided key-value pair to the custom fields of `meta` using the `WithCustomFieldKV` method.
//
// Parameters:
//   - `key`: A string representing the custom field key to set.
//   - `value`: The value associated with the custom field key.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithCustomFieldKV(key string, value any) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithCustomFieldKV(key, value)
	return w
}

// WithCustomFieldKVf sets a specific custom field key-value pair in the `meta` field of the `wrapper` instance
// using a formatted value.
//
// This function constructs a formatted string value using the provided `format` string and arguments (`args`).
// It then calls the `WithCustomFieldKV` method to add or update the custom field with the specified key and
// the formatted value. If the `meta` field of the `wrapper` instance is not initialized, it is created
// before setting the custom field.
//
// Parameters:
//   - key: A string representing the key for the custom field.
//   - format: A format string to construct the value.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `wrapper` instance, enabling method chaining.
func (w *wrapper) WithCustomFieldKVf(key string, format string, args ...any) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = Meta()
	}
	w.meta.WithCustomFieldKVf(key, format, args...)
	return w
}

// WithPage sets the current page number in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified page number is then
// applied to the pagination instance.
//
// Parameters:
//   - v: The page number to set.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithPage(v int) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = Pages()
	}
	w.pagination.WithPage(v)
	return w
}

// WithPerPage sets the number of items per page in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified items-per-page value
// is then applied to the pagination instance.
//
// Parameters:
//   - v: The number of items per page to set.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithPerPage(v int) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = Pages()
	}
	w.pagination.WithPerPage(v)
	return w
}

// WithTotalPages sets the total number of pages in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified total pages value
// is then applied to the pagination instance.
//
// Parameters:
//   - v: The total number of pages to set.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithTotalPages(v int) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = Pages()
	}
	w.pagination.WithTotalPages(v)
	return w
}

// WithTotalItems sets the total number of items in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified total items value
// is then applied to the pagination instance.
//
// Parameters:
//   - v: The total number of items to set.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithTotalItems(v int) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = Pages()
	}
	w.pagination.WithTotalItems(v)
	return w
}

// WithIsLast sets whether the current page is the last one in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified boolean value is then
// applied to indicate whether the current page is the last.
//
// Parameters:
//   - v: A boolean indicating whether the current page is the last.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithIsLast(v bool) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = Pages()
	}
	w.pagination.WithIsLast(v)
	return w
}

// Hash256 generates a hash string for the `wrapper` instance.
//
// This method concatenates the values of the `statusCode`, `message`, `data`, and `meta` fields
// into a single string and then computes a hash of that string using the `strutil.Hash256` function.
// The resulting hash string can be used for various purposes, such as caching or integrity checks.
func (w *wrapper) Hash256() (string, *wrapper) {
	if !w.Available() {
		return "", w
	}
	h, err := hashy.Hash256(w.StatusCode(), w.message, w.data, w.meta.Respond())
	if err != nil {
		return "", New().
			WithHeader(InternalServerError).
			WithErrSck(err).
			WithMessage("Failed to generate hash")
	}
	return h, New().
		WithHeader(OK).
		WithMessage("Successfully generated hash")
}

// Hash256Safe generates a hash string for the `wrapper` instance.
//
// This method generates a hash string for the `wrapper` instance using the `Hash256` method.
// If the `wrapper` instance is not available or the hash generation fails, it returns an empty string.
//
// Returns:
//   - A string representing the hash value.
//   - An empty string if the `wrapper` instance is not available or the hash generation fails.
func (w *wrapper) Hash256Safe() string {
	hash, _w := w.Hash256()
	if _w.IsError() {
		return ""
	}
	return hash
}

// Hash generates a hash value for the `wrapper` instance.
//
// This method generates a hash value for the `wrapper` instance using the `Hash` method.
// If the `wrapper` instance is not available or the hash generation fails, it returns an error.
//
// Returns:
//   - A uint64 representing the hash value.
//   - An error if the `wrapper` instance is not available or the hash generation fails.
func (w *wrapper) Hash() (uint64, *wrapper) {
	if !w.Available() {
		return 0, w
	}
	h, err := hashy.Hash(w.StatusCode(), w.message, w.data, w.meta.Respond())
	if err != nil {
		return 0, New().
			WithHeader(InternalServerError).
			WithErrSck(err).
			WithMessage("Failed to generate hash")
	}
	return h, New().
		WithHeader(OK).
		WithMessage("Successfully generated hash")
}

// HashSafe generates a hash value for the `wrapper` instance.

// This method generates a hash value for the `wrapper` instance using the `Hash` method.
// If the `wrapper` instance is not available or the hash generation fails, it returns an empty string.
//
// Returns:
//   - A string representing the hash value.
//   - An empty string if the `wrapper` instance is not available or the hash generation fails.
func (w *wrapper) HashSafe() uint64 {
	hash, _w := w.Hash()
	if _w.IsError() {
		return 0
	}
	return hash
}

// WithStreaming enables streaming mode for the wrapper and returns a streaming wrapper for enhanced data transfer capabilities.
//
// This function is the primary entry point for activating streaming functionality on an existing wrapper instance.
// It creates a new StreamingWrapper that preserves the metadata and context of the original wrapper while adding
// streaming-specific features such as chunk-based transfer, compression, progress tracking, and bandwidth throttling.
// The returned StreamingWrapper allows for method chaining to configure streaming parameters before initiating transfer.
//
// Parameters:
//   - reader: An io.Reader implementation providing the source data stream (e.g., *os.File, *http.Response.Body, *bytes.Buffer).
//     Cannot be nil; streaming will fail if no valid reader is provided.
//   - config: A *StreamConfig containing streaming configuration options (chunk size, compression, strategy, concurrency).
//     If nil, a default configuration is automatically created with sensible defaults:
//   - ChunkSize: 65536 bytes (64KB)
//   - Strategy: STRATEGY_BUFFERED (balanced throughput and memory)
//   - Compression: COMP_NONE
//   - MaxConcurrentChunks: 4
//
// Returns:
//   - A pointer to a new StreamingWrapper instance that wraps the original wrapper.
//   - The StreamingWrapper preserves all metadata from the original wrapper.
//   - If the receiver wrapper is nil, creates a new default wrapper before enabling streaming.
//   - The returned StreamingWrapper can be chained with configuration methods before calling Start().
//
// Example:
//
//	file, _ := os.Open("large_file.bin")
//	defer file.Close()
//
//	// Simple streaming with defaults
//	result := replify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/file").
//	    WithStreaming(file, nil).
//	    WithChunkSize(1024 * 1024).
//	    WithCompressionType(COMP_GZIP).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            fmt.Printf("Transferred: %.2f MB / %.2f MB\n",
//	                float64(p.TransferredBytes) / 1024 / 1024,
//	                float64(p.TotalBytes) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background()).
//	    WithMessage("File transfer completed")
//
// See Also:
//   - AsStreaming: Simplified version with default configuration
//   - Start: Initiates the streaming operation
//   - WithChunkSize: Configures chunk size
//   - WithCompressionType: Enables data compression
func (w *wrapper) WithStreaming(reader io.Reader, config *StreamConfig) *StreamingWrapper {
	if w == nil {
		return NewStreaming(reader, config)
	}
	if config == nil {
		config = NewStreamConfig()
	}

	sw := NewStreaming(reader, config)
	// Copy existing wrapper metadata
	sw.wrapper = w
	sw.wrapper.WithMessage("Streaming mode enabled")

	return sw
}

// AsStreaming converts a regular wrapper instance into a streaming-enabled response with default configuration.
//
// This function provides a simplified, one-line alternative to WithStreaming for common streaming scenarios.
// It automatically creates a new wrapper if the receiver is nil and applies default streaming configuration,
// eliminating the need for manual configuration object creation. This is ideal for quick implementations where
// standard settings (64KB chunks, buffered strategy, no compression) are acceptable.
//
// Parameters:
//   - reader: An io.Reader implementation providing the source data stream (e.g., *os.File, *http.Response.Body, *bytes.Buffer).
//     Cannot be nil; streaming will fail if no valid reader is provided.
//
// Returns:
//   - A pointer to a new StreamingWrapper instance configured with default settings:
//   - ChunkSize: 65536 bytes (64KB)
//   - Strategy: STRATEGY_BUFFERED
//   - Compression: COMP_NONE
//   - MaxConcurrentChunks: 4
//   - UseBufferPool: true
//   - ReadTimeout: 30 seconds
//   - WriteTimeout: 30 seconds
//   - If the receiver wrapper is nil, automatically creates a new wrapper before enabling streaming.
//   - Returns a StreamingWrapper ready for optional configuration before calling Start().
//
// Example:
//
//	// Minimal streaming setup with defaults - best for simple file downloads
//	file, _ := os.Open("document.pdf")
//	defer file.Close()
//
//	result := replify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/document").
//	    AsStreaming(file).
//	    WithTotalBytes(fileSize).
//	    Start(context.Background())
//
//	// Or without creating a new wrapper first
//	result := (*wrapper)(nil).
//	    AsStreaming(file).
//	    Start(context.Background())
//
// Comparison:
//
//	// Using AsStreaming (simple, defaults only)
//	streaming := response.AsStreaming(reader)
//
//	// Using WithStreaming (more control)
//	streaming := response.WithStreaming(reader, &StreamConfig{
//	    ChunkSize:           512 * 1024,
//	    Compression:         COMP_GZIP,
//	    MaxConcurrentChunks: 8,
//	})
//
// See Also:
//   - WithStreaming: For custom streaming configuration
//   - NewStreamConfig: To create custom configuration objects
//   - Start: Initiates the streaming operation
//   - WithCallback: Adds progress tracking after AsStreaming
func (w *wrapper) AsStreaming(reader io.Reader) *StreamingWrapper {
	if w == nil {
		w = New()
	}
	return w.WithStreaming(reader, NewStreamConfig())
}

// Respond generates a map representation of the `wrapper` instance.
//
// This method collects various fields of the `wrapper` (e.g., `data`, `header`, `meta`, etc.)
// and organizes them into a key-value map. Only non-nil or meaningful fields are added
// to the resulting map to ensure a clean and concise response structure.
//
// Fields included in the response:
//   - `data`: The primary data payload, if present.
//   - `headers`: The structured header details, if present.
//   - `meta`: Metadata about the response, if present.
//   - `pagination`: Pagination details, if applicable.
//   - `debug`: Debugging information, if provided.
//   - `total`: Total number of items, if set to a valid non-negative value.
//   - `status_code`: The HTTP status code, if greater than 0.
//   - `message`: A descriptive message, if not empty.
//   - `path`: The request path, if not empty.
//
// Returns:
//   - A `map[string]interface{}` containing the structured response data.
func (w *wrapper) Respond() map[string]any {
	if !w.Available() {
		return nil
	}
	w.cacheMutex.RLock()
	hash := w.Hash256Safe()

	if w.cacheHash == hash && w.cachedWrap != nil {
		defer w.cacheMutex.RUnlock()
		return w.cachedWrap
	}
	w.cacheMutex.RUnlock()

	w.cacheMutex.Lock()
	defer w.cacheMutex.Unlock()

	response := w.build()
	w.cachedWrap = response
	w.cacheHash = hash

	return response
}

// R represents a wrapper around the main `wrapper` struct. It is used as a high-level
// abstraction to provide a simplified interface for handling API responses.
// The `R` type allows for easier manipulation of the wrapped data, metadata, and other
// response components, while maintaining the flexibility of the underlying `wrapper` structure.
//
// Example usage:
//
//	var response replify.R = replify.NewWrap().Reply()
//	fmt.Println(response.Json())  // Prints the wrapped response details, including data, headers, and metadata.
func (w *wrapper) Reply() R {
	return R{wrapper: w}
}

// ReplyPtr returns a pointer to a new R instance that wraps the current `wrapper`.
//
// This method creates a new `R` struct, initializing it with the current `wrapper` instance,
// and returns a pointer to this new `R` instance. This allows for easier manipulation
// of the wrapped data and metadata through the `R` abstraction.
//
// Returns:
//   - A pointer to an `R` struct that wraps the current `wrapper` instance.
//
// Example usage:
//
//	var responsePtr *replify.R = replify.NewWrap().ReplyPtr()
//	fmt.Println(responsePtr.Json())  // Prints the wrapped response details, including data, headers, and metadata.
func (w *wrapper) ReplyPtr() *R {
	return &R{wrapper: w}
}

// Json serializes the `wrapper` instance into a compact JSON string.
//
// This function uses the `encoding.Json` utility to generate a JSON representation
// of the `wrapper` instance. The output is a compact JSON string with no additional
// whitespace or formatting.
//
// Returns:
//   - A compact JSON string representation of the `wrapper` instance.
func (w *wrapper) Json() string {
	return jsonpass(w.Respond())
}

// JsonPretty serializes the `wrapper` instance into a prettified JSON string.
//
// This function uses the `encoding.JsonPretty` utility to generate a JSON representation
// of the `wrapper` instance. The output is a human-readable JSON string with
// proper indentation and formatting for better readability.
//
// Returns:
//   - A prettified JSON string representation of the `wrapper` instance.
func (w *wrapper) JsonPretty() string {
	return jsonpretty(w.Respond())
}

// build generates a map representation of the `wrapper` instance.
// This method collects various fields of the `wrapper` (e.g., `data`, `header`, `meta`, etc.)
// and organizes them into a key-value map. It ensures that only non-empty or meaningful fields
// are included in the resulting map, providing a clean and structured response.
// The following fields are included in the response:
//   - `data`: The primary data payload, if present.
//   - `headers`: The structured header details, if present.
//   - `meta`: Metadata about the response, if present.
//   - `pagination`: Pagination details, if applicable.
//   - `debug`: Debugging information, if provided.
//   - `total`: Total number of items, if set to a valid non-negative value.
//   - `status_code`: The HTTP status code, if greater than 0.
//   - `message`: A descriptive message, if not empty.
//   - `path`: The request path, if not empty.
//
// Returns:
//   - A `map[string]interface{}` containing the structured response data.
func (w *wrapper) build() map[string]any {
	m := make(map[string]any)
	if w.IsBodyPresent() {
		m["data"] = w.data
	}
	if w.IsHeaderPresent() {
		m["headers"] = w.header.Respond()
	}
	if w.IsMetaPresent() {
		m["meta"] = w.meta.Respond()
	}
	if w.IsPagingPresent() {
		m["pagination"] = w.pagination.Respond()
	}
	if w.IsDebuggingPresent() {
		m["debug"] = w.debug
	}
	if w.IsTotalPresent() {
		m["total"] = w.total
	}
	if w.IsStatusCodePresent() {
		m["status_code"] = w.statusCode
	}
	if strutil.IsNotEmpty(w.message) {
		m["message"] = w.message
	}
	if strutil.IsNotEmpty(w.path) {
		m["path"] = w.path
	}
	return m
}
