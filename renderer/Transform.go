package renderer

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/vectormath"
)

type Transform interface {
	ApplyTransform(transform Transform)
	Set(transform Transform)
	From(scale, translation vectormath.Vector3, orientation vectormath.Quaternion)
	TransformCoordinate(v vectormath.Vector3) vectormath.Vector3
	TransformNormal(n vectormath.Vector3) vectormath.Vector3
}

type GlTransform struct {
	Mat mgl32.Mat4
}

func CreateTransform() Transform {
	glTx := &GlTransform{Mat: mgl32.Ident4()}
	return glTx
}

func FacingTransform(transform Transform, rotation float64, newNormal, normal, tangent vectormath.Vector3) {
	angleCorrection := -tangent.AngleBetween(newNormal.Subtract(newNormal.Project(normal)))
	if normal.Cross(tangent).Dot(newNormal) < 0 {
		angleCorrection = -angleCorrection
	}
	angleQ := vectormath.AngleAxis(rotation+angleCorrection, normal)
	betweenVectorsQ := vectormath.BetweenVectors(normal, newNormal)
	orientation := betweenVectorsQ.Multiply(angleQ)
	transform.From(vectormath.Vector3{1, 1, 1}, vectormath.Vector3{0, 0, 0}, orientation)
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

func (glTx *GlTransform) From(scale, translation vectormath.Vector3, orientation vectormath.Quaternion) {
	quat := convertQuaternion(orientation)
	tx := convertVector(translation)
	s := convertVector(scale)
	glTx.Mat = mgl32.Translate3D(tx[0], tx[1], tx[2]).Mul4(mgl32.Scale3D(s[0], s[1], s[2])).Mul4(quat.Mat4())
}

func (glTx *GlTransform) TransformCoordinate(v vectormath.Vector3) vectormath.Vector3 {
	result := mgl32.TransformCoordinate(convertVector(v), glTx.Mat)
	return vectormath.Vector3{float64(result[0]), float64(result[1]), float64(result[2])}
}

func (glTx *GlTransform) TransformNormal(n vectormath.Vector3) vectormath.Vector3 {
	result := mgl32.TransformCoordinate(convertVector(n), glTx.Mat.Inv().Transpose())
	return vectormath.Vector3{float64(result[0]), float64(result[1]), float64(result[2])}
}

func convertVector(v vectormath.Vector3) mgl32.Vec3 {
	return mgl32.Vec3{float32(v.X), float32(v.Y), float32(v.Z)}
}

func convertQuaternion(q vectormath.Quaternion) mgl32.Quat {
	return mgl32.Quat{W: float32(q.W), V: mgl32.Vec3{float32(q.X), float32(q.Y), float32(q.Z)}}
}
