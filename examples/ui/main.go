package main

import (
	"bytes"
	"fmt"
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

/*
This example renders 2 windows
The left window uses the ui API (func populateContent)
The right window uses the html/css parser (func htmlContent)
Both are exactly the same but show the 2 methods
*/
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

		// create windows with content containers
		window1, content1 := createWindow(controllerManager)
		window2, content2 := createWindow(controllerManager)

		// populate content and set window tab order
		window1.Tabs = populateContent(content1)
		window2.Tabs = htmlContent(content2)

		// position windows
		window1.SetTranslation(vmath.Vector2{X: 50, Y: 50}.ToVector3())
		window2.SetTranslation(vmath.Vector2{X: 450, Y: 50}.ToVector3())

		// Add the windows to the engine
		gameEngine.AddOrtho(window1)
		gameEngine.AddOrtho(window2)

		// render all windows
		window1.Render()
		window2.Render()
	})
}

func createWindow(controllerManager *glfwController.ControllerManager) (window *ui.Window, content *ui.Container) {
	// Create window with size
	window = ui.NewWindow()
	window.SetScale(vmath.Vector2{X: 300}.ToVector3())

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

	return
}

func populateContent(content *ui.Container) []ui.Activatable {
	// example text title
	textElement := ui.NewTextElement("UI EXAMPLE!", color.Black, 16, nil)

	// example image element
	img, _ := assets.ImportImageCached("resources/cubemap.png")
	imageElement := ui.NewImageElement(img)
	imageElement.SetWidth(200)

	// example text field
	tf := ui.NewTextField("", color.Black, 16, nil)
	tf.SetPlaceholder("this is a placeholder")
	tf.SetBackgroundColor(255, 255, 255, 255)
	tf.SetMargin(ui.Margin{10, 0, 10, 0})

	// example hidden text field
	passwordTf := ui.NewTextField("", color.Black, 16, nil)
	passwordTf.SetHidden(true)
	passwordTf.SetBackgroundColor(255, 255, 255, 255)
	passwordTf.SetMargin(ui.Margin{0, 0, 10, 0})

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
	content.AddChildren(textElement, imageElement, tf, passwordTf, button)

	// return everything that should be included in the Tabs order
	return []ui.Activatable{tf, passwordTf}
}

func htmlContent(content *ui.Container) []ui.Activatable {
	html := bytes.NewBufferString(`
		<body>
			<div class=content>
				<h1 id=heading>UI EXAMPLE!</h1>
				<img src=testImage></img>
				<input type=text placeholder="this is a placeholder"></input>
				<input type=password></input>
				<button onclick=clickButton></button>
			<div>
		</body>
	`)

	css := bytes.NewBufferString(`
		.content img {
			width: 200;
		}

		.content input {
			background-color: #fff;
			margin: 10 0 0 0;
		}

		.content button {
			padding: 20;
			margin: 10 0 0 0;
			background-color: #a00;
		}
		
		.content button:hover {
			background-color: #e99;
		}
	`)

	// create assets
	htmlAssets := ui.NewHtmlAssets()

	// image
	img, _ := assets.ImportImageCached("resources/cubemap.png")
	htmlAssets.AddImage("testImage", img)

	// button click callback
	htmlAssets.AddCallback("clickButton", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // on release
			content.TextElementById("heading").SetText("release").SetTextColor(color.NRGBA{254, 0, 0, 254}).ReRender()
		} else {
			content.TextElementById("heading").SetText("press").SetTextColor(color.NRGBA{0, 254, 0, 254}).ReRender()
		}
	})

	// Render the html/css code to the content container
	activatables, err := ui.LoadHTML(content, html, css, htmlAssets)
	if err != nil {
		fmt.Println("Error loading html: ", err)
	}

	return activatables
}
