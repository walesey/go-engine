package physics

import (
	"fmt"
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
}

func NewPhysicsSpace() *PhysicsSpace {
	return &PhysicsSpace{
		StepDt:       0.018,
		GlobalForces: NewForceStore(),
		objects:      make([]*PhysicsObject, 0, 500),
		workerQueue:  make([]workerQueueItem, 0, 500),
		workerPool:   NewWorkerPool(),
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
		ps.GlobalForces.DoStep(ps.StepDt, object)
		object.doStep(ps.StepDt)
	}

	ps.contactCache.MarkContactsAsOld()

	//do broadphase overlaps and spawn narrow phase workers for each
	queueIndex := 0
	for i, object1 := range ps.objects {
		for j, object2 := range ps.objects {
			if i != j {
				if (*object1).BroadPhaseOverlap(*object2) {
					worker := ps.workerPool.GetWorker()
					worker.Write(PhysicsPair{*object1, *object2})
					ps.workerQueue = append(ps.workerQueue[:queueIndex], workerQueueItem{worker: worker, object1: object1, object2: object2})
					queueIndex = queueIndex + 1
				}
			}
		}
	}

	//read narrow phase results from workers
	for i := 0; i < queueIndex; i++ {
		if ps.workerQueue[i].worker.Read() {

			obj1 := ps.workerQueue[i].object1
			obj2 := ps.workerQueue[i].object2
			pair := PhysicsPair{*obj1, *obj2}
			inContact := ps.contactCache.Contains(pair)
			if !inContact {
				fmt.Println("TODO: Contact EVENT")
				ps.contactCache.Add(pair)
			}

			/*	obj1.doStep(-ps.StepDt * 0.5)
				obj2.doStep(-ps.StepDt * 0.5)

				worker := ps.workerPool.GetWorker()
				worker.Write(PhysicsPair{*object1, *object2})
				ps.workerQueue[queueIndex] = workerQueueItem{worker: worker, object1: object1, object2: object2}
				queueIndex = queueIndex + 1 */
		}
	}

	// do a step at half velocity in reverse
	// recheck collisions
	// handle each collision

	ps.contactCache.CleanOldContacts()
}
