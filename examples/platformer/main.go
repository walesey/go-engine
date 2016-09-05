package main

import (
	"image/color"
	"runtime"
	"time"

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
	walking    bool
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
	characterMat := assets.CreateMaterial(characterImg, nil, nil, nil)
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
	gameEngine.InitFpsDial()
	gameEngine.SetFpsCap(120)

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

		// create and manage a new particle system
		spawnParticles := func(load func() *effects.ParticleSystem, position vmath.Vector2) *effects.ParticleSystem {
			particles := load()
			particles.SetTranslation(position.ToVector3())
			gameEngine.AddOrtho(particles)
			age := 0.0
			gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
				particles.Update(dt)
				age += dt
				if age >= 0.2 {
					particles.DisableSpawning = true
				}
				if age >= 10 {
					gameEngine.RemoveOrtho(particles, true)
				}
			}))
			return particles
		}

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
			character.body.SetVelocity(character.body.GetVelocity().Add(vmath.Vector2{Y: -300}))
			spawnParticles(dustParticles, character.body.GetPosition().Add(vmath.Vector2{Y: 0.5 * characterSize}))
		}, controller.KeySpace, controller.Press)

		// shoot
		customController.BindMouseAction(func() {
			spawnParticles(majicParticles, character.body.GetPosition())
		}, controller.MouseButtonLeft, controller.Press)

		// create sparks when player collides with somthing
		var collisionTimer time.Time
		physicsSpace.SetOnCollision(func(shapeA, shapeB *chipmunk.Shape) {
			if time.Since(collisionTimer) > 200*time.Millisecond {
				spawnParticles(sparkParticles, character.body.GetPosition().Add(vmath.Vector2{Y: 0.5 * characterSize}))
			}
			collisionTimer = time.Now()
		})
	})
}

func dustParticles() *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("resources/smoke.png")
	material := assets.CreateMaterial(img, nil, nil, nil)
	material.LightingMode = renderer.MODE_UNLIT
	particleSystem := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     5,
		ParticleEmitRate: 20,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		Material:         material,
		TotalFrames:      64,
		FramesX:          8,
		FramesY:          8,
		MaxLife:          3.3,
		MinLife:          2.7,
		StartColor:       color.NRGBA{254, 254, 254, 120},
		EndColor:         color.NRGBA{254, 254, 254, 0},
		StartSize:        vmath.Vector2{40, 40}.ToVector3(),
		EndSize:          vmath.Vector2{180, 180}.ToVector3(),
		MinTranslation:   vmath.Vector2{-2, -2}.ToVector3(),
		MaxTranslation:   vmath.Vector2{2, 2}.ToVector3(),
		MinStartVelocity: vmath.Vector2{-5.0, -5.0}.ToVector3(),
		MaxStartVelocity: vmath.Vector2{5.0, 5.0}.ToVector3(),
		MaxRotation:      -3.14,
		MinRotation:      3.14,
	})
	particleSystem.FaceCamera = false
	return particleSystem
}

func majicParticles() *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("resources/majic.png")
	material := assets.CreateMaterial(img, nil, nil, nil)
	material.LightingMode = renderer.MODE_UNLIT
	particleSystem := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     300,
		ParticleEmitRate: 700,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		Material:         material,
		TotalFrames:      9,
		FramesX:          3,
		FramesY:          3,
		MaxLife:          1.3,
		MinLife:          0.7,
		StartColor:       color.NRGBA{254, 254, 254, 254},
		EndColor:         color.NRGBA{254, 254, 254, 254},
		StartSize:        vmath.Vector2{10, 10}.ToVector3(),
		EndSize:          vmath.Vector2{10, 10}.ToVector3(),
		MinTranslation:   vmath.Vector2{-1, -1}.ToVector3(),
		MaxTranslation:   vmath.Vector2{1, 1}.ToVector3(),
		MinStartVelocity: vmath.Vector2{-170, -170}.ToVector3(),
		MaxStartVelocity: vmath.Vector2{170, 170}.ToVector3(),
		MaxRotation:      -3.14,
		MinRotation:      3.14,
	})
	particleSystem.FaceCamera = false
	return particleSystem
}

func sparkParticles() *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("resources/spark.png")
	material := assets.CreateMaterial(img, nil, nil, nil)
	material.LightingMode = renderer.MODE_UNLIT
	particleSystem := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     80,
		ParticleEmitRate: 400,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		Material:         material,
		TotalFrames:      1,
		FramesX:          1,
		FramesY:          1,
		MaxLife:          3.3,
		MinLife:          1.7,
		StartColor:       color.NRGBA{254, 160, 90, 254},
		EndColor:         color.NRGBA{254, 160, 90, 254},
		StartSize:        vmath.Vector2{2, 2}.ToVector3(),
		EndSize:          vmath.Vector2{2, 2}.ToVector3(),
		MinTranslation:   vmath.Vector2{-5, -5}.ToVector3(),
		MaxTranslation:   vmath.Vector2{5, 5}.ToVector3(),
		MinStartVelocity: vmath.Vector2{-100, -100}.ToVector3(),
		MaxStartVelocity: vmath.Vector2{100, 100}.ToVector3(),
		Acceleration:     vmath.Vector2{Y: 400}.ToVector3(),
		OnParticleUpdate: func(p *effects.Particle) {
			p.Scale = p.Scale.Multiply(vmath.Vector2{X: 1 + p.Velocity.Length()*0.05, Y: 1}.ToVector3())
			p.Orientation = vmath.BetweenVectors(vmath.Vector3{X: 1}, p.Velocity)
		},
	})
	particleSystem.FaceCamera = false
	return particleSystem
}
