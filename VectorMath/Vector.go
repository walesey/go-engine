package VectorMath

import "math"

type Vector3 struct {
    X,Y,Z float64
}

type Vector interface {
    Length() float64
    LengthSquared() float64
}

func SetVector3(value Vector3, store *Vector3){
    SetValueVector3( value.X, value.Y, value.Z, store )
}

func SetValueVector3(X,Y,Z float64, store *Vector3){
    (*store).X = X
    (*store).Y = Y
    (*store).Z = Z
}

func (v Vector3) LengthSquared() float64 {
    return ( v.X * v.X ) + ( v.Y * v.Y ) + ( v.Z * v.Z );
}

func (v Vector3) Length() float64 {
    return math.Sqrt( v.LengthSquared() );
}

func (v Vector3) Copy( other Vector3 ) Vector3 {
    result := Vector3{}
    SetVector3( other, &result )
    return result
}

func (v Vector3) Add( other Vector3 ) Vector3 {
    result := Vector3{}
    v.AddStore( other, &result )
    return result
}

func (v Vector3) AddStore( other Vector3, store *Vector3 ) {
    SetValueVector3( v.X + other.X, v.Y + other.Y, v.Z + other.Z, store )
}

func (v Vector3) Subtract( other Vector3 ) Vector3 {
    result := Vector3{}
    v.SubtractStore( other, &result )
    return result
}

func (v Vector3) SubtractStore( other Vector3, store *Vector3 ) {
    SetValueVector3( v.X - other.X, v.Y - other.Y, v.Z - other.Z, store )
}

func (v Vector3) Multiply( other Vector3 ) Vector3 {
    result := Vector3{}
    v.MultiplyStore( other, &result )
    return result
}

func (v Vector3) MultiplyStore( other Vector3, store *Vector3 ) {
    SetValueVector3( v.X * other.X, v.Y * other.Y, v.Z * other.Z, store )
}

func (v Vector3) MultiplyScalar( scalar float64 ) Vector3 {
    result := Vector3{}
    v.MultiplyScalarStore( scalar, &result )
    return result
}

func (v Vector3) MultiplyScalarStore( scalar float64, store *Vector3 ) {
    SetValueVector3( v.X * scalar, v.Y * scalar, v.Z * scalar, store )
}

func (v Vector3) Divide( other Vector3 ) Vector3 {
    result := Vector3{}
    v.DivideStore( other, &result )
    return result
}

func (v Vector3) DivideStore( other Vector3, store *Vector3 ) {
    SetValueVector3( v.X / other.X, v.Y / other.Y, v.Z / other.Z, store )
}

func (v Vector3) DivideScalar( scalar float64 ) Vector3 {
    result := Vector3{}
    v.DivideScalarStore( scalar, &result )
    return result
}

func (v Vector3) DivideScalarStore( scalar float64, store *Vector3 ) {
    SetValueVector3( v.X / scalar, v.Y / scalar, v.Z / scalar, store )
}

func (v Vector3) Cross( other Vector3 ) Vector3 {
    result := Vector3{}
    v.CrossStore( other, &result )
    return result
}

func (v Vector3) CrossStore( other Vector3, store *Vector3 ) {
    SetValueVector3(
        (other.Y*v.Z) - (other.Z*v.Y),
        (other.Z*v.X) - (other.X*v.Z),
        (other.X*v.Y) - (other.Y*v.X), store )
}

func (v Vector3) Dot( other Vector3 ) float64 {
    return (v.X * other.X) + (v.Y * other.Y) + (v.Z * other.Z)
}
