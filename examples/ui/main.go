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

		// create a window with a content container
		window, content := createWindow(controllerManager)

		// populate content and set window tab order
		window.Tabs = populateContent(content)

		// Add the window to the engine
		gameEngine.AddOrtho(window)

		// render all window content
		window.Render()
	})
}

func createWindow(controllerManager *glfwController.ControllerManager) (window *ui.Window, content *ui.Container) {
	// Create window with size and position
	window = ui.NewWindow()
	window.SetScale(vmath.Vector2{X: 300}.ToVector3())
	window.SetTranslation(vmath.Vector2{X: 50, Y: 50}.ToVector3())

	// create a click and drag tab
	tab := ui.NewContainer()
	tab.SetBackgroundColor(70, 70, 170, 255)
	tab.SetHeight(40)

	// create a content container
	content = ui.NewContainer()
	content.SetBackgroundColor(200, 200, 200, 255)
	content.SetPadding(ui.NewMargin(10))

	// Add all the containers to the window
	mainContainer := ui.NewContainer()
	mainContainer.AddChildren(tab, content)
	window.SetElement(mainContainer)

	// create uiController
	uiController := ui.NewUiController(window)
	controllerManager.AddController(uiController.(glfwController.Controller))
	ui.ClickAndDragWindow(window, tab.Hitbox, uiController)

	// Make sure all text fields get deselected on click
	uiController.BindMouseAction(func() {
		ui.DeactivateAllTextElements(content)
	}, controller.MouseButton1, controller.Press)
	return
}

func populateContent(c *ui.Container) []ui.Activatable {
	// example text title
	textElement := ui.NewTextElement("UI EXAMPLE!", color.Black, 16, nil)

	// example image element
	img, _ := assets.ImportImageCached("resources/cubemap.png")
	imageElement := ui.NewImageElement(img)
	imageElement.SetWidth(200)

	// example text field
	tf := ui.NewTextElement("", color.Black, 16, nil)
	tf.SetPlaceholder("this is a placeholder")

	// example hidden text field
	passwordTf := ui.NewTextElement("", color.Black, 16, nil)
	passwordTf.SetHidden(true)

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
	c.AddChildren(textElement, imageElement, textField(tf), textField(passwordTf), button)

	// return everything that should be included in the Tabs order
	return []ui.Activatable{tf, passwordTf}
}

func textField(textElem *ui.TextElement) *ui.Container {
	tf := ui.NewContainer()
	tf.SetHeight(26)
	tf.SetBackgroundColor(200, 200, 200, 254)
	tf.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
		if !release {
			textElem.Activate()
		}
	})
	return tf
}
