package editor

import (
	"bytes"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/effects"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
	"github.com/walesey/go-engine/util"
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
	selectSprite          renderer.Entity
}

func New() *Editor {
	return &Editor{
		uiAssets:    ui.NewHtmlAssets(),
		rootMapNode: renderer.NewNode(),
		currentMap: &editorModels.MapModel{
			Name: "default",
			Root: editorModels.NewNodeModel("root"),
		},
		mouseMode: "scale",
	}
}

func (e *Editor) Start() {
	glRenderer := &opengl.OpenglRenderer{
		WindowTitle: "GoEngine Editor",
	}
	e.renderer = glRenderer
	e.gameEngine = engine.NewEngine(e.renderer)

	e.gameEngine.Start(func() {

		shader, err := assets.ImportShader("build/shaders/pbr.vert", "build/shaders/pbr.frag")
		if err != nil {
			panic("error importing shader")
		}

		e.gameEngine.DefaultShader(shader)

		// Sky cubemap
		skyImg, err := assets.ImportImage("resources/cubemap.png")
		if err == nil {
			geom := renderer.CreateSkyBox()
			geom.Transform(mgl32.Scale3D(10000, 10000, 10000))
			skyNode := renderer.NewNode()
			skyNode.Material = renderer.NewMaterial(renderer.NewTexture("diffuseMap", skyImg))
			skyNode.RendererParams = renderer.NewRendererParams()
			skyNode.RendererParams.CullBackface = false
			skyNode.RendererParams.Unlit = true
			skyNode.Add(geom)
			e.gameEngine.AddSpatial(skyNode)
		}

		environmentCubemap := &renderer.Texture{
			TextureName: "environmentMap",
			CubeMap:     renderer.CreateCubemap(skyImg),
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

		e.initSelectSprite()
		e.gameEngine.AddUpdatable(engine.UpdatableFunc(e.updateSelectSprite))

		//editor controller
		e.controllerManager.AddController(NewEditorController(e).(glfwController.Controller))

		//custom controller
		e.customController = controller.CreateController()
		e.controllerManager.AddController(e.customController.(glfwController.Controller))

		e.setupUI()
	})
}

func (e *Editor) initSelectSprite() {
	img, _ := assets.DecodeImage(bytes.NewBuffer(util.Base64ToBytes(GeometryIconData)))
	mat := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img))
	selectSprite := effects.CreateSprite(1, 1, 1, mat)
	spriteNode := renderer.NewNode()
	spriteNode.RendererParams = &renderer.RendererParams{
		Unlit:        true,
		Transparency: renderer.EMISSIVE,
	}
	spriteNode.Add(selectSprite)
	e.selectSprite = spriteNode
	e.gameEngine.AddSpatialTransparent(selectSprite)
}

func (e *Editor) updateSelectSprite(dt float64) {
	selectedModel, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root)
	if selectedModel != nil {
		size := selectedModel.Translation.Sub(e.gameEngine.Camera().Translation).Len() * 0.02
		translation, err := e.rootMapNode.RelativePosition(selectedModel.GetNode())
		if err == nil {
			e.selectSprite.SetScale(mgl32.Vec3{size, size, size})
			e.selectSprite.SetTranslation(translation)
		} else {
			e.selectSprite.SetScale(mgl32.Vec3{})
		}
	}
}
