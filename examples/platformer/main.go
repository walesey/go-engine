package main

import (
	"image/color"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/go-gl/mathgl/mgl32"
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
	c.sprite.SetTranslation(c.body.GetPosition().Vec3(0))

	// character animation
	c.frameTimer += dt
	if c.frameTimer >= timePerAnimationFrame {
		c.frameTimer -= timePerAnimationFrame
		c.sprite.NextFrame()
	}
}

func NewCharacter(shader *renderer.Shader) *Character {
	characterImg, _ := assets.ImportImageCached("resources/stickman.png")
	characterMat := renderer.NewMaterial(renderer.NewTexture("diffuseMap", characterImg))
	sprite := effects.CreateSprite(4, 4, 1, characterMat)
	sprite.FaceCamera = false
	sprite.SetScale(mgl32.Vec2{characterSize, characterSize}.Vec3(0))
	sprite.SetTranslation(mgl32.Vec2{400, 400}.Vec3(0))

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
	physicsSpace.SetGravity(mgl32.Vec2{0, 400})
	gameEngine.AddUpdatable(physicsSpace)

	gameEngine.Start(func() {

		shader := renderer.NewShader()
		vertsrc, err := ioutil.ReadFile("build/shaders/basic.vert")
		if err != nil {
			panic(err)
		}
		shader.VertSrc = string(vertsrc)

		fragsrc, err := ioutil.ReadFile("build/shaders/basic.frag")
		if err != nil {
			panic(err)
		}
		shader.FragSrc = string(fragsrc)
		glRenderer.SetDefaultShader(shader)

		skyImg, err := assets.ImportImage("resources/cubemap.png")
		if err == nil {
			geom := renderer.CreateSkyBox()
			geom.Material = renderer.NewMaterial(renderer.NewTexture("diffuseMap", skyImg))
			geom.CullBackface = false
			geom.Transform(mgl32.Scale3D(100, 100, 100))
			gameEngine.AddSpatial(geom)
		}

		// The player object
		character := NewCharacter(shader)
		character.body.SetPosition(mgl32.Vec2{400, 400})

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
		spawnParticles := func(load func(shader *renderer.Shader) *effects.ParticleSystem, position mgl32.Vec2) *effects.ParticleSystem {
			particles := load(shader)
			particles.SetTranslation(position.Vec3(0))
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
		customController.BindKeyAction(func() {
			character.body.SetVelocity(mgl32.Vec2{-200, character.body.GetVelocity().Y()})
		}, controller.KeyA, controller.Press)
		customController.BindKeyAction(func() {
			character.body.SetVelocity(mgl32.Vec2{200, character.body.GetVelocity().Y()})
		}, controller.KeyD, controller.Press)

		//Stop walking
		customController.BindKeyAction(func() {
			character.body.SetVelocity(mgl32.Vec2{0, character.body.GetVelocity().Y()})
		}, controller.KeyA, controller.Release)
		customController.BindKeyAction(func() {
			character.body.SetVelocity(mgl32.Vec2{0, character.body.GetVelocity().Y()})
		}, controller.KeyD, controller.Release)

		// Jump
		customController.BindKeyAction(func() {
			character.body.SetVelocity(character.body.GetVelocity().Add(mgl32.Vec2{0, -300}))
			spawnParticles(dustParticles, character.body.GetPosition().Add(mgl32.Vec2{0, 0.5 * characterSize}))
		}, controller.KeySpace, controller.Press)

		// shoot
		customController.BindMouseAction(func() {
			spawnParticles(majicParticles, character.body.GetPosition())
		}, controller.MouseButtonLeft, controller.Press)

		// create sparks when player collides with somthing
		var collisionTimer time.Time
		physicsSpace.SetOnCollision(func(shapeA, shapeB *chipmunk.Shape) {
			if time.Since(collisionTimer) > 200*time.Millisecond {
				spawnParticles(sparkParticles, character.body.GetPosition().Add(mgl32.Vec2{0, 0.5 * characterSize}))
			}
			collisionTimer = time.Now()
		})
	})
}

func dustParticles(shader *renderer.Shader) *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("resources/smoke.png")
	material := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img))
	material.Transparency = renderer.EMISSIVE
	material.DepthMask = false
	particleSystem := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     5,
		ParticleEmitRate: 20,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		Material:         material,
		Shader:           shader,
		TotalFrames:      64,
		FramesX:          8,
		FramesY:          8,
		MaxLife:          3.3,
		MinLife:          2.7,
		StartColor:       color.NRGBA{254, 254, 254, 120},
		EndColor:         color.NRGBA{254, 254, 254, 0},
		StartSize:        mgl32.Vec2{40, 40}.Vec3(0),
		EndSize:          mgl32.Vec2{180, 180}.Vec3(0),
		MinTranslation:   mgl32.Vec2{-2, -2}.Vec3(0),
		MaxTranslation:   mgl32.Vec2{2, 2}.Vec3(0),
		MinStartVelocity: mgl32.Vec2{-5.0, -5.0}.Vec3(0),
		MaxStartVelocity: mgl32.Vec2{5.0, 5.0}.Vec3(0),
		MaxRotation:      -3.14,
		MinRotation:      3.14,
	})
	particleSystem.FaceCamera = false
	return particleSystem
}

func majicParticles(shader *renderer.Shader) *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("resources/majic.png")
	material := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img))
	material.Transparency = renderer.EMISSIVE
	material.DepthMask = false
	particleSystem := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     300,
		ParticleEmitRate: 700,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		Material:         material,
		Shader:           shader,
		TotalFrames:      9,
		FramesX:          3,
		FramesY:          3,
		MaxLife:          1.3,
		MinLife:          0.7,
		StartColor:       color.NRGBA{254, 254, 254, 254},
		EndColor:         color.NRGBA{254, 254, 254, 254},
		StartSize:        mgl32.Vec2{10, 10}.Vec3(0),
		EndSize:          mgl32.Vec2{10, 10}.Vec3(0),
		MinTranslation:   mgl32.Vec2{-1, -1}.Vec3(0),
		MaxTranslation:   mgl32.Vec2{1, 1}.Vec3(0),
		MinStartVelocity: mgl32.Vec2{-170, -170}.Vec3(0),
		MaxStartVelocity: mgl32.Vec2{170, 170}.Vec3(0),
		MaxRotation:      -3.14,
		MinRotation:      3.14,
	})
	particleSystem.FaceCamera = false
	return particleSystem
}

func sparkParticles(shader *renderer.Shader) *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("resources/spark.png")
	material := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img))
	material.Transparency = renderer.EMISSIVE
	material.DepthMask = false
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
		StartSize:        mgl32.Vec2{2, 2}.Vec3(0),
		EndSize:          mgl32.Vec2{2, 2}.Vec3(0),
		MinTranslation:   mgl32.Vec2{-5, -5}.Vec3(0),
		MaxTranslation:   mgl32.Vec2{5, 5}.Vec3(0),
		MinStartVelocity: mgl32.Vec2{-100, -100}.Vec3(0),
		MaxStartVelocity: mgl32.Vec2{100, 100}.Vec3(0),
		Acceleration:     mgl32.Vec2{0, 400}.Vec3(0),
		OnParticleUpdate: func(p *effects.Particle) {
			p.Scale[0] = p.Scale[0] * (1 + p.Velocity.Len()*0.05)
			p.Orientation = mgl32.QuatBetweenVectors(mgl32.Vec3{1, 0, 0}, p.Velocity)
		},
	})
	particleSystem.FaceCamera = false
	return particleSystem
}
