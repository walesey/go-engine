package vectormath

import "math"

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
