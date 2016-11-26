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
	UseRendererParams(params RendererParams)

	CreateGeometry(geometry *Geometry)
	DestroyGeometry(geometry *Geometry)
	DrawGeometry(geometry *Geometry, transform mgl32.Mat4)

	CreateMaterial(material *Material)
	DestroyMaterial(material *Material)
	UseMaterial(material *Material)

	CreateCubeMap(cubeMap *CubeMap)
	DestroyCubeMap(cubeMap *CubeMap)
	UseCubeMap(cubeMap *CubeMap)

	CreateShader(shader *Shader)
	UseShader(shader *Shader)

	CreatePostEffect(shader *Shader)
	DestroyPostEffects(shader *Shader)

	AddLight(light *Light)
	RemoveLight(light *Light)
}

//A Spatial is something that can be Drawn by a Renderer
type Spatial interface {
	Draw(renderer Renderer, transform mgl32.Mat4)
	Optimize(geometry *Geometry, transform mgl32.Mat4)
	Destroy(renderer Renderer)
	Centre() mgl32.Vec3
	SetParent(parent *Node)
}

//An Entity is something that can be scaled, positioned and rotated (orientation)
type Entity interface {
	SetScale(scale mgl32.Vec3)
	SetTranslation(translation mgl32.Vec3)
	SetOrientation(orientation mgl32.Quat)
}
