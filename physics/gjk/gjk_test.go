package gjk

import (
	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
	"testing"
)

func TestGJKOverlap(t *testing.T) {
}

func TestGJKDistance(t *testing.T) {
}

func TestEPACollisionInfo(t *testing.T) {
}

func TestFarthestPointInDirection(t *testing.T) {
	points := []vmath.Vector3{
		vmath.Vector3{2, 2, 0},
		vmath.Vector3{1, 2, 0},
		vmath.Vector3{1, 1, 0},
		vmath.Vector3{2, 1, 0},
	}
	convexSet, ok := NewConvexSet(points).(*ConvexSet)
	assert.True(t, ok, "NewConvexSet should create a *ConvexSet")
	assert.EqualValues(t, vmath.Vector3{2, 2, 0}, convexSet.farthestPointInDirection(vmath.Vector3{1, 1, 0}), "Farthest point incorrect")
	assert.EqualValues(t, vmath.Vector3{1, 2, 0}, convexSet.farthestPointInDirection(vmath.Vector3{-1, 1, 0}), "Farthest point incorrect")
	assert.EqualValues(t, vmath.Vector3{1, 1, 0}, convexSet.farthestPointInDirection(vmath.Vector3{-1, -1, 0}), "Farthest point incorrect")
	assert.EqualValues(t, vmath.Vector3{2, 1, 0}, convexSet.farthestPointInDirection(vmath.Vector3{1, -1, 0}), "Farthest point incorrect")
}
