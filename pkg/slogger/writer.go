package slogger

import (
	"io"
	"os"
)

// MultiWriter fans log output to multiple io.Writer targets simultaneously.
type MultiWriter struct {
	writers []io.Writer
}

// NewMultiWriter returns a MultiWriter that writes to all provided writers.
//
// Parameters:
//   - `writers`: one or more destination writers
//
// Returns:
//
// a *MultiWriter targeting every supplied writer.
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	dst := make([]io.Writer, len(writers))
	copy(dst, writers)
	return &MultiWriter{writers: dst}
}

// Write writes p to every registered writer.
// All writers are attempted even when an earlier write fails.
//
// Parameters:
//   - `p`: the bytes to write
//
// Returns:
//
// the byte count returned by the first writer and the first error encountered.
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	var firstN int
	var firstErr error
	for i, w := range mw.writers {
		nn, werr := w.Write(p)
		if i == 0 {
			firstN = nn
			firstErr = werr
		}
	}
	return firstN, firstErr
}

// Stdout returns os.Stdout as an io.Writer.
//
// Returns:
//
// the standard output writer.
func Stdout() io.Writer { return os.Stdout }

// Stderr returns os.Stderr as an io.Writer.
//
// Returns:
//
// the standard error writer.
func Stderr() io.Writer { return os.Stderr }
