package actor

import (
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

//an actor that can move around freely in space using wasd and mouse
type FreeMoveFPSActor struct {
	Entity                  renderer.Entity
	Location                vmath.Vector3
	forwardMove, strafeMove float64
	lookPitch, lookAngle    float64
	MoveSpeed, LookSpeed    float64
}

func NewFreeMoveFPSActor(entity renderer.Entity) *FreeMoveFPSActor {
	return &FreeMoveFPSActor{
		Entity:    entity,
		lookAngle: 0.0,
		lookPitch: 0.0,
		MoveSpeed: 10.0,
		LookSpeed: 0.001,
	}
}

func (actor *FreeMoveFPSActor) Update(dt float64) {

	//orientation
	vertRot := vmath.AngleAxis(actor.lookAngle, vmath.Vector3{0, 1, 0})
	axis := vertRot.Apply(vmath.Vector3{1, 0, 0}).Cross(vmath.Vector3{0, 1, 0})
	horzRot := vmath.AngleAxis(actor.lookPitch, axis)
	orientation := horzRot.Multiply(vertRot)
	actor.Entity.SetOrientation(orientation)

	//location
	velocity := orientation.Apply(vmath.Vector3{actor.forwardMove, 0, actor.strafeMove})
	actor.Location = actor.Location.Add(velocity.MultiplyScalar(dt))
	actor.Entity.SetTranslation(actor.Location)
}

func (actor *FreeMoveFPSActor) Look(dx, dy float64) {
	actor.lookAngle = actor.lookAngle - actor.LookSpeed*dx
	actor.lookPitch = actor.lookPitch + actor.LookSpeed*dy
	if actor.lookPitch > 1.5 {
		actor.lookPitch = 1.5
	}
	if actor.lookPitch < -1.5 {
		actor.lookPitch = -1.5
	}
}

func (actor *FreeMoveFPSActor) StartMovingForward() {
	actor.forwardMove = actor.MoveSpeed
}

func (actor *FreeMoveFPSActor) StartMovingBackward() {
	actor.forwardMove = -actor.MoveSpeed
}

func (actor *FreeMoveFPSActor) StartStrafingLeft() {
	actor.strafeMove = -actor.MoveSpeed
}

func (actor *FreeMoveFPSActor) StartStrafingRight() {
	actor.strafeMove = actor.MoveSpeed
}

func (actor *FreeMoveFPSActor) StopMovingForward() {
	actor.forwardMove = 0
}

func (actor *FreeMoveFPSActor) StopMovingBackward() {
	actor.forwardMove = 0
}

func (actor *FreeMoveFPSActor) StopStrafingLeft() {
	actor.strafeMove = 0
}

func (actor *FreeMoveFPSActor) StopStrafingRight() {
	actor.strafeMove = 0
}

func (actor *FreeMoveFPSActor) Jump()           {}
func (actor *FreeMoveFPSActor) StandUp()        {}
func (actor *FreeMoveFPSActor) Crouch()         {}
func (actor *FreeMoveFPSActor) Prone()          {}
func (actor *FreeMoveFPSActor) StartSprinting() {}
func (actor *FreeMoveFPSActor) StopSprinting()  {}
