package effects

import (
	"image/color"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/util"
)

type ParticleSettings struct {
	MaxParticles                             int
	ParticleEmitRate                         float32
	MaxLife, MinLife                         float32
	StartSize, EndSize                       mgl32.Vec3
	StartColor, EndColor                     color.NRGBA
	MaxTranslation, MinTranslation           mgl32.Vec3
	MaxStartVelocity, MinStartVelocity       mgl32.Vec3
	Acceleration                             mgl32.Vec3
	MaxRotation, MinRotation                 float32
	MaxRotationVelocity, MinRotationVelocity float32
	BaseGeometry                             *renderer.Geometry
	TotalFrames, FramesX, FramesY            int
	OnParticleUpdate                         func(p *Particle)
}

type ParticleSystem struct {
	Location        mgl32.Vec3
	DisableSpawning bool
	FaceCamera      bool
	Settings        ParticleSettings
	Node            *renderer.Node
	geometry        *renderer.Geometry
	particles       []*Particle
	life            float32
	cameraPosition  mgl32.Vec3
}

type Particle struct {
	active              bool
	geometry            *renderer.Geometry
	Life, LifeRemaining float32
	Scale               mgl32.Vec3
	Translation         mgl32.Vec3
	Orientation         mgl32.Quat
	Rotation            float32
	Velocity            mgl32.Vec3
	RotationVelocity    float32
	Color               color.Color
	Frame               int
}

func CreateParticleSystem(settings ParticleSettings) *ParticleSystem {
	geometry := renderer.CreateGeometry(make([]uint32, 0, 0), make([]float32, 0, 0))
	node := renderer.NewNode()
	node.Add(geometry)

	ps := ParticleSystem{
		Settings:   settings,
		FaceCamera: true,
		Node:       node,
		geometry:   geometry,
		particles:  make([]*Particle, settings.MaxParticles),
	}
	ps.initParitcles()
	return &ps
}

func (ps *ParticleSystem) initParitcles() {
	for i := 0; i < ps.Settings.MaxParticles; i = i + 1 {
		ps.particles[i] = &Particle{
			active:   false,
			geometry: ps.Settings.BaseGeometry.Copy(),
		}
	}
}

func (ps *ParticleSystem) Draw(renderer renderer.Renderer, transform mgl32.Mat4) {
	ps.Node.Draw(renderer, transform)
}

func (ps *ParticleSystem) Destroy(renderer renderer.Renderer) {
	ps.Node.Destroy(renderer)
}

func (ps *ParticleSystem) Centre() mgl32.Vec3 {
	return ps.Location
}

func (ps *ParticleSystem) SetParent(parent *renderer.Node) {
	ps.Node.SetParent(parent)
}

func (ps *ParticleSystem) Optimize(geometry *renderer.Geometry, transform mgl32.Mat4) {
	ps.geometry.Optimize(geometry, transform)
}

func (ps *ParticleSystem) BoundingRadius(transform mgl32.Mat4) float32 {
	return ps.Node.BoundingRadius(transform)
}

func (ps *ParticleSystem) SetCameraLocation(cameraLocation mgl32.Vec3) {
	ps.cameraPosition = cameraLocation
}

func (ps *ParticleSystem) Update(dt float64) {
	//number of new particles to spawn
	previousLife := ps.life
	ps.life = ps.life + float32(dt)
	previouseSpawnCount := int(previousLife * ps.Settings.ParticleEmitRate)
	newSpawnCount := int(ps.life * ps.Settings.ParticleEmitRate)
	spawnCount := newSpawnCount - previouseSpawnCount
	if !ps.DisableSpawning {
		for i := 0; i < spawnCount; i = i + 1 {
			ps.spawnParticle()
		}
	}
	//update all particles:
	for _, p := range ps.particles {
		ps.updateParticle(p, dt)
	}
	//sort particles and build geometry
	ps.geometry.ClearBuffers()
	for _, p := range ps.particles {
		if p.active {
			ps.loadParticle(p)
		}
	}
}

func (ps *ParticleSystem) spawnParticle() {
	//get first available inactive particle
	for i, particle := range ps.particles {
		if !particle.active {
			//spawn partice
			ps.particles[i].active = true
			randomNb := rand.Float32()
			ps.particles[i].Life = ps.Settings.MaxLife*(1.0-randomNb) + ps.Settings.MinLife*randomNb
			randomNb = rand.Float32()
			ps.particles[i].Rotation = ps.Settings.MaxRotation*(1.0-randomNb) + ps.Settings.MinRotation*randomNb
			ps.particles[i].LifeRemaining = ps.particles[i].Life
			ps.particles[i].Translation = ps.Location.Add(randomVector(ps.Settings.MinTranslation, ps.Settings.MaxTranslation))
			ps.particles[i].Velocity = randomVector(ps.Settings.MinStartVelocity, ps.Settings.MaxStartVelocity)
			// ps.particles[i].angularVelocity = ps.Settings.MaxStartVelocity.Slerp( ps.Settings.MinStartVelocity, rand.Float64() )
			randomNb = rand.Float32()
			ps.particles[i].RotationVelocity = ps.Settings.MaxRotationVelocity*(1.0-randomNb) + ps.Settings.MinRotationVelocity*randomNb
			break
		}
	}

}

func (ps *ParticleSystem) updateParticle(p *Particle, dt float64) {
	//set translation
	p.Translation = p.Translation.Add(p.Velocity.Mul(float32(dt)))
	p.Velocity = p.Velocity.Add(ps.Settings.Acceleration.Mul(float32(dt)))
	//set orientation / rotation
	p.Rotation = p.Rotation + (p.RotationVelocity * float32(dt))
	p.LifeRemaining = p.LifeRemaining - float32(dt)
	// set valuew based on life remaining
	lifeRatio := p.LifeRemaining / p.Life
	p.Scale = util.Vec3Lerp(ps.Settings.EndSize, ps.Settings.StartSize, lifeRatio)
	p.Color = lerpColor(ps.Settings.EndColor, ps.Settings.StartColor, lifeRatio)
	p.Frame = int((1.0 - lifeRatio) * float32(ps.Settings.TotalFrames))
	//face the camera
	if ps.FaceCamera {
		p.Orientation = util.FacingOrientation(p.Rotation, ps.cameraPosition.Sub(p.Translation), mgl32.Vec3{0, 0, 1}, mgl32.Vec3{-1, 0, 0})
	} else {
		p.Orientation = mgl32.QuatRotate(p.Rotation, mgl32.Vec3{0, 0, 1})
	}
	//is particle dead
	if p.LifeRemaining <= 0 {
		p.active = false
	}
	if ps.Settings.OnParticleUpdate != nil && p.active {
		ps.Settings.OnParticleUpdate(p)
	}
}

func (ps *ParticleSystem) loadParticle(p *Particle) {
	//set color
	p.geometry.SetColor(p.Color)
	//set flipbook uv
	BoxFlipbook(p.geometry, p.Frame, ps.Settings.FramesX, ps.Settings.FramesY)
	//rotate and move
	particleTransform := util.Mat4From(p.Scale, p.Translation, p.Orientation)
	//add geometry to particle system
	p.geometry.Optimize(ps.geometry, particleTransform)
}

//sets the location where the particles with be emitted from
func (ps *ParticleSystem) SetTranslation(translation mgl32.Vec3) {
	ps.Location = translation
}

func (ps *ParticleSystem) SetScale(scale mgl32.Vec3)             {} //na
func (ps *ParticleSystem) SetOrientation(orientation mgl32.Quat) {} //na

func lerpColor(color1, color2 color.NRGBA, amount float32) color.NRGBA {
	r := int(float32(color1.R)*(1.0-amount)) + int(float32(color2.R)*amount)
	g := int(float32(color1.G)*(1.0-amount)) + int(float32(color2.G)*amount)
	b := int(float32(color1.B)*(1.0-amount)) + int(float32(color2.B)*amount)
	a := int(float32(color1.A)*(1.0-amount)) + int(float32(color2.A)*amount)
	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func randomVector(min, max mgl32.Vec3) mgl32.Vec3 {
	r1, r2, r3 := rand.Float32(), rand.Float32(), rand.Float32()
	return mgl32.Vec3{min.X()*(1.0-r1) + max.X()*r1, min.Y()*(1.0-r2) + max.Y()*r2, min.Z()*(1.0-r3) + max.Z()*r3}
}
