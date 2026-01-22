package replify

import (
	"net/http"
	"time"
)

// Normalize performs a comprehensive normalization of the wrapper instance.
//
// It sequentially calls the `NormHSC` method to normalize the relationship
// between the header and status code, followed by the `NormPaging` method
// to normalize the pagination information.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) Normalize() *wrapper {
	return w.NormHSC().NormPaging().NormMeta()
}

// NormPaging normalizes the pagination information in the wrapper.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. It then calls the `Normalize` method
// on the pagination instance to ensure its values are consistent.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) NormPaging() *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = Pages()
	} else {
		w.pagination.Normalize()
		w.RandDeltaValue() // Indicate that a change has occurred
	}
	return w
}

// NormHSC normalizes the relationship between the header and status code.
//
// If the status code is not present but the header is, it sets the status code
// from the header's code. If the header is not present but the status code is,
// it creates a new header with the status code and its corresponding text.
//
// If both the status code and header are present, it ensures the status code
// matches the header's code.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) NormHSC() *wrapper {
	hasChanges := false
	switch {
	case !w.IsStatusCodePresent() && w.IsHeaderPresent():
		w.statusCode = w.header.Code()
		hasChanges = true
	case !w.IsHeaderPresent() && w.IsStatusCodePresent():
		w.header = Header().WithCode(w.statusCode).WithText(http.StatusText(w.statusCode))
		hasChanges = true
	case w.IsStatusCodePresent() && w.IsHeaderPresent():
		if w.statusCode == w.header.Code() {
			break
		}
		w.statusCode = w.header.Code()
		hasChanges = true
	}

	if hasChanges {
		w.RandDeltaValue() // Indicate that a change has occurred
	}
	return w
}

// NormMeta normalizes the metadata in the wrapper.
//
// If the meta object is not already initialized, it creates a new one
// using the `Meta` function. It then ensures that essential fields such as
// locale, API version, request ID, and requested time are set to default
// values if they are not already present.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) NormMeta() *wrapper {
	hasChanges := false

	if !w.IsMetaPresent() {
		w.meta = Meta()
		hasChanges = true
	}

	// Set default locale if not present
	if !w.meta.IsLocalePresent() {
		w.meta.locale = string(LocaleEnUS)
		hasChanges = true
	}

	// Set default API version if not present
	if !w.meta.IsApiVersionPresent() {
		w.meta.apiVersion = "v0.0.1"
		hasChanges = true
	}

	// Generate request ID if not present
	if !w.meta.IsRequestIDPresent() {
		w.meta.RandRequestID()
		hasChanges = true
	}

	// Set requested time if not present
	if !w.meta.IsRequestedTimePresent() {
		w.meta.requestedTime = time.Now()
		hasChanges = true
	}

	if hasChanges {
		w.RandDeltaValue() // Indicate that a change has occurred
	}
	return w
}
