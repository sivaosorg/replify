package slogger

// Add registers hook for every level it reports.
//
// Parameters:
//   - `hook`: the Hook to register
func (h *Hooks) Add(hook Hook) {
	if h == nil {
		return
	}
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
	if h == nil {
		return nil
	}
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
	if h == nil {
		return 0
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.hooks[level])
}

// HooksFor returns a copy of the hooks registered for a specific level.
//
// Parameters:
//   - `level`: the level whose hooks should be returned
//
// Returns:
//
// a copy of the []Hook slice for that level.
func (h *Hooks) HooksFor(level Level) []Hook {
	if h == nil {
		return nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.hooks[level] == nil {
		return nil
	}
	result := make([]Hook, len(h.hooks[level]))
	copy(result, h.hooks[level])
	return result
}
