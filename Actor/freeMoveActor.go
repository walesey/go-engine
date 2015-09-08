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
	fma.velocity = vectorMath.Vector3{ fma.orientation.Apply(vectorMath.Vector3{fma.MoveSpeed,0,0}).X, 0, 0 }
}

func (fma *FreeMoveActor) StartMovingDown() {
	fma.velocity = vectorMath.Vector3{ fma.orientation.Apply(vectorMath.Vector3{-fma.MoveSpeed,0,0}).X, 0, 0 }
}

func (fma *FreeMoveActor) StartMovingLeft() {
	fma.velocity = vectorMath.Vector3{ 0, 0, fma.orientation.Apply( vectorMath.Vector3{0,0,-fma.MoveSpeed}).Z }
}

func (fma *FreeMoveActor) StartMovingRight() {
	fma.velocity = vectorMath.Vector3{ 0, 0, fma.orientation.Apply( vectorMath.Vector3{0,0,fma.MoveSpeed}).Z }
}

func (fma *FreeMoveActor) StopMovingUp() {
	fma.velocity = vectorMath.Vector3{0,fma.velocity.Y,fma.velocity.Z}
}

func (fma *FreeMoveActor) StopMovingDown() {
	fma.velocity = vectorMath.Vector3{0,fma.velocity.Y,fma.velocity.Z}
}

func (fma *FreeMoveActor) StopMovingLeft() {
	fma.velocity = vectorMath.Vector3{fma.velocity.X,fma.velocity.Y,0}
}

func (fma *FreeMoveActor) StopMovingRight() {
	fma.velocity = vectorMath.Vector3{fma.velocity.X,fma.velocity.Y,0}
}
