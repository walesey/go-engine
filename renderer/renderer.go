package renderer

import "github.com/go-gl/mathgl/mgl32"

type Renderer interface {
	SetInit(callback func())
	SetUpdate(callback func())
	SetRender(callback func())
	SetCamera(camera *Camera)
	Start()

	BackGroundColor(r, g, b, a float32)
	WindowDimensions() mgl32.Vec2
	LockCursor(lock bool)

	CreateGeometry(geometry *Geometry)
	DestroyGeometry(geometry *Geometry)
	DrawGeometry(geometry *Geometry, transform mgl32.Mat4)

	CreateMaterial(material *Material)
	DestroyMaterial(material *Material)

	CreateShader(shader *Shader)
	CreatePostEffect(shader *Shader)
	DestroyPostEffects(shader *Shader)
}
