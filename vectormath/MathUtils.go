package vectormath

import (
	"fmt"
	"math"
)

func Round(val float64, roundOn float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}

func RoundHalfUp(val float64) (newVal int) {
	return (int)(Round(val, .5, 0))
}

func ApproxEqual(value1, value2, epsilon float64) bool {
	return math.Abs(value1-value2) <= epsilon
}

//PointToPlaneDist distance from plane (a,b,c) to point
func PointToPlaneDist(a, b, c, point Vector3) float64 {
	ab := b.Subtract(a)
	ac := c.Subtract(a)
	ap := point.Subtract(a)
	normal := ab.Cross(ac).Normalize()
	return math.Abs(ap.Dot(normal))
}

//PointToLineDist distance from line (a,b) to point
func PointToLineDist(a, b, point Vector3) float64 {
	ab := b.Subtract(a)
	ap := point.Subtract(a)
	prj := ap.Dot(ab)
	lenSq := ab.Dot(ab)
	t := prj / lenSq
	return ab.MultiplyScalar(t).Add(a).Subtract(point).Length()
}

//PointLiesInsideTriangle - return true if the point lies within the triangle formed by points (a,b,c)
func PointLiesInsideTriangle(a, b, c, point Vector3) bool {
	ab := a.Subtract(b)
	bc := b.Subtract(c)
	ca := c.Subtract(a)
	ap := a.Subtract(point)
	bp := b.Subtract(point)
	cp := c.Subtract(point)
	cross1 := ab.Cross(ap)
	cross2 := bc.Cross(bp)
	cross3 := ca.Cross(cp)
	dot12 := cross1.Dot(cross2)
	dot13 := cross1.Dot(cross3)
	return dot12 > 0 && dot13 > 0
}

//PointLiesInsideAABB - return true if the point lies within the rectan formed by points a and b
func PointLiesInsideAABB(a, b, point Vector2) bool {
	if (point.X > a.X && point.X > b.X) || (point.X < a.X && point.X < b.X) {
		return false
	}
	if (point.Y > a.Y && point.Y > b.Y) || (point.Y < a.Y && point.Y < b.Y) {
		return false
	}
	return true
}

func CramerSolve3(mat Matrix3, col Vector3) (Vector3, error) {
	det := mat.Determinant()
	if det != 0 {
		matx := Matrix3{
			col.X, mat.M01, mat.M02,
			col.Y, mat.M11, mat.M12,
			col.Z, mat.M21, mat.M22,
		}
		maty := Matrix3{
			mat.M00, col.X, mat.M02,
			mat.M10, col.Y, mat.M12,
			mat.M20, col.Z, mat.M22,
		}
		matz := Matrix3{
			mat.M00, mat.M01, col.X,
			mat.M10, mat.M11, col.Y,
			mat.M20, mat.M21, col.Z,
		}
		detx := matx.Determinant()
		dety := maty.Determinant()
		detz := matz.Determinant()
		return Vector3{X: detx / det, Y: dety / det, Z: detz / det}, nil
	}
	return Vector3{}, fmt.Errorf("No solution")
}

func CramerSolve4(mat Matrix4, col Vector4) (Vector4, error) {
	det := mat.Determinant()
	if det != 0 {
		matx := Matrix4{
			col.X, mat.M01, mat.M02, mat.M03,
			col.Y, mat.M11, mat.M12, mat.M13,
			col.Z, mat.M21, mat.M22, mat.M23,
			col.W, mat.M31, mat.M32, mat.M33,
		}
		maty := Matrix4{
			mat.M00, col.X, mat.M02, mat.M03,
			mat.M10, col.Y, mat.M12, mat.M13,
			mat.M20, col.Z, mat.M22, mat.M23,
			mat.M30, col.W, mat.M32, mat.M33,
		}
		matz := Matrix4{
			mat.M00, mat.M01, col.X, mat.M03,
			mat.M10, mat.M11, col.Y, mat.M13,
			mat.M20, mat.M21, col.Z, mat.M23,
			mat.M30, mat.M31, col.W, mat.M33,
		}
		matw := Matrix4{
			mat.M00, mat.M01, mat.M02, col.X,
			mat.M10, mat.M11, mat.M12, col.Y,
			mat.M20, mat.M21, mat.M22, col.Z,
			mat.M30, mat.M31, mat.M32, col.W,
		}
		detx := matx.Determinant()
		dety := maty.Determinant()
		detz := matz.Determinant()
		detw := matw.Determinant()
		return Vector4{X: detx / det, Y: dety / det, Z: detz / det, W: detw / det}, nil
	}
	return Vector4{}, fmt.Errorf("No solution")
}

// RowMat3ColumnProduct - Calculate the product of a row vec, a 3x3 mat and a 1x3 column vec
func RowMat3ColumnProduct(mat Matrix3, row, col Vector3) float64 {
	a := row.X*mat.M00*col.X + row.Y*mat.M10*col.X + row.Z*mat.M20*col.X
	b := row.X*mat.M01*col.Y + row.Y*mat.M11*col.Y + row.Z*mat.M21*col.Y
	c := row.X*mat.M02*col.Z + row.Y*mat.M12*col.Z + row.Z*mat.M22*col.Z
	return a + b + c
}

//FacingOrientation - return an orientation that always faces the given direction with rotation
func FacingOrientation(rotation float64, direction, normal, tangent Vector3) Quaternion {
	betweenVectorsQ := BetweenVectors(normal, direction)
	angleQ := AngleAxis(rotation, normal)
	orientation := betweenVectorsQ.Multiply(angleQ)
	return orientation
}

// TwoSegmentIntersect - find the intersection point of two line segments <p11-p12> and <p21-p22>
func TwoSegmentIntersect(p11, p12, p21, p22 Vector2) (Vector2, error) {
	p := p11
	q := p21
	r := p12.Subtract(p11)
	s := p22.Subtract(p21)
	if math.Abs(r.Cross(s)) < 0.0000001 {
		return Vector2{}, fmt.Errorf("No intersections: lines parallel")
	}
	t := q.Subtract(p).Cross(s) / r.Cross(s)
	u := p.Subtract(q).Cross(r) / s.Cross(r)
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		return p.Add(r.MultiplyScalar(t)), nil
	}
	return Vector2{}, fmt.Errorf("No intersections")
}

func SegmentCircleIntersect(radius float64, center, start, finish Vector2) (Vector2, error) {
	d := finish.Subtract(start)
	f := start.Subtract(center)

	a := d.Dot(d)
	b := f.MultiplyScalar(2).Dot(d)
	c := f.Dot(f) - radius*radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return Vector2{}, fmt.Errorf("No intersection")
	} else {
		discriminant = math.Sqrt(discriminant)

		t1 := (-b - discriminant) / (2 * a)
		t2 := (-b + discriminant) / (2 * a)

		if t1 >= 0 && t1 <= 1 {
			return Vector2{X: start.X + t1*d.X, Y: start.Y + t1*d.Y}, nil
		}
		if t2 >= 0 && t2 <= 1 {
			return Vector2{X: start.X + t2*d.X, Y: start.Y + t2*d.Y}, nil
		}
	}

	return Vector2{}, fmt.Errorf("No intersections")
}
