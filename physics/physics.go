package physics

import (
	"fmt"
	vmath "github.com/walesey/go-engine/vectormath"
)

type PhysicsSpace struct {
	objects      []*PhysicsObject
	workerQueue  []workerQueueItem
	workerPool   *WorkerPool
	objectPool   *PhysicsObjectPool
	contactCache ContactCache
	StepDt       float64
	GlobalForces *ForceStore
}

type workerQueueItem struct {
	worker           *PhysicsWorker
	object1, object2 *PhysicsObject
	index1, index2   int
}

func NewPhysicsSpace() *PhysicsSpace {
	return &PhysicsSpace{
		StepDt:       0.018,
		GlobalForces: NewForceStore(),
		objects:      make([]*PhysicsObject, 0, 500),
		objectPool:   NewPhysicsObjectPool(),
		contactCache: NewContactCache(),
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

func (ps *PhysicsSpace) DoStep() {

	//do standard movement step
	for _, object := range ps.objects {
		if !object.Static {
			ps.GlobalForces.DoStep(ps.StepDt, object)
			object.doStep(ps.StepDt)
		}
	}

	ps.contactCache.MarkContactsAsOld()

	//do broadphase overlaps and narrow phase checks for each
	for i, object1 := range ps.objects {
		for j, object2 := range ps.objects {
			if i != j {
				if object1.BroadPhaseOverlap(object2) {
					if object1.NarrowPhaseOverlap(object2) {
						//check contact cache
						inContact := ps.contactCache.Contains(i, j)
						if !inContact {
							fmt.Println("TODO: Contact EVENT")
						}
						ps.contactCache.Add(i, j)

						//Collision response
						if !object1.Static {
							object1.doStep(-ps.StepDt * 0.5)
							object1.Velocity = vmath.Vector3{0, 0, 0}
						}
						if !object2.Static {
							object2.doStep(-ps.StepDt * 0.5)
							object2.Velocity = vmath.Vector3{0, 0, 0}
						}
					}
				}
			}
		}
	}

	ps.contactCache.CleanOldContacts()
}
