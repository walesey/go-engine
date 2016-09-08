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
	vmath "github.com/walesey/go-engine/vectormath"
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

		createWindow(uiController, populateContent)
	})
}

func createWindow(uiController controller.Controller, populate func(c *ui.Container) []ui.Activatable) {
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

func populateContent(c *ui.Container) []ui.Activatable {
	// example text title
	textElement := ui.NewTextElement("UI EXAMPLE!", color.Black, 16, nil)

	// example image element
	img, _ := assets.ImportImageCached("resources/cubemap.png")
	imageElement := ui.NewImageElement(img)

	// example text field
	text := ui.NewTextElement("", color.Black, 16, nil)
	text.SetPlaceholder("this is a placeholder")
	tf := textField(text)

	passwordText := ui.NewTextElement("", color.Black, 16, nil)
	passwordText.SetHidden(true)
	passwordTf := textField(passwordText)

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
	c.AddChildren(textElement, imageElement, tf, passwordTf, button)

	// return everything that should be included in the Tabs order
	return []ui.Activatable{text, passwordText}
}

func textField(textElem *ui.TextElement) *ui.Container {
	tf := ui.NewContainer()
	tf.SetHeight(16)
	tf.SetBackgroundColor(200, 200, 200, 254)
	tf.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
		if !release {
			textElem.Activate()
		}
	})
	return tf
}
