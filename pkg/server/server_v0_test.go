package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evertras/rcc/pkg/badge"
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

			mockRepo := newMockCoverageRepo().withValue1000(test.key, test.storedValue1000)
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

	assert.Contains(t, body, "??%", "Should have ??% to denote unknown value in badge")
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

	mockRepo := newMockCoverageRepo().withValue1000(key, value1000)

	s := server.New(server.NewDefaultConfig(), mockRepo)

	s.Handler().ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Unexpected HTTP status code")

	body := w.Body.String()

	// TODO: Better test of SVG correctness, but this is pretty visual...
	assert.Contains(t, body, expectedText, "Missing expected percent coverage")
	assert.Contains(t, body, badge.ColorRed, "Color is wrong")
}

func TestV0BadgeCoverageReturnsSVGWithColorThresholds(t *testing.T) {
	const key = "github.com/Evertras/rcc"

	tests := []struct {
		name               string
		thresholdOrange100 int
		thresholdRed100    int
		actualValue100     int
		expectedColor      badge.Color
		expected400        bool
		label              string
	}{
		{
			name:           "GreenWithDefaults",
			actualValue100: 80,
			expectedColor:  badge.ColorGreen,
		},
		{
			name:           "OrangeWithDefaults",
			actualValue100: 79,
			expectedColor:  badge.ColorOrange,
		},
		{
			name:           "OrangeWithDefaults",
			actualValue100: 79,
			expectedColor:  badge.ColorOrange,
		},
		{
			name:               "OrangeWithOnlyOrangeAdjusted",
			thresholdOrange100: 82,
			actualValue100:     81,
			expectedColor:      badge.ColorOrange,
		},
		{
			name:            "RedOnlyRedAdjustedAboveOrangeThreshold",
			thresholdRed100: 82,
			actualValue100:  81,
			expectedColor:   badge.ColorRed,
		},
		{
			name:           "RedWithDefaults",
			actualValue100: 49,
			expectedColor:  badge.ColorRed,
			label:          "redredred",
		},
		{
			name:               "GreenWithHighThresholds",
			thresholdOrange100: 98,
			thresholdRed100:    97,
			actualValue100:     99,
			expectedColor:      badge.ColorGreen,
		},
		{
			name:               "GreenWithLowThresholds",
			thresholdOrange100: 9,
			thresholdRed100:    5,
			actualValue100:     10,
			expectedColor:      badge.ColorGreen,
		},
		{
			name:               "OrangeWithLowThresholds",
			thresholdOrange100: 9,
			thresholdRed100:    5,
			actualValue100:     8,
			expectedColor:      badge.ColorOrange,
		},
		{
			name:               "RedWithLowThresholds",
			thresholdOrange100: 9,
			thresholdRed100:    5,
			actualValue100:     3,
			expectedColor:      badge.ColorRed,
		},
		{
			name:               "RedWithHighThresholds",
			thresholdOrange100: 99,
			thresholdRed100:    98,
			actualValue100:     97,
			expectedColor:      badge.ColorRed,
		},
		{
			name:               "OrangeWithHighThresholds",
			thresholdOrange100: 99,
			thresholdRed100:    98,
			actualValue100:     98,
			expectedColor:      badge.ColorOrange,
		},
		{
			name:               "OrangeTooHigh",
			thresholdOrange100: 101,
			actualValue100:     98,
			expected400:        true,
			label:              "orange2high",
		},
		{
			name:               "OrangeTooLow",
			thresholdOrange100: -1,
			actualValue100:     98,
			expected400:        true,
		},
		{
			name:            "RedTooHigh",
			thresholdRed100: 101,
			expected400:     true,
		},
		{
			name:            "RedTooLow",
			thresholdRed100: -1,
			expected400:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedLabel := "coverage"

			path := fmt.Sprintf("/api/v0/badge/coverage?key=%s", key)

			if test.thresholdOrange100 != 0 {
				path += fmt.Sprintf("&thresholdOrange100=%d%%25", test.thresholdOrange100)
			}

			if test.thresholdRed100 != 0 {
				path += fmt.Sprintf("&thresholdRed100=%d%%25", test.thresholdRed100)
			}

			if test.label != "" {
				path += fmt.Sprintf("&label=%s", test.label)
				expectedLabel = test.label
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", path, nil)

			mockRepo := newMockCoverageRepo().withValue1000(key, test.actualValue100*10)

			s := server.New(server.NewDefaultConfig(), mockRepo)

			s.Handler().ServeHTTP(w, r)

			if test.expected400 {
				assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 failure")
				return
			}

			assert.Equal(t, http.StatusOK, w.Code, "Unexpected HTTP status code")

			body := w.Body.String()

			assert.Contains(t, body, test.expectedColor, "Missing expected color")
			assert.Contains(t, body, expectedLabel, "Missing expected label")
		})
	}
}
