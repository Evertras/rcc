package server_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/evertras/rcc/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestServerHealthz(t *testing.T) {
	s := server.New(server.NewDefaultConfig(), newMockCoverageRepo())

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/healthz", nil)

	s.Handler().ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Health check should return 200")
}

func TestServerStopsOnContextDone(t *testing.T) {
	s := server.New(server.NewDefaultConfig(), newMockCoverageRepo())

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
