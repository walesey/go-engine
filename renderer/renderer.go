package renderer

import "github.com/go-gl/mathgl/mgl32"

type Renderer interface {
	SetInit(callback func())
	SetUpdate(callback func())
	SetRender(callback func())
	SetCamera(camera *Camera)
	Camera() *Camera
	Start()

	BackGroundColor(r, g, b, a float32)
	WindowDimensions() mgl32.Vec2
	LockCursor(lock bool)
	UseRendererParams(params RendererParams)

	DrawGeometry(geometry *Geometry, transform mgl32.Mat4)
	DestroyGeometry(geometry *Geometry)

	UseMaterial(material *Material)
	DestroyMaterial(material *Material)

	UseCubeMap(cubeMap *CubeMap)
	DestroyCubeMap(cubeMap *CubeMap)

	UseShader(shader *Shader)

	CreatePostEffect(shader *Shader)
	DestroyPostEffects(shader *Shader)

	AddLight(light *Light)
	RemoveLight(light *Light)
}

//A Spatial is something that can be Added to scenegraph nodes
type Spatial interface {
	Draw(renderer Renderer, transform mgl32.Mat4)
	Optimize(geometry *Geometry, transform mgl32.Mat4)
	Destroy(renderer Renderer)
	Center() mgl32.Vec3
	BoundingRadius() float32
	OrthoOrder() int
	SetParent(parent *Node)
}

//An Entity is something that can be scaled, positioned and rotated (orientation)
type Entity interface {
	SetScale(scale mgl32.Vec3)
	SetTranslation(translation mgl32.Vec3)
	SetOrientation(orientation mgl32.Quat)
}
