package dynamics

import (
	"github.com/walesey/go-engine/physics/physicsAPI"
	vmath "github.com/walesey/go-engine/vectormath"
)

type PhysicsSpaceImpl struct {
	objects          []*PhysicsObjectImpl
	objectPool       *PhysicsObjectPool
	contactCache     ContactCache
	constraintSolver physicsAPI.ConstraintSolver
	constraints      []physicsAPI.Constraint
	OnEvent          func(Event)
	gravity          vmath.Vector3
}

func NewPhysicsSpace() physicsAPI.PhysicsSpace {
	return &PhysicsSpaceImpl{
		objects:          make([]*PhysicsObjectImpl, 0, 500),
		objectPool:       NewPhysicsObjectPool(),
		contactCache:     NewContactCache(),
		constraints:      make([]physicsAPI.Constraint, 0, 500),
		gravity:          vmath.Vector3{0, -10, 0},
		constraintSolver: NewSequentialImpulseSolver(),
	}
}

func (ps *PhysicsSpaceImpl) Delete() {}

// CreateObject create a new object and add it to the world
func (ps *PhysicsSpaceImpl) CreateObject() physicsAPI.PhysicsObject {
	object := ps.objectPool.GetPhysicsObject()
	ps.objects = append(ps.objects, object)
	return object
}

// AddObject Add objects to the world
func (ps *PhysicsSpaceImpl) AddObject(objects ...physicsAPI.PhysicsObject) {
	for _, object := range objects {
		pObject, ok := object.(*PhysicsObjectImpl)
		if ok {
			ps.objects = append(ps.objects, pObject)
		}
	}
}

// RemoveObject remove objects from the world
func (ps *PhysicsSpaceImpl) RemoveObject(objects ...physicsAPI.PhysicsObject) {
	for _, object := range objects {
		for index, phyObj := range ps.objects {
			if phyObj == object {
				ps.objects = append(ps.objects[:index], ps.objects[index+1:]...)
				ps.objectPool.ReleasePhysicsObject(phyObj)
			}
		}
	}
	ps.contactCache.Clear()
}

func (ps *PhysicsSpaceImpl) SetConstraintSolver(solver physicsAPI.ConstraintSolver) {
	ps.constraintSolver = solver
}

func (ps *PhysicsSpaceImpl) AddConstraint(constraint physicsAPI.Constraint) {
	ps.constraints = append(ps.constraints, constraint)
}

func (ps *PhysicsSpaceImpl) RemoveConstraints(constraint ...physicsAPI.Constraint) {
	//TODO
}

func (ps *PhysicsSpaceImpl) SetGravity(gravity vmath.Vector3) {
	ps.gravity = gravity
}

func (ps *PhysicsSpaceImpl) GetGravity() vmath.Vector3 {
	return ps.gravity
}

func (ps *PhysicsSpaceImpl) applyGravity() {
	for _, object := range ps.objects {
		object.ApplyGravity(ps.gravity)
	}
}

func (ps *PhysicsSpaceImpl) clearForces() {
	for _, object := range ps.objects {
		object.totalForce = vmath.Vector3{0, 0, 0}
		object.totalTorque = vmath.Vector3{0, 0, 0}
	}
}

// SimulateStep DoStep update simulate a period of time with a number of subSteps
func (ps *PhysicsSpaceImpl) SimulateStep(stepTime float64, subSteps int) {
	stepDt := stepTime / float64(subSteps)

	ps.applyGravity()

	for iteration := 0; iteration < subSteps; iteration = iteration + 1 {
		ps.simulateSingleStep(stepDt)
	}

	ps.clearForces()
}

func (ps *PhysicsSpaceImpl) simulateSingleStep(stepTime float64) {

	//predict motion
	for _, object := range ps.objects {
		if !object.IsStatic() && object.IsActive() {
			object.ApplyGravity(ps.gravity)
			object.DoStep(stepTime)
		}
	}

	ps.contactCache.MarkContactsAsOld()

	//TODO: find a way to reuse constraints
	ps.constraints = ps.constraints[:0]

	//do broadphase overlaps and narrow phase checks for each
	for i, object1 := range ps.objects {
		if !object1.IsStatic() && object1.IsActive() {
			for j, object2 := range ps.objects {
				if i != j {
					if object1.BroadPhaseOverlap(object2) {
						if object1.NarrowPhaseOverlap(object2) {

							// activate object
							object2.active = true

							//check contact cache
							inContact := ps.contactCache.Contains(i, j)
							if inContact {

							} else {
								ps.contactCache.Add(i, j)
							}

							penV := object1.PenetrationVector(object2)
							globalContact := object1.ContactPoint(object2)
							localContact1 := globalContact.Subtract(object1.position).Add(penV)
							localContact2 := globalContact.Subtract(object2.position)

							//collision normal
							norm := vmath.Vector3{1, 0, 0}
							if penV.LengthSquared() > 0 {
								norm = penV.Normalize()
							}

							//Constraint info
							contactConstraint := &ContactConstraint{
								BodyIndex1:    i,
								BodyIndex2:    j,
								Body1:         object1,
								Body2:         object2,
								LocalContact1: localContact1,
								LocalContact2: localContact2,
								Penetration:   penV,
								Normal:        norm,
								Restitution:   0.0,
							}
							ps.AddConstraint(contactConstraint)

							//TODO: deactivate if moving too slow

							// event handler
							if !inContact && ps.OnEvent != nil {
								ps.OnEvent(CollisionEvent(globalContact))
							}
						}
					}
				}
			}
		}
	}

	ps.constraintSolver.SolveGroup(stepTime, &ps.constraints)

	ps.contactCache.CleanOldContacts()
}
