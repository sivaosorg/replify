package replify

import "github.com/sivaosorg/replify/pkg/sysx"

// Resource returns the underlying [sysx.Resource], which exposes the
// serialized payload as a seekable, closeable stream ([sysx.ReadSeekCloser]).
// Valid until Close is called; returns nil when the Dump itself is nil.
func (d *Dump) Resource() *sysx.Resource {
	if d == nil {
		return nil
	}
	return d.syr
}

// Filepath returns the destination path of the permanent on-disk file when
// the Dump was produced by [wrapper.DumpTo]. Returns an empty string for
// Dumps produced by [wrapper.Dump].
func (d *Dump) Filepath() string {
	if d == nil {
		return ""
	}
	return d.filepath
}

// Size returns the byte length of the serialized payload as reported by the
// underlying [sysx.Resource]. Returns 0 when the Dump is nil.
func (d *Dump) Size() int64 {
	if d == nil || d.syr == nil {
		return 0
	}
	return d.syr.Size()
}

// Rewind seeks the content stream back to offset 0, allowing the payload to
// be read again without creating a new Dump. Useful in retry paths or when
// the same body must be forwarded to multiple destinations.
//
// Returns [sysx.ErrNilResource] when the Resource content is nil.
func (d *Dump) Rewind() error {
	if d == nil || d.syr == nil {
		return nil
	}
	return d.syr.Rewind()
}

// Close removes the backing temporary file and releases all held resources.
// It is safe to call from multiple goroutines simultaneously; the cleanup
// runs exactly once via [sync.Once] — subsequent calls are no-ops that
// return nil. The first caller receives any I/O error from the cleanup.
//
// When produced by [wrapper.DumpBodyToFile], only the in-process temp copy
// is removed. The permanent file at FilePath is never touched.
func (d *Dump) Close() error {
	if d == nil {
		return nil
	}
	d.once.Do(func() {
		if d.syr != nil {
			d.closeErr = d.syr.Close()
		}
	})
	return d.closeErr
}
