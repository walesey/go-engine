package editor

import (
	"github.com/luxengine/gobullet"
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/physics/bullet"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Editor struct {
	renderer          renderer.Renderer
	gameEngine        engine.Engine
	customController  *controller.ActionMap
	controllerManager *controller.ControllerManager
	uiAssets          ui.HtmlAssets
	mainMenu          *ui.Window
	mainMenuOpen      bool
}

func New() *Editor {
	return &Editor{}
}

func (e *Editor) Start() {
	glRenderer := &renderer.OpenglRenderer{
		WindowTitle: "GoEngine editor",
	}
	e.renderer = glRenderer
	e.gameEngine = engine.NewEngine(e.renderer)

	//physics engine
	sdk := gobullet.NewBulletSDK()
	defer sdk.Delete()
	physicsWorld := bullet.NewBtDynamicsWorld(sdk)
	defer physicsWorld.Delete()
	physicsWorld.SetGravity(vmath.Vector3{0, -10, 0})
	e.gameEngine.AddUpdatable(physicsWorld)

	e.gameEngine.Start(func() {
		//lighting
		e.renderer.CreateLight(0.0, 0.0, 0.0, 0.5, 0.5, 0.5, 0.7, 0.7, 0.7, true, vmath.Vector3{0.3, -1, 0.2}, 0)

		//input/controller manager
		e.controllerManager = controller.NewControllerManager(glRenderer.Window)

		//camera + player
		camera := e.gameEngine.Camera()
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.MoveSpeed = 3.0
		mainController := controller.NewBasicMovementController(freeMoveActor)
		e.controllerManager.AddController(mainController)

		//custom controller
		e.customController = controller.NewActionMap()
		e.controllerManager.AddController(e.customController)

		e.setupUI()

		//event loop
		e.gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			//TODO:
		}))
	})
}
