package renderer

import "github.com/go-gl/mathgl/mgl32"

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
	WindowDimensions() mgl32.Vec2
	Perspective(location, lookat, up mgl32.Vec3, angle, near, far float32)
	Ortho()
	CameraLocation() mgl32.Vec3
	FrustrumContainsSphere(radius float32) bool
	EnableDepthTest(depthTest bool)
	EnableDepthMask(depthMast bool)
	PushTransform(transform mgl32.Mat4)
	PopTransform()
	CreateGeometry(geometry *Geometry)
	DestroyGeometry(geometry *Geometry)
	CreateMaterial(material *Material)
	DestroyMaterial(material *Material)
	DrawGeometry(geometry *Geometry)
	CreateLight(ar, ag, ab, dr, dg, db, sr, sg, sb float32, directional bool, position mgl32.Vec3, i int)
	DestroyLight(i int)
	ReflectionMap(cm *CubeMap)
	CreatePostEffect(shader Shader)
	DestroyPostEffects(shader Shader)
	LockCursor(lock bool)
}
