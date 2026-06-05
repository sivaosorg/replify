package sysx

import (
	"bytes"
	"io"
)

// Read implements io.Reader.
func (s *spillBuffer) Read(p []byte) (int, error) {
	return s.reader.Read(p)
}

// Seek implements io.Seeker.
func (s *spillBuffer) Seek(offset int64, whence int) (int64, error) {
	return s.reader.Seek(offset, whence)
}

// Close releases the memory buffer and closes the spill file when one
// was created.
//
// Returns:
//
// Any error returned by the spill file's Close, or nil otherwise.
func (s *spillBuffer) Close() error {
	s.mem = nil
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}

// drain copies src into the spill buffer, switching to disk as soon as
// the in-memory portion would exceed threshold.
func (s *spillBuffer) drain(src io.Reader, tempDir, tempPattern string, remove bool) error {
	const chunk = 32 << 10
	buf := make([]byte, chunk)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if wErr := s.write(buf[:n], tempDir, tempPattern, remove); wErr != nil {
				return wErr
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// write appends p to the spill buffer, transparently spilling to disk
// when the in-memory portion would exceed the threshold.
func (s *spillBuffer) write(p []byte, tempDir, tempPattern string, remove bool) error {
	if s.file != nil {
		n, err := s.file.Write(p)
		s.size += int64(n)
		return err
	}
	if int64(s.mem.Len())+int64(len(p)) <= s.threshold {
		n, _ := s.mem.Write(p)
		s.size += int64(n)
		return nil
	}
	tf, err := newTempFileWith(tempDir, tempPattern, remove)
	if err != nil {
		return err
	}
	if _, err := tf.Write(s.mem.Bytes()); err != nil {
		_ = tf.Close()
		return err
	}
	s.mem.Reset()
	s.file = tf
	n, err := s.file.Write(p)
	s.size += int64(n)
	return err
}

// seal positions the underlying reader at offset 0 and locks in size.
func (s *spillBuffer) seal() error {
	if s.file != nil {
		if _, err := s.file.Seek(0, io.SeekStart); err != nil {
			return err
		}
		s.reader = s.file
		return nil
	}
	s.reader = bytes.NewReader(s.mem.Bytes())
	return nil
}

// compile-time interface assertion.
var _ ReadSeekCloser = (*spillBuffer)(nil)
