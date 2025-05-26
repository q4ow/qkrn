package store

import (
	"testing"

	"github.com/q4ow/qkrn/pkg/types"
)

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore()

	err := store.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	value, err := store.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%s'", value)
	}

	_, err = store.Get("nonexistent")
	if err != types.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}

	err = store.Delete("key1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get("key1")
	if err != types.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound after delete, got %v", err)
	}

	store.Set("a", "1")
	store.Set("b", "2")
	store.Set("c", "3")

	keys := store.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	if store.Size() != 3 {
		t.Errorf("Expected size 3, got %d", store.Size())
	}

	err = store.Set("", "value")
	if err != types.ErrEmptyKey {
		t.Errorf("Expected ErrEmptyKey, got %v", err)
	}
}
