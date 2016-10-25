package main

import (
	"image/color"
	
	"github.com/gopherjs/gopherjs/js"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/webgl"
	"github.com/walesey/go-engine/renderer"
)

//
func main() {
	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)

	//renderer and game engine
	webRenderer := webgl.NewWebRenderer(canvas)
	gameEngine := engine.NewEngine(webRenderer)

	gameEngine.Start(func() {
		gameEngine.InitFpsDial()

		//lighting
		webRenderer.CreateLight(
			0.3, 0.3, 0.3, //ambient
			0.5, 0.5, 0.5, //diffuse
			0.7, 0.7, 0.7, //specular
			false, mgl32.Vec3{0.7, 0.2, 0.7}, //position
			0, //index
		)

		// Create a red box geometry, attach to a node, add the node to the scenegraph
		boxGeometry := renderer.CreateBox(10, 10)
		boxGeometry.Material = renderer.CreateMaterial()
		boxGeometry.SetColor(color.NRGBA{254, 0, 0, 254})
		boxGeometry.CullBackface = false
		boxNode := renderer.CreateNode()
		boxNode.SetTranslation(mgl32.Vec3{30, 0})
		boxNode.Add(boxGeometry)
		gameEngine.AddSpatial(boxNode)

		// make the box spin
		var angle float32
		gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			angle += float32(dt)
			boxNode.SetOrientation(mgl32.QuatRotate(angle, mgl32.Vec3{0, 1, 0}))
		}))
	})
}

