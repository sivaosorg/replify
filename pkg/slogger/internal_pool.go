package slogger

import "sync"

var entryPool = sync.Pool{New: func() interface{} { return &Entry{Fields: make([]Field, 0, 8)} }}

// acquireEntry retrieves an Entry from the pool and sets its Logger field.
//
// Parameters:
//   - `l`: the logger that owns this entry
//
// Returns:
//
// a ready-to-use *Entry with Logger set.
func acquireEntry(l *Logger) *Entry {
	e := entryPool.Get().(*Entry)
	e.Logger = l
	return e
}

// releaseEntry resets e and returns it to the pool for reuse.
//
// Parameters:
//   - `e`: the entry to release
func releaseEntry(e *Entry) {
	e.reset()
	entryPool.Put(e)
}
