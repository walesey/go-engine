package physics

import (
	"github.com/walesey/go-engine/physics/dynamics"
	vmath "github.com/walesey/go-engine/vectormath"
)

type PhysicsSpace struct {
	objects      []*dynamics.PhysicsObject
	objectPool   *dynamics.PhysicsObjectPool
	contactCache ContactCache
	StepDt       float64
	Iterations   int
	GlobalForces *dynamics.ForceStore
}

func NewPhysicsSpace() *PhysicsSpace {
	return &PhysicsSpace{
		StepDt:       0.018,
		Iterations:   1,
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

	for iteration := 0; iteration < ps.Iterations; iteration = iteration + 1 {

		stepDt := ps.StepDt / float64(ps.Iterations)
		for _, object := range ps.objects {
			if !object.Static && object.Active {
				ps.GlobalForces.DoStep(stepDt, object)
				object.DoStep(stepDt)
			}
		}

		ps.contactCache.MarkContactsAsOld()

		//do broadphase overlaps and narrow phase checks for each
		for i, object1 := range ps.objects {
			if !object1.Static && object1.Active {
				for j, object2 := range ps.objects {
					if i != j {
						if object1.BroadPhaseOverlap(object2) {
							if object1.NarrowPhaseOverlap(object2) {

								// activate object
								object2.Active = true

								//check contact cache
								inContact := ps.contactCache.Contains(i, j)
								if !inContact {
									// fmt.Printf("TODO: Contact EVENT %v - %v\n", i, j)
								}
								ps.contactCache.Add(i, j)

								//Collision info
								penV := object1.PenetrationVector(object2)

								//position correction
								object1.Position = object1.Position.Subtract(penV)

								globalContact := object1.ContactPoint(object2)
								localContact1 := globalContact.Subtract(object1.Position)
								localContact2 := globalContact.Subtract(object2.Position)

								//collision normal
								var norm vmath.Vector3
								if penV.LengthSquared() > 0 {
									norm = penV.Normalize()
								} else if !object2.Position.ApproxEqual(object1.Position, 0.00001) {
									norm = object2.Position.Subtract(object1.Position).Normalize()
								} else {
									norm = vmath.Vector3{1, 0, 0}
								}

								//process constraint
								contactConstraint := dynamics.ContactConstraint{
									Normal:        norm,
									LocalContact1: localContact1,
									LocalContact2: localContact2,
									Object1:       object1,
									Object2:       object2,
									InContact:     inContact,
								}
								contactConstraint.Solve()

								//deactivate if moving too slow
								if inContact &&
									object1.Velocity.LengthSquared() <= object1.ActiveVelocity*object1.ActiveVelocity &&
									object1.AngularVelocity.W <= object1.ActiveVelocity {
									object1.Active = true
								}
							}
						}
					}
				}
			}
		}

		ps.contactCache.CleanOldContacts()
	}
}
