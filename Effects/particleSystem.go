package effects

import(
	"image"
	"image/color"
	"math/rand"

    "github.com/Walesey/goEngine/vectorMath"
	"github.com/Walesey/goEngine/renderer"

	"github.com/disintegration/imaging"
)

type ParticleSettings struct {
	MaxParticles int
	ParticleEmitRate float64
	Sprite Sprite
	FaceCamera bool
	MaxLife, MinLife float64
	StartSize, EndSize vectorMath.Vector3
	StartColor, EndColor color.NRGBA
	MaxTranslation, MinTranslation vectorMath.Vector3
	MaxStartVelocity, MinStartVelocity vectorMath.Vector3
	Acceleration vectorMath.Vector3
	MaxAngularVelocity, MinAngularVelocity vectorMath.Quaternion
	MaxRotationVelocity, MinRotationVelocity float64
}

type ParticleSystem struct {
	Node renderer.Node
	particles []Particle
	settings ParticleSettings
	Location vectorMath.Vector3 
	life float64
}

type Particle struct {
	node renderer.Node
	sprite *Sprite
	active bool
	life, lifeRemaining float64
	translation, size vectorMath.Vector3 
	orientation vectorMath.Quaternion
	rotation float64
	velocity vectorMath.Vector3
	angularVelocity vectorMath.Quaternion
	rotationVelocity float64
}

func CreateParticleSystem( settings ParticleSettings ) ParticleSystem {
	ps := ParticleSystem{ 
		settings: settings, 
		Node: renderer.CreateNode(), 
		particles: make([]Particle, settings.MaxParticles ),
	}
	ps.initParitcles()
	return ps
}

func (ps *ParticleSystem) initParitcles() {
	for i:=0; i<ps.settings.MaxParticles; i=i+1 {
		node := renderer.CreateNode()
		sprite := ps.settings.Sprite
		node.Add( &sprite )
		ps.particles[i] = Particle{
			node: node,
			active: false,
			size: vectorMath.Vector3{0,0,0},
		}
		ps.particles[i].sprite = &sprite
		ps.Node.Add(&ps.particles[i].node)
	}
}

func (ps *ParticleSystem) Update( dt float64, renderer renderer.Renderer ){
	//update all particles
	for index,_ := range ps.particles {
		ps.UpdateParticle(&ps.particles[index], renderer.CameraLocation(), dt)
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
	index := 0
	for i,particle := range ps.particles {
		if !particle.active {
			index = i
			break
		}
	}
	//spawn partice
	ps.particles[index].active = true
	randomNb := rand.Float64()
	ps.particles[index].life = ps.settings.MaxLife * (1.0-randomNb) + ps.settings.MinLife * randomNb
	ps.particles[index].lifeRemaining = ps.particles[index].life
	ps.particles[index].translation = ps.Location.Add( ps.settings.MaxTranslation.Lerp( ps.settings.MinTranslation, rand.Float64() ) )
	ps.particles[index].orientation = vectorMath.IdentityQuaternion()
	ps.particles[index].rotation = 3.14
	ps.particles[index].velocity = ps.settings.MaxStartVelocity.Lerp( ps.settings.MinStartVelocity, rand.Float64() )
	// ps.particles[index].angularVelocity = ps.settings.MaxStartVelocity.Slerp( ps.settings.MinStartVelocity, rand.Float64() )
	randomNb = rand.Float64()
	ps.particles[index].rotationVelocity = ps.settings.MaxRotationVelocity * (1.0-randomNb) + ps.settings.MinRotationVelocity * randomNb
}

func (ps *ParticleSystem) UpdateParticle( p *Particle, camera vectorMath.Vector3, dt float64 ){
	if p.active {
		//set translation
		p.translation = p.translation.Add(p.velocity.MultiplyScalar(dt))
		p.node.SetTranslation(p.translation)
		//set orientation / rotation
		if ps.settings.FaceCamera {
			p.rotation = p.rotation + (p.rotationVelocity * dt)
			p.node.SetFacing( p.rotation, camera.Subtract(p.translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )
		} else {
			//TODO
			p.node.SetOrientation(p.orientation)
		}
		//TODO set color

		p.lifeRemaining = p.lifeRemaining - dt
		lifeRatio := p.lifeRemaining / p.life
		p.node.SetScale(ps.settings.EndSize.Lerp(ps.settings.StartSize, lifeRatio))
		//set flipbook frame
		p.sprite.FrameLerp(lifeRatio)
		//is particle dead
		if p.lifeRemaining <= 0 {
			p.active = false
			p.node.SetScale(vectorMath.Vector3{0,0,0})
		} 
	} else {
		p.node.SetScale(vectorMath.Vector3{0,0,0})
	}
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

func colorImage( img image.Image, colorRGBA color.NRGBA ) image.Image{
	return imaging.AdjustFunc( img,
		func(c color.NRGBA) color.NRGBA {
			r := (int(c.R) * int(colorRGBA.R)) / 255
			g := (int(c.G) * int(colorRGBA.G)) / 255
			b := (int(c.B) * int(colorRGBA.B)) / 255
			a := (int(c.A) * int(colorRGBA.A)) / 255
			return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		},
	)
}

