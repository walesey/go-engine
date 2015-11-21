package gjk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
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

	convexSet1.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	assert.True(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{1.01, 1.01, 1.01}, vmath.IdentityQuaternion())
	assert.False(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{0.91, 0.91, 0.91}, vmath.IdentityQuaternion())
	assert.True(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{-0.91, -0.91, -0.91}, vmath.IdentityQuaternion())
	assert.True(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{-1.01, -1.01, -1.01}, vmath.IdentityQuaternion())
	assert.False(t, convexSet1.Overlap(convexSet2))
}

func TestGJKOverlap_rotation(t *testing.T) {
	verts := []vmath.Vector3{
		vmath.Vector3{-0.1, 0, -1},
		vmath.Vector3{0.1, 0, -1},
		vmath.Vector3{0.1, 0.2, -1},
		vmath.Vector3{-0.1, 0.2, -1},
		vmath.Vector3{-0.1, 0, 1},
		vmath.Vector3{0.1, 0, 1},
		vmath.Vector3{0.1, 0.2, 1},
		vmath.Vector3{-0.1, 0.2, 1},
	}
	convexSet1 := NewConvexSet(verts)
	convexSet2 := NewConvexSet(verts)

	convexSet1.Offset(vmath.Vector3{0, 0.3, 0}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	assert.False(t, convexSet1.Overlap(convexSet2))

	convexSet1.Offset(vmath.Vector3{0, 0.3, 0}, vmath.AngleAxis(-0.3, vmath.Vector3{1, 0, 0}))
	convexSet2.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	assert.True(t, convexSet1.Overlap(convexSet2))
}

func TestGJKContact(t *testing.T) {
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

	convexSet1.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{1.01, 1.01, 1.01}, vmath.IdentityQuaternion())
	contact := convexSet1.ContactPoint(convexSet2)
	assert.True(t, contact.ApproxEqual(vmath.Vector3{1, 1, 1}, 0.01), "contact point incorrect")

	convexSet1.Offset(vmath.Vector3{-1, -1, -1}, vmath.IdentityQuaternion())
	convexSet2.Offset(vmath.Vector3{-2.01, -2.01, -2.01}, vmath.IdentityQuaternion())
	contact = convexSet1.ContactPoint(convexSet2)
	assert.True(t, contact.ApproxEqual(vmath.Vector3{-1, -1, -1}, 0.01), "contact point incorrect")
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

func TestClosestPointOnTriangleToOrigin(t *testing.T) {
	assert.EqualValues(t, vmath.Vector3{0, 0, 0}, closestPointOnTriangleToOrigin(vmath.Vector3{1, 1, 1}, vmath.Vector3{1, 0, 1}, vmath.Vector3{0, 0, 0}), "incorrect closest point to plane")
	assert.EqualValues(t, vmath.Vector3{-1, -1, -2}, closestPointOnTriangleToOrigin(vmath.Vector3{-4, 0, -4}, vmath.Vector3{0, -4, -4}, vmath.Vector3{-1, -1, -2}), "closest point to triangle should lie within the triangle")
}
