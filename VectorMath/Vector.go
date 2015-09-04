package vectorMath

import "math"

type Vector3 struct {
    X,Y,Z float64
}

type Vector interface {
    Length() float64
    LengthSquared() float64
}

func CreateVector3() *Vector3{
    return &Vector3{0,0,0}
}

func (v *Vector3) Set(value Vector3){
    v.SetValue(value.X, value.Y, value.Z)
}

func (v *Vector3) SetValue(X,Y,Z float64){
    v.X = X
    v.Y = Y
    v.Z = Z
}

func (v *Vector3) LengthSquared() float64 {
    return ( v.X * v.X ) + ( v.Y * v.Y ) + ( v.Z * v.Z );
}

func (v *Vector3) Length() float64 {
    return math.Sqrt( v.LengthSquared() );
}

func (v *Vector3) Copy() *Vector3 {
    result := CreateVector3()
    result.SetValue(v.X, v.Y, v.Z)
    return result
}

func (v *Vector3) Add( other Vector3 ) {
    v.X = v.X + other.X
    v.Y = v.Y + other.Y
    v.Z = v.Z + other.Z
}

func (v *Vector3) Subtract( other Vector3 ) {
    v.X = v.X - other.X
    v.Y = v.Y - other.Y
    v.Z = v.Z - other.Z
}

func (v *Vector3) Multiply( other Vector3 ) {
    v.X = v.X * other.X
    v.Y = v.Y * other.Y
    v.Z = v.Z * other.Z
}

func (v *Vector3) MultiplyScalar( scalar float64 ) {
    v.X = v.X * scalar
    v.Y = v.Y * scalar
    v.Z = v.Z * scalar
}

func (v *Vector3) Divide( other Vector3 ) {
    v.X = v.X / other.X
    v.Y = v.Y / other.Y
    v.Z = v.Z / other.Z
}

func (v *Vector3) DivideScalar( scalar float64 ) {
    v.X = v.X / scalar
    v.Y = v.Y / scalar
    v.Z = v.Z / scalar
}

func (v *Vector3) Cross( other Vector3 ) *Vector3 {
    result := v.Copy()
    result.X = (other.Y*v.Z) - (other.Z*v.Y)
    result.Y = (other.Z*v.X) - (other.X*v.Z)
    result.Z = (other.X*v.Y) - (other.Y*v.X)
    return result
}

func (v *Vector3) Dot( other Vector3 ) float64 {
    return (v.X * other.X) + (v.Y * other.Y) + (v.Z * other.Z)
}
