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
	if mw == nil {
		return 0, nil
	}
	mw.mu.Lock()
	writers := make([]io.Writer, len(mw.writers))
	copy(writers, mw.writers)
	mw.mu.Unlock()

	var firstN int
	var firstErr error
	for i, w := range writers {
		nn, werr := w.Write(p)
		if i == 0 {
			firstN = nn
			firstErr = werr
		}
	}
	return firstN, firstErr
}

// Writers returns a copy of the writer list.
//
// Returns:
//
// a copy of the []io.Writer slice registered with this MultiWriter.
func (mw *MultiWriter) Writers() []io.Writer {
	if mw == nil || mw.writers == nil {
		return nil
	}
	mw.mu.Lock()
	defer mw.mu.Unlock()
	result := make([]io.Writer, len(mw.writers))
	copy(result, mw.writers)
	return result
}

// AddWriter appends a writer to the list.
// This method is safe for concurrent use.
//
// Parameters:
//   - `w`: the writer to add
func (mw *MultiWriter) AddWriter(w io.Writer) {
	if mw == nil {
		return
	}
	mw.mu.Lock()
	mw.writers = append(mw.writers, w)
	mw.mu.Unlock()
}

// Stdout returns os.Stdout as an io.Writer.
//
// Returns:
//
// os.Stdout wrapped as an io.Writer.
func Stdout() io.Writer { return os.Stdout }

// Stderr returns os.Stderr as an io.Writer.
//
// Returns:
//
// os.Stderr wrapped as an io.Writer.
func Stderr() io.Writer { return os.Stderr }
