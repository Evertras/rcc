package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evertras/rcc/pkg/server"
)

func TestV0PutWithNoKeyReturns400(t *testing.T) {
	const route = "/api/v0/coverage/value100"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", route, nil)

	s := server.New()

	s.Handle(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Unexpected HTTP response code")
}

func TestV0PutWithSimpleKeyWorks(t *testing.T) {
	const key = "simple"

	const route = "/api/v0/coverage/value100?key=" + key

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", route, nil)

	s := server.New()

	s.Handle(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Unexpected HTTP response code")
}

func TestV0PutWithGithubRepoKeyWorks(t *testing.T) {
	const key = "github.com/Evertras/rcc"

	const route = "/api/v0/coverage/value100?key=" + key

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", route, nil)

	s := server.New()

	s.Handle(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Unexpected HTTP response code")
}
