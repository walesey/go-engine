package util

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func Vec3_64_from(v mgl32.Vec3) mgl64.Vec3 {
	return mgl64.Vec3{float64(v[0]), float64(v[1]), float64(v[2])}
}

func Vec3_32_from(v mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(v[0]), float32(v[1]), float32(v[2])}
}

func Quat_64_from(q mgl32.Quat) mgl64.Quat {
	return mgl64.Quat{V: Vec3_64_from(q.V), W: float64(q.W)}
}

func Quat_32_from(q mgl64.Quat) mgl32.Quat {
	return mgl32.Quat{V: Vec3_32_from(q.V), W: float32(q.W)}
}

func Vec3LenSq(v1 mgl32.Vec3) float32 {
	return v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2]
}

func Vec3LenSq64(v1 mgl64.Vec3) float64 {
	return v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2]
}

func Vec2LenSq(v1 mgl32.Vec2) float32 {
	return v1[0]*v1[0] + v1[1]*v1[1]
}

func Vec2LenSq64(v1 mgl64.Vec2) float64 {
	return v1[0]*v1[0] + v1[1]*v1[1]
}

func Vec3Lerp(start, end mgl32.Vec3, amount float32) mgl32.Vec3 {
	return start.Mul(1.0 - amount).Add(end.Mul(amount))
}

func Vec3Lerp64(start, end mgl64.Vec3, amount float64) mgl64.Vec3 {
	return start.Mul(1.0 - amount).Add(end.Mul(amount))
}

func Vec2Cross(v1, v2 mgl32.Vec2) float32 {
	return v1[0]*v2[1] - v1[1]*v2[0]
}

func Vec2Rotate(v mgl32.Vec2, angle float32) mgl32.Vec2 {
	sn, cs := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))
	return mgl32.Vec2{v[0]*cs - v[1]*sn, v[0]*sn + v[1]*cs}
}

func Vec2AngleBetween(v1, v2 mgl32.Vec2) float32 {
	return float32(math.Acos(float64(v1.Normalize().Dot(v2.Normalize()))))
}

func Mat4From(scale, translation mgl32.Vec3, orientation mgl32.Quat) mgl32.Mat4 {
	s, tx, q := scale, translation, orientation
	return mgl32.Translate3D(tx[0], tx[1], tx[2]).Mul4(q.Mat4()).Mul4(mgl32.Scale3D(s[0], s[1], s[2]))
}

func Mat4From64(scale, translation mgl64.Vec3, orientation mgl64.Quat) mgl64.Mat4 {
	s, tx, q := scale, translation, orientation
	return mgl64.Translate3D(tx[0], tx[1], tx[2]).Mul4(q.Mat4()).Mul4(mgl64.Scale3D(s[0], s[1], s[2]))
}

func Round(val float32, roundOn float32, places int) float32 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * float64(val)
	_, div := math.Modf(digit)
	if div >= float64(roundOn) {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return float32(round / pow)
}

func RoundHalfUp(val float32) (newVal int) {
	return (int)(Round(val, .5, 0))
}

//PointToPlaneDist distance from plane (a,b,c) to point
func PointToPlaneDist(a, b, c, point mgl32.Vec3) float32 {
	ab := b.Sub(a)
	ac := c.Sub(a)
	ap := point.Sub(a)
	normal := ac.Cross(ab).Normalize()
	return float32(math.Abs(float64(ap.Dot(normal))))
}

//PointToLineDist distance from line (a,b) to point
func PointToLineDist(a, b, point mgl32.Vec3) float32 {
	ab := b.Sub(a)
	ap := point.Sub(a)
	prj := ap.Dot(ab)
	lenSq := ab.Dot(ab)
	t := prj / lenSq
	return ab.Mul(t).Add(a).Sub(point).Len()
}

//PointLiesInsideTriangle - return true if the point lies within the triangle formed by points (a,b,c)
func PointLiesInsideTriangle(a, b, c, point mgl32.Vec3) bool {
	ab := a.Sub(b)
	bc := b.Sub(c)
	ca := c.Sub(a)
	ap := a.Sub(point)
	bp := b.Sub(point)
	cp := c.Sub(point)
	cross1 := ap.Cross(ab)
	cross2 := bp.Cross(bc)
	cross3 := cp.Cross(ca)
	dot12 := cross1.Dot(cross2)
	dot13 := cross1.Dot(cross3)
	return dot12 > 0 && dot13 > 0
}

//RayTriangleIntersect - return true if the ray formed by start and direction intersects with the triangle formed by points (a,b,c)
func RayTriangleIntersect(a, b, c, start, direction mgl32.Vec3) (point mgl32.Vec3, ok bool) {
	normal := a.Sub(b).Cross(a.Sub(c))
	u := direction
	w := start.Sub(a)
	d := normal.Dot(u)
	n := -normal.Dot(w)

	// is parallel
	if math.Abs(float64(d)) < 0.00000001 {
		ok = false
		return
	}

	// intersect param
	s := n / d
	if s < 0 {
		ok = false
		return
	}
	point = start.Add(u.Mul(s))
	ok = PointLiesInsideTriangle(a, b, c, point)
	return
}

//PointLiesInsideAABB - return true if the point lies within the rectan formed by points a and b
func PointLiesInsideAABB(a, b, point mgl32.Vec2) bool {
	if (point.X() > a.X() && point.X() > b.X()) || (point.X() < a.X() && point.X() < b.X()) {
		return false
	}
	if (point.Y() > a.Y() && point.Y() > b.Y()) || (point.Y() < a.Y() && point.Y() < b.Y()) {
		return false
	}
	return true
}

//FacingOrientation - return an orientation that always faces the given direction with rotation
func FacingOrientation(rotation float32, direction, normal, tangent mgl32.Vec3) mgl32.Quat {
	result := FacingOrientation64(float64(rotation), Vec3_64_from(direction), Vec3_64_from(normal), Vec3_64_from(tangent))
	return Quat_32_from(result)
}

func FacingOrientation64(rotation float64, direction, normal, tangent mgl64.Vec3) mgl64.Quat {
	betweenVectorsQ := mgl64.QuatBetweenVectors(normal, direction)
	angleQ := mgl64.QuatRotate(rotation, normal)
	orientation := betweenVectorsQ.Mul(angleQ)
	return orientation
}

// TwoSegmentIntersect - find the intersection point of two line segments <p11-p12> and <p21-p22>
func TwoSegmentIntersect(p11, p12, p21, p22 mgl32.Vec2) (mgl32.Vec2, error) {
	p := p11
	q := p21
	r := p12.Sub(p11)
	s := p22.Sub(p21)
	if math.Abs(float64(Vec2Cross(r, s))) < 0.0000001 {
		return mgl32.Vec2{}, fmt.Errorf("No intersections: lines parallel")
	}
	t := Vec2Cross(q.Sub(p), s) / Vec2Cross(r, s)
	u := Vec2Cross(p.Sub(q), r) / Vec2Cross(s, r)
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		return p.Add(r.Mul(t)), nil
	}
	return mgl32.Vec2{}, fmt.Errorf("No intersections")
}

func SegmentCircleIntersect(radius float32, center, start, finish mgl32.Vec2) (mgl32.Vec2, error) {
	d := finish.Sub(start)
	f := start.Sub(center)

	a := d.Dot(d)
	b := f.Mul(2).Dot(d)
	c := f.Dot(f) - radius*radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return mgl32.Vec2{}, fmt.Errorf("No intersection")
	} else {
		discriminant = float32(math.Sqrt(float64(discriminant)))

		t1 := (-b - discriminant) / (2 * a)
		t2 := (-b + discriminant) / (2 * a)

		if t1 >= 0 && t1 <= 1 {
			return mgl32.Vec2{start.X() + t1*d.X(), start.Y() + t1*d.Y()}, nil
		}
		if t2 >= 0 && t2 <= 1 {
			return mgl32.Vec2{start.X() + t2*d.X(), start.Y() + t2*d.Y()}, nil
		}
	}

	return mgl32.Vec2{}, fmt.Errorf("No intersections")
}

// Max - returns the largest value of the values provided.
func Max(values ...float64) float64 {
	var result float64
	for i, value := range values {
		if i == 0 || value > result {
			result = value
		}
	}
	return result
}

// Max - returns the largest value of the values provided.
func MaxF32(values ...float32) float32 {
	var result float32
	for i, value := range values {
		if i == 0 || value > result {
			result = value
		}
	}
	return result
}

// Min - returns the smallest value of the values provided.
func Min(values ...float64) float64 {
	var result float64
	for i, value := range values {
		if i == 0 || value < result {
			result = value
		}
	}
	return result
}

// Min - returns the smallest value of the values provided.
func MinF32(values ...float32) float32 {
	var result float32
	for i, value := range values {
		if i == 0 || value < result {
			result = value
		}
	}
	return result
}
