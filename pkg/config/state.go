package config

import (
	"fmt"
	"sync"
	
	"github.com/jacobdavidalcock/consolekit/pkg/utils"
)

// State manages global application state with thread safety
type State struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

// NewState creates a new state manager
func NewState() *State {
	return &State{
		data: make(map[string]interface{}),
	}
}

// Set sets a state value (thread-safe)
func (s *State) Set(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = value
}

// Get gets a state value (thread-safe)
func (s *State) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, exists := s.data[key]
	return value, exists
}

// GetString gets a string state value (thread-safe)
func (s *State) GetString(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, exists := s.data[key]
	if !exists {
		return "", false
	}
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// GetInt gets an integer state value (thread-safe)
func (s *State) GetInt(key string) (int, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, exists := s.data[key]
	if !exists {
		return 0, false
	}
	if i, ok := value.(int); ok {
		return i, true
	}
	return 0, false
}

// GetBool gets a boolean state value (thread-safe)
func (s *State) GetBool(key string) (bool, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, exists := s.data[key]
	if !exists {
		return false, false
	}
	if b, ok := value.(bool); ok {
		return b, true
	}
	return false, false
}

// Delete removes a state value (thread-safe)
func (s *State) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, key)
}

// Clear removes all state values (thread-safe)
func (s *State) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = make(map[string]interface{})
}

// Keys returns all state keys (thread-safe)
func (s *State) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var keys []string
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

// ShowAll displays all state values (thread-safe)
func (s *State) ShowAll() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	fmt.Println("\n--- Current State ---")
	for key, value := range s.data {
		// Mask sensitive values
		if isSensitive(key) {
			if str, ok := value.(string); ok {
				fmt.Printf("  %-15s : %s\n", key, utils.MaskString(str, 4, 4))
				continue
			}
		}
		fmt.Printf("  %-15s : %v\n", key, value)
	}
	fmt.Println("--------------------")
}

// isSensitive checks if a key contains sensitive information
func isSensitive(key string) bool {
	sensitiveKeys := []string{
		"password", "token", "secret", "key", "credential",
		"apikey", "api_key", "auth", "jwt", "bearer",
	}
	
	for _, sensitive := range sensitiveKeys {
		if key == sensitive {
			return true
		}
	}
	return false
}