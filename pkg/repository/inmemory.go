package repository

import (
	"context"
	"fmt"
	"sync"
)

type InMemory struct {
	mu sync.RWMutex

	vals1000 map[string]int
}

func NewInMemory() *InMemory {
	return &InMemory{
		vals1000: make(map[string]int),
	}
}
func (r *InMemory) StoreValue1000(ctx context.Context, key string, value int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.vals1000[key] = value

	return nil
}

func (r *InMemory) GetValue1000(ctx context.Context, key string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if val, ok := r.vals1000[key]; ok {
		return val, nil
	}

	return 0, fmt.Errorf("key %q not found", key)
}
