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
	bb1.AttachTo(&obj1)
	bb2.AttachTo(&obj2)
	assert.True(t, obj1.BroadPhaseOverlap(obj2), "BroadPhase")
	obj2.Position = vmath.Vector3{3, 0, 1}
	assert.False(t, obj1.BroadPhaseOverlap(obj2), "BroadPhase")
}

func TestBoundingBox(t *testing.T) {
	offset1 := vmath.Vector3{0, 0, 0}
	offset2 := vmath.Vector3{0, 0, 0}
	bb1 := NewBoundingBox(vmath.Vector3{1, 1, 1})
	bb2 := NewBoundingBox(vmath.Vector3{1, 1, 1})
	bb1.offset = &offset1
	bb2.offset = &offset2
	assert.True(t, bb1.Overlap(bb2), "Bounding box overlap")
	offset2 = vmath.Vector3{0, 0, 1.1}
	assert.False(t, bb1.Overlap(bb2), "Bounding box not overlap")
}
