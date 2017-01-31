package main

import (
	"bytes"
	"fmt"
	"go/build"
	"image/color"
	"os"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
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
	// set working dir to access assets
	p, _ := build.Import("github.com/walesey/go-engine", "", build.FindOnly)
	os.Chdir(p.Dir)
}

func main() {

	//renderer and game engine
	glRenderer := opengl.NewOpenglRenderer("ui", 800, 800, false)
	gameEngine := engine.NewEngine(glRenderer)
	gameEngine.InitFpsDial()

	gameEngine.Start(func() {
		// load in ui shader
		shader := renderer.NewShader()
		shader.FragSrc = uiFragmentShader
		shader.VertSrc = uiVertexShader
		gameEngine.DefaultShader(shader)

		// input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// create windows with content containers
		window1, content1 := createWindow(controllerManager)
		window2, content2 := createWindow(controllerManager)

		// populate content and set window tab order
		window1.Tabs = populateContent(content1)
		window2.Tabs = htmlContent(content2)

		// position windows
		window1.SetTranslation(mgl32.Vec2{50, 50}.Vec3(0))
		window2.SetTranslation(mgl32.Vec2{450, 50}.Vec3(0))

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
	window.SetScale(mgl32.Vec2{300, 0}.Vec3(0))

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
	button.Hitbox.AddOnClick(func(button int, release bool, position mgl32.Vec2) {
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

const uiVertexShader = `
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 texCoord;
in vec4 color;

out vec2 fragTexCoord;
out vec4 fragColor;

void main() {
	fragTexCoord = texCoord;
	fragColor = color;
	gl_Position = projection * camera * model * vec4(vert, 1);
}
`

const uiFragmentShader = `
#version 330

uniform bool useTextures;
uniform sampler2D diffuseMap;

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 outputColor;

void main() {
	
	if (useTextures) {
  	outputColor = texture(diffuseMap, fragTexCoord) * fragColor;
	} else {
		outputColor = fragColor;
	}
}
`
