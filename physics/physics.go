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

// CreateObject create a new object and add it to the world
func (ps *PhysicsSpace) CreateObject() *PhysicsObject {
	object := ps.objectPool.GetPhysicsObject()
	ps.objects = append(ps.objects, object)
	return object
}

// Remove remove objects from the world
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
			object.doStep(ps.StepDt)
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

							//collision normal
							var norm vmath.Vector3
							if penV.LengthSquared() > 0 {
								norm = penV.Normalize()
							} else if !object2.Position.ApproxEqual(object1.Position, 0.00001) {
								norm = object2.Position.Subtract(object1.Position).Normalize()
							} else {
								norm = vmath.Vector3{1, 0, 0}
							}

							globalContact := object1.ContactPoint(object2)
							localContact1 := globalContact.Subtract(object1.Position)
							localContact2 := globalContact.Subtract(object2.Position)

							//velocities
							angularV1 := object1.AngularVelocityVector()
							angularV2 := object2.AngularVelocityVector()
							radialV1 := localContact1.Cross(angularV1)
							radialV2 := localContact2.Cross(angularV2)
							contactV1 := radialV1.Add(object1.Velocity)
							contactV2 := radialV2.Add(object2.Velocity)

							if object2.Static {
								object1.Position = object1.Position.Subtract(penV)

							} else {
								halfPen := penV.MultiplyScalar(0.5)
								object1.Position = object1.Position.Subtract(halfPen)
								object2.Position = object2.Position.Add(halfPen)

								mR1 := 0.4 * object1.Mass * object1.Radius
								mR2 := 0.4 * object2.Mass * object2.Radius
								tensor1 := vmath.Matrix3{
									mR1, 0.0, 0.0,
									0.0, mR1, 0.0,
									0.0, 0.0, mR1,
								}
								tensor2 := vmath.Matrix3{
									mR2, 0.0, 0.0,
									0.0, mR2, 0.0,
									0.0, 0.0, mR2,
								}

								contactV := contactV1.Subtract(contactV2)
								relativeV := norm.Dot(contactV)
								velocityImpulse := -relativeV

								vel1 := tensor1.Inverse().Transform(localContact1.Cross(norm))
								vel1 = vel1.Cross(localContact1)
								impulseDenom1 := (1.0 / object1.Mass) + norm.Dot(vel1)
								vel2 := tensor2.Inverse().Transform(localContact2.Cross(norm))
								vel2 = vel2.Cross(localContact2)
								impulseDenom2 := (1.0 / object2.Mass) + norm.Dot(vel2)
								impulseDenom := impulseDenom1 + impulseDenom2
								normalImpulse := velocityImpulse / impulseDenom
								impulseVector1 := norm.MultiplyScalar(normalImpulse)
								impulseVector2 := impulseVector1.MultiplyScalar(-1)
								torqueImpulse1 := localContact1.Cross(impulseVector1).MultiplyScalar(-1)
								torqueImpulse2 := localContact2.Cross(impulseVector2).MultiplyScalar(-1)

								object1.Velocity = object1.Velocity.Add(impulseVector1.DivideScalar(object1.Mass))
								object2.Velocity = object2.Velocity.Add(impulseVector2.DivideScalar(object2.Mass))

								newAngularV1 := angularV1.Add(tensor1.Inverse().Transform(torqueImpulse1))
								newAngularV2 := angularV2.Add(tensor2.Inverse().Transform(torqueImpulse2))
								object1.SetAngularVelocityVector(newAngularV1)
								object2.SetAngularVelocityVector(newAngularV2)
							}
						}
					}
				}
			}
		}
	}

	ps.contactCache.CleanOldContacts()
}
