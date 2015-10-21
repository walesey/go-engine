package physics

import (
	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
	"testing"
)

func TestTriangleOverlap(t *testing.T) {
	t1 := Triangle{vmath.Vector3{-0.504109, 0.540383, 1.441891}, vmath.Vector3{1.004138, -0.547662, 0.706127}, vmath.Vector3{0.499970, 0.007280, -1.148018}}
	t2 := Triangle{vmath.Vector3{0.615065, -0.037162, -1.208676}, vmath.Vector3{-0.355948, 0.614617, 0.413767}, vmath.Vector3{0.892249, -0.426572, 1.579071}}
	assert.EqualValues(t, false, t1.Overlap(t2), "Triangle Overlap should be false")

	t1 = Triangle{vmath.Vector3{-1, 0, 1}, vmath.Vector3{1, 0, 1}, vmath.Vector3{1, 0, -1}}
	t2 = Triangle{vmath.Vector3{-1, 0.3, 1}, vmath.Vector3{1, 0.3, 1}, vmath.Vector3{1, 0.3, -1}}
	assert.EqualValues(t, false, t1.Overlap(t2), "Triangle Overlap should be false (Parallel triangles)")

	t1 = Triangle{vmath.Vector3{9, 0, 1}, vmath.Vector3{10, 0, 1}, vmath.Vector3{10, 0, -1}}
	t2 = Triangle{vmath.Vector3{-1, 0, 1}, vmath.Vector3{1, 0, 1}, vmath.Vector3{1, 0, -1}}
	assert.EqualValues(t, false, t1.Overlap(t2), "Triangle Overlap should be false (coplanar triangles)")

	t1 = Triangle{vmath.Vector3{-0.504109, 0.540383, 1.441891}, vmath.Vector3{1.004138, -0.547662, 0.706127}, vmath.Vector3{0.499970, 0.007280, -1.148018}}
	t2 = Triangle{vmath.Vector3{0.048924, -1.436378, 0.544679}, vmath.Vector3{-0.679617, 0.418397, 0.715096}, vmath.Vector3{0.775831, 1.099488, -0.475613}}
	assert.EqualValues(t, true, t1.Overlap(t2), "Triangle Overlap should be true")

	t1 = Triangle{vmath.Vector3{-0.504109, 0.540383, 1.441891}, vmath.Vector3{1.004138, -0.547662, 0.706127}, vmath.Vector3{0.499970, 0.007280, -1.148018}}
	t2 = Triangle{vmath.Vector3{0.081168, 0.966970, 1.599979}, vmath.Vector3{0.320969, -0.553781, 0.323337}, vmath.Vector3{0.749229, 0.741485, -1.139154}}
	assert.EqualValues(t, true, t1.Overlap(t2), "Triangle Overlap should be true")
}

func TestConvexHull(t *testing.T) {
	t1 := Triangle{vmath.Vector3{0, 0, 0}, vmath.Vector3{0, 1, 0}, vmath.Vector3{0, 0, 1}}
	t2 := Triangle{vmath.Vector3{0, 0, 0}, vmath.Vector3{0, 1, 0}, vmath.Vector3{1, 0, 0}}
	t3 := Triangle{vmath.Vector3{0, 1, 0}, vmath.Vector3{1, 0, 0}, vmath.Vector3{0, 0, 1}}
	triangles1 := []Triangle{t1, t2, t3}
	convexHull1 := NewConvexHull(triangles1)
	obj1 := NewPhysicsObject()
	convexHull1.AttachTo(&obj1)

	t4 := Triangle{vmath.Vector3{0, 0, 0}, vmath.Vector3{0, 1, 0}, vmath.Vector3{0, 0, -1}}
	t5 := Triangle{vmath.Vector3{0, 0, 0}, vmath.Vector3{0, 1, 0}, vmath.Vector3{-1, 0, 0}}
	t6 := Triangle{vmath.Vector3{0, 1, 0}, vmath.Vector3{-1, 0, 0}, vmath.Vector3{0, 0, -1}}
	triangles2 := []Triangle{t4, t5, t6}
	convexHull2 := NewConvexHull(triangles2)
	obj2 := NewPhysicsObject()
	convexHull2.AttachTo(&obj2)

	obj2.Position = vmath.Vector3{-0.3, 0, -0.3}
	assert.False(t, obj1.NarrowPhaseOverlap(obj2), "ConvexHull Overlap should be false")

	obj2.Orientation = vmath.AngleAxis(3.14, vmath.Vector3{0, 1, 0})
	assert.True(t, obj1.NarrowPhaseOverlap(obj2), "ConvexHull Overlap should be true")
}
