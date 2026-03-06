package slogger

import (
"io"
"os"
)

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
func Stdout() io.Writer { return os.Stdout }

// Stderr returns os.Stderr as an io.Writer.
func Stderr() io.Writer { return os.Stderr }
