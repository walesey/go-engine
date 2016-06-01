package examples

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/physics/chipmunk"
	"github.com/walesey/go-engine/renderer"

	"github.com/codegangsta/cli"
	vmath "github.com/walesey/go-engine/vectormath"
)

//
func Chipmunk(c *cli.Context) {
	glRenderer := &opengl.OpenglRenderer{
		WindowTitle:  "GoEngine",
		WindowWidth:  1900,
		WindowHeight: 1000,
	}
	gameEngine := engine.NewEngine(glRenderer)
	physicsSpace := chipmunkPhysics.NewChipmonkSpace()
	physicsSpace.SetGravity(vmath.Vector2{0, -10})
	gameEngine.AddUpdatable(physicsSpace)

	geomMonkey, _ := assets.ImportObj("TestAssets/Files/physicsMonkey/phyMonkey.obj")
	segmentsMonkey := assets.IntersectGeometry(geomMonkey)

	gameEngine.Start(func() {
		//lighting
		glRenderer.CreateLight(0.0, 0.0, 0.0, 0.5, 0.5, 0.5, 0.7, 0.7, 0.7, true, vmath.Vector3{0.3, -1, 0.2}, 0)

		//Sky
		skyImg, err := assets.ImportImage("TestAssets/Files/skybox/cubemap.png")
		if err == nil {
			gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
		}

		//input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		//camera + player
		camera := gameEngine.Camera()
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.MoveSpeed = 20.0
		freeMoveActor.LookSpeed = 0.002
		mainController := controller.NewBasicMovementController(freeMoveActor, true)
		controllerManager.AddController(mainController.(glfwController.Controller))
		gameEngine.AddUpdatable(freeMoveActor)

		// Static objects
		monkeyNode := renderer.CreateNode()
		monkeyNode.Add(geomMonkey)
		gameEngine.AddSpatial(monkeyNode)

		monkeyBody := chipmunkPhysics.NewChipmunkBodyStatic()
		for _, seg := range segmentsMonkey {
			segStart := vect.Vect{vect.Float(seg[0].X), vect.Float(seg[0].Y)}
			segFinish := vect.Vect{vect.Float(seg[1].X), vect.Float(seg[1].Y)}
			segment := chipmunk.NewSegment(segStart, segFinish, 0)
			segment.SetElasticity(0.6)
			monkeyBody.Body.AddShape(segment)
		}
		physicsSpace.AddBody(monkeyBody)
		monkeyBody.SetPosition(vmath.Vector2{0, -20})

		physicsActor := actor.NewPhysicsActor2D(monkeyNode, monkeyBody, vmath.Vector3{1, 0, 1})
		gameEngine.AddUpdatable(physicsActor)

		// Dynamic objects
		spawn := func() {
			monkeyNode := renderer.CreateNode()
			monkeyNode.Add(geomMonkey)
			gameEngine.AddSpatial(monkeyNode)

			monkeyBody := chipmunkPhysics.NewChipmunkBody(1, 1)
			for _, seg := range segmentsMonkey {
				segStart := vect.Vect{vect.Float(seg[0].X), vect.Float(seg[0].Y)}
				segFinish := vect.Vect{vect.Float(seg[1].X), vect.Float(seg[1].Y)}
				segment := chipmunk.NewSegment(segStart, segFinish, 0)
				segment.SetElasticity(0.6)
				monkeyBody.Body.AddShape(segment)
			}
			physicsSpace.AddBody(monkeyBody)

			physicsActor := actor.NewPhysicsActor2D(monkeyNode, monkeyBody, vmath.Vector3{1, 0, 1})
			gameEngine.AddUpdatable(physicsActor)

		}

		spawn()
	})
}
