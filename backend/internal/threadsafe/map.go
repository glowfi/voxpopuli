package threadsafe

import "sync"

// ThreadSafeMap is a thread-safe map implementation using generics.
type ThreadSafeMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewThreadSafeMap creates a new instance of ThreadSafeMap.
func NewThreadSafeMap[K comparable, V any]() *ThreadSafeMap[K, V] {
	return &ThreadSafeMap[K, V]{data: make(map[K]V)}
}

// Put adds a key-value pair to the map.
func (m *ThreadSafeMap[K, V]) Put(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get retrieves a value by key.
func (m *ThreadSafeMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.data[key]
	return value, exists
}

// Remove deletes a key-value pair from the map.
func (m *ThreadSafeMap[K, V]) Remove(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// ContainsKey checks if the map contains a key.
func (m *ThreadSafeMap[K, V]) ContainsKey(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.data[key]
	return exists
}

// Size returns the number of key-value pairs in the map.
func (m *ThreadSafeMap[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Clear removes all key-value pairs from the map.
func (m *ThreadSafeMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[K]V)
}

// Load returns a copy of the entire map.
func (m *ThreadSafeMap[K, V]) Load() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Create a new map to return a copy of the data
	copy := make(map[K]V, len(m.data))
	for k, v := range m.data {
		copy[k] = v
	}
	return copy
}
