package gjk

import (
	"fmt"
	"github.com/walesey/go-engine/physics"
	vmath "github.com/walesey/go-engine/vectormath"
)

const epaMaxIterations = 40

type ConvexSet struct {
	verticies   []vmath.Vector3
	offset      vmath.Vector3
	orientation vmath.Quaternion
	simplex     *Simplex
	searchDir   vmath.Vector3
}

// NewConvexSet
func NewConvexSet(verticies []vmath.Vector3) physics.Collider {
	return &ConvexSet{
		verticies:   verticies,
		offset:      vmath.Vector3{0, 0, 0},
		orientation: vmath.IdentityQuaternion(),
		simplex:     NewSimplex(),
		searchDir:   vmath.Vector3{0, 1, 0},
	}
}

func (cs *ConvexSet) Offset(offset vmath.Vector3, orientation vmath.Quaternion) {
	cs.offset = offset
	cs.orientation = orientation
}

// Calculate the overlap of a convex set and another collider
func (cs *ConvexSet) Overlap(other physics.Collider) bool {
	switch t := other.(type) {
	default:
		fmt.Printf("unsupported type for other collider: %T\n", t)
	case *ConvexSet:
		return cs.OverlapConvexSet(other.(*ConvexSet))
	}
	return false
}

//PenetrationVector - get the vector of penetration between the two colliders
func (cs *ConvexSet) PenetrationVector(other physics.Collider) vmath.Vector3 {
	switch t := other.(type) {
	default:
		fmt.Printf("unsupported type for other collider: %T\n", t)
	case *ConvexSet:
		return cs.CollisionInfo(other.(*ConvexSet))
	}
	return vmath.Vector3{0, 0, 0}
}

// OverlapConvexSet Return true if the two convex sets overlap
func (cs *ConvexSet) OverlapConvexSet(other *ConvexSet) bool {
	cs.simplex.Clear()
	d := cs.searchDir
	cs.simplex.Add(SimplexPoint{mPoint: cs.support(other, d)})
	d = d.MultiplyScalar(-1)
	for true {
		cs.simplex.Add(SimplexPoint{mPoint: cs.support(other, d)})
		if cs.simplex.GetLast().mPoint.Dot(d) <= 0 {
			return false
		} else {
			if cs.simplex.ContainsOrigin(&d) {
				return true
			}
		}
	}
	return false
}

//CollisionInfo returns the penetration vector
func (cs *ConvexSet) CollisionInfo(other *ConvexSet) vmath.Vector3 {
	if cs.simplex.Len() < 4 {
		//build new 4 point simplex
		cs.simplex.Clear()
		cs.simplex.Add(SimplexPoint{mPoint: cs.support(other, vmath.Vector3{1, 0, 0})})
		cs.simplex.Add(SimplexPoint{mPoint: cs.support(other, vmath.Vector3{0, 1, 0})})
		cs.simplex.Add(SimplexPoint{mPoint: cs.support(other, vmath.Vector3{0, 0, 1})})
		cs.simplex.Add(SimplexPoint{mPoint: cs.support(other, vmath.Vector3{-1, -1, -1})})
	}
	// setup tetrahedron faces
	cs.simplex.ClearFaces()
	cs.simplex.AddFace(SimplexFace{p1: 0, p2: 1, p3: 2})
	cs.simplex.AddFace(SimplexFace{p1: 0, p2: 2, p3: 3})
	cs.simplex.AddFace(SimplexFace{p1: 0, p2: 3, p3: 1})
	cs.simplex.AddFace(SimplexFace{p1: 1, p2: 2, p3: 3})
	for iters := 0; iters < epaMaxIterations; iters = iters + 1 {
		face, dist := cs.simplex.ClosestFace()
		norm := cs.simplex.FaceNormal(face.index)
		p := cs.support(other, norm)
		if cs.simplex.containsPoint(p, 0.00001) {
			return norm.Normalize().MultiplyScalar(dist)
		} else {
			cs.simplex.AddPointToFace(SimplexPoint{mPoint: p}, face.index)
		}
	}
	return vmath.Vector3{0, 0, 0}
}

func (cs *ConvexSet) support(other *ConvexSet, direction vmath.Vector3) vmath.Vector3 {
	p1 := cs.farthestPointInDirection(direction)
	p2 := other.farthestPointInDirection(direction.MultiplyScalar(-1))
	result := p1.Subtract(p2)
	return result
}

func (cs *ConvexSet) farthestPointInDirection(direction vmath.Vector3) vmath.Vector3 {
	var farthest vmath.Vector3
	max := 0.0
	for i, point := range cs.verticies {
		p := cs.transformPoint(point)
		d := direction.Dot(p)
		if d >= max || i == 0 {
			max = d
			farthest = p
		}
	}
	return farthest
}

func (cs *ConvexSet) transformPoint(point vmath.Vector3) vmath.Vector3 {
	return cs.orientation.Apply(point).Add(cs.offset)
}
