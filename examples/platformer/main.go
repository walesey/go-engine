package main

import (
	"image/color"
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

func NewCharacter() *Character {
	characterImg, _ := assets.ImportImageCached("resources/stickman.png")
	characterMat := renderer.NewMaterial(renderer.NewTexture("diffuseMap", characterImg, false))
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
	glRenderer := opengl.NewOpenglRenderer("Platformer", 800, 800, false)
	gameEngine := engine.NewEngine(glRenderer)
	gameEngine.InitFpsDial()

	// physics engine (Chipmonk)
	physicsSpace := chipmunkPhysics.NewChipmonkSpace()
	physicsSpace.SetGravity(mgl32.Vec2{0, 400})
	gameEngine.AddUpdatable(physicsSpace)

	gameEngine.Start(func() {
		glRenderer.BackGroundColor(0.7, 0.7, 0.9, 0.0)

		// load in default shader
		shader := renderer.NewShader()
		shader.FragSrc = fragShader
		shader.VertSrc = vertShader
		gameEngine.DefaultShader(shader)

		// The player object
		character := NewCharacter()
		character.body.SetPosition(mgl32.Vec2{400, 400})

		// Add the character to all the things
		physicsSpace.AddBody(character.body)
		gameEngine.AddOrtho(character.sprite)
		gameEngine.AddUpdatable(character)

		// terrain
		terrainGeometry := renderer.CreateBoxWithOffset(800, 800-floorHeight, 0, floorHeight)
		terrainGeometry.SetColor(color.NRGBA{0, 80, 0, 254})
		gameEngine.AddOrtho(terrainGeometry)

		// terrain physics body
		terrainBody := chipmunkPhysics.NewChipmunkBodyStatic()
		segment := chipmunk.NewSegment(vect.Vect{0, floorHeight}, vect.Vect{4000, floorHeight}, 0)
		terrainBody.Body.AddShape(segment)
		physicsSpace.AddBody(terrainBody)

		particleNode := renderer.NewNode()
		gameEngine.AddOrtho(particleNode)

		// create a new timed particle system
		spawnParticles := func(load func() *effects.ParticleGroup, position mgl32.Vec2) {
			particleGroup := load()
			particleGroup.SetTranslation(position.Vec3(0))
			effects.TriggerTimedParticleGroup(
				effects.TimedParticleGroup{
					Particles:  particleGroup,
					GameEngine: gameEngine,
					TargetNode: particleNode,
					Life:       0.2,
					Cleanup:    10,
				},
			)
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

func dustParticles() *effects.ParticleGroup {
	img, _ := assets.ImportImageCached("resources/smoke.png")
	material := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img, false))
	particleSystem := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     5,
		ParticleEmitRate: 20,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
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
	particleGroup := effects.NewParticleGroup(nil, particleSystem)
	particleGroup.Node.Material = material
	return particleGroup
}

func majicParticles() *effects.ParticleGroup {
	img, _ := assets.ImportImageCached("resources/majic.png")
	material := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img, false))
	particleSystem := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     300,
		ParticleEmitRate: 700,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
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
	particleGroup := effects.NewParticleGroup(nil, particleSystem)
	particleGroup.Node.Material = material
	return particleGroup
}

func sparkParticles() *effects.ParticleGroup {
	img, _ := assets.ImportImageCached("resources/spark.png")
	material := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img, false))
	particleSystem := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     80,
		ParticleEmitRate: 400,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
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
	particleGroup := effects.NewParticleGroup(nil, particleSystem)
	particleGroup.Node.Material = material
	return particleGroup
}

const vertShader = `
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

const fragShader = `
#version 330

uniform bool useTextures;
uniform sampler2D diffuseMap;

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 outputColor;

void main() {
	//repeat textures
	float textureX = fragTexCoord.x - int(fragTexCoord.x);
	float textureY = fragTexCoord.y - int(fragTexCoord.y);
	if (fragTexCoord.x < 0) {textureX = textureX + 1.0;}
	if (fragTexCoord.y < 0) {textureY = textureY + 1.0;}
	vec2 textCoord = vec2(textureX, textureY);

	if (useTextures) {
  	outputColor = texture(diffuseMap, textCoord) * fragColor;
	} else {
		outputColor = fragColor;
	}
}
`
