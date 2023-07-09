package server_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/evertras/rcc/pkg/server"
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

func (r *mockCoverageRepo) getValue1000(key string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if val, ok := r.vals1000[key]; ok {
		return val, nil
	}

	return 0, fmt.Errorf("key %q not found", key)
}

func (r *mockCoverageRepo) StoreValue1000(ctx context.Context, key string, value int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.vals1000[key] = value

	return nil
}

func (r *mockCoverageRepo) GetValue1000(ctx context.Context, key string) (int, error) {
	return r.getValue1000(key)
}

func TestServerStopsOnContextDone(t *testing.T) {
	s := server.New(newMockCoverageRepo())

	const expectedTimeout = time.Millisecond * 10

	timeout := time.NewTimer(expectedTimeout * 2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	listenResult := make(chan error)

	go func() {
		listenResult <- s.ListenAndServe(ctx)
	}()

	select {
	case err := <-listenResult:
		if !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Unexpected error returned from listen: %s", err.Error())
		}

	case <-timeout.C:
		t.Error("Reached timout before server finished")
	}
}
