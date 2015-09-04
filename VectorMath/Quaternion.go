package vectorMath

import "math"

type Quaternion struct {
    X,Y,Z,W float64
}

func IdentityQuaternion() *Quaternion {
	return &Quaternion{0,0,0,1}
}

func AngleAxisQuaternion( angle float64, axis Vector3 ) *Quaternion {
	quat := IdentityQuaternion()
	quat.AngleAxis(angle, axis)
	return quat
}

func (q *Quaternion) Set( value Quaternion ) {
    q.SetValue(value.X, value.Y, value.Z, value.W)
}

func (q *Quaternion) SetValue(X,Y,Z,W float64){
    q.X = X
    q.Y = Y
    q.Z = Z
    q.W = W
}



//
func (q *Quaternion) AngleAxis( angle float64, axis Vector3 ) {
	s := math.Sin( angle/2 )
	q.X = axis.X * s
	q.Y = axis.Y * s
	q.Z = axis.Z * s
	q.W = math.Cos(angle/2)
}
