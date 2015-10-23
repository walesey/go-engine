package physics

import (
//	"github.com/walesey/go-engine/vectormath"
)

type PhysicsSpace struct {
	objects []*PhysicsObject
	pool    *WorkerPool
	StepDt  float64
}

func NewPhysicsSpace() *PhysicsSpace {
	return &PhysicsSpace{
		StepDt:  0.018,
		objects: make([]*PhysicsObject, 0, 500),
		pool:    NewWorkerPool(),
	}
}

func (ps *PhysicsSpace) PhysicsStep() {
	for _, object := range ps.objects {
		object.doStep(ps.StepDt)
	}
}

func (ps *PhysicsSpace) Add(objects ...*PhysicsObject) {
	ps.objects = append(ps.objects, objects...)
}

func (ps *PhysicsSpace) Remove(objects ...*PhysicsObject) {
	//find the address in the slice
	for _, remove := range objects {
		for index, object := range ps.objects {
			if object == remove {
				if index+1 == len(ps.objects) {
					ps.objects = ps.objects[:index]
				} else {
					ps.objects = append(ps.objects[:index], ps.objects[index+1:]...)
				}
				break
			}
		}
	}
}
