package vectormath

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestPointToLineDist(t *testing.T) {
	assert.EqualValues(t, 1, PointToLineDist(Vector3{0, 0, 0}, Vector3{1, 0, 1}, Vector3{0, 1, 0}))
}
