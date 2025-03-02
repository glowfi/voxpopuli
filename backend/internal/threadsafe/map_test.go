package threadsafe_test

import (
	"testing"

	"github.com/glowfi/voxpopuli/backend/internal/threadsafe"
)

func TestThreadSafeMap(t *testing.T) {
	tests := []struct {
		name     string
		actions  func(m *threadsafe.ThreadSafeMap[string, int])
		expected map[string]int
	}{
		{
			name: "Put and Get",
			actions: func(m *threadsafe.ThreadSafeMap[string, int]) {
				m.Put("One", 1)
				m.Put("Two", 2)
			},
			expected: map[string]int{"One": 1, "Two": 2},
		},
		{
			name: "Remove",
			actions: func(m *threadsafe.ThreadSafeMap[string, int]) {
				m.Put("One", 1)
				m.Put("Two", 2)
				m.Remove("One")
			},
			expected: map[string]int{"Two": 2},
		},
		{
			name: "ContainsKey",
			actions: func(m *threadsafe.ThreadSafeMap[string, int]) {
				m.Put("One", 1)
				m.Put("Two", 2)
			},
			expected: map[string]int{"One": 1, "Two": 2},
		},
		{
			name: "Size",
			actions: func(m *threadsafe.ThreadSafeMap[string, int]) {
				m.Put("One", 1)
				m.Put("Two", 2)
				m.Remove("One")
			},
			expected: map[string]int{"Two": 2},
		},
		{
			name: "Clear",
			actions: func(m *threadsafe.ThreadSafeMap[string, int]) {
				m.Put("One", 1)
				m.Put("Two", 2)
				m.Clear()
			},
			expected: map[string]int{},
		},
		{
			name: "Load",
			actions: func(m *threadsafe.ThreadSafeMap[string, int]) {
				m.Put("One", 1)
				m.Put("Two", 2)
			},
			expected: map[string]int{"One": 1, "Two": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := threadsafe.NewThreadSafeMap[string, int]()
			tt.actions(m)

			// Verify the expected values
			for k, v := range tt.expected {
				if value, exists := m.Get(k); !exists || value != v {
					t.Errorf("Expected %v for key %v, got %v (exists: %v)", v, k, value, exists)
				}
			}

			// Check the size
			if size := m.Size(); size != len(tt.expected) {
				t.Errorf("Expected size %v, got %v", len(tt.expected), size)
			}

			// Check if the map is empty after clear
			if tt.name == "Clear" {
				if size := m.Size(); size != 0 {
					t.Errorf("Expected size 0 after clear, got %v", size)
				}
			}

			// Check if all keys are present in Load
			if tt.name == "Load" {
				loadedData := m.Load()
				for k, v := range tt.expected {
					if value, exists := loadedData[k]; !exists || value != v {
						t.Errorf("Load() expected %v for key %v, got %v (exists: %v)", v, k, value, exists)
					}
				}
			}
		})
	}
}
