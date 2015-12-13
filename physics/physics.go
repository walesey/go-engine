package physics

import (
	"fmt"

	"github.com/walesey/go-engine/physics/dynamics"
)

type PhysicsSpace struct {
	objects      []*dynamics.PhysicsObject
	objectPool   *dynamics.PhysicsObjectPool
	contactCache ContactCache
	StepDt       float64
	GlobalForces *dynamics.ForceStore
}

func NewPhysicsSpace() *PhysicsSpace {
	return &PhysicsSpace{
		StepDt:       0.018,
		GlobalForces: dynamics.NewForceStore(),
		objects:      make([]*dynamics.PhysicsObject, 0, 500),
		objectPool:   dynamics.NewPhysicsObjectPool(),
		contactCache: NewContactCache(),
	}
}

// CreateObject create a new object and add it to the world
func (ps *PhysicsSpace) CreateObject() *dynamics.PhysicsObject {
	object := ps.objectPool.GetPhysicsObject()
	ps.objects = append(ps.objects, object)
	return object
}

// Remove remove objects from the world
func (ps *PhysicsSpace) Remove(objects ...*dynamics.PhysicsObject) {
	//find the address in the slice
	for _, remove := range objects {
		for index, object := range ps.objects {
			if object == remove {
				if index+1 == len(ps.objects) {
					ps.objects = ps.objects[:index]
				} else {
					ps.objects = append(ps.objects[:index], ps.objects[index+1:]...)
				}
				ps.objectPool.ReleasePhysicsObject(object)
				break
			}
		}
	}
}

// DoStep update all objects
func (ps *PhysicsSpace) DoStep() {

	//do standard movement step
	for _, object := range ps.objects {
		if !object.Static {
			ps.GlobalForces.DoStep(ps.StepDt, object)
			object.DoStep(ps.StepDt)
		}
	}

	ps.contactCache.MarkContactsAsOld()

	//do broadphase overlaps and narrow phase checks for each
	for i, object1 := range ps.objects {
		if !object1.Static {
			for j, object2 := range ps.objects {
				if i != j {
					if object1.BroadPhaseOverlap(object2) {
						if object1.NarrowPhaseOverlap(object2) {

							//check contact cache
							inContact := ps.contactCache.Contains(i, j)
							if !inContact {
								fmt.Printf("TODO: Contact EVENT %v - %v\n", i, j)
							}
							ps.contactCache.Add(i, j)

							//Collision info
							penV := object1.PenetrationVector(object2)

							if object2.Static {
								object1.Position = object1.Position.Subtract(penV)
							} else {
								halfPen := penV.MultiplyScalar(0.5)
								object1.Position = object1.Position.Subtract(halfPen)
								object2.Position = object2.Position.Add(halfPen)
							}

							globalContact := object1.ContactPoint(object2)
							localContact1 := globalContact.Subtract(object1.Position)
							localContact2 := globalContact.Subtract(object2.Position)

							contactConstraint := dynamics.ContactConstraint{
								PenetrationVector: penV,
								LocalContact1:     localContact1,
								LocalContact2:     localContact2,
								Object1:           object1,
								Object2:           object2,
							}
							contactConstraint.Solve()
						}
					}
				}
			}
		}
	}

	ps.contactCache.CleanOldContacts()
}
