package actor

import (
	"github.com/walesey/go-engine/physics/dynamics"
	"github.com/walesey/go-engine/renderer"
)

type PhysicsActor struct {
	Entity renderer.Entity
	Object *dynamics.PhysicsObject
}

func NewPhysicsActor(Entity renderer.Entity, Object *dynamics.PhysicsObject) *PhysicsActor {
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
