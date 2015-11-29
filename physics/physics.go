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

	pointForce := new(PointForce)

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

							// norm := object1.Position.Subtract(object2.Position)
							// if norm.LengthSquared() < 0.0000001 {
							// 	norm = vmath.Vector3{1, 0, 0}
							// } else {
							// 	norm = norm.Normalize()
							// }

							globalContact := object1.ContactPoint(object2)
							// angular velocity
							localContact1 := globalContact.Subtract(object1.Position)
							localContact2 := globalContact.Subtract(object2.Position)
							angularV1 := object1.AngularVelocityVector()
							angularV2 := object2.AngularVelocityVector()

							if object2.Static {
								object1.Position = object1.Position.Subtract(penV)
								e := 0.1
								impulse := SphericalCollisionResponse(
									e, object1.Mass, 999999999999999999.9,
									localContact1, localContact2, norm,
									object1.Velocity, vmath.Vector3{0, 0, 0}, angularV1, vmath.Vector3{0, 0, 0},
								)

								pointForce.Value = norm.MultiplyScalar(-impulse)
								pointForce.Position = localContact1
								pointForce.DoStep(1.0, object1)

							} else {
								halfPen := penV.MultiplyScalar(0.5)
								object1.Position = object1.Position.Subtract(halfPen)
								object2.Position = object2.Position.Add(halfPen)

								e := 0.1
								impulse := SphericalCollisionResponse(
									e, object1.Mass, object2.Mass,
									localContact1, localContact2, norm,
									object1.Velocity, object2.Velocity, angularV1, angularV2,
								)

								pointForce.Value = norm.MultiplyScalar(-impulse)
								pointForce.Position = localContact1
								pointForce.DoStep(1.0, object1)

								pointForce.Value = norm.MultiplyScalar(impulse)
								pointForce.Position = localContact2
								pointForce.DoStep(1.0, object2)
							}
						}
					}
				}
			}
		}
	}

	ps.contactCache.CleanOldContacts()
}

// SphericalCollisionResponse - collision response assuming both bodies are spherical
func SphericalCollisionResponse(e, ma, mb float64, ra, rb, n, vai, vbi, wai, wbi vmath.Vector3) float64 {
	Ia := vmath.Matrix3{
		0.4 * ma * ra.LengthSquared(), 0.0, 0.0,
		0.0, 0.4 * ma * ra.LengthSquared(), 0.0,
		0.0, 0.0, 0.4 * ma * ra.LengthSquared(),
	}
	Ib := vmath.Matrix3{
		0.4 * mb * rb.LengthSquared(), 0.0, 0.0,
		0.0, 0.4 * mb * rb.LengthSquared(), 0.0,
		0.0, 0.0, 0.4 * mb * rb.LengthSquared(),
	}
	return CollisionResponse(e, ma, mb, Ia, Ib, ra, rb, n, vai, vbi, wai, wbi)
}

// CollisionResponse calculates the collision impulse of 2 objects
// e coefficient of restitution which depends on the nature of the two colliding materials ( e<1.0 && e>0.0 )
// ma/mb mass of the objects
// Ia/Ib inertial tensor in global coordinates
// ra/rb position of contact relative to the centre of mass
// n normal of collision
// vai/vbi initial velocity
// wai/wbi initial angular velocity
func CollisionResponse(e, ma, mb float64, Ia, Ib vmath.Matrix3, ra, rb, n, vai, vbi, wai, wbi vmath.Vector3) float64 {
	IaInverse := Ia.Inverse()
	IbInverse := Ib.Inverse()
	deltaVa := ra.Cross(n).Dot(IaInverse.Transform(ra.Cross(n)))
	deltaVb := rb.Cross(n).Dot(IbInverse.Transform(rb.Cross(n)))

	impulse := (e + 1.0) * (vai.Subtract(vbi).Dot(n) + ra.Cross(n).Dot(wai) - rb.Cross(n).Dot(wbi))
	impulse = impulse / (1.0/ma + 1.0/mb + deltaVa + deltaVb)

	return impulse
}
