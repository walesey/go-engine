package actor

import (
	"github.com/walesey/go-engine/renderer"
)

//an actor that can move around bound by physics (gravity), can jump, can walk/run
type FPSActor struct {
	entity renderer.Entity
}

func (actor *FPSActor) Look(dx, dy float64) {

}
