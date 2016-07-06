package renderer

import (
	"github.com/go-gl/mathgl/mgl32"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Transform interface {
	ApplyTransform(transform Transform)
	Set(transform Transform)
	FromMatrix(matrix vmath.Matrix3)
	From(scale, translation vmath.Vector3, orientation vmath.Quaternion)
	TransformCoordinate(v vmath.Vector3) vmath.Vector3
	TransformNormal(n vmath.Vector3) vmath.Vector3
}

type GlTransform struct {
	Mat mgl32.Mat4
}

func CreateTransform() Transform {
	glTx := &GlTransform{Mat: mgl32.Ident4()}
	return glTx
}

func (glTx *GlTransform) ApplyTransform(transform Transform) {
	othertx, found := transform.(*GlTransform)
	if found {
		glTx.Mat = glTx.Mat.Mul4(othertx.Mat)
	}
}

func (glTx *GlTransform) Set(transform Transform) {
	othertx, found := transform.(*GlTransform)
	if found {
		glTx.Mat = othertx.Mat
	}
}

func (glTx *GlTransform) FromMatrix(matrix vmath.Matrix3) {
	glTx.Mat = convertMatrix(matrix)
}

func (glTx *GlTransform) From(scale, translation vmath.Vector3, orientation vmath.Quaternion) {
	quat := convertQuaternion(orientation)
	tx := convertVector(translation)
	s := convertVector(scale)
	glTx.Mat = mgl32.Translate3D(tx[0], tx[1], tx[2]).Mul4(quat.Mat4()).Mul4(mgl32.Scale3D(s[0], s[1], s[2]))
}

func (glTx *GlTransform) TransformCoordinate(v vmath.Vector3) vmath.Vector3 {
	result := mgl32.TransformCoordinate(convertVector(v), glTx.Mat)
	return vmath.Vector3{float64(result[0]), float64(result[1]), float64(result[2])}
}

func (glTx *GlTransform) TransformNormal(n vmath.Vector3) vmath.Vector3 {
	result := mgl32.TransformCoordinate(convertVector(n), glTx.Mat.Inv().Transpose())
	return vmath.Vector3{float64(result[0]), float64(result[1]), float64(result[2])}
}

func convertVector(v vmath.Vector3) mgl32.Vec3 {
	return mgl32.Vec3{float32(v.X), float32(v.Y), float32(v.Z)}
}

func convertQuaternion(q vmath.Quaternion) mgl32.Quat {
	return mgl32.Quat{W: float32(q.W), V: mgl32.Vec3{float32(q.X), float32(q.Y), float32(q.Z)}}
}

func convertMatrix(m vmath.Matrix3) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(m.M00), float32(m.M01), float32(m.M02), 0,
		float32(m.M10), float32(m.M11), float32(m.M12), 0,
		float32(m.M20), float32(m.M21), float32(m.M22), 0,
		0, 0, 0, 1,
	}
}
