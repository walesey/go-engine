package effects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

type ParticleGroup struct {
	node      *renderer.Node
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
		particle.SetCameraLocation(pg.camera.GetTranslation())
		particle.Update(dt)
	}
}

func (pg *ParticleGroup) SetTranslation(translation mgl32.Vec3) {
	for _, particle := range pg.particles {
		particle.Location = translation
	}
}

func (pg *ParticleGroup) Draw(renderer renderer.Renderer) {
	pg.node.Draw(renderer)
}

func (pg *ParticleGroup) Destroy(renderer renderer.Renderer) {
	pg.node.Destroy(renderer)
}

func (pg *ParticleGroup) Centre() mgl32.Vec3 {
	return pg.node.Centre()
}

func (pg *ParticleGroup) Optimize(geometry *renderer.Geometry, transform mgl32.Mat4) {
	pg.node.Optimize(geometry, transform)
}

func NewParticleGroup(camera *renderer.Camera, particles ...*ParticleSystem) *ParticleGroup {
	node := renderer.CreateNode()
	for _, particle := range particles {
		node.Add(particle)
	}
	return &ParticleGroup{
		node:      node,
		camera:    camera,
		particles: particles,
	}
}
