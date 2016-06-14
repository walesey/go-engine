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
	gameEngine.InitFpsDial()

	physicsSpace := chipmunkPhysics.NewChipmonkSpace()
	physicsSpace.SetGravity(vmath.Vector2{0, -10})
	gameEngine.AddUpdatable(physicsSpace)

	geomMonkey, _ := assets.ImportObj("TestAssets/Files/physicsMonkey/phyMonkey.obj")

	geomMap, _ := assets.ImportObj("TestAssets/Files/CollisionMesh/Collision_Mesh.obj")
	tx := renderer.CreateTransform()
	tx.From(vmath.Vector3{1.0, 1.0, 1.0}, vmath.Vector3{0.0, -1.0, 0.0}, vmath.IdentityQuaternion())
	geomMap.Transform(tx)
	segmentsMap := assets.IntersectGeometry(geomMap)

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
		freeMoveActor.Location = vmath.Vector3{0, 140, 0}
		mainController := controller.NewBasicMovementController(freeMoveActor, true)
		controllerManager.AddController(mainController.(glfwController.Controller))
		gameEngine.AddUpdatable(freeMoveActor)

		applySegmentShapes := func(body *chipmunk.Body) {
			for _, seg := range segmentsMap {
				segStart := vect.Vect{vect.Float(seg[0].X), vect.Float(seg[0].Y)}
				segFinish := vect.Vect{vect.Float(seg[1].X), vect.Float(seg[1].Y)}
				segment := chipmunk.NewSegment(segStart, segFinish, 0)
				segment.SetElasticity(0.6)
				body.AddShape(segment)
			}
		}

		// Static objects
		mapNode := renderer.CreateNode()
		mapNode.Add(geomMap)
		gameEngine.AddSpatial(mapNode)

		mapBody := chipmunkPhysics.NewChipmunkBodyStatic()
		applySegmentShapes(mapBody.Body)
		physicsSpace.AddBody(mapBody)

		physicsActor := actor.NewPhysicsActor2D(mapNode, mapBody, vmath.Vector3{1, 0, 1})
		gameEngine.AddUpdatable(physicsActor)

		// Dynamic objects
		spawn := func(pos float64) {
			monkeyNode := renderer.CreateNode()
			monkeyNode.Add(geomMonkey)
			gameEngine.AddSpatial(monkeyNode)

			monkeyBody := chipmunkPhysics.NewChipmunkBody(1, 1)
			circle := chipmunk.NewCircle(vect.Vector_Zero, 1.0)
			circle.SetElasticity(0.95)
			monkeyBody.Body.AddShape(circle)
			physicsSpace.AddBody(monkeyBody)
			monkeyBody.SetPosition(vmath.Vector2{pos, 20 + pos})

			physicsActor := actor.NewPhysicsActor2D(monkeyNode, monkeyBody, vmath.Vector3{1, 0, 1})
			gameEngine.AddUpdatable(physicsActor)

		}

		for i := 0; i < 200; i++ {
			spawn(float64(i))
		}

	})
}
