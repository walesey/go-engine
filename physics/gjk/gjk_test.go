package gjk

import (
	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
	"testing"
)

func TestGJKOverlap(t *testing.T) {
	verts := []vmath.Vector3{
		vmath.Vector3{0, 0, 0},
		vmath.Vector3{1, 0, 0},
		vmath.Vector3{1, 1, 0},
		vmath.Vector3{0, 1, 0},
		vmath.Vector3{0, 0, 1},
		vmath.Vector3{1, 0, 1},
		vmath.Vector3{1, 1, 1},
		vmath.Vector3{0, 1, 1},
	}
	convexSet1 := NewConvexSet(verts)
	convexSet2 := NewConvexSet(verts)

	convexSet1.Offset(vmath.Vector3{0.5, 0.5, 0.5}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{0.5, 0.5, 0.5}, vmath.IdentityQuaternion())
	assert.True(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{0.5, 0.5, 0.5}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{1.51, 1.51, 1.51}, vmath.IdentityQuaternion())
	assert.False(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{0.5, 0.5, 0.5}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{1.41, 1.41, 1.41}, vmath.IdentityQuaternion())
	assert.True(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{-0.5, -0.5, -0.5}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{-1.41, -1.41, -1.41}, vmath.IdentityQuaternion())
	assert.True(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{-0.5, -0.5, -0.5}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{-1.51, -1.51, -1.51}, vmath.IdentityQuaternion())
	assert.False(t, convexSet1.Overlap(convexSet2))
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
