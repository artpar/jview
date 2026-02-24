package darwin

import "sync"

// CallbackRegistry maps uint64 IDs to Go callback functions.
// Used by ObjC target-action to route events back to Go.
type CallbackRegistry struct {
	mu      sync.RWMutex
	next    uint64
	entries map[uint64]func(string)
}

var globalRegistry = &CallbackRegistry{
	next:    1,
	entries: make(map[uint64]func(string)),
}

// Register stores a callback and returns its ID.
func (r *CallbackRegistry) Register(fn func(string)) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.next
	r.next++
	r.entries[id] = fn
	return id
}

// Invoke calls the callback with the given ID.
func (r *CallbackRegistry) Invoke(id uint64, data string) {
	r.mu.RLock()
	fn, ok := r.entries[id]
	r.mu.RUnlock()
	if ok {
		fn(data)
	}
}

// Unregister removes a callback.
func (r *CallbackRegistry) Unregister(id uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, id)
}
