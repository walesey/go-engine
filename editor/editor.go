package editor

import (
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Editor struct {
	renderer              renderer.Renderer
	gameEngine            engine.Engine
	currentMap            *editorModels.MapModel
	rootMapNode           *renderer.Node
	customController      controller.Controller
	controllerManager     *glfwController.ControllerManager
	uiAssets              ui.HtmlAssets
	mainMenu              *ui.Window
	mainMenuController    glfwController.Controller
	overviewMenu          *Overview
	mainMenuOpen          bool
	progressBar           *ui.Window
	fileBrowser           *FileBrowser
	fileBrowserOpen       bool
	fileBrowserController glfwController.Controller
	mouseMode             string
}

func New() *Editor {
	return &Editor{
		uiAssets:    ui.NewHtmlAssets(),
		rootMapNode: renderer.CreateNode(),
		currentMap: &editorModels.MapModel{
			Name: "default",
			Root: editorModels.NewNodeModel("root"),
		},
		mouseMode: "scale",
	}
}

func (e *Editor) Start() {
	glRenderer := &opengl.OpenglRenderer{
		WindowTitle: "GoEngine editor",
	}
	e.renderer = glRenderer
	e.gameEngine = engine.NewEngine(e.renderer)

	e.gameEngine.Start(func() {
		//lighting
		e.renderer.CreateLight(0.0, 0.0, 0.0, 0.5, 0.5, 0.5, 0.7, 0.7, 0.7, true, vmath.Vector3{0.5, -1, 0.3}, 0)

		//Sky
		skyImg, err := assets.ImportImage("TestAssets/Files/skybox/cubemap.png")
		if err == nil {
			e.gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
		}

		//root node
		e.gameEngine.AddSpatial(e.rootMapNode)

		//input/controller manager
		e.controllerManager = glfwController.NewControllerManager(glRenderer.Window)

		//camera + player
		camera := e.gameEngine.Camera()
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.MoveSpeed = 20.0
		freeMoveActor.LookSpeed = 0.002
		mainController := controller.NewBasicMovementController(freeMoveActor, true)
		e.controllerManager.AddController(mainController.(glfwController.Controller))
		e.gameEngine.AddUpdatable(freeMoveActor)

		//editor controller
		e.controllerManager.AddController(NewEditorController(e).(glfwController.Controller))

		//custom controller
		e.customController = controller.CreateController()
		e.controllerManager.AddController(e.customController.(glfwController.Controller))

		e.setupUI()
	})
}
