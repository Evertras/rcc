package server_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/evertras/rcc/pkg/server"
)

func TestServerStopsOnContextDone(t *testing.T) {
	s := server.New()

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
