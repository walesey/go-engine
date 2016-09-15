package actor

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/physics/physicsAPI"
	"github.com/walesey/go-engine/renderer"
)

type PhysicsActor2D struct {
	Entity renderer.Entity
	Object physicsAPI.PhysicsObject2D
	Mask   mgl32.Vec3
}

func NewPhysicsActor2D(entity renderer.Entity, object physicsAPI.PhysicsObject2D, mask mgl32.Vec3) *PhysicsActor2D {
	return &PhysicsActor2D{
		Entity: entity,
		Object: object,
		Mask:   mask,
	}
}

func (actor *PhysicsActor2D) Update(dt float64) {
	objPos := actor.Object.GetPosition()
	var position mgl32.Vec3
	var orientation mgl32.Quat
	if actor.Mask.X() < actor.Mask.Y() && actor.Mask.X() < actor.Mask.Z() { // YZ plane
		position = mgl32.Vec3{0, actor.Mask.Y() * objPos.X(), actor.Mask.Z() * objPos.Y()}
		orientation = mgl32.QuatRotate(actor.Object.GetAngle(), mgl32.Vec3{1, 0, 0})
	} else if actor.Mask.Y() < actor.Mask.X() && actor.Mask.Y() < actor.Mask.Z() { // XZ plane
		position = mgl32.Vec3{actor.Mask.X() * objPos.X(), 0, actor.Mask.Z() * objPos.Y()}
		orientation = mgl32.QuatRotate(actor.Object.GetAngle(), mgl32.Vec3{0, 1, 0})
	} else { // XY plane
		position = mgl32.Vec3{actor.Mask.X() * objPos.X(), actor.Mask.Y() * objPos.Y(), 0}
		orientation = mgl32.QuatRotate(actor.Object.GetAngle(), mgl32.Vec3{0, 0, 1})
	}
	actor.Entity.SetTranslation(position)
	actor.Entity.SetOrientation(orientation)
}
