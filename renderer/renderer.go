package renderer

import "github.com/walesey/go-engine/vectormath"

type Uniform struct {
	Name  string
	Value interface{}
}

type Shader struct {
	Name     string
	Uniforms []Uniform
}

type Renderer interface {
	Init(callback func())
	Update(callback func())
	Render(callback func())
	Start()
	BackGroundColor(r, g, b, a float32)
	Projection(angle, near, far float32)
	Camera(location, lookat, up vectormath.Vector3)
	Ortho()
	CameraLocation() vectormath.Vector3
	PopTransform()
	PushTransform()
	EnableDepthTest(depthTest bool)
	ApplyTransform(transform Transform)
	CreateGeometry(geometry *Geometry)
	DestroyGeometry(geometry *Geometry)
	CreateMaterial(material *Material)
	DestroyMaterial(material *Material)
	DrawGeometry(geometry *Geometry)
	CreateLight(ar, ag, ab, dr, dg, db, sr, sg, sb float32, directional bool, position vectormath.Vector3, i int)
	DestroyLight(i int)
	ReflectionMap(cm *CubeMap)
	CreatePostEffect(shader Shader)
	DestroyPostEffects(shader Shader)
	LockCursor(lock bool)
}
