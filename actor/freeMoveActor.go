package actor

import (
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/vectormath"
)

//an actor that can move around freely in space
type FreeMoveActor struct {
	Entity                  renderer.Entity
	Location                vectormath.Vector3
	forwardMove, strafeMove float64
	lookPitch, lookAngle    float64
	MoveSpeed, LookSpeed    float64
}

func CreateFreeMoveActor(entity renderer.Entity) *FreeMoveActor {
	return &FreeMoveActor{
		Entity:    entity,
		lookAngle: 0.0,
		lookPitch: 0.0,
		MoveSpeed: 10.0,
		LookSpeed: 0.001,
	}
}

func (actor *FreeMoveActor) Update(dt float64) {

	//orientation
	vertRot := vectormath.AngleAxis(actor.lookAngle, vectormath.Vector3{0, 1, 0})
	axis := vertRot.Apply(vectormath.Vector3{1, 0, 0}).Cross(vectormath.Vector3{0, 1, 0})
	horzRot := vectormath.AngleAxis(actor.lookPitch, axis)
	orientation := horzRot.Multiply(vertRot)
	velocity := orientation.Apply(vectormath.Vector3{actor.forwardMove, 0, actor.strafeMove})
	actor.Location = actor.Location.Add(velocity.MultiplyScalar(dt))

	//update entity
	actor.Entity.SetTranslation(actor.Location)
	actor.Entity.SetOrientation(orientation)
}

func (actor *FreeMoveActor) Look(dx, dy float64) {
	actor.lookAngle = actor.lookAngle - actor.LookSpeed*dx
	actor.lookPitch = actor.lookPitch + actor.LookSpeed*dy
	if actor.lookPitch > 1.5 {
		actor.lookPitch = 1.5
	}
	if actor.lookPitch < -1.5 {
		actor.lookPitch = -1.5
	}
}

func (actor *FreeMoveActor) StartMovingUp() {
	actor.forwardMove = actor.MoveSpeed
}

func (actor *FreeMoveActor) StartMovingDown() {
	actor.forwardMove = -actor.MoveSpeed
}

func (actor *FreeMoveActor) StartMovingLeft() {
	actor.strafeMove = -actor.MoveSpeed
}

func (actor *FreeMoveActor) StartMovingRight() {
	actor.strafeMove = actor.MoveSpeed
}

func (actor *FreeMoveActor) StopMovingUp() {
	actor.forwardMove = 0
}

func (actor *FreeMoveActor) StopMovingDown() {
	actor.forwardMove = 0
}

func (actor *FreeMoveActor) StopMovingLeft() {
	actor.strafeMove = 0
}

func (actor *FreeMoveActor) StopMovingRight() {
	actor.strafeMove = 0
}
