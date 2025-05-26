package store

import (
	"sync"

	"github.com/q4ow/qkrn/pkg/types"
)

type MemoryStore struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (s *MemoryStore) Get(key string) (string, error) {
	if key == "" {
		return "", types.ErrEmptyKey
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	if !exists {
		return "", types.ErrKeyNotFound
	}

	return value, nil
}

func (s *MemoryStore) Set(key, value string) error {
	if key == "" {
		return types.ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return nil
}

func (s *MemoryStore) Delete(key string) error {
	if key == "" {
		return types.ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return types.ErrKeyNotFound
	}

	delete(s.data, key)
	return nil
}

func (s *MemoryStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}

func (s *MemoryStore) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
