package vectorMath


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
