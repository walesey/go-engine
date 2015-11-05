package gjk

import (
	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
	"testing"
)

func TestSimplex(t *testing.T) {
	simplex := NewSimplex()
	simplex.Add(SimplexPoint{vmath.Vector3{1, 1, 1}})
	simplex.Add(SimplexPoint{vmath.Vector3{2, 2, 2}})
	simplex.Add(SimplexPoint{vmath.Vector3{3, 3, 3}})
	simplex.Add(SimplexPoint{vmath.Vector3{4, 4, 4}})
	assert.EqualValues(t, 4, simplex.Len(), "simplex should have 4 entries")

	simplex.Remove(2)
	assert.EqualValues(t, 3, simplex.Len(), "simplex should have 3 entries")
	assert.EqualValues(t, vmath.Vector3{2, 2, 2}, simplex.Get(1).mPoint, "incorrect value in simplex")
	assert.EqualValues(t, vmath.Vector3{4, 4, 4}, simplex.Get(2).mPoint, "incorrect value in simplex")

	simplex.Remove(2)
	assert.EqualValues(t, 2, simplex.Len(), "simplex should have 2 entries")
	assert.EqualValues(t, vmath.Vector3{2, 2, 2}, simplex.Get(1).mPoint, "incorrect value in simplex")
}
