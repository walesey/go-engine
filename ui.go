package main

import (
	"runtime"

	"github.com/walesey/go-engine/assets"
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

	gameEngine.Start(func() {
		// glRenderer.BackGroundColor(1.0, 1.0, 1.0, 1.0)

		window := ui.NewWindow()
		container := ui.NewContainer()

		imgElement := ui.NewImageElement(alienwareImg)
		imgElement.SetSize(500, 0)
		container.AddChildren(imgElement)
		imgElement = ui.NewImageElement(alienwareImg)
		imgElement.SetSize(300, 0)
		container.AddChildren(imgElement)
		container.HorizontalAlign = true
		window.SetElement(container)
		window.SetScale(vmath.Vector3{1000, 1000, 1})
		gameEngine.AddOrtho(window)

	})
}
