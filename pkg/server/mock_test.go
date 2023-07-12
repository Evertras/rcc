package server_test

import (
	"context"
	"fmt"
	"sync"
)

type mockCoverageRepo struct {
	mu sync.RWMutex

	vals1000 map[string]int
}

func newMockCoverageRepo() *mockCoverageRepo {
	return &mockCoverageRepo{
		vals1000: make(map[string]int),
	}
}

func (r *mockCoverageRepo) withValue1000(key string, value1000 int) *mockCoverageRepo {
	err := r.StoreValue1000(context.Background(), key, value1000)

	if err != nil {
		panic(fmt.Sprintf("Somehow failed to store, bad mock code: %s", err.Error()))
	}

	return r
}

func (r *mockCoverageRepo) getValue1000(key string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if val, ok := r.vals1000[key]; ok {
		return val, nil
	}

	return 0, fmt.Errorf("key %q not found", key)
}

func (r *mockCoverageRepo) StoreValue1000(ctx context.Context, key string, value1000 int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.vals1000[key] = value1000

	return nil
}

func (r *mockCoverageRepo) GetValue1000(ctx context.Context, key string) (int, error) {
	return r.getValue1000(key)
}
