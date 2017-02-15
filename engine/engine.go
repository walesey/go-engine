package engine

import (
	"fmt"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
)

type Engine interface {
	Start(Init func())
	DefaultShader(shader *renderer.Shader)
	DefaultCubeMap(cubeMap *renderer.CubeMap)
	AddOrtho(spatial renderer.Spatial)
	AddSpatial(spatial renderer.Spatial)
	AddSpatialTransparent(spatial renderer.Spatial)
	RemoveSpatial(spatial renderer.Spatial, destroy bool)
	AddUpdatable(updatable Updatable)
	RemoveUpdatable(updatable Updatable)
	AddLight(light *renderer.Light)
	RemoveLight(light *renderer.Light)
	RequestAnimationFrame(cb func())
	Camera() *renderer.Camera
	Renderer() renderer.Renderer
	SetFpsCap(FpsCap float64)
	FPS() float64
	InitFpsDial()
	Update()
}

type EngineImpl struct {
	fpsMeter       *renderer.FPSMeter
	renderer       renderer.Renderer
	sceneGraph     *renderer.SceneGraph
	camera         *renderer.Camera
	updatableStore *UpdatableStore
	stepCounter    int64

	opaqueNode, transparentNode, orthoNode *renderer.Node
}

func (engine *EngineImpl) Start(Init func()) {
	if engine.renderer != nil {
		engine.renderer.SetInit(Init)
		engine.renderer.SetUpdate(engine.Update)
		engine.renderer.SetRender(engine.Render)
		engine.renderer.SetCamera(engine.camera)
		engine.renderer.Start()
	} else {
		Init()
		for {
			engine.Update()
		}
	}
}

func (engine *EngineImpl) DefaultShader(shader *renderer.Shader) {
	engine.opaqueNode.Shader, engine.transparentNode.Shader, engine.orthoNode.Shader = shader, shader, shader
}

func (engine *EngineImpl) DefaultCubeMap(cubeMap *renderer.CubeMap) {
	engine.opaqueNode.CubeMap, engine.transparentNode.CubeMap, engine.orthoNode.CubeMap = cubeMap, cubeMap, cubeMap
}

func (engine *EngineImpl) Update() {
	engine.stepCounter++
	dt := engine.fpsMeter.UpdateFPSMeter()
	engine.updatableStore.UpdateAll(dt)
}

func (engine *EngineImpl) Render() {
	engine.camera.Ortho = false
	engine.sceneGraph.RenderScene(engine.renderer, engine.camera.Translation)
	engine.camera.Ortho = true
	engine.orthoNode.Draw(engine.renderer, mgl32.Ident4())
}

func (engine *EngineImpl) AddOrtho(spatial renderer.Spatial) {
	engine.orthoNode.Add(spatial)
}

func (engine *EngineImpl) AddSpatial(spatial renderer.Spatial) {
	engine.opaqueNode.Add(spatial)
}

func (engine *EngineImpl) AddSpatialTransparent(spatial renderer.Spatial) {
	engine.transparentNode.Add(spatial)
}

func (engine *EngineImpl) RemoveSpatial(spatial renderer.Spatial, destroy bool) {
	engine.opaqueNode.Remove(spatial, destroy)
	engine.transparentNode.Remove(spatial, destroy)
	engine.orthoNode.Remove(spatial, destroy)
}

func (engine *EngineImpl) AddUpdatable(updatable Updatable) {
	engine.updatableStore.Add(updatable)
}

func (engine *EngineImpl) RemoveUpdatable(updatable Updatable) {
	engine.updatableStore.Remove(updatable)
}

func (engine *EngineImpl) AddLight(light *renderer.Light) {
	engine.renderer.AddLight(light)
}

func (engine *EngineImpl) RemoveLight(light *renderer.Light) {
	engine.renderer.RemoveLight(light)
}

func (engine *EngineImpl) RequestAnimationFrame(cb func()) {
	step := engine.stepCounter
	var updater Updatable
	updater = UpdatableFunc(func(dt float64) {
		if engine.stepCounter != step {
			engine.RemoveUpdatable(updater)
			cb()
		}
	})
	engine.AddUpdatable(updater)
}

func (engine *EngineImpl) Camera() *renderer.Camera {
	return engine.camera
}

func (engine *EngineImpl) Renderer() renderer.Renderer {
	return engine.renderer
}

func (engine *EngineImpl) SetFpsCap(FpsCap float64) {
	engine.fpsMeter.FpsCap = FpsCap
}

func (engine *EngineImpl) FPS() float64 {
	return engine.fpsMeter.Value()
}

func (engine *EngineImpl) InitFpsDial() {
	window := ui.NewWindow()
	window.SetTranslation(mgl32.Vec3{10, 10, 1})
	window.SetScale(mgl32.Vec3{400, 0, 1})
	window.SetBackgroundColor(0, 0, 0, 0)

	container := ui.NewContainer()
	container.SetBackgroundColor(0, 0, 0, 0)
	window.SetElement(container)

	text := ui.NewTextElement("0", color.RGBA{255, 0, 0, 255}, 18, nil)
	container.AddChildren(text)
	engine.AddUpdatable(UpdatableFunc(func(dt float64) {
		fps := engine.FPS()
		text.SetText(fmt.Sprintf("%v", int(fps)))
		switch {
		case fps < 20:
			text.SetTextColor(color.RGBA{255, 0, 0, 255})
		case fps < 30:
			text.SetTextColor(color.RGBA{255, 90, 0, 255})
		case fps < 50:
			text.SetTextColor(color.RGBA{255, 255, 0, 255})
		default:
			text.SetTextColor(color.RGBA{0, 255, 0, 255})
		}
		text.ReRender()
	}))

	window.Render()
	engine.AddOrtho(window)
}

func (engine *EngineImpl) initNodes() {
	engine.opaqueNode, engine.transparentNode, engine.orthoNode = renderer.NewNode(), renderer.NewNode(), renderer.NewNode()
}

func NewEngine(r renderer.Renderer) Engine {
	fpsMeter := renderer.CreateFPSMeter(1.0)
	fpsMeter.FpsCap = 144

	sceneGraph := renderer.CreateSceneGraph()
	updatableStore := NewUpdatableStore()
	camera := renderer.CreateCamera()

	engine := &EngineImpl{
		fpsMeter:       fpsMeter,
		sceneGraph:     sceneGraph,
		updatableStore: updatableStore,
		renderer:       r,
		camera:         camera,
	}

	engine.initNodes()
	sceneGraph.Add(engine.opaqueNode)

	engine.transparentNode.RendererParams = renderer.NewRendererParams()
	engine.transparentNode.RendererParams.CullBackface = false
	engine.transparentNode.RendererParams.DepthMask = false
	sceneGraph.AddTransparent(engine.transparentNode)

	engine.orthoNode.RendererParams = &renderer.RendererParams{
		Unlit:        true,
		DepthMask:    true,
		Transparency: renderer.NON_EMISSIVE,
	}
	return engine
}

func NewHeadlessEngine() Engine {
	updatableStore := NewUpdatableStore()
	fpsMeter := renderer.CreateFPSMeter(1.0)
	fpsMeter.FpsCap = 144

	engine := &EngineImpl{
		fpsMeter:       fpsMeter,
		updatableStore: updatableStore,
	}

	engine.initNodes()
	return engine
}
