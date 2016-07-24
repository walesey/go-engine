package effects

import (
	"image/color"
	"math/rand"

	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ParticleSettings struct {
	MaxParticles                             int
	ParticleEmitRate                         float64
	MaxLife, MinLife                         float64
	StartSize, EndSize                       vmath.Vector3
	StartColor, EndColor                     color.NRGBA
	MaxTranslation, MinTranslation           vmath.Vector3
	MaxStartVelocity, MinStartVelocity       vmath.Vector3
	Acceleration                             vmath.Vector3
	MaxRotation, MinRotation                 float64
	MaxRotationVelocity, MinRotationVelocity float64
	BaseGeometry                             *renderer.Geometry
	Material                                 *renderer.Material
	TotalFrames, FramesX, FramesY            int
	OnParticleUpdate                         func(p *Particle)
}

type ParticleSystem struct {
	Location          vmath.Vector3
	DisableSpawning   bool
	Settings          ParticleSettings
	geometry          *renderer.Geometry
	particles         []*Particle
	particleTransform renderer.Transform
	life              float64
	cameraPosition    vmath.Vector3
}

type Particle struct {
	active              bool
	geometry            *renderer.Geometry
	Life, LifeRemaining float64
	Scale               vmath.Vector3
	Translation         vmath.Vector3
	Orientation         vmath.Quaternion
	Rotation            float64
	Velocity            vmath.Vector3
	RotationVelocity    float64
	Color               color.Color
	Frame               int
}

func CreateParticleSystem(settings ParticleSettings) *ParticleSystem {
	geometry := renderer.CreateGeometry(make([]uint32, 0, 0), make([]float32, 0, 0))
	geometry.Material = settings.Material
	geometry.CullBackface = false
	ps := ParticleSystem{
		Settings:          settings,
		geometry:          geometry,
		particles:         make([]*Particle, settings.MaxParticles),
		particleTransform: renderer.CreateTransform(),
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

func (ps *ParticleSystem) Draw(renderer renderer.Renderer) {
	ps.geometry.Draw(renderer)
}

func (ps *ParticleSystem) Destroy(renderer renderer.Renderer) {
	ps.geometry.Destroy(renderer)
}

func (ps *ParticleSystem) Centre() vmath.Vector3 {
	return ps.Location
}

func (ps *ParticleSystem) Optimize(geometry *renderer.Geometry, transform renderer.Transform) {
	ps.geometry.Optimize(geometry, transform)
}

func (ps *ParticleSystem) SetCameraLocation(cameraLocation vmath.Vector3) {
	ps.cameraPosition = cameraLocation
}

func (ps *ParticleSystem) Update(dt float64) {
	//number of new particles to spawn
	previousLife := ps.life
	ps.life = ps.life + dt
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
			randomNb := rand.Float64()
			ps.particles[i].Life = ps.Settings.MaxLife*(1.0-randomNb) + ps.Settings.MinLife*randomNb
			randomNb = rand.Float64()
			ps.particles[i].Rotation = ps.Settings.MaxRotation*(1.0-randomNb) + ps.Settings.MinRotation*randomNb
			ps.particles[i].LifeRemaining = ps.particles[i].Life
			ps.particles[i].Translation = ps.Location.Add(randomVector(ps.Settings.MinTranslation, ps.Settings.MaxTranslation))
			ps.particles[i].Velocity = randomVector(ps.Settings.MinStartVelocity, ps.Settings.MaxStartVelocity)
			// ps.particles[i].angularVelocity = ps.Settings.MaxStartVelocity.Slerp( ps.Settings.MinStartVelocity, rand.Float64() )
			randomNb = rand.Float64()
			ps.particles[i].RotationVelocity = ps.Settings.MaxRotationVelocity*(1.0-randomNb) + ps.Settings.MinRotationVelocity*randomNb
			break
		}
	}

}

func (ps *ParticleSystem) updateParticle(p *Particle, dt float64) {
	//set translation
	p.Translation = p.Translation.Add(p.Velocity.MultiplyScalar(dt))
	p.Velocity = p.Velocity.Add(ps.Settings.Acceleration.MultiplyScalar(dt))
	//set orientation / rotation
	p.Rotation = p.Rotation + (p.RotationVelocity * dt)
	p.LifeRemaining = p.LifeRemaining - dt
	// set valuew based on life remaining
	lifeRatio := p.LifeRemaining / p.Life
	p.Scale = ps.Settings.EndSize.Lerp(ps.Settings.StartSize, lifeRatio)
	p.Color = lerpColor(ps.Settings.EndColor, ps.Settings.StartColor, lifeRatio)
	p.Frame = int((1.0 - lifeRatio) * float64(ps.Settings.TotalFrames))
	//face the camera
	p.Orientation = vmath.FacingOrientation(p.Rotation, ps.cameraPosition.Subtract(p.Translation), vmath.Vector3{0, 0, 1}, vmath.Vector3{-1, 0, 0})
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
	ps.particleTransform.From(p.Scale, p.Translation, p.Orientation)
	//add geometry to particle system
	p.geometry.Optimize(ps.geometry, ps.particleTransform)
}

//sets the location where the particles with be emitted from
func (ps *ParticleSystem) SetTranslation(translation vmath.Vector3) {
	ps.Location = translation
}

func (ps *ParticleSystem) SetScale(scale vmath.Vector3)                {} //na
func (ps *ParticleSystem) SetOrientation(orientation vmath.Quaternion) {} //na

func lerpColor(color1, color2 color.NRGBA, amount float64) color.NRGBA {
	r := int(float64(color1.R)*(1.0-amount)) + int(float64(color2.R)*amount)
	g := int(float64(color1.G)*(1.0-amount)) + int(float64(color2.G)*amount)
	b := int(float64(color1.B)*(1.0-amount)) + int(float64(color2.B)*amount)
	a := int(float64(color1.A)*(1.0-amount)) + int(float64(color2.A)*amount)
	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func randomVector(min, max vmath.Vector3) vmath.Vector3 {
	r1, r2, r3 := rand.Float64(), rand.Float64(), rand.Float64()
	return vmath.Vector3{min.X*(1.0-r1) + max.X*r1, min.Y*(1.0-r2) + max.Y*r2, min.Z*(1.0-r3) + max.Z*r3}
}
