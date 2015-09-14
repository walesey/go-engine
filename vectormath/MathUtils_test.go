package vectormath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRound(t *testing.T) {
	n := 1.999
	i := Round(n, .5, 0)
	assert.EqualValues(t, 2.0, i, "LengthSquared")
}

func TestRoundHalfUp(t *testing.T) {
	n := 1.501
	i := RoundHalfUp(n)
	assert.EqualValues(t, 2, i, "LengthSquared")
}

func TestApproxEqual(t *testing.T) {
	assert.True(t, ApproxEqual(1.53121, 1.5, 0.1), "ApproxEqual")
	assert.False(t, ApproxEqual(1.53121, 1.5, 0.01), "ApproxEqual")
	assert.True(t, ApproxEqual(1.49121, 1.5, 0.01), "ApproxEqual")
}
