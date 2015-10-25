package actor

import (
	"github.com/walesey/go-engine/physics"
	"github.com/walesey/go-engine/renderer"
)

type PhysicsActor struct {
	Entity renderer.Entity
	Object *physics.PhysicsObject
}

func NewPhysicsActor(Entity renderer.Entity, Object *physics.PhysicsObject) *PhysicsActor {
	return &PhysicsActor{
		Entity: Entity,
		Object: Object,
	}
}

func (actor *PhysicsActor) Update(dt float64) {
	//update entity
	actor.Entity.SetTranslation(actor.Object.Position)
	actor.Entity.SetOrientation(actor.Object.Orientation)
}
