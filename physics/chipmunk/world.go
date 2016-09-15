package chipmunkPhysics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/walesey/go-engine/physics/physicsAPI"
)

type ChipmonkSpace struct {
	Space       *chipmunk.Space
	OnCollision func(shapeA, shapeB *chipmunk.Shape)
}

func NewChipmonkSpace() *ChipmonkSpace {
	return &ChipmonkSpace{
		Space: chipmunk.NewSpace(),
	}
}

func (cSpace *ChipmonkSpace) Update(dt float64) {
	cSpace.Space.Step(vect.Float(dt))
	if cSpace.OnCollision != nil {
		for _, a := range cSpace.Space.Arbiters {
			cSpace.OnCollision(a.ShapeA, a.ShapeB)
		}
	}
}

func (cSpace *ChipmonkSpace) SetOnCollision(onCollision func(shapeA, shapeB *chipmunk.Shape)) {
	cSpace.OnCollision = onCollision
}

func (cSpace *ChipmonkSpace) AddBody(body physicsAPI.PhysicsObject2D) {
	cBody, ok := body.(*ChipmunkBody)
	if ok {
		cSpace.Space.AddBody(cBody.Body)
	}
}

func (cSpace *ChipmonkSpace) RemoveBody(body physicsAPI.PhysicsObject2D) {
	cBody, ok := body.(*ChipmunkBody)
	if ok {
		cSpace.Space.RemoveBody(cBody.Body)
	}
}

func (cSpace *ChipmonkSpace) SetGravity(gravity mgl32.Vec2) {
	cSpace.Space.Gravity = convertToVect(gravity)
}

func (cSpace *ChipmonkSpace) GetGravity() mgl32.Vec2 {
	return convertFromVect(cSpace.Space.Gravity)
}

func convertFromVect(v vect.Vect) mgl32.Vec2 {
	return mgl32.Vec2{float32(v.X), float32(v.Y)}
}

func convertToVect(v mgl32.Vec2) vect.Vect {
	return vect.Vect{X: vect.Float(v.X()), Y: vect.Float(v.Y())}
}
