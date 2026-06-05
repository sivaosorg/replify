package sysx

import (
	"bytes"
	"io"
)

// NewMemBlob wraps the supplied byte slice in a MemBlob backing.
//
// The slice is retained, not copied; callers must not mutate it for the
// lifetime of the returned value. The Reader is positioned at offset 0.
//
// Parameters:
//   - `data`: the payload to expose as a seekable, closeable stream.
//
// Returns:
//
// A pointer to a MemBlob ready for consumption by Resource.
//
// Example:
//
//	blob := sysx.NewMemBlob([]byte("hello"))
//	defer blob.Close()
func NewMemBlob(data []byte) *MemBlob {
	return &MemBlob{data: data, reader: bytes.NewReader(data)}
}

// Bytes returns the underlying byte slice. The returned slice is the
// same one passed to NewMemBlob; callers must not mutate it.
//
// Returns:
//
// The payload bytes.
func (m *MemBlob) Bytes() []byte {
	if m == nil {
		return nil
	}
	return m.data
}

// Len returns the total number of bytes in the blob.
//
// Returns:
//
// The blob length in bytes.
func (m *MemBlob) Len() int64 {
	if m == nil {
		return 0
	}
	return int64(len(m.data))
}

// Read implements io.Reader by delegating to the embedded *bytes.Reader.
func (m *MemBlob) Read(p []byte) (int, error) {
	return m.reader.Read(p)
}

// Seek implements io.Seeker by delegating to the embedded *bytes.Reader.
func (m *MemBlob) Seek(offset int64, whence int) (int64, error) {
	return m.reader.Seek(offset, whence)
}

// Close implements io.Closer. It is a no-op; the underlying byte slice is
// released by the garbage collector when the MemBlob becomes unreachable.
//
// Returns:
//
// nil.
func (m *MemBlob) Close() error {
	return nil
}

// compile-time interface assertion.
var _ ReadSeekCloser = (*MemBlob)(nil)
var _ io.Seeker = (*MemBlob)(nil)
