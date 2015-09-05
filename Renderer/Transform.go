package renderer

import (
	"fmt"

	"github.com/Walesey/goEngine/vectorMath"
	"github.com/go-gl/mathgl/mgl32"
)

type Transform interface {
	ApplyTransform( transform Transform )
	From( scale, translation vectorMath.Vector3, orientation vectorMath.Quaternion )
}

type GlTransform struct {
	Mat mgl32.Mat4
}

func CreateTransform() Transform {
	glTx := &GlTransform{ Mat: mgl32.Ident4() }
	return glTx
}

func (glTx *GlTransform) ApplyTransform( transform Transform ) {
	switch v := transform.(type) {
    default:
        fmt.Printf("unexpected type for ApplyTransform GlTransform: %T", v)
    case *GlTransform:
		glTx.Mat = glTx.Mat.Mul4( transform.(*GlTransform).Mat )
    }
}

//
func (glTx *GlTransform) From( scale, translation vectorMath.Vector3, orientation vectorMath.Quaternion ) {
	quat := convertQuaternion(orientation)
	tx := convertVector(translation)
	s := convertVector(scale)
	glTx.Mat = mgl32.Scale3D( s[0], s[1], s[2] ).Mul4( quat.Mat4() ).Mul4( mgl32.Translate3D(  tx[0], tx[1], tx[2] ) )
}

//
func convertVector( v vectorMath.Vector3 ) mgl32.Vec3{
	return mgl32.Vec3{(float32)(v.X), (float32)(v.Y), (float32)(v.Z)}
}

func convertQuaternion( q vectorMath.Quaternion ) mgl32.Quat{
	return mgl32.Quat{ W:(float32)(q.W), V:mgl32.Vec3{(float32)(q.X), (float32)(q.Y), (float32)(q.Z)}}
}