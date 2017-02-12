package effects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

type ParticleGroup struct {
	Node      *renderer.Node
	camera    *renderer.Camera
	particles []*ParticleSystem
}

func (pg *ParticleGroup) Disable(disable bool) {
	for _, particle := range pg.particles {
		particle.DisableSpawning = disable
	}
}

func (pg *ParticleGroup) Update(dt float64) {
	for _, particle := range pg.particles {
		if pg.camera != nil {
			particle.SetCameraLocation(pg.camera.Translation)
		}
		particle.Update(dt)
	}
}

func (pg *ParticleGroup) SetTranslation(translation mgl32.Vec3) {
	for _, particle := range pg.particles {
		particle.Location = translation
	}
}

func (pg *ParticleGroup) Draw(renderer renderer.Renderer, transform mgl32.Mat4) {
	pg.Node.Draw(renderer, transform)
}

func (pg *ParticleGroup) Destroy(renderer renderer.Renderer) {
	pg.Node.Destroy(renderer)
}

func (pg *ParticleGroup) Center() mgl32.Vec3 {
	return pg.Node.Center()
}

func (pg *ParticleGroup) SetParent(parent *renderer.Node) {
	pg.Node.SetParent(parent)
}

func (pg *ParticleGroup) Optimize(geometry *renderer.Geometry, transform mgl32.Mat4) {
	pg.Node.Optimize(geometry, transform)
}

func (pg *ParticleGroup) BoundingRadius() float32 {
	return pg.Node.BoundingRadius()
}

func (pg *ParticleGroup) OrthoOrder() int {
	return pg.Node.OrthoOrder()
}

func NewParticleGroup(camera *renderer.Camera, particles ...*ParticleSystem) *ParticleGroup {
	node := renderer.NewNode()
	for _, particle := range particles {
		node.Add(particle)
	}
	return &ParticleGroup{
		Node:      node,
		camera:    camera,
		particles: particles,
	}
}
