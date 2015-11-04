package physics

import (
	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
	"testing"
)

func TestBoundingBoxWithPhysicsObject(t *testing.T) {
	obj1 := NewPhysicsObject()
	obj2 := NewPhysicsObject()
	bb1 := NewBoundingBox(vmath.Vector3{1, 1, 1})
	bb2 := NewBoundingBox(vmath.Vector3{1, 1, 1})
	obj1.BroadPhase = bb1
	obj2.BroadPhase = bb2
	assert.True(t, obj1.BroadPhaseOverlap(obj2), "BroadPhase")
	obj2.Position = vmath.Vector3{3, 0, 1}
	assert.False(t, obj1.BroadPhaseOverlap(obj2), "BroadPhase")
}

func TestBoundingBox(t *testing.T) {
	bb1 := NewBoundingBox(vmath.Vector3{1, 1, 1})
	bb2 := NewBoundingBox(vmath.Vector3{1, 1, 1})
	bb1.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	bb2.Offset(vmath.Vector3{0, 0, 0}, vmath.IdentityQuaternion())
	assert.True(t, bb1.Overlap(bb2), "Bounding box overlap")
	bb1.Offset(vmath.Vector3{0, 0, 1.1}, vmath.IdentityQuaternion())
	assert.False(t, bb1.Overlap(bb2), "Bounding box not overlap")
}
