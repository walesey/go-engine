package main

import (
	"image/color"
	"runtime"

	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/ui"
)

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
}

func main() {

	//renderer and game engine
	glRenderer := &opengl.OpenglRenderer{
		WindowTitle:  "ui",
		WindowWidth:  800,
		WindowHeight: 800,
	}
	gameEngine := engine.NewEngine(glRenderer)

	gameEngine.Start(func() {

		// input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)
		uiController := controller.CreateController()
		controllerManager.AddController(uiController.(glfwController.Controller))

		createWindow(populateContent)
	})
}

func createWindow(populate func(c *ui.Container) []Activatable) {
	// Create a ui window with a click and drag tab
	window := ui.NewWindow()
	mainContainer := ui.NewContainer()
	window.SetElement(mainContainer)

	tab := ui.NewContainer()
	tab.SetBackgroundColor(70, 70, 170, 255)
	tab.SetHeight(40)

	content := ui.NewContainer()
	content.SetBackgroundColor(200, 200, 200, 255)
	content.SetPadding(ui.NewMargin(10))

	mainContainer.AddChildren(tab, content)
	ui.ClickAndDragWindow(window, tab.Hitbox, uiController)

	// Make sure all text fields get deselected on click
	uiController.BindMouseAction(func() {
		ui.DeactivateAllTextElements(content)
	}, controller.MouseButton1, controller.Press)

	window.Tabs = populate(content)
}

func populateContent(c *ui.Container) []Activatable {
	// example text title
	textElement := ui.NewTextElement("UI EXAMPLE!", color.Black, 16, nil)

	// example image element
	skyImg, _ := assets.ImportImageCached("resources/cubemap.png")
	imageElement := ui.NewImageElement(img)

	// onClickHandler for text fields
	activateOnClick := func(button int, release bool, position vmath.Vector2) {
		if !release {
			textField.Activate()
		}
	}

	// example text field
	textField := ui.NewTextElement("", color.Black, 16, nil)
	textField.SetPlaceholder("this is a placeholder")
	textField.Hitbox.AddOnClick(activateOnClick)

	passwordField := ui.NewTextElement("", color.Black, 16, nil)
	passwordField.SetHidden(true)
	passwordField.Hitbox.AddOnClick(activateOnClick)

	// example button
	button := ui.NewContainer()
	button.SetBackgroundColor(160, 0, 0, 254)
	button.SetPadding(ui.NewMargin(20))
	// button on click event
	button.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
		if release {
			textElement.SetText("release").SetTextColor(color.NRGBA{254, 0, 0, 254}).ReRender()
		} else {
			textElement.SetText("click").SetTextColor(color.NRGBA{0, 254, 0, 254}).ReRender()
		}
	})
	// button on hover event
	button.Hitbox.AddOnHover(func() {
		button.SetBackgroundColor(210, 100, 100, 254)
	})
	button.Hitbox.AddOnUnHover(func() {
		button.SetBackgroundColor(160, 0, 0, 254)
	})

	// add everything to the content container
	c.AddChildren(textElement, imageElement, textField, passwordField, button)

	// return everything that should be included in the Tabs order
	return []Activatable{textField, passwordField}
}
