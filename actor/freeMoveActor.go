package actor

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

//an actor that can move around freely in space
type FreeMoveActor struct {
	Entity                  renderer.Entity
	Location                mgl32.Vec3
	forwardMove, strafeMove float32
	lookPitch, lookAngle    float32
	MoveSpeed, LookSpeed    float32
}

func NewFreeMoveActor(entity renderer.Entity) *FreeMoveActor {
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
	vertRot := mgl32.QuatRotate(actor.lookAngle, mgl32.Vec3{0, 1, 0})
	axis := vertRot.Rotate(mgl32.Vec3{1, 0, 0}).Cross(mgl32.Vec3{0, 1, 0})
	horzRot := mgl32.QuatRotate(actor.lookPitch, axis)
	orientation := horzRot.Mul(vertRot)
	velocity := orientation.Rotate(mgl32.Vec3{actor.forwardMove, 0, actor.strafeMove})
	actor.Location = actor.Location.Add(velocity.Mul(float32(dt)))

	//update entity
	actor.Entity.SetTranslation(actor.Location)
	actor.Entity.SetOrientation(orientation)
}

func (actor *FreeMoveActor) Look(dx, dy float32) {
	actor.lookAngle = actor.lookAngle - actor.LookSpeed*dx
	actor.lookPitch = actor.lookPitch - actor.LookSpeed*dy
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
