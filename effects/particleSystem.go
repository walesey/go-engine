package effects

import (
	"image/color"
	"math/rand"
	"sort"

	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/vectormath"
)

type ParticleSettings struct {
	MaxParticles                             int
	ParticleEmitRate                         float64
	FaceCamera                               bool
	MaxLife, MinLife                         float64
	StartSize, EndSize                       vectormath.Vector3
	StartColor, EndColor                     color.NRGBA
	MaxTranslation, MinTranslation           vectormath.Vector3
	MaxStartVelocity, MinStartVelocity       vectormath.Vector3
	Acceleration                             vectormath.Vector3
	MaxAngularVelocity, MinAngularVelocity   vectormath.Quaternion
	MaxRotationVelocity, MinRotationVelocity float64
	BaseGeometry                             *renderer.Geometry
	Material                                 *renderer.Material
	TotalFrames, FramesX, FramesY            int
}

type ParticleSystem struct {
	Location          vectormath.Vector3
	geometry          *renderer.Geometry
	particles         particleList
	settings          ParticleSettings
	particleTransform renderer.Transform
	life              float64
}

type Particle struct {
	active              bool
	geometry            *renderer.Geometry
	life, lifeRemaining float64
	translation         vectormath.Vector3
	orientation         vectormath.Quaternion
	rotation            float64
	velocity            vectormath.Vector3
	angularVelocity     vectormath.Quaternion
	rotationVelocity    float64
}

type particleList struct {
	particles      []Particle
	cameraPosition vectormath.Vector3
}

func (slice particleList) Len() int {
	return len(slice.particles)
}

func (slice particleList) Less(i, j int) bool {
	cameraDeltai := slice.cameraPosition.Subtract(slice.particles[i].translation).LengthSquared()
	cameraDeltaj := slice.cameraPosition.Subtract(slice.particles[j].translation).LengthSquared()
	return cameraDeltai > cameraDeltaj
}

func (slice particleList) Swap(i, j int) {
	slice.particles[i], slice.particles[j] = slice.particles[j], slice.particles[i]
}

func CreateParticleSystem(settings ParticleSettings) *ParticleSystem {
	geometry := renderer.CreateGeometry(make([]uint32, 0, 0), make([]float32, 0, 0))
	geometry.Material = settings.Material
	geometry.CullBackface = false
	ps := ParticleSystem{
		settings:          settings,
		geometry:          geometry,
		particles:         particleList{particles: make([]Particle, settings.MaxParticles)},
		particleTransform: renderer.CreateTransform(),
	}
	ps.initParitcles()
	return &ps
}

func (ps *ParticleSystem) initParitcles() {
	indicies := append([]uint32(nil), ps.settings.BaseGeometry.Indicies...)
	verticies := append([]float32(nil), ps.settings.BaseGeometry.Verticies...)
	for i := 0; i < ps.settings.MaxParticles; i = i + 1 {
		ps.particles.particles[i] = Particle{
			active:   false,
			geometry: renderer.CreateGeometry(indicies, verticies),
		}
	}
}

func (ps *ParticleSystem) Draw(renderer renderer.Renderer) {
	ps.geometry.Draw(renderer)
}

func (ps *ParticleSystem) Centre() vectormath.Vector3 {
	return ps.Location
}

func (ps *ParticleSystem) Optimize(geometry *renderer.Geometry, transform renderer.Transform) {
	ps.geometry.Optimize(geometry, transform)
}

func (ps *ParticleSystem) Update(dt float64, renderer renderer.Renderer) {
	ps.geometry.ClearBuffers()
	//sort particles
	ps.particles.cameraPosition = renderer.CameraLocation()
	sort.Sort(ps.particles)
	//update all particles:
	for index, _ := range ps.particles.particles {
		ps.updateParticle(&ps.particles.particles[index], renderer.CameraLocation(), dt)
	}
	//number of new particles to spawn
	previousLife := ps.life
	ps.life = ps.life + dt
	previouseSpawnCount := int(previousLife * ps.settings.ParticleEmitRate)
	newSpawnCount := int(ps.life * ps.settings.ParticleEmitRate)
	spawnCount := newSpawnCount - previouseSpawnCount
	for i := 0; i < spawnCount; i = i + 1 {
		ps.spawnParticle()
	}
}

func (ps *ParticleSystem) spawnParticle() {
	//get first available inactive particle
	for i, particle := range ps.particles.particles {
		if !particle.active {
			//spawn partice
			ps.particles.particles[i].active = true
			randomNb := rand.Float64()
			ps.particles.particles[i].life = ps.settings.MaxLife*(1.0-randomNb) + ps.settings.MinLife*randomNb
			ps.particles.particles[i].lifeRemaining = ps.particles.particles[i].life
			ps.particles.particles[i].translation = ps.Location.Add(randomVector(ps.settings.MinTranslation, ps.settings.MaxTranslation))
			ps.particles.particles[i].orientation = vectormath.IdentityQuaternion()
			ps.particles.particles[i].velocity = randomVector(ps.settings.MinStartVelocity, ps.settings.MaxStartVelocity)
			// ps.particles[i].angularVelocity = ps.settings.MaxStartVelocity.Slerp( ps.settings.MinStartVelocity, rand.Float64() )
			randomNb = rand.Float64()
			ps.particles.particles[i].rotationVelocity = ps.settings.MaxRotationVelocity*(1.0-randomNb) + ps.settings.MinRotationVelocity*randomNb
			break
		}
	}

}

func (ps *ParticleSystem) updateParticle(p *Particle, camera vectormath.Vector3, dt float64) {

	//set translation
	p.translation = p.translation.Add(p.velocity.MultiplyScalar(dt))
	p.velocity = p.velocity.Add(ps.settings.Acceleration.MultiplyScalar(dt))
	//set orientation / rotation
	if ps.settings.FaceCamera {
		p.rotation = p.rotation + (p.rotationVelocity * dt)
	} else {
		//TODO
	}
	p.lifeRemaining = p.lifeRemaining - dt
	//is particle dead
	if p.lifeRemaining <= 0 {
		p.active = false
	}
	lifeRatio := p.lifeRemaining / p.life
	scale := ps.settings.EndSize.Lerp(ps.settings.StartSize, lifeRatio)
	color := lerpColor(ps.settings.EndColor, ps.settings.StartColor, lifeRatio)
	frame := int((1.0 - lifeRatio) * float64(ps.settings.TotalFrames))
	if !p.active {
		scale = vectormath.Vector3{0, 0, 0}
	}
	//build geometry
	ps.setBaseParticle(p)
	//set color
	p.geometry.SetColor(color)
	//set flipbook uv
	BoxFlipbook(p.geometry, frame, ps.settings.FramesX, ps.settings.FramesY)
	//face the camera
	renderer.FacingTransform(ps.particleTransform, p.rotation, camera.Subtract(p.translation), vectormath.Vector3{0, 0, 1}, vectormath.Vector3{-1, 0, 0})
	p.geometry.Transform(ps.particleTransform)
	//rotate and move
	ps.particleTransform.From(scale, p.translation, p.orientation)
	//add geometry to particle system
	p.geometry.Optimize(ps.geometry, ps.particleTransform)

}

//Restore the particle's geometry back to the unTranformed state
func (ps *ParticleSystem) setBaseParticle(p *Particle) {
	for i, _ := range ps.settings.BaseGeometry.Verticies {
		p.geometry.Verticies[i] = ps.settings.BaseGeometry.Verticies[i]
	}
}

//sets the location where the particles with be emitted from
func (ps *ParticleSystem) SetTranslation(translation vectormath.Vector3) {
	ps.Location = translation
}

func (ps *ParticleSystem) SetScale(scale vectormath.Vector3)                {} //na
func (ps *ParticleSystem) SetOrientation(orientation vectormath.Quaternion) {} //na

func lerpColor(color1, color2 color.NRGBA, amount float64) color.NRGBA {
	r := int(float64(color1.R)*(1.0-amount)) + int(float64(color2.R)*amount)
	g := int(float64(color1.G)*(1.0-amount)) + int(float64(color2.G)*amount)
	b := int(float64(color1.B)*(1.0-amount)) + int(float64(color2.B)*amount)
	a := int(float64(color1.A)*(1.0-amount)) + int(float64(color2.A)*amount)
	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func randomVector(min, max vectormath.Vector3) vectormath.Vector3 {
	r1, r2, r3 := rand.Float64(), rand.Float64(), rand.Float64()
	return vectormath.Vector3{min.X*(1.0-r1) + max.X*r1, min.Y*(1.0-r2) + max.Y*r2, min.Z*(1.0-r3) + max.Z*r3}
}
