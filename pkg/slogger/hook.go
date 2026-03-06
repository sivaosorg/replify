package slogger

import "sync"

// Hook fires side-effects (e.g. alerting, metrics) for matching log levels.
// Implementations must be safe for concurrent use.
type Hook interface {
	// Levels returns the set of log levels this hook handles.
	Levels() []Level
	// Fire is called for each matching entry; it must not retain the entry.
	Fire(*Entry) error
}

// Hooks is a level-indexed registry of Hook instances.
type Hooks struct {
	mu    sync.RWMutex
	hooks map[Level][]Hook
}

// NewHooks returns an empty Hooks registry.
//
// Returns:
//
// a ready-to-use *Hooks.
func NewHooks() *Hooks {
	h := &Hooks{
		hooks: make(map[Level][]Hook, 7),
	}
	return h
}

// Add registers hook for every level it reports.
//
// Parameters:
//   - `hook`: the Hook to register
func (h *Hooks) Add(hook Hook) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, lvl := range hook.Levels() {
		h.hooks[lvl] = append(h.hooks[lvl], hook)
	}
}

// Fire calls every hook registered for level, passing entry to each.
// All hooks are called even if earlier ones return errors; errors are
// aggregated into a single combined error.
//
// Parameters:
//   - `level`: the level whose hooks should fire
//   - `entry`: the log entry to pass to each hook
//
// Returns:
//
// a combined error from all hooks, or nil if all succeeded.
func (h *Hooks) Fire(level Level, entry *Entry) error {
	h.mu.RLock()
	hooks := h.hooks[level]
	h.mu.RUnlock()

	var first error
	for _, hook := range hooks {
		if err := hook.Fire(entry); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// Len returns the number of hooks registered for level.
//
// Parameters:
//   - `level`: the level to query
//
// Returns:
//
// the count of hooks registered for that level.
func (h *Hooks) Len(level Level) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.hooks[level])
}
