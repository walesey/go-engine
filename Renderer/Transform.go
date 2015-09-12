package renderer

import (
	"github.com/Walesey/goEngine/vectorMath"
	"github.com/go-gl/mathgl/mgl32"
)

type Transform interface {
	ApplyTransform( transform Transform )
	Set( transform Transform )
	From( scale, translation vectorMath.Vector3, orientation vectorMath.Quaternion )
	TransformCoordinate( v vectorMath.Vector3 ) vectorMath.Vector3
	TransformNormal( n vectorMath.Vector3 ) vectorMath.Vector3
}

type GlTransform struct {
	Mat mgl32.Mat4
}

func CreateTransform() Transform {
	glTx := &GlTransform{ Mat: mgl32.Ident4() }
	return glTx
}

func FacingTransform( transform Transform, rotation float64, newNormal, normal, tangent vectorMath.Vector3 ) {
    angleCorrection := -tangent.AngleBetween( newNormal.Subtract(newNormal.Project(normal)) )
    if normal.Cross(tangent).Dot(newNormal) < 0 {
        angleCorrection = -angleCorrection
    }
    angleQ := vectorMath.AngleAxis( rotation + angleCorrection, normal )
    betweenVectorsQ := vectorMath.BetweenVectors( normal, newNormal ) 
    orientation := betweenVectorsQ.Multiply(angleQ)
    transform.From( vectorMath.Vector3{1,1,1}, vectorMath.Vector3{0,0,0}, orientation )
}

func (glTx *GlTransform) ApplyTransform( transform Transform ) {
	othertx, found := transform.(*GlTransform)
	if found {
		glTx.Mat = glTx.Mat.Mul4( othertx.Mat )
	}
}

func (glTx *GlTransform) Set( transform Transform ) {
	othertx, found := transform.(*GlTransform)
	if found {
		glTx.Mat = othertx.Mat
	}
}

func (glTx *GlTransform) From( scale, translation vectorMath.Vector3, orientation vectorMath.Quaternion ) {
	quat := convertQuaternion(orientation)
	tx := convertVector(translation)
	s := convertVector(scale)
	glTx.Mat = mgl32.Translate3D( tx[0], tx[1], tx[2] ).Mul4( mgl32.Scale3D( s[0], s[1], s[2] ) ).Mul4( quat.Mat4() )
}

func (glTx *GlTransform) TransformCoordinate( v vectorMath.Vector3 ) vectorMath.Vector3 {
	result := mgl32.TransformCoordinate(convertVector(v), glTx.Mat)
	return vectorMath.Vector3{ float64(result[0]), float64(result[1]), float64(result[2]) }
}

func (glTx *GlTransform) TransformNormal( n vectorMath.Vector3 ) vectorMath.Vector3 {
	result := mgl32.TransformCoordinate(convertVector(n), glTx.Mat.Inv().Transpose() )
	return vectorMath.Vector3{ float64(result[0]), float64(result[1]), float64(result[2]) }
}

func convertVector( v vectorMath.Vector3 ) mgl32.Vec3{
	return mgl32.Vec3{float32(v.X), float32(v.Y), float32(v.Z)}
}

func convertQuaternion( q vectorMath.Quaternion ) mgl32.Quat{
	return mgl32.Quat{ W:float32(q.W), V:mgl32.Vec3{float32(q.X), float32(q.Y), float32(q.Z)}}
}