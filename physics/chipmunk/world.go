package chipmunkPhysics

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/walesey/go-engine/physics/physicsAPI"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ChipmonkSpace struct {
	Space *chipmunk.Space
}

func NewChipmonkSpace() *ChipmonkSpace {
	return &ChipmonkSpace{
		Space: chipmunk.NewSpace(),
	}
}

func (cSpace *ChipmonkSpace) Update(dt float64) {
	cSpace.Space.Step(vect.Float(dt))
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

func (cSpace *ChipmonkSpace) SetGravity(gravity vmath.Vector2) {
	cSpace.Space.Gravity = convertToVect(gravity)
}

func (cSpace *ChipmonkSpace) GetGravity() vmath.Vector2 {
	return convertFromVect(cSpace.Space.Gravity)
}

func convertFromVect(v vect.Vect) vmath.Vector2 {
	return vmath.Vector2{X: float64(v.X), Y: float64(v.Y)}
}

func convertToVect(v vmath.Vector2) vect.Vect {
	return vect.Vect{X: vect.Float(v.X), Y: vect.Float(v.Y)}
}
