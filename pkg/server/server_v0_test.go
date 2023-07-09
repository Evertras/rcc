package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evertras/rcc/pkg/server"
)

func TestV0PutWithNoKeyReturns400(t *testing.T) {
	const route = "/api/v0/coverage"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", route, nil)

	s := server.New(newMockCoverageRepo())

	s.Handle(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Unexpected HTTP response code")
}

func TestV0PutWithNoValueReturns400(t *testing.T) {
	const route = "/api/v0/coverage?key=abc"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", route, nil)

	s := server.New(newMockCoverageRepo())

	s.Handle(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Unexpected HTTP response code")
}

func TestV0Put(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		valueRaw string

		expectedStore1000 int
	}{
		{
			name:     "Simple",
			key:      "abc",
			valueRaw: "10",

			expectedStore1000: 100,
		},
		{
			name:     "GithubRepo",
			key:      "github.com/Evertras/rcc",
			valueRaw: "90.7",

			expectedStore1000: 907,
		},
		{
			name:     "NoCoverageInt",
			key:      "idk",
			valueRaw: "0",

			expectedStore1000: 0,
		},
		{
			name:     "NoCoverageFloatPercent",
			key:      "idk",
			valueRaw: "0.00000%",

			expectedStore1000: 0,
		},
		{
			name:     "FullCoverageFloatPercent",
			key:      "idk",
			valueRaw: "100.00%",

			expectedStore1000: 1000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			route := fmt.Sprintf("/api/v0/coverage?key=%s&value100=%s", url.QueryEscape(test.key), url.QueryEscape(test.valueRaw))

			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", route, nil)
			mockRepo := newMockCoverageRepo()

			s := server.New(mockRepo)

			s.Handle(w, r)

			assert.Equal(t, http.StatusOK, w.Code, "Unexpected HTTP response code")

			val, err := mockRepo.getValue1000(test.key)

			assert.NoError(t, err, "Failed to check value that should exist")

			if t.Failed() {
				t.FailNow()
			}

			assert.Equal(t, val, test.expectedStore1000, "Unexpected value stored")
		})
	}
}
