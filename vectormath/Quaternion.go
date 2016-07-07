package vectormath

import "math"

type Quaternion struct {
	X, Y, Z, W float64
}

func IdentityQuaternion() Quaternion {
	return Quaternion{0, 0, 0, 1}
}

func AngleAxis(angle float64, axis Vector3) Quaternion {
	s := math.Sin(angle / 2)
	return Quaternion{
		axis.X * s,
		axis.Y * s,
		axis.Z * s,
		math.Cos(angle / 2),
	}
}

func BetweenVectors(start, finish Vector3) Quaternion {
	axis := finish.Cross(start).Normalize()
	angle := start.AngleBetween(finish)
	if finish.Dot(start) < 0 {
		angle := start.MultiplyScalar(-1).AngleBetween(finish)
		return AngleAxis(3.14-angle, axis)
	}
	return AngleAxis(angle, axis)
}

func (q Quaternion) Multiply(other Quaternion) Quaternion {
	return Quaternion{
		q.X*other.W + q.Y*other.Z - q.Z*other.Y + q.W*other.X,
		-q.X*other.Z + q.Y*other.W + q.Z*other.X + q.W*other.Y,
		q.X*other.Y - q.Y*other.X + q.Z*other.W + q.W*other.Z,
		-q.X*other.X - q.Y*other.Y - q.Z*other.Z + q.W*other.W,
	}
}

func (q Quaternion) MagnitudeSquared() float64 {
	return q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z
}

func (q Quaternion) Magnitude() float64 {
	magSq := q.MagnitudeSquared()
	if ApproxEqual(magSq, 1.0, 0.000001) {
		return 1.0
	}
	return math.Sqrt(magSq)
}

func (q Quaternion) Normalize() Quaternion {
	n := 1.0 / q.Magnitude()
	return Quaternion{
		q.X * n,
		q.Y * n,
		q.Z * n,
		q.W * n,
	}
}

func (q Quaternion) RotationMatrix4() Matrix4 {
	return Matrix4{
		1 - 2*q.Y*q.Y - 2*q.Z*q.Z, 2*q.X*q.Y + 2*q.W*q.Z, 2*q.X*q.Z - 2*q.W*q.Y, 0,
		2*q.X*q.Y - 2*q.W*q.Z, 1 - 2*q.X*q.X - 2*q.Z*q.Z, 2*q.Y*q.Z + 2*q.W*q.X, 0,
		2*q.X*q.Z + 2*q.W*q.Y, 2*q.Y*q.Z - 2*q.W*q.X, 1 - 2*q.X*q.X - 2*q.Y*q.Y, 0,
		0, 0, 0, 1,
	}
}

func (q Quaternion) RotationMatrix3() Matrix3 {
	return Matrix3{
		1 - 2*q.Y*q.Y - 2*q.Z*q.Z, 2*q.X*q.Y + 2*q.W*q.Z, 2*q.X*q.Z - 2*q.W*q.Y,
		2*q.X*q.Y - 2*q.W*q.Z, 1 - 2*q.X*q.X - 2*q.Z*q.Z, 2*q.Y*q.Z + 2*q.W*q.X,
		2*q.X*q.Z + 2*q.W*q.Y, 2*q.Y*q.Z - 2*q.W*q.X, 1 - 2*q.X*q.X - 2*q.Y*q.Y,
	}
}

func (q Quaternion) Apply(v Vector3) Vector3 {
	return Vector3{
		q.W*q.W*v.X + 2*q.Y*q.W*v.Z - 2*q.Z*q.W*v.Y + q.X*q.X*v.X + 2*q.Y*q.X*v.Y + 2*q.Z*q.X*v.Z - q.Z*q.Z*v.X - q.Y*q.Y*v.X,
		2*q.X*q.Y*v.X + q.Y*q.Y*v.Y + 2*q.Z*q.Y*v.Z + 2*q.W*q.Z*v.X - q.Z*q.Z*v.Y + q.W*q.W*v.Y - 2*q.X*q.W*v.Z - q.X*q.X*v.Y,
		2*q.X*q.Z*v.X + 2*q.Y*q.Z*v.Y + q.Z*q.Z*v.Z - 2*q.W*q.Y*v.X - q.Y*q.Y*v.Z + 2*q.W*q.X*v.Y - q.X*q.X*v.Z + q.W*q.W*v.Z,
	}
}

func (q Quaternion) ApproxEqual(other Quaternion, epsilon float64) bool {
	return ApproxEqual(q.X, other.X, epsilon) && ApproxEqual(q.Y, other.Y, epsilon) && ApproxEqual(q.Z, other.Z, epsilon) && ApproxEqual(q.W, other.W, epsilon)
}
