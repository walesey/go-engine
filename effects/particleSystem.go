package effects

import (
	"image/color"
	"math/rand"
	"sort"

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
	MaxAngularVelocity, MinAngularVelocity   vmath.Quaternion
	MaxRotationVelocity, MinRotationVelocity float64
	BaseGeometry                             *renderer.Geometry
	Material                                 *renderer.Material
	TotalFrames, FramesX, FramesY            int
}

type ParticleSystem struct {
	Location          vmath.Vector3
	DisableSpawning   bool
	geometry          *renderer.Geometry
	particles         []*Particle
	settings          ParticleSettings
	particleTransform renderer.Transform
	life              float64
	cameraPosition    vmath.Vector3
}

type Particle struct {
	active              bool
	geometry            *renderer.Geometry
	life, lifeRemaining float64
	translation         vmath.Vector3
	rotation            float64
	velocity            vmath.Vector3
	rotationVelocity    float64
}

func (ps *ParticleSystem) Len() int {
	return len(ps.particles)
}

func (ps *ParticleSystem) Less(i, j int) bool {
	cameraDeltai := ps.particles[i].translation.Subtract(ps.cameraPosition).LengthSquared()
	cameraDeltaj := ps.particles[j].translation.Subtract(ps.cameraPosition).LengthSquared()
	return cameraDeltai < cameraDeltaj
}

func (ps *ParticleSystem) Swap(i, j int) {
	ps.particles[i], ps.particles[j] = ps.particles[j], ps.particles[i]
}

func CreateParticleSystem(settings ParticleSettings) *ParticleSystem {
	geometry := renderer.CreateGeometry(make([]uint32, 0, 0), make([]float32, 0, 0))
	geometry.Material = settings.Material
	geometry.CullBackface = false
	ps := ParticleSystem{
		settings:          settings,
		geometry:          geometry,
		particles:         make([]*Particle, settings.MaxParticles),
		particleTransform: renderer.CreateTransform(),
	}
	ps.initParitcles()
	return &ps
}

func (ps *ParticleSystem) initParitcles() {
	for i := 0; i < ps.settings.MaxParticles; i = i + 1 {
		ps.particles[i] = &Particle{
			active:   false,
			geometry: ps.settings.BaseGeometry.Copy(),
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
	//update all particles:
	for _, p := range ps.particles {
		ps.updateParticle(p, dt)
	}
	//number of new particles to spawn
	previousLife := ps.life
	ps.life = ps.life + dt
	previouseSpawnCount := int(previousLife * ps.settings.ParticleEmitRate)
	newSpawnCount := int(ps.life * ps.settings.ParticleEmitRate)
	spawnCount := newSpawnCount - previouseSpawnCount
	if !ps.DisableSpawning {
		for i := 0; i < spawnCount; i = i + 1 {
			ps.spawnParticle()
		}
	}
	//sort particles and build geometry
	sort.Sort(ps)
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
			ps.particles[i].life = ps.settings.MaxLife*(1.0-randomNb) + ps.settings.MinLife*randomNb
			ps.particles[i].lifeRemaining = ps.particles[i].life
			ps.particles[i].translation = ps.Location.Add(randomVector(ps.settings.MinTranslation, ps.settings.MaxTranslation))
			ps.particles[i].velocity = randomVector(ps.settings.MinStartVelocity, ps.settings.MaxStartVelocity)
			// ps.particles[i].angularVelocity = ps.settings.MaxStartVelocity.Slerp( ps.settings.MinStartVelocity, rand.Float64() )
			randomNb = rand.Float64()
			ps.particles[i].rotationVelocity = ps.settings.MaxRotationVelocity*(1.0-randomNb) + ps.settings.MinRotationVelocity*randomNb
			break
		}
	}

}

func (ps *ParticleSystem) updateParticle(p *Particle, dt float64) {
	//set translation
	p.translation = p.translation.Add(p.velocity.MultiplyScalar(dt))
	p.velocity = p.velocity.Add(ps.settings.Acceleration.MultiplyScalar(dt))
	//set orientation / rotation
	p.rotation = p.rotation + (p.rotationVelocity * dt)
	p.lifeRemaining = p.lifeRemaining - dt
	//is particle dead
	if p.lifeRemaining <= 0 {
		p.active = false
	}
}

func (ps *ParticleSystem) loadParticle(p *Particle) {
	lifeRatio := p.lifeRemaining / p.life
	scale := ps.settings.EndSize.Lerp(ps.settings.StartSize, lifeRatio)
	color := lerpColor(ps.settings.EndColor, ps.settings.StartColor, lifeRatio)
	frame := int((1.0 - lifeRatio) * float64(ps.settings.TotalFrames))
	//set color
	p.geometry.SetColor(color)
	//set flipbook uv
	BoxFlipbook(p.geometry, frame, ps.settings.FramesX, ps.settings.FramesY)
	//face the camera
	orientation := vmath.FacingOrientation(p.rotation, ps.cameraPosition.Subtract(p.translation), vmath.Vector3{0, 0, 1}, vmath.Vector3{-1, 0, 0})
	//rotate and move
	ps.particleTransform.From(scale, p.translation, orientation)
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
