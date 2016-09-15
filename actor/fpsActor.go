package actor

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/physics/physicsAPI"
	"github.com/walesey/go-engine/renderer"
)

//an actor that can move around bound by physics (gravity), can jump, can walk/run
type FPSActor struct {
	Entity                  renderer.Entity
	Character               physicsAPI.CharacterController
	MoveSpeed, SprintSpeed  float32
	LookSpeed               float32
	lookPitch, lookAngle    float32
	forwardMove, strafeMove float32
	walkDirection           mgl32.Vec3
}

func NewFPSActor(entity renderer.Entity, character physicsAPI.CharacterController) *FPSActor {
	return &FPSActor{
		Entity:      entity,
		Character:   character,
		MoveSpeed:   0.3,
		SprintSpeed: 0.5,
		LookSpeed:   0.001,
	}
}

func (actor *FPSActor) Update(dt float64) {
	// orientation
	vertRot := mgl32.QuatRotate(actor.lookAngle, mgl32.Vec3{0, 1, 0})
	axis := vertRot.Rotate(mgl32.Vec3{1, 0, 0}).Cross(mgl32.Vec3{0, 1, 0})
	horzRot := mgl32.QuatRotate(actor.lookPitch, axis)
	orientation := horzRot.Mul(vertRot)
	actor.Entity.SetOrientation(orientation)

	// walking direction
	if actor.Character.OnGround() {
		actor.walkDirection = orientation.Rotate(mgl32.Vec3{actor.forwardMove, 0, actor.strafeMove})
	}
	actor.Character.SetWalkDirection(actor.walkDirection)
	actor.Entity.SetTranslation(actor.Character.GetPosition())
}

func (actor *FPSActor) Look(dx, dy float32) {
	actor.lookAngle = actor.lookAngle - actor.LookSpeed*dx
	actor.lookPitch = actor.lookPitch - actor.LookSpeed*dy
	if actor.lookPitch > 1.5 {
		actor.lookPitch = 1.5
	}
	if actor.lookPitch < -1.5 {
		actor.lookPitch = -1.5
	}
}

func (actor *FPSActor) StartMovingForward() {
	actor.forwardMove = actor.MoveSpeed
}

func (actor *FPSActor) StartMovingBackward() {
	actor.forwardMove = -actor.MoveSpeed
}

func (actor *FPSActor) StartStrafingLeft() {
	actor.strafeMove = -actor.MoveSpeed
}

func (actor *FPSActor) StartStrafingRight() {
	actor.strafeMove = actor.MoveSpeed
}

func (actor *FPSActor) StopMovingForward() {
	actor.forwardMove = 0
}

func (actor *FPSActor) StopMovingBackward() {
	actor.forwardMove = 0
}

func (actor *FPSActor) StopStrafingLeft() {
	actor.strafeMove = 0
}

func (actor *FPSActor) StopStrafingRight() {
	actor.strafeMove = 0
}

func (actor *FPSActor) Jump() {
	if actor.Character.CanJump() {
		actor.Character.Jump()
	}
}

func (actor *FPSActor) StandUp() {

}

func (actor *FPSActor) Crouch() {

}

func (actor *FPSActor) Prone() {

}

func (actor *FPSActor) StartSprinting() {

}

func (actor *FPSActor) StopSprinting() {

}
