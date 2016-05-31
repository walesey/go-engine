package engine

import (
	"fmt"
	"time"

	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

var ids = 0

// Engine is a wrapper for all the go-engine boilerblate code.
// It sets up a basic render / Update loop and provides a nice interface for writing games.
type Engine interface {
	Start(Init func())
	AddOrtho(spatial renderer.Spatial)
	AddSpatial(spatial renderer.Spatial)
	RemoveSpatial(spatial renderer.Spatial, destroy bool)
	RemoveOrtho(spatial renderer.Spatial, destroy bool)
	AddUpdatable(updatable Updatable) string
	RemoveUpdatable(updatable Updatable)
	AddUpdatableKey(key string, updatable Updatable)
	RemoveUpdatableKey(key string)
	Sky(material *renderer.Material, size float64)
	Camera() *renderer.Camera
	Update()
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
	if engine.renderer != nil {
		engine.renderer.Init(Init)
		engine.renderer.Update(engine.Update)
		engine.renderer.Render(engine.Render)
		engine.renderer.Start()
	} else {
		Init()
		for {
			engine.Update()
			time.Sleep(18 * time.Millisecond)
		}
	}
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

func (engine *EngineImpl) AddOrtho(spatial renderer.Spatial) {
	engine.orthoNode.Add(spatial)
}

func (engine *EngineImpl) RemoveOrtho(spatial renderer.Spatial, destroy bool) {
	engine.orthoNode.Remove(spatial, destroy)
}

func (engine *EngineImpl) AddSpatial(spatial renderer.Spatial) {
	engine.sceneGraph.Add(spatial)
}

func (engine *EngineImpl) RemoveSpatial(spatial renderer.Spatial, destroy bool) {
	engine.sceneGraph.Remove(spatial, destroy)
}

func (engine *EngineImpl) AddUpdatable(updatable Updatable) string {
	ids++
	key := fmt.Sprintf("%v", ids)
	engine.updatableStore.Add(key, updatable)
	return key
}

func (engine *EngineImpl) RemoveUpdatable(updatable Updatable) {
	engine.updatableStore.RemoveUpdatable(updatable)
}

func (engine *EngineImpl) AddUpdatableKey(key string, updatable Updatable) {
	engine.updatableStore.Add(key, updatable)
}

func (engine *EngineImpl) RemoveUpdatableKey(key string) {
	engine.updatableStore.Remove(key)
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

func NewHeadlessEngine() Engine {
	updatableStore := NewUpdatableStore()
	fpsMeter := renderer.CreateFPSMeter(1.0)
	fpsMeter.FpsCap = 60

	return &EngineImpl{
		fpsMeter:       fpsMeter,
		updatableStore: updatableStore,
	}
}
