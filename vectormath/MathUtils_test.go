package vectormath

import (
	"fmt"
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

func TestFastInvSqrt(t *testing.T) {
	assert.True(t, ApproxEqual(FastInvSqrt(4.0), 1.0/2.0, 0.01), fmt.Sprintf("FastInvSqrt of 4: expect:%v to equal:%v", FastInvSqrt(4.0), 1.0/2.0))
	assert.True(t, ApproxEqual(FastInvSqrt(89.0), 1.0/9.0, 0.01), fmt.Sprintf("FastInvSqrt of 89: expect:%v to equal:%v", FastInvSqrt(89.0), 1.0/9.0))
	assert.True(t, ApproxEqual(FastInvSqrt(732736.0), 1.0/856.0, 0.01), fmt.Sprintf("FastInvSqrt of 732736: expect:%v to equal:%v", FastInvSqrt(732736.0), 1.0/856.0))
}
