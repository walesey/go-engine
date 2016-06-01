package actor

import (
	"github.com/walesey/go-engine/physics/physicsAPI"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type PhysicsActor2D struct {
	Entity renderer.Entity
	Object physicsAPI.PhysicsObject2D
	Mask   vmath.Vector3
}

func NewPhysicsActor2D(entity renderer.Entity, object physicsAPI.PhysicsObject2D, mask vmath.Vector3) *PhysicsActor2D {
	return &PhysicsActor2D{
		Entity: entity,
		Object: object,
		Mask:   mask,
	}
}

func (actor *PhysicsActor2D) Update(dt float64) {
	objPos := actor.Object.GetPosition()
	var position vmath.Vector3
	var orientation vmath.Quaternion
	if actor.Mask.X < actor.Mask.Y && actor.Mask.X < actor.Mask.Z { // YZ plane
		position = vmath.Vector3{0, actor.Mask.Y * objPos.X, actor.Mask.Z * objPos.Y}
		orientation = vmath.AngleAxis(actor.Object.GetAngle(), vmath.Vector3{1, 0, 0})
	} else if actor.Mask.Y < actor.Mask.X && actor.Mask.Y < actor.Mask.Z { // XZ plane
		position = vmath.Vector3{actor.Mask.X * objPos.X, 0, actor.Mask.Z * objPos.Y}
		orientation = vmath.AngleAxis(actor.Object.GetAngle(), vmath.Vector3{0, 1, 0})
	} else { // XY plane
		position = vmath.Vector3{actor.Mask.X * objPos.X, actor.Mask.Y * objPos.Y, 0}
		orientation = vmath.AngleAxis(actor.Object.GetAngle(), vmath.Vector3{0, 0, 1})
	}
	actor.Entity.SetTranslation(position)
	actor.Entity.SetOrientation(orientation)
}
