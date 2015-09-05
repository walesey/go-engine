package vectorMath

import "math"

type Quaternion struct {
    X,Y,Z,W float64
}

func IdentityQuaternion() Quaternion {
	return Quaternion{0,0,0,1}
}

func AngleAxis( angle float64, axis Vector3 ) Quaternion {
	s := math.Sin( angle/2 )
	return Quaternion{
		axis.X * s,
		axis.Y * s,
		axis.Z * s,
		math.Cos(angle/2),
	}
}

func BetweenVectors( start, finish Vector3 ) Quaternion {
	start = start.Normalize()
	finish = finish.Normalize()
	axis := start.Cross(finish)
	angle := math.Acos( start.Dot(finish) )
	return AngleAxis(angle, axis)
}

func (q *Quaternion) Set( value Quaternion ) Quaternion {
    return q.SetValue(value.X, value.Y, value.Z, value.W)
}

func (q *Quaternion) SetValue(X,Y,Z,W float64) Quaternion {
    q.X = X
    q.Y = Y
    q.Z = Z
    q.W = W
	return *q
}

func (q Quaternion) Multiply( other Quaternion ) Quaternion {	
	return Quaternion{
		q.X * other.W + q.Y * other.Z - q.Z * other.Y + q.W * other.X,
		-q.X * other.Z + q.Y * other.W + q.Z * other.X + q.W * other.Y,
		q.X * other.Y - q.Y * other.X + q.Z * other.W + q.W * other.Z,
		-q.X * other.X - q.Y * other.Y - q.Z * other.Z + q.W * other.W,
	}
}

func (q Quaternion) Apply( v Vector3 ) Vector3 {	
	return Vector3{
		q.W * q.W * v.X + 2 * q.Y * q.W * v.Z - 2 * q.Z * q.W * v.Y + q.X * q.X * v.X + 2 * q.Y * q.X * v.Y + 2 * q.Z * q.X * v.Z - q.Z * q.Z * v.X - q.Y * q.Y * v.X,
		2 * q.X * q.Y * v.X + q.Y * q.Y * v.Y + 2 * q.Z * q.Y * v.Z + 2 * q.W * q.Z * v.X - q.Z * q.Z * v.Y + q.W * q.W * v.Y - 2 * q.X * q.W * v.Z - q.X * q.X * v.Y,
		2 * q.X * q.Z * v.X + 2 * q.Y * q.Z * v.Y + q.Z * q.Z * v.Z - 2 * q.W * q.Y * v.X - q.Y * q.Y * v.Z + 2 * q.W * q.X * v.Y - q.X * q.X * v.Z + q.W * q.W * v.Z,
	}
}