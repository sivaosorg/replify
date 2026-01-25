package replify

import (
	"fmt"
	"io"
)

// NewError returns an error with the supplied message and records the stack trace
// at the point it was called. The error contains the message and the stack trace
// which can be used for debugging or logging the error along with the call stack.
//
// Usage example:
//
//	err := NewError("Something went wrong")
//	fmt.Println(err) // "Something went wrong" along with stack trace
func NewError(message string) error {
	return &underlying{
		msg:   message,
		stack: Callers(),
	}
}

// NewErrorf formats the given arguments according to the format specifier and
// returns the formatted string as an error. It also records the stack trace at
// the point it was called.
//
// Usage example:
//
//	err := NewErrorf("Failed to load file %s", filename)
//	fmt.Println(err) // "Failed to load file <filename>" along with stack trace
func NewErrorf(format string, args ...any) error {
	return &underlying{
		msg:   fmt.Sprintf(format, args...),
		stack: Callers(),
	}
}

// NewErrorAck annotates an existing error with a stack trace at the point
// NewErrorAck was called. If the provided error is nil, it simply returns nil.
//
// Usage example:
//
//	err := errors.New("original error")
//	errWithStack := NewErrorAck(err)
//	fmt.Println(errWithStack) // original error with stack trace
func NewErrorAck(err error) error {
	if err == nil {
		return nil
	}
	return &underlyingStack{
		err,
		Callers(),
	}
}

// WithErrWrap returns an error that annotates the provided error with a new message
// and a stack trace at the point WithErrWrap was called. If the provided error is nil,
// WithErrWrap returns nil.
//
// Usage example:
//
//	err := errors.New("file not found")
//	wrappedErr := WithErrWrap(err, "Failed to read the file")
//	fmt.Println(wrappedErr) // "Failed to read the file: file not found" with stack trace
func WithErrWrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &underlyingMessage{
		cause: err,
		msg:   message,
	}
	return &underlyingStack{
		err,
		Callers(),
	}
}

// WithErrWrapf returns an error that annotates the provided error with a formatted message
// and a stack trace at the point WithErrWrapf was called. If the provided error is nil,
// WithErrWrapf returns nil.
//
// Usage example:
//
//	err := errors.New("file not found")
//	wrappedErr := WithErrWrapf(err, "Failed to load file %s", filename)
//	fmt.Println(wrappedErr) // "Failed to load file <filename>: file not found" with stack trace
func WithErrWrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	err = &underlyingMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &underlyingStack{
		err,
		Callers(),
	}
}

// WithMessage annotates an existing error with a new message. If the error is nil,
// it returns nil.
//
// Usage example:
//
//	err := errors.New("original error")
//	errWithMessage := WithMessage(err, "Additional context")
//	fmt.Println(errWithMessage) // "Additional context: original error"
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &underlyingMessage{
		cause: err,
		msg:   message,
	}
}

// WithMessagef annotates an existing error with a formatted message. If the error
// is nil, it returns nil.
//
// Usage example:
//
//	err := errors.New("original error")
//	errWithMessage := WithMessagef(err, "Context: %s", "something went wrong")
//	fmt.Println(errWithMessage) // "Context: something went wrong: original error"
func WithMessagef(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &underlyingMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

// Cause traverses the error chain and returns the underlying cause of the error
// if it implements the `Cause()` method. If the error doesn't implement `Cause()`,
// it simply returns the original error. If the error is nil, nil is returned.
//
// Usage example:
//
//	err := Wrap(errors.New("file not found"), "Failed to open file")
//	causeErr := Cause(err)
//	fmt.Println(causeErr) // "file not found"
//
// An error value has a cause if it implements the following
// interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Error implements the error interface for the `underlying` type, returning the
// message stored in the error object. It is used to retrieve the error message.
//
// Usage example:
//
//	err := WithError("Something went wrong")
//	fmt.Println(err.Error()) // "Something went wrong"
func (u *underlying) Error() string {
	return u.msg
}

// Format formats the error according to the specified format verb. If the `+`
// flag is provided, it will format both the message and the stack trace. Otherwise,
// it will format just the message.
//
// Usage example:
//
//	err := WithError("Something went wrong")
//	fmt.Printf("%+v\n", err) // "Something went wrong" with stack trace
func (u *underlying) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, u.msg)
			u.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, u.msg)
	case 'q':
		fmt.Fprintf(s, "%q", u.msg)
	}
}

// Cause returns the underlying error cause for the `underlyingStack` type.
func (u *underlyingStack) Cause() error { return u.error }

// Unwrap returns the underlying error for the `underlyingStack` type.
func (u *underlyingStack) Unwrap() error { return u.error }

// Format formats the error with the stack trace and the error message. If the
// `+` flag is set, it will include the stack trace as well.
//
// Usage example:
//
//	err := Wrap(errors.New("file not found"), "Failed to open file")
//	fmt.Printf("%+v\n", err) // "Failed to open file: file not found" with stack trace
func (u *underlyingStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", u.Cause())
			u.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, u.Error())
	case 'q':
		fmt.Fprintf(s, "%q", u.Error())
	}
}

// Cause returns the underlying cause of the error for the `underlyingMessage` type.
func (u *underlyingMessage) Cause() error { return u.cause }

// Unwrap returns the underlying cause of the error for the `underlyingMessage` type.
func (u *underlyingMessage) Unwrap() error { return u.cause }

// Error returns the error message concatenated with the underlying error message
// for the `underlyingMessage` type.
//
// Usage example:
//
//	err := WithMessage(errors.New("file not found"), "Failed to open file")
//	fmt.Println(err) // "Failed to open file: file not found"
func (u *underlyingMessage) Error() string {
	return u.msg + ": " + u.cause.Error()
}
