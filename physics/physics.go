package physics

import (
	"fmt"
	"math"

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
								fmt.Println("TODO: Contact EVENT")
							}
							ps.contactCache.Add(i, j)

							//Collision info
							penV := object1.PenetrationVector(object2)

							//collision normal
							var norm vmath.Vector3
							if penV.LengthSquared() > 0 {
								norm = penV.Normalize()
							} else {
								norm = object2.Position.Subtract(object1.Position).Normalize()
							}

							globalContact := object1.ContactPoint(object2)
							fmt.Println(globalContact)
							//velocities
							linearV1 := object1.Velocity.Dot(norm)
							tangentV1 := object1.Velocity.Subtract(norm.MultiplyScalar(linearV1))
							linearV2 := object2.Velocity.Dot(norm)
							tangentV2 := object2.Velocity.Subtract(norm.MultiplyScalar(linearV2))
							//radii / tangent velocity at contact due to angular velocity
							localContact1 := globalContact.Subtract(object1.Position)
							localContact2 := globalContact.Subtract(object2.Position)
							r1 := localContact1.Length()
							r2 := localContact2.Length()
							angularV1 := object1.AngularVelocity
							angularV2 := object2.AngularVelocity
							angularV1axis := vmath.Vector3{angularV1.X, angularV1.Y, angularV1.Z}
							angularV2axis := vmath.Vector3{angularV2.X, angularV2.Y, angularV2.Z}
							contactTangentV1 := localContact1.Cross(angularV1axis)
							contactTangentV2 := localContact2.Cross(angularV2axis)
							if contactTangentV1.LengthSquared() > 0.00000001 {
								contactTangentV1 = contactTangentV1.Normalize().MultiplyScalar(r1 * angularV1.W).Add(tangentV1)
							}
							if contactTangentV2.LengthSquared() > 0.00000001 {
								contactTangentV2 = contactTangentV2.Normalize().MultiplyScalar(r2 * angularV2.W).Add(tangentV2)
							}

							if object2.Static {
								object1.Position = object1.Position.Subtract(penV)
								//linear impulse
								dVl1 := -math.Sqrt(0.5 * linearV1 * linearV1 * object1.Friction)
								object1.Velocity = norm.MultiplyScalar(dVl1).Add(tangentV1)
							} else {
								halfPen := penV.MultiplyScalar(0.5)
								object1.Position = object1.Position.Subtract(halfPen)
								object2.Position = object2.Position.Add(halfPen)
								//linear momentum
								linearVf := (linearV1*object1.Mass + linearV2*object2.Mass) / (object1.Mass + object2.Mass)
								//impulse (bounce) (1/2 mv^2)
								impulse1 := 0.5 * object1.Mass * (linearV1 - linearVf) * (linearV1 - linearVf)
								impulse2 := 0.5 * object2.Mass * (linearV2 - linearVf) * (linearV2 - linearVf)
								impulse := impulse1 + impulse2
								dVl1 := -math.Sqrt(impulse * object1.Friction / object1.Mass)
								dVl2 := math.Sqrt(impulse * object2.Friction / object2.Mass)
								//combine and apply velocities
								object1.Velocity = norm.MultiplyScalar(linearVf + dVl1).Add(tangentV1)
								object2.Velocity = norm.MultiplyScalar(linearVf + dVl2).Add(tangentV2)

								//angular
								if r1 > 0 && r2 > 0 {
									AngMomentum1 := contactTangentV1.MultiplyScalar(r1).MultiplyScalar(object1.Mass)
									AngMomentum2 := contactTangentV2.MultiplyScalar(r2).MultiplyScalar(object2.Mass)
									// angular momentum
									contactTangentVf := AngMomentum1.Add(AngMomentum2).DivideScalar(r1*object1.Mass + r2*object2.Mass)
									// get final angular velocities
									localTangentV1f := contactTangentVf.Subtract(tangentV1)
									localTangentV2f := contactTangentVf.Subtract(tangentV2)
									angularV1AxisF := localContact1.Cross(localTangentV1f)
									angularV2AxisF := localContact2.Cross(localTangentV2f)
									object1.AngularVelocity = vmath.Quaternion{angularV1AxisF.X, angularV1AxisF.Y, angularV1AxisF.Z, -localTangentV1f.Length() / r1}
									object2.AngularVelocity = vmath.Quaternion{angularV2AxisF.X, angularV2AxisF.Y, angularV2AxisF.Z, localTangentV2f.Length() / r2}
								}
							}
						}
					}
				}
			}
		}
	}

	ps.contactCache.CleanOldContacts()
}
