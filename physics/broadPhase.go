package physics

import (
	"fmt"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Collider interface {
	Overlap(other Collider) bool
	AttachTo(object *PhysicsObject)
}

type BoundingBox struct {
	bounds vmath.Vector3
	offset *vmath.Vector3
}

func NewBoundingBox(bounds vmath.Vector3) *BoundingBox {
	return &BoundingBox{
		bounds: vmath.Vector3{0, 0, 0},
	}
}

func (bb *BoundingBox) Overlap(other Collider) bool {
	switch t := other.(type) {
	default:
		fmt.Printf("unsupported type for other collider: %T\n", t)
	case *BoundingBox:
		return bb.OverlapBoundingBox(other.(*BoundingBox))
	}
	return false
}

//AttachTo - this bounding box is now bound to the object
func (bb *BoundingBox) AttachTo(object *PhysicsObject) {
	object.broadPhase = bb
	bb.offset = &object.Position
}

//OverlapBoundingBox - classic AABB overlap test
func (bb *BoundingBox) OverlapBoundingBox(other *BoundingBox) bool {
	if OneDimensionOverlap(bb.offset.X+bb.bounds.X/2, bb.offset.X-bb.bounds.X/2, other.offset.X+other.bounds.X/2, other.offset.X-other.bounds.X/2) {
		return false
	}
	if OneDimensionOverlap(bb.offset.Y+bb.bounds.Y/2, bb.offset.Y-bb.bounds.Y/2, other.offset.Y+other.bounds.Y/2, other.offset.Y-other.bounds.Y/2) {
		return false
	}
	if OneDimensionOverlap(bb.offset.Z+bb.bounds.Z/2, bb.offset.Z-bb.bounds.Z/2, other.offset.Z+other.bounds.Z/2, other.offset.Z-other.bounds.Z/2) {
		return false
	}
	return true
}

func OneDimensionOverlap(high1, low1, high2, low2 float64) bool {
	return (high1 < high2 && high1 < low2) || (low1 > low2 && low1 > high2)
}

//TODO: sphere base broadphase collisions
