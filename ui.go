package main

import (
	"image/color"
	"runtime"

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
		container := ui.NewContainer()

		imgElement := ui.NewImageElement(alienwareImg)
		container.AddChildren(imgElement)
		container.SetSize(400, 0)
		container.SetMargin(ui.NewMargin(15))
		container.SetPadding(ui.NewMargin(15))
		container.SetBackgroundColor(0, 255, 0, 255)
		imgElement = ui.NewImageElement(alienwareImg)
		imgElement.SetSize(300, 0)
		imgElement.Hitbox.AddOnHover(func() {
			imgElement.SetSize(350, 0)
			window.Render()
		})
		imgElement.Hitbox.AddOnUnHover(func() {
			imgElement.SetSize(300, 0)
			window.Render()
		})
		container.AddChildren(imgElement)
		container.AddChildren(ui.NewTextElement("test text text test text test text test text test text test text test text test text", color.Black, 32, textFont))
		window.SetElement(container)
		window.SetTranslation(vmath.Vector3{100, 100, 1})
		window.SetScale(vmath.Vector3{600, 700, 1})
		window.SetBackgroundColor(90, 0, 255, 255)
		gameEngine.AddOrtho(window)

		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)

		uiController := ui.NewUiController(window)
		controllerManager.AddController(uiController)

	})
}
