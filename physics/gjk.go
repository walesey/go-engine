package physics

import (
	"fmt"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ConvexSet struct {
	verticies   []vmath.Vector3
	offset      vmath.Vector3
	orientation vmath.Quaternion
}

type Simplex struct {
}

// NewConvexSet
func NewConvexSet(verticies []vmath.Vector3) Collider {
	return &ConvexSet{
		verticies:   verticies,
		offset:      vmath.Vector3{0, 0, 0},
		orientation: vmath.IdentityQuaternion(),
	}
}

func (cs *ConvexSet) Offset(offset vmath.Vector3, orientation vmath.Quaternion) {
	cs.offset = offset
	cs.orientation = orientation
}

// Calculate the overlap of a convex set and another collider
func (cs *ConvexSet) Overlap(other Collider) bool {
	switch t := other.(type) {
	default:
		fmt.Printf("unsupported type for other collider: %T\n", t)
	case *ConvexSet:
		return cs.OverlapConvexSet(other.(*ConvexSet))
	}
	return false
}

// OverlapConvexSet Return true if the two convex sets overlap
func (cs *ConvexSet) OverlapConvexSet(other *ConvexSet) bool {

	return false
}

//CollisionInfo
func (cs *ConvexSet) CollisionInfo(other *ConvexSet) (norm vmath.Vector3, pen float64) {
	//TODO:
	return vmath.Vector3{}, 0.0
}

// Get the two closest points on two convex sets
func (cs *ConvexSet) ClosestPoints(other *ConvexSet) (p1, p2 vmath.Vector3) {
	//TODO:
	return vmath.Vector3{}, vmath.Vector3{}
}

func (cs *ConvexSet) support(other *ConvexSet, direction vmath.Vector3) vmath.Vector3 {
	p1 := cs.farthestPointInDirection(direction)
	p2 := other.farthestPointInDirection(direction.MultiplyScalar(-1))
	return p1.Subtract(p2)
}

func (cs *ConvexSet) farthestPointInDirection(direction vmath.Vector3) vmath.Vector3 {
	var farthest vmath.Vector3
	max := 0.0
	for i, point := range cs.verticies {
		p := cs.transformPoint(point)
		d := direction.Dot(p)
		if d >= max || i == 0 {
			max = d
			farthest = point
		}
	}
	return farthest
}

func (cs *ConvexSet) transformPoint(point vmath.Vector3) vmath.Vector3 {
	return cs.orientation.Apply(point).Add(cs.offset)
}
