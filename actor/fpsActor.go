package actor

import (
	"fmt"

	"github.com/walesey/go-engine/physics/physicsAPI"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

//an actor that can move around bound by physics (gravity), can jump, can walk/run
type FPSActor struct {
	Entity                  renderer.Entity
	Character               physicsAPI.CharacterController
	MoveSpeed, SprintSpeed  float64
	LookSpeed               float64
	forwardMove, strafeMove float64
	lookPitch, lookAngle    float64
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
	vertRot := vmath.AngleAxis(actor.lookAngle, vmath.Vector3{0, 1, 0})
	axis := vertRot.Apply(vmath.Vector3{1, 0, 0}).Cross(vmath.Vector3{0, 1, 0})
	horzRot := vmath.AngleAxis(actor.lookPitch, axis)
	orientation := horzRot.Multiply(vertRot)
	actor.Entity.SetOrientation(orientation)

	// walking direction
	direction := orientation.Apply(vmath.Vector3{actor.forwardMove, 0, actor.strafeMove})
	if !actor.Character.OnGround() {
		direction = vmath.Vector3{0, 0, 0}
	}
	actor.Character.SetWalkDirection(direction)
	actor.Entity.SetTranslation(actor.Character.GetPosition())
}

func (actor *FPSActor) Look(dx, dy float64) {
	actor.lookAngle = actor.lookAngle - actor.LookSpeed*dx
	actor.lookPitch = actor.lookPitch + actor.LookSpeed*dy
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
		fmt.Println()
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
