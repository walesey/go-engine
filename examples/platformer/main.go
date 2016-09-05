package main

import (
	"image/color"
	"runtime"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/effects"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/physics/chipmunk"
	"github.com/walesey/go-engine/physics/physicsAPI"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

const characterSize = 40.0
const timePerAnimationFrame = 0.3
const floorHeight = 600

type Character struct {
	body       physicsAPI.PhysicsObject2D
	sprite     *effects.Sprite
	frameTimer float64
}

func (c *Character) Update(dt float64) {
	c.sprite.SetTranslation(c.body.GetPosition().ToVector3())

	// character animation
	c.frameTimer += dt
	if c.frameTimer >= timePerAnimationFrame {
		c.frameTimer -= timePerAnimationFrame
		c.sprite.NextFrame()
	}
}

func NewCharacter() *Character {
	characterImg, _ := assets.ImportImageCached("resources/stickman.png")
	characterMat := renderer.CreateMaterial()
	characterMat.Diffuse = characterImg
	characterMat.LightingMode = renderer.MODE_UNLIT
	sprite := effects.CreateSprite(4, 4, 1, characterMat)
	sprite.SetScale(vmath.Vector2{characterSize, characterSize}.ToVector3())
	sprite.SetTranslation(vmath.Vector2{X: 400, Y: 400}.ToVector3())

	body := chipmunkPhysics.NewChipmunkBody(1, 1)
	circle := chipmunk.NewCircle(vect.Vector_Zero, float32(characterSize*0.5))
	body.Body.AddShape(circle)

	return &Character{
		body:   body,
		sprite: sprite,
	}
}

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
}

func main() {

	//renderer and game engine
	glRenderer := &opengl.OpenglRenderer{
		WindowTitle:  "Platformer",
		WindowWidth:  800,
		WindowHeight: 800,
	}
	gameEngine := engine.NewEngine(glRenderer)

	// physics engine (Chipmonk)
	physicsSpace := chipmunkPhysics.NewChipmonkSpace()
	physicsSpace.SetGravity(vmath.Vector2{Y: 400})
	gameEngine.AddUpdatable(physicsSpace)

	gameEngine.Start(func() {

		// Sky cubemap
		skyImg, err := assets.ImportImageCached("resources/cubemap.png")
		if err == nil {
			gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
		}

		// The player object
		character := NewCharacter()
		character.body.SetPosition(vmath.Vector2{400, 400})

		// Add the character to all the things
		physicsSpace.AddBody(character.body)
		gameEngine.AddOrtho(character.sprite)
		gameEngine.AddUpdatable(character)

		// terrain
		terrainGeometry := renderer.CreateBoxWithOffset(800, 800-floorHeight, 0, floorHeight)
		terrainGeometry.SetColor(color.NRGBA{0, 254, 0, 254})
		gameEngine.AddOrtho(terrainGeometry)

		// terrain physics body
		terrainBody := chipmunkPhysics.NewChipmunkBodyStatic()
		segment := chipmunk.NewSegment(vect.Vect{0, floorHeight}, vect.Vect{4000, floorHeight}, 0)
		terrainBody.Body.AddShape(segment)
		physicsSpace.AddBody(terrainBody)

		// input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// key bindings
		customController := controller.CreateController()
		controllerManager.AddController(customController.(glfwController.Controller))

		// Walk
		customController.BindAction(func() {
			character.body.SetVelocity(vmath.Vector2{X: -200, Y: character.body.GetVelocity().Y})
		}, controller.KeyA, controller.Press)
		customController.BindAction(func() {
			character.body.SetVelocity(vmath.Vector2{X: 200, Y: character.body.GetVelocity().Y})
		}, controller.KeyD, controller.Press)

		//Stop walking
		customController.BindAction(func() {
			character.body.SetVelocity(vmath.Vector2{Y: character.body.GetVelocity().Y})
		}, controller.KeyA, controller.Release)
		customController.BindAction(func() {
			character.body.SetVelocity(vmath.Vector2{Y: character.body.GetVelocity().Y})
		}, controller.KeyD, controller.Release)

		// Jump
		customController.BindAction(func() {
			character.body.SetVelocity(character.body.GetVelocity().Add(vmath.Vector2{Y: -400}))
		}, controller.KeySpace, controller.Press)

	})
}
