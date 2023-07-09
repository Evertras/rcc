package badge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evertras/rcc/pkg/badge"
)

func TestGenerateCoverageSVG(t *testing.T) {
	const (
		coverageValue1000 = 948

		// Round down
		expectedCoverageAmount = "94%"
	)

	svg, err := badge.GenerateCoverageSVG(coverageValue1000, badge.ColorGreen)

	assert.NoError(t, err, "Unexpected error in generation")
	assert.Contains(t, svg, expectedCoverageAmount, "Didn't find coverage percent text")
}
