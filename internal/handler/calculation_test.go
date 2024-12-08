package handler

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizeInt(t *testing.T) {
	var size float64 = 92.545454
	sizeInt := int(size)
	assert.Equal(t, 92, sizeInt)
	roundedNum := math.Round(size)
	assert.Equal(t, float64(93), roundedNum)

	size_2 := 92.50
	assert.Equal(t, float64(93), math.Round(size_2))

	size_3 := 94.49
	assert.Equal(t, float64(94), math.Round(size_3))
}
