package repository_test

import (
	"context"
	"strings"
	"testing"

	"github.com/evertras/rcc/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryReturnsNotFoundCorrectly(t *testing.T) {
	r := repository.NewInMemory()

	_, err := r.GetValue1000(context.Background(), "nope")

	assert.Error(t, err, "Expected an error")

	assert.True(t, strings.Contains(err.Error(), "not found"), "Error should contain 'not found'")
}

func TestInMemoryStoresAndRetrieves(t *testing.T) {
	r := repository.NewInMemory()

	const (
		key   = "github.com/Evertras/rcc"
		value = 1000
	)

	err := r.StoreValue1000(context.Background(), key, value)

	assert.NoError(t, err, "Failed to store value")

	returnedValue, err := r.GetValue1000(context.Background(), key)

	assert.NoError(t, err, "Failed to get value")
	assert.Equal(t, value, returnedValue, "Wrong value returned")
}
