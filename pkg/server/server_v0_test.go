package server_test

import (
	"context"
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

	s := server.New(server.NewDefaultConfig(), newMockCoverageRepo())

	s.Handler().ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Unexpected HTTP response code")
}

func TestV0PutWithNoValueReturns400(t *testing.T) {
	const route = "/api/v0/coverage?key=abc"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", route, nil)

	s := server.New(server.NewDefaultConfig(), newMockCoverageRepo())

	s.Handler().ServeHTTP(w, r)

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

			s := server.New(server.NewDefaultConfig(), mockRepo)

			s.Handler().ServeHTTP(w, r)

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

func TestV0GetNoKey(t *testing.T) {
	const route = "/api/v0/coverage"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", route, nil)

	s := server.New(server.NewDefaultConfig(), newMockCoverageRepo())

	s.Handler().ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestV0Get(t *testing.T) {
	tests := []struct {
		name            string
		key             string
		storedValue1000 int
		expectedReturn  string
	}{
		{
			name:            "Simple",
			key:             "abc",
			storedValue1000: 123,
			expectedReturn:  "12.3",
		},
		{
			name:            "GithubRepo",
			key:             "github.com/Evertras/rcc",
			storedValue1000: 100,
			expectedReturn:  "10.0",
		},
		{
			name:            "ZeroCoverage",
			key:             "idk",
			storedValue1000: 0,
			expectedReturn:  "0.0",
		},
		{
			name:            "FullCoverage",
			key:             "idk",
			storedValue1000: 1000,
			expectedReturn:  "100.0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			route := "/api/v0/coverage?key=" + url.QueryEscape(test.key)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", route, nil)

			mockRepo := newMockCoverageRepo()
			err := mockRepo.StoreValue1000(context.Background(), test.key, test.storedValue1000)
			assert.NoError(t, err, "Failed to set the initial mock value, bad test setup")
			s := server.New(server.NewDefaultConfig(), mockRepo)

			s.Handler().ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code, "Unexpected HTTP response code")

			body := w.Body.String()

			assert.Equal(t, test.expectedReturn, body, "Unexpected body value")
		})
	}
}

func TestV0BadgeCoverageNoKey(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v0/badge/coverage", nil)

	s := server.New(server.NewDefaultConfig(), newMockCoverageRepo())

	s.Handler().ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Unexpected HTTP status code")
}

func TestV0BadgeCoverageNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v0/badge/coverage?key=idk", nil)

	s := server.New(server.NewDefaultConfig(), newMockCoverageRepo())

	s.Handler().ServeHTTP(w, r)

	// Should return 200 OK but the badge itself should indicate not found
	assert.Equal(t, http.StatusOK, w.Code, "Unexpected HTTP status code")

	body := w.Body.String()

	assert.Contains(t, body, "??.?%", "Should have ?? to denote unknown value in badge")
}

func TestV0BadgeCoverageReturnsSVG(t *testing.T) {
	const (
		key       = "github.com/Evertras/rcc"
		value1000 = 489

		// Round down for now
		expectedText = "48%"
	)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v0/badge/coverage?key="+key, nil)

	mockRepo := newMockCoverageRepo()
	err := mockRepo.StoreValue1000(context.Background(), key, value1000)

	assert.NoError(t, err, "Unexpected error setting test value, bad test setup")

	s := server.New(server.NewDefaultConfig(), mockRepo)

	s.Handler().ServeHTTP(w, r)

	// Should return 200 OK but the badge itself should indicate not found
	assert.Equal(t, http.StatusOK, w.Code, "Unexpected HTTP status code")

	body := w.Body.String()

	// TODO: Better test of SVG correctness
	assert.Contains(t, body, expectedText, "Missing expected percent coverage")
}
