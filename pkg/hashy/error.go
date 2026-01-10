package hashy

import "fmt"

// Error returns the error message for the `ErrNotStringer` type.
//
// Returns:
//   - A string representing the error message.
//
// Example:
//
//	err := &ErrNotStringer{Field: "name"}
//	fmt.Println(err.Error()) // "pkg.hash: field \"name\" has hash:\"string\" tag but does not implement fmt.Stringer"
func (e *ErrNotStringer) Error() string {
	return fmt.Sprintf("pkg.hash: field %q has hash:\"string\" tag but does not implement fmt.Stringer", e.Field)
}
