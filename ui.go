package main

import (
	"runtime"

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
	glRenderer := &renderer.OpenglRenderer{WindowTitle: "GoEngine"}
	gameEngine := engine.NewEngine(glRenderer)

	// alienwareImg := assets.ImportImage("TestAssets/test.png")

	gameEngine.Start(func() {

		window := ui.NewWindow()
		mainContainer := ui.NewContainer()

		tab := ui.NewContainer()
		tab.SetBackgroundColor(120, 120, 120, 255)
		tab.SetHeight(40)

		container := ui.NewContainer()
		// container.SetMargin(ui.NewMargin(30))
		container.SetPadding(ui.NewMargin(30))
		container.SetBackgroundColor(0, 255, 0, 255)
		mainContainer.AddChildren(tab, container)

		innerContainer := ui.NewContainer()
		innerContainer.SetPadding(ui.NewMargin(30))
		innerContainer.SetBackgroundColor(255, 0, 0, 255)
		innerContainer.SetWidth(350)
		innerContainer.SetHeight(350)
		innerContainer.Hitbox.AddOnHover(func() {
			innerContainer.SetWidth(350)
		})
		innerContainer.Hitbox.AddOnUnHover(func() {
			innerContainer.SetWidth(300)
		})
		container.AddChildren(innerContainer)

		// imgElement := ui.NewImageElement(alienwareImg)
		// imgElement.SetWidth(200)
		// innerContainer.AddChildren(imgElement)
		// imgElement = ui.NewImageElement(alienwareImg)
		// imgElement.SetWidth(100)
		// innerContainer.AddChildren(imgElement)

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
		ui.ClickAndDragWindow(window, tab.Hitbox, customController)
	})
}
