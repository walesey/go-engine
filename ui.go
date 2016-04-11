package main

import (
	"image/color"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	glRenderer := renderer.NewOpenglRenderer("GoEngine", 1900, 1000)
	gameEngine := engine.NewEngine(glRenderer)

	assetLib, err := assets.LoadAssetLibrary("TestAssets/ui.asset")
	if err != nil {
		panic(err)
	}
	alienwareImg := assetLib.GetImage("alienware")
	textFont, _ := ui.LoadFont("TestAssets/Audiowide-Regular.ttf")

	gameEngine.Start(func() {

		window := ui.NewWindow()
		mainContainer := ui.NewContainer()

		tab := ui.NewContainer()
		tab.SetBackgroundColor(120, 120, 120, 255)
		tab.SetSize(0, 40)

		container := ui.NewContainer()
		container.SetMargin(ui.NewMargin(15))
		container.SetPadding(ui.NewMargin(15))
		container.SetBackgroundColor(0, 255, 0, 255)
		mainContainer.AddChildren(tab, container)

		imgElement := ui.NewImageElement(alienwareImg)
		container.AddChildren(imgElement)
		imgElement = ui.NewImageElement(alienwareImg)
		imgElement.SetSize(300, 0)
		imgElement.Hitbox.AddOnHover(func() {
			imgElement.SetSize(350, 0)
		})
		imgElement.Hitbox.AddOnUnHover(func() {
			imgElement.SetSize(300, 0)
		})
		container.AddChildren(imgElement)

		textElement := ui.NewTextElement("test", color.Black, 32, textFont)
		textElement.Activate()
		container.AddChildren(textElement)

		window.SetElement(mainContainer)
		window.SetScale(vmath.Vector3{600, 700, 1})
		window.SetBackgroundColor(90, 0, 255, 255)
		gameEngine.AddOrtho(window)

		gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			window.Render()
		}))

		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)

		uiController := ui.NewUiController(window)
		controllerManager.AddController(uiController)

		//custom controller
		customController := controller.NewActionMap()
		controllerManager.AddController(customController)

		//Click and drag window
		grabbed := false
		grabOffset := vmath.Vector2{}
		tab.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
			grabOffset = position
			grabbed = !release
		})
		customController.BindMouseAction(func() {
			grabbed = false
		}, glfw.MouseButton1, glfw.Release)
		customController.BindAxisAction(func(xpos, ypos float64) {
			if grabbed {
				position := vmath.Vector2{xpos, ypos}
				window.SetTranslation(position.Subtract(grabOffset).ToVector3())
			}
		})

	})
}
