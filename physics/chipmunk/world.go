package chipmunkPhysics

import (
	"github.com/vova616/chipmunk"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ChipmonkSpace struct {
	space *chipmunk.Space
}

func NewChipmonkSpace() *ChipmonkSpace {
	return &ChipmonkSpace{}
}

func (cSpace *ChipmonkSpace) Update(dt float64) {

}

func (cSpace *ChipmonkSpace) AddBody(body physicsAPI.PhysicsObject2D) {

}

func (cSpace *ChipmonkSpace) RemoveBody(body physicsAPI.PhysicsObject2D) {

}
