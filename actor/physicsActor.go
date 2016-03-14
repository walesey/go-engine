package actor

import (
	"github.com/walesey/go-engine/physics/physicsAPI"
	"github.com/walesey/go-engine/renderer"
)

type PhysicsActor struct {
	Entity renderer.Entity
	Object physicsAPI.PhysicsObject
}

func NewPhysicsActor(entity renderer.Entity, object physicsAPI.PhysicsObject) *PhysicsActor {
	return &PhysicsActor{
		Entity: entity,
		Object: object,
	}
}

func (actor *PhysicsActor) Update(dt float64) {
	//update entity\
	actor.Entity.SetTranslation(actor.Object.GetPosition())
	actor.Entity.SetOrientation(actor.Object.GetOrientation())
}
