package physics

import (
	"fmt"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Triangle struct {
	Point1, Point2, Point3 vmath.Vector3
}

type ConvexHull struct {
	triangles   []Triangle
	offset      *vmath.Vector3
	orientation *vmath.Quaternion
}

// NewConvexHull
func NewConvexHull(triangles []Triangle) *ConvexHull {
	return &ConvexHull{
		triangles: triangles,
	}
}

func (ch *ConvexHull) AttachTo(obj *PhysicsObject) {
	ch.offset = &obj.Position
	ch.orientation = &obj.Orientation
	obj.narrowPhase = ch
}

// Calculate the overlap of a convex hull and another collider
func (ch *ConvexHull) Overlap(other Collider) bool {
	switch t := other.(type) {
	default:
		fmt.Printf("unsupported type for other collider: %T\n", t)
	case *ConvexHull:
		return ch.OverlapConvexHull(other.(*ConvexHull))
	}
	return false
}

// Calculate the overlap of two convex hulls
func (ch *ConvexHull) OverlapConvexHull(other *ConvexHull) bool {
	for _, t := range ch.triangles {
		t = t.Rotate(*ch.orientation)
		t = t.Translate(*ch.offset)
		for _, tt := range other.triangles {
			tt = tt.Rotate(*other.orientation)
			tt = tt.Translate(*other.offset)
			if t.Overlap(tt) {
				return true
			}
		}
	}
	return false
}

// Calculates the overlap of two triangles
func (triangle Triangle) Overlap(other Triangle) bool {
	//get triangle plane equation
	N1 := triangle.Point2.Subtract(triangle.Point1).Cross(triangle.Point2.Subtract(triangle.Point1))
	d1 := -N1.Dot(triangle.Point1)
	p21, p22, p23, err := other.planeOverlap(N1, d1)
	if err != nil {
		return false
	}

	//get other plane equation
	N2 := other.Point1.Subtract(other.Point2).Cross(other.Point1.Subtract(other.Point3))
	d2 := -N2.Dot(other.Point1)
	p11, p12, p13, err := triangle.planeOverlap(N2, d2)
	if err != nil {
		return false
	}

	//line/plane equation
	lineD := N1.Cross(N2)

	//project points onto the intersection line
	prj11 := lineD.Dot(p11)
	prj12 := lineD.Dot(p12)
	prj13 := lineD.Dot(p13)
	prj21 := lineD.Dot(p21)
	prj22 := lineD.Dot(p22)
	prj23 := lineD.Dot(p23)

	//get distances from points to planes
	dist11 := N2.Dot(p11) + d2
	dist12 := N2.Dot(p12) + d2
	dist13 := N2.Dot(p13) + d2
	dist21 := N1.Dot(p21) + d1
	dist22 := N1.Dot(p22) + d1
	dist23 := N1.Dot(p23) + d1

	//co-planar
	if dist11 == 0 && dist12 == 0 && dist13 == 0 && dist21 == 0 && dist22 == 0 && dist23 == 0 {
		//TODO: handle coplanar
		return false
	}

	//find parametric intervals
	t11 := prj11 + ((prj13 - prj11) * (dist11 / (dist11 - dist13)))
	t12 := prj12 + ((prj13 - prj12) * (dist12 / (dist12 - dist13)))
	t21 := prj21 + ((prj23 - prj21) * (dist21 / (dist21 - dist23)))
	t22 := prj22 + ((prj23 - prj22) * (dist22 / (dist22 - dist23)))

	//return true if the intervals overlap
	return !((t11 > t21 && t11 > t22 && t12 > t21 && t12 > t22) || (t11 < t21 && t11 < t22 && t12 < t21 && t12 < t22))
}

//Intersection of triangle and plane (given by normal/d plane equation)
//Returns error if no intersection else returns the two points on one side of the plane followed by the point on the other side
func (triangle Triangle) planeOverlap(normal vmath.Vector3, d float64) (vmath.Vector3, vmath.Vector3, vmath.Vector3, error) {
	//calculate the distance from each vertex of triangle to the plane
	dist1 := normal.Dot(triangle.Point1) + d
	dist2 := normal.Dot(triangle.Point2) + d
	dist3 := normal.Dot(triangle.Point3) + d

	//triangle is completely on one side of triangle plane
	//All dist must be != 0, all must have the same sign
	if dist1 != 0 && dist2 != 0 && dist3 != 0 {
		if (dist1 > 0 && dist2 > 0 && dist3 > 0) || (dist1 < 0 && dist2 < 0 && dist3 < 0) {
			return triangle.Point1, triangle.Point2, triangle.Point3, fmt.Errorf("No intersection")
		}
	}

	//which vertex is on it's own
	if (dist1 > 0 && dist2 < 0 && dist3 < 0) || (dist1 < 0 && dist2 > 0 && dist3 > 0) {
		return triangle.Point2, triangle.Point3, triangle.Point1, nil
	}
	if (dist1 > 0 && dist2 < 0 && dist3 > 0) || (dist1 < 0 && dist2 > 0 && dist3 < 0) {
		return triangle.Point1, triangle.Point3, triangle.Point2, nil
	}
	return triangle.Point1, triangle.Point2, triangle.Point3, nil
}

//Rotate
func (t Triangle) Rotate(q vmath.Quaternion) Triangle {
	return Triangle{
		Point1: q.Apply(t.Point1),
		Point2: q.Apply(t.Point2),
		Point3: q.Apply(t.Point3),
	}
}

//Translate
func (t Triangle) Translate(v vmath.Vector3) Triangle {
	return Triangle{
		Point1: t.Point1.Add(v),
		Point2: t.Point2.Add(v),
		Point3: t.Point3.Add(v),
	}
}
