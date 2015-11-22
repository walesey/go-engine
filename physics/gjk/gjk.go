package gjk

import (
	"fmt"

	"github.com/walesey/go-engine/physics"
	vmath "github.com/walesey/go-engine/vectormath"
)

const epaMaxIterations = 40
const gjkMaxIterations = 20

// ConvexSet A set of points representing a convex geometry
type ConvexSet struct {
	verticies   []vmath.Vector3
	offset      vmath.Vector3
	orientation vmath.Quaternion
	simplex     *Simplex
	searchDir   vmath.Vector3
}

// NewConvexSet create a new convex set
func NewConvexSet(verticies []vmath.Vector3) physics.Collider {
	return &ConvexSet{
		verticies:   verticies,
		offset:      vmath.Vector3{0, 0, 0},
		orientation: vmath.IdentityQuaternion(),
		simplex:     NewSimplex(),
		searchDir:   vmath.Vector3{0, 1, 0},
	}
}

// Offset test offset of this collider
func (cs *ConvexSet) Offset(offset vmath.Vector3, orientation vmath.Quaternion) {
	cs.offset = offset
	cs.orientation = orientation
}

// Overlap Calculate the overlap of a convex set and another collider
func (cs *ConvexSet) Overlap(other physics.Collider) bool {
	switch t := other.(type) {
	default:
		fmt.Printf("unsupported type for other collider: %T\n", t)
	case *ConvexSet:
		return cs.OverlapConvexSet(other.(*ConvexSet))
	}
	return false
}

// ContactPoints - gets the global contact point for this shape and the other shape
func (cs *ConvexSet) ContactPoint(other physics.Collider) vmath.Vector3 {
	switch t := other.(type) {
	default:
		fmt.Printf("unsupported type for other collider: %T\n", t)
	case *ConvexSet:
		return cs.ContactPointConvexSet(other.(*ConvexSet))
	}
	return vmath.Vector3{}
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

// OverlapConvexSet Uses GJK and Returns true if the two convex sets overlap
func (cs *ConvexSet) OverlapConvexSet(other *ConvexSet) bool {
	cs.simplex.Clear()
	d := cs.searchDir
	cs.simplex.Add(cs.support(other, d))
	d = d.MultiplyScalar(-1)
	for iters := 0; iters < gjkMaxIterations; iters = iters + 1 {
		cs.simplex.Add(cs.support(other, d))
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

// ContactPointsConvexSet Uses GJK and Returns the global contact point
func (cs *ConvexSet) ContactPointConvexSet(other *ConvexSet) vmath.Vector3 {
	cs.simplex.Clear()
	d := cs.searchDir
	cs.simplex.Add(cs.support(other, d))
	d = d.MultiplyScalar(-1)
	cs.simplex.Add(cs.support(other, d))
	d = closestPointOnSegmentToOrigin(cs.simplex.Get(0).mPoint, cs.simplex.Get(1).mPoint)
	d = d.MultiplyScalar(-1)
	if d.LengthSquared() < 0.0001 {
		return cs.sourcePoints()
	}
	cs.simplex.Add(cs.support(other, d))
	d = closestPointOnTriangleToOrigin(cs.simplex.Get(0).mPoint, cs.simplex.Get(1).mPoint, cs.simplex.Get(2).mPoint)
	for iters := 0; iters < gjkMaxIterations; iters = iters + 1 {
		d = d.MultiplyScalar(-1)
		if d.LengthSquared() < 0.00001 {
			return cs.sourcePoints()
		}
		cs.simplex.Add(cs.support(other, d))
		//work out which points to keep
		p012 := closestPointOnTriangleToOrigin(cs.simplex.Get(0).mPoint, cs.simplex.Get(1).mPoint, cs.simplex.Get(2).mPoint)
		p013 := closestPointOnTriangleToOrigin(cs.simplex.Get(0).mPoint, cs.simplex.Get(1).mPoint, cs.simplex.Get(3).mPoint)
		p023 := closestPointOnTriangleToOrigin(cs.simplex.Get(0).mPoint, cs.simplex.Get(2).mPoint, cs.simplex.Get(3).mPoint)
		p123 := closestPointOnTriangleToOrigin(cs.simplex.Get(1).mPoint, cs.simplex.Get(2).mPoint, cs.simplex.Get(3).mPoint)

		d012 := p012.LengthSquared() - 0.00001 // tolerance
		d013 := p013.LengthSquared()
		d023 := p023.LengthSquared()
		d123 := p123.LengthSquared()
		if d013 < d012 && d013 < d023 && d013 < d123 {
			cs.simplex.Remove(2)
			d = p013
		} else if d023 < d013 && d023 < d012 && d023 < d123 {
			cs.simplex.Remove(1)
			d = p023
		} else if d123 < d012 && d123 < d013 && d123 < d023 {
			cs.simplex.Remove(0)
			d = p123
		} else {
			cs.simplex.Remove(3)
			return cs.sourcePoints()
		}
	}
	return vmath.Vector3{}
}

//CollisionInfo uses EPA and returns the penetration vector
func (cs *ConvexSet) CollisionInfo(other *ConvexSet) vmath.Vector3 {
	if cs.simplex.Len() < 4 {
		//build new 4 point simplex
		cs.simplex.Clear()
		cs.simplex.Add(cs.support(other, vmath.Vector3{1, 0, 0}))
		cs.simplex.Add(cs.support(other, vmath.Vector3{0, 1, 0}))
		cs.simplex.Add(cs.support(other, vmath.Vector3{0, 0, 1}))
		cs.simplex.Add(cs.support(other, vmath.Vector3{-1, -1, -1}))
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
		if cs.simplex.containsPoint(p.mPoint, 0.00001) {
			return norm.Normalize().MultiplyScalar(dist)
		} else {
			cs.simplex.AddPointToFace(p, face.index)
		}
	}
	return vmath.Vector3{0, 0, 0}
}

// sourcePoint - return the source point given the two closes triangles of the two shapes
func (cs *ConvexSet) sourcePoints() vmath.Vector3 {
	if cs.simplex.Len() == 2 {

		p1 := cs.simplex.Get(0).point
		p2 := cs.simplex.Get(1).point
		other1 := cs.simplex.Get(0).otherPoint
		other2 := cs.simplex.Get(1).otherPoint

		//distances
		d1 := vmath.PointToLineDist(other1, other2, p1)
		d2 := vmath.PointToLineDist(other1, other2, p2)
		od1 := vmath.PointToLineDist(p1, p2, other1)
		od2 := vmath.PointToLineDist(p1, p2, other2)

		if d1 < d2 && d1 < od1 && d1 < od2 {
			return p1
		} else if d2 < d1 && d2 < od1 && d2 < od2 {
			return p2
		} else if od1 < d1 && od1 < d2 && od1 < od2 {
			return other1
		} else {
			return other2
		}

	} else if cs.simplex.Len() == 3 {

		//points
		p1 := cs.simplex.Get(0).point
		p2 := cs.simplex.Get(1).point
		p3 := cs.simplex.Get(2).point
		other1 := cs.simplex.Get(0).otherPoint
		other2 := cs.simplex.Get(1).otherPoint
		other3 := cs.simplex.Get(2).otherPoint

		//distances
		d1 := vmath.PointToPlaneDist(other1, other2, other3, p1)
		d2 := vmath.PointToPlaneDist(other1, other2, other3, p2)
		d3 := vmath.PointToPlaneDist(other1, other2, other3, p3)
		od1 := vmath.PointToPlaneDist(p1, p2, p3, other1)
		od2 := vmath.PointToPlaneDist(p1, p2, p3, other2)
		od3 := vmath.PointToPlaneDist(p1, p2, p3, other3)

		if d1 < d2 && d1 < d3 && d1 < od1 && d1 < od2 && d1 < od3 {
			return p1
		} else if d2 < d1 && d2 < d3 && d2 < od1 && d2 < od2 && d2 < od3 {
			return p2
		} else if d3 < d1 && d3 < d2 && d3 < od1 && d3 < od2 && d3 < od3 {
			return p3
		} else if od1 < d1 && od1 < d2 && od1 < d3 && od1 < od2 && od1 < od3 {
			return other1
		} else if od2 < d1 && od2 < d2 && od2 < d3 && od2 < od1 && od2 < od3 {
			return other2
		} else {
			return other3
		}
	}
	return vmath.Vector3{}
}

func closestPointOnSegmentToOrigin(sPoint1, sPoint2 vmath.Vector3) vmath.Vector3 {
	p12 := sPoint1.Subtract(sPoint2)
	prj := sPoint1.MultiplyScalar(-1).Dot(p12)
	lenSq := p12.Dot(p12)
	t := prj / lenSq
	return p12.MultiplyScalar(t).Add(sPoint1)
}

func closestPointOnTriangleToOrigin(tPoint1, tPoint2, tPoint3 vmath.Vector3) vmath.Vector3 {
	p12 := tPoint1.Subtract(tPoint2)
	p13 := tPoint1.Subtract(tPoint3)
	norm := p12.Cross(p13)
	d := tPoint1.Dot(norm)
	c := d / norm.LengthSquared()
	closestPoint := norm.MultiplyScalar(c)
	if vmath.PointLiesInsideTriangle(tPoint1, tPoint2, tPoint3, closestPoint) {
		return closestPoint
	}
	//return the closest triangle point
	if tPoint1.LengthSquared() < tPoint2.LengthSquared() && tPoint1.LengthSquared() < tPoint3.LengthSquared() {
		return tPoint1
	} else if tPoint2.LengthSquared() < tPoint3.LengthSquared() && tPoint2.LengthSquared() < tPoint1.LengthSquared() {
		return tPoint2
	} else {
		return tPoint3
	}
}

func (cs *ConvexSet) support(other *ConvexSet, direction vmath.Vector3) SimplexPoint {
	p1 := cs.farthestPointInDirection(direction)
	p2 := other.farthestPointInDirection(direction.MultiplyScalar(-1))
	mPoint := p1.Subtract(p2)
	return SimplexPoint{
		mPoint:     mPoint,
		point:      p1,
		otherPoint: p2,
	}
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
