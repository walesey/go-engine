package engine

import (
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

// Engine is a wrapper for all the go-engine boilerblate code.
// It sets up a basic render / Update loop and provides a nice interface for writing games.
type Engine interface {
	Start(Init func())
	AddOrtho(spatials ...renderer.Spatial)
	AddSpatial(spatial ...renderer.Spatial)
	RemoveSpatial(spatial ...renderer.Spatial)
	AddUpdatable(updatables ...Updatable)
	RemoveUpdatable(updatables ...Updatable)
	Sky(material *renderer.Material, size float64)
	Camera() *renderer.Camera
}

type EngineImpl struct {
	fpsMeter       *renderer.FPSMeter
	renderer       renderer.Renderer
	sceneGraph     renderer.SceneGraph
	orthoNode      *renderer.Node
	camera         *renderer.Camera
	updatableStore *UpdatableStore
}

func (engine *EngineImpl) Start(Init func()) {
	engine.renderer.Init(Init)
	engine.renderer.Update(engine.Update)
	engine.renderer.Render(engine.Render)
	engine.renderer.Start()
}

func (engine *EngineImpl) Update() {
	engine.fpsMeter.UpdateFPSMeter()
	engine.updatableStore.UpdateAll(0.018)
}

func (engine *EngineImpl) Render() {
	engine.camera.Perspective()
	engine.sceneGraph.RenderScene(engine.renderer)
	engine.camera.Ortho()
	engine.orthoNode.Draw(engine.renderer)
}

func (engine *EngineImpl) AddOrtho(spatials ...renderer.Spatial) {
	for _, s := range spatials {
		engine.orthoNode.Add(s)
	}
}

func (engine *EngineImpl) AddSpatial(spatials ...renderer.Spatial) {
	for _, s := range spatials {
		engine.sceneGraph.Add(s)
	}
}

func (engine *EngineImpl) RemoveSpatial(spatials ...renderer.Spatial) {
	for _, s := range spatials {
		engine.sceneGraph.Remove(s)
	}
}

func (engine *EngineImpl) AddUpdatable(updatables ...Updatable) {
	engine.updatableStore.Add(updatables...)
}

func (engine *EngineImpl) RemoveUpdatable(updatables ...Updatable) {
	engine.updatableStore.Remove(updatables...)
}

func (engine *EngineImpl) Sky(material *renderer.Material, size float64) {
	geom := renderer.CreateSkyBox()
	geom.Material = material
	geom.Material.LightingMode = renderer.MODE_UNLIT
	geom.CullBackface = false
	skyNode := renderer.CreateNode()
	skyNode.Add(geom)
	skyNode.SetRotation(1.57, vmath.Vector3{0, 1, 0})
	skyNode.SetScale(vmath.Vector3{1, 1, 1}.MultiplyScalar(size))
	engine.AddSpatial(skyNode)
	cubeMap := renderer.CreateCubemap(material.Diffuse)
	engine.renderer.ReflectionMap(cubeMap)
}

func (engine *EngineImpl) Camera() *renderer.Camera {
	return engine.camera
}

func NewEngine(r renderer.Renderer) Engine {
	fpsMeter := renderer.CreateFPSMeter(1.0)
	fpsMeter.FpsCap = 60

	sceneGraph := renderer.CreateSceneGraph()
	orthoNode := renderer.CreateNode()
	updatableStore := NewUpdatableStore()
	camera := renderer.CreateCamera(r)

	return &EngineImpl{
		fpsMeter:       fpsMeter,
		sceneGraph:     sceneGraph,
		orthoNode:      orthoNode,
		updatableStore: updatableStore,
		renderer:       r,
		camera:         camera,
	}
}
