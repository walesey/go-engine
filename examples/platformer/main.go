package main

import (
	"runtime"

	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/effects"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
}

func main() {

	glRenderer := &opengl.OpenglRenderer{
		WindowTitle:  "Platformer",
		WindowWidth:  800,
		WindowHeight: 800,
	}
	gameEngine := engine.NewEngine(glRenderer)

	gameEngine.Start(func() {

		// Sky cubemap
		skyImg, err := assets.ImportImageCached("resources/cubemap.png")
		if err == nil {
			gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
		}

		// input/controller manager
		// controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// The player sprite
		characterImg, _ := assets.ImportImageCached("resources/stickman.png")
		characterMat := renderer.CreateMaterial()
		characterMat.Diffuse = characterImg
		characterMat.LightingMode = renderer.MODE_UNLIT
		characterSprite := effects.CreateSprite(4, 4, 1, characterMat)
		characterSprite.SetScale(vmath.Vector2{40, 40}.ToVector3())
		characterSprite.SetTranslation(vmath.Vector2{X: 400, Y: 400}.ToVector3())
		gameEngine.AddOrtho(characterSprite)

		// character animation
		timePerFrame := 0.3
		var frameTimer float64
		gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			frameTimer += dt
			if frameTimer >= timePerFrame {
				frameTimer -= timePerFrame
				characterSprite.NextFrame()
			}
		}))
	})
}
