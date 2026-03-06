package slogger

// acquireEntry retrieves an Entry from the pool and sets its logger field.
//
// Parameters:
//   - `l`: the logger that owns this entry
//
// Returns:
//
// a ready-to-use *Entry with logger set.
func acquireEntry(l *Logger) *Entry {
e := entryPool.Get().(*Entry)
e.logger = l
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
