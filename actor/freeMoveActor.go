package actor

import (
    "github.com/Walesey/goEngine/vectorMath"
    "github.com/Walesey/goEngine/renderer"
)

//an actor that can move around freely in space (up/down/left/right/forward/backward)
type FreeMoveActor struct {
	Entity renderer.Entity
	location vectorMath.Vector3
	velocity vectorMath.Vector3
	orientation vectorMath.Quaternion
	MoveSpeed float64
}

func CreateFreeMoveActor( entity renderer.Entity ) *FreeMoveActor {
	return &FreeMoveActor{
		Entity: entity,
		orientation: vectorMath.IdentityQuaternion(),
		MoveSpeed: 1.0,
	}
}

func (fma *FreeMoveActor) Update( dt float64 ) {
	fma.location = fma.location.Add( fma.velocity.MultiplyScalar(dt) )
    fma.Entity.SetTranslation(fma.location)
    fma.Entity.SetOrientation(fma.orientation)
}

func (fma *FreeMoveActor) StartMovingUp() {
	fma.velocity.X = fma.orientation.Apply(vectorMath.Vector3{fma.MoveSpeed,0,0}).X
}

func (fma *FreeMoveActor) StartMovingDown() {
	fma.velocity.X = fma.orientation.Apply(vectorMath.Vector3{-fma.MoveSpeed,0,0}).X
}

func (fma *FreeMoveActor) StartMovingLeft() {
	fma.velocity.Z = fma.orientation.Apply( vectorMath.Vector3{0,0,-fma.MoveSpeed}).Z
}

func (fma *FreeMoveActor) StartMovingRight() {
	fma.velocity.Z = fma.orientation.Apply( vectorMath.Vector3{0,0,fma.MoveSpeed}).Z
}

func (fma *FreeMoveActor) StopMovingUp() {
	fma.velocity.X = 0
}

func (fma *FreeMoveActor) StopMovingDown() {
	fma.velocity.X = 0
}

func (fma *FreeMoveActor) StopMovingLeft() {
	fma.velocity.Z = 0
}

func (fma *FreeMoveActor) StopMovingRight() {
	fma.velocity.Z = 0
}
