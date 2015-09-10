package effects

import(
	"image/color"
	"math/rand"

    "github.com/Walesey/goEngine/vectorMath"
	"github.com/Walesey/goEngine/renderer"
)

type ParticleSettings struct {
	MaxParticles int
	ParticleEmitRate float64
	FaceCamera bool
	MaxLife, MinLife float64
	StartSize, EndSize vectorMath.Vector3
	StartColor, EndColor color.NRGBA
	MaxTranslation, MinTranslation vectorMath.Vector3
	MaxStartVelocity, MinStartVelocity vectorMath.Vector3
	Acceleration vectorMath.Vector3
	MaxAngularVelocity, MinAngularVelocity vectorMath.Quaternion
	MaxRotationVelocity, MinRotationVelocity float64
	Material renderer.Material
	TotalFrames, FramesX, FramesY int
}

type ParticleSystem struct {
	geometry renderer.Geometry
	particles []Particle
	settings ParticleSettings
	Location vectorMath.Vector3 
	particleTransform renderer.Transform
	life float64
}

type Particle struct {
	active bool
	life, lifeRemaining float64
	translation vectorMath.Vector3 
	orientation vectorMath.Quaternion
	rotation float64
	velocity vectorMath.Vector3
	angularVelocity vectorMath.Quaternion
	rotationVelocity float64
}

func CreateParticleSystem( settings ParticleSettings ) ParticleSystem {
	geometry := renderer.CreateGeometry( make([]uint32,0,0), make([]float32,0,0) )
	geometry.Material = &settings.Material
	geometry.CullBackface = false
	ps := ParticleSystem{ 
		settings: settings, 
		geometry: geometry,
		particles: make([]Particle, settings.MaxParticles ),
		particleTransform: renderer.CreateTransform(),
	}
	ps.initParitcles()
	return ps
}

func (ps *ParticleSystem) initParitcles() {
	for i:=0; i<ps.settings.MaxParticles; i=i+1 {
		ps.particles[i] = Particle{ active: false }
	}
}

func (ps *ParticleSystem) Draw( renderer renderer.Renderer ) {
	ps.geometry.Draw(renderer)
}

func (ps *ParticleSystem) Optimize( geometry *renderer.Geometry, transform renderer.Transform ) {
    ps.geometry.Optimize(geometry, transform)
}

func (ps *ParticleSystem) Update( dt float64, renderer renderer.Renderer ){
	ps.geometry.ClearBuffers()
	//update all particles
	for index,_ := range ps.particles {
		ps.updateParticle(&ps.particles[index], renderer.CameraLocation(), dt)
	}
	//number of new particles to spawn
	previousLife := ps.life
	ps.life = ps.life + dt
	previouseSpawnCount := int(previousLife * ps.settings.ParticleEmitRate)
	newSpawnCount := int(ps.life * ps.settings.ParticleEmitRate)
	spawnCount := newSpawnCount - previouseSpawnCount
	for i:=0; i<spawnCount; i=i+1 {
		ps.spawnParticle()
	}
}

func (ps *ParticleSystem) spawnParticle(){
	//get first available inactive particle
	for i,particle := range ps.particles {
		if !particle.active {
			//spawn partice
			ps.particles[i].active = true
			randomNb := rand.Float64()
			ps.particles[i].life = ps.settings.MaxLife * (1.0-randomNb) + ps.settings.MinLife * randomNb
			ps.particles[i].lifeRemaining = ps.particles[i].life
			ps.particles[i].translation = ps.Location.Add( randomVector(ps.settings.MinTranslation, ps.settings.MaxTranslation) )
			ps.particles[i].orientation = vectorMath.IdentityQuaternion()
			ps.particles[i].rotation = 3.14
			ps.particles[i].velocity = randomVector(ps.settings.MinStartVelocity, ps.settings.MaxStartVelocity)
			// ps.particles[i].angularVelocity = ps.settings.MaxStartVelocity.Slerp( ps.settings.MinStartVelocity, rand.Float64() )
			randomNb = rand.Float64()
			ps.particles[i].rotationVelocity = ps.settings.MaxRotationVelocity * (1.0-randomNb) + ps.settings.MinRotationVelocity * randomNb
			break
		}
	}
	
}

func (ps *ParticleSystem) updateParticle( p *Particle, camera vectorMath.Vector3, dt float64 ){

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
	color := lerpColor( ps.settings.EndColor, ps.settings.StartColor, lifeRatio )
	// frame := int( lifeRatio * float64(ps.settings.TotalFrames) )//TODO
	if !p.active {
		scale = vectorMath.Vector3{0,0,0}
	}
	//build geometry
	geometry := renderer.CreateBox(float32(scale.X), float32(scale.Y))
	//set color
	geometry.SetColor(color)
	//face the camera
	renderer.FacingTransform( ps.particleTransform, p.rotation, camera.Subtract(p.translation), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,1} )
	geometry.Transform( ps.particleTransform )
	//rotate and move
	ps.particleTransform.From( scale, p.translation, p.orientation )
	//add geometry to particle system
	geometry.Optimize( &ps.geometry, ps.particleTransform )		

}

//sets the location where the particles with be emitted from
func (ps *ParticleSystem) SetTranslation( translation vectorMath.Vector3 ) {
	ps.Location = translation
}

func (ps *ParticleSystem) SetScale( scale vectorMath.Vector3 ) {} //na
func (ps *ParticleSystem) SetOrientation( orientation vectorMath.Quaternion  ) {} //na

func lerpColor( color1, color2 color.NRGBA, amount float64 ) color.NRGBA{
	r := int(float64(color1.R)*(1.0-amount)) + int(float64(color2.R)*amount)
	g := int(float64(color1.G)*(1.0-amount)) + int(float64(color2.G)*amount)
	b := int(float64(color1.B)*(1.0-amount)) + int(float64(color2.B)*amount)
	a := int(float64(color1.A)*(1.0-amount)) + int(float64(color2.A)*amount)
	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func randomVector( min, max vectorMath.Vector3 ) vectorMath.Vector3 {
	r1, r2, r3 := rand.Float64(), rand.Float64(), rand.Float64()
	return vectorMath.Vector3{ min.X*(1.0-r1) + max.X*r1, min.Y*(1.0-r2) + max.Y*r2, min.Z*(1.0-r3) + max.Z*r3 }
}