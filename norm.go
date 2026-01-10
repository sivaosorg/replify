package replify

import (
	"net/http"
)

// NormPagination normalizes the pagination information in the wrapper.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. It then calls the `Normalize` method
// on the pagination instance to ensure its values are consistent.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) NormPagination() *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = Pages()
	} else {
		w.pagination.Normalize()
		w.RandRequestID() // Indicate that a change has occurred
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
		w.RandRequestID() // Indicate that a change has occurred
	}
	return w
}
