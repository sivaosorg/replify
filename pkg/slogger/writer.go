package slogger

import (
"io"
"os"
)
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
