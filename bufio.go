package replify

// Get returns buffer from pool
//
// If pool is empty, it creates a new buffer
// of predefined size
//
// Returns:
//   - A byte slice buffer
func (bp *BufferPool) Get() []byte {
	select {
	case buf := <-bp.buffers:
		return buf
	default:
		return make([]byte, bp.size)
	}
}

// Put returns buffer to pool
//
// # If pool is full, the buffer is discarded
//
// Parameters:
//   - buf: A byte slice buffer to be returned to the pool
//
// Returns:
//   - None
func (bp *BufferPool) Put(buf []byte) {
	if buf == nil {
		return
	}
	select {
	case bp.buffers <- buf:
	default:
		// Pool is full, discard
	}
}
