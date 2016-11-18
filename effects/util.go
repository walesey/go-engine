package effects

import (
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/renderer"
)

type TimedParticleGroup struct {
	GameEngine    engine.Engine
	TargetNode    *renderer.Node
	Particles     *ParticleGroup
	Life, Cleanup float64
}

func TriggerTimedParticleGroup(tpg TimedParticleGroup) {
	tpg.TargetNode.Add(tpg.Particles)
	tpg.GameEngine.AddUpdatable(tpg.Particles)
	particleLife := tpg.Life

	var updater engine.Updatable
	updater = engine.UpdatableFunc(func(dt float64) {
		particleLife -= dt
		if particleLife < 0 {
			tpg.Particles.Disable(true)
		}
		if particleLife < -tpg.Cleanup {
			tpg.TargetNode.Remove(tpg.Particles, true)
			tpg.GameEngine.RemoveUpdatable(tpg.Particles)
			tpg.GameEngine.RemoveUpdatable(updater)
		}
	})
	tpg.GameEngine.AddUpdatable(updater)
}
