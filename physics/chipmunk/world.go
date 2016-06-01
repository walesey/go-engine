package chipmunkPhysics

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/walesey/go-engine/physics/physicsAPI"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ChipmonkSpace struct {
	space *chipmunk.Space
}

func NewChipmonkSpace() *ChipmonkSpace {
	return &ChipmonkSpace{
		space: chipmunk.NewSpace(),
	}
}

func (cSpace *ChipmonkSpace) Update(dt float64) {
	cSpace.space.Step(vect.Float(dt))
}

func (cSpace *ChipmonkSpace) AddBody(body physicsAPI.PhysicsObject2D) {
	cBody, ok := body.(*ChipmunkBody)
	if ok {
		cSpace.space.AddBody(cBody.Body)
	}
}

func (cSpace *ChipmonkSpace) RemoveBody(body physicsAPI.PhysicsObject2D) {
	cBody, ok := body.(*ChipmunkBody)
	if ok {
		cSpace.space.RemoveBody(cBody.Body)
	}
}

func (cSpace *ChipmonkSpace) SetGravity(gravity vmath.Vector2) {
	cSpace.space.Gravity = convertToVect(gravity)
}

func (cSpace *ChipmonkSpace) GetGravity() vmath.Vector2 {
	return convertFromVect(cSpace.space.Gravity)
}

func convertFromVect(v vect.Vect) vmath.Vector2 {
	return vmath.Vector2{X: float64(v.X), Y: float64(v.Y)}
}

func convertToVect(v vmath.Vector2) vect.Vect {
	return vect.Vect{X: vect.Float(v.X), Y: vect.Float(v.Y)}
}
