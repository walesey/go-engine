package dynamics

import vmath "github.com/walesey/go-engine/vectormath"

type PhysicsSpace interface {
	SimulateStep(stepTime float64, subSteps int)
	CreateObject() PhysicsObject
	RemoveObject(objects ...PhysicsObject)
	SetConstraintSolver(solver ConstraintSolver)
	AddConstraint(constraint Constraint)
	RemoveConstraints(constraint ...Constraint)
	SetGravity(gravity vmath.Vector3)
	GetGravity() vmath.Vector3
}

type PhysicsSpaceImpl struct {
	objects          []*PhysicsObjectImpl
	objectPool       *PhysicsObjectPool
	contactCache     ContactCache
	constraintSolver ConstraintSolver
	constraints      []Constraint
	OnEvent          func(Event)
	gravity          vmath.Vector3
}

func NewPhysicsSpace() PhysicsSpace {
	return &PhysicsSpaceImpl{
		objects:      make([]*PhysicsObjectImpl, 0, 500),
		objectPool:   NewPhysicsObjectPool(),
		contactCache: NewContactCache(),
		constraints:  make([]Constraint, 0, 500),
		gravity:      vmath.Vector3{0, -10, 0},
	}
}

// CreateObject create a new object and add it to the world
func (ps *PhysicsSpaceImpl) CreateObject() PhysicsObject {
	object := ps.objectPool.GetPhysicsObject()
	ps.objects = append(ps.objects, object)
	return object
}

// Remove remove objects from the world
func (ps *PhysicsSpaceImpl) RemoveObject(objects ...PhysicsObject) {
	//TODO
}

func (ps *PhysicsSpaceImpl) SetConstraintSolver(solver ConstraintSolver) {
	ps.constraintSolver = solver
}

func (ps *PhysicsSpaceImpl) AddConstraint(constraint Constraint) {
	ps.constraints = append(ps.constraints, constraint)
}

func (ps *PhysicsSpaceImpl) RemoveConstraints(constraint ...Constraint) {
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

// DoStep update simulate a period of time with a number of subSteps
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
		if !object.static && object.active {
			object.DoStep(stepTime)
		}
	}

	ps.contactCache.MarkContactsAsOld()

	//do broadphase overlaps and narrow phase checks for each
	for i, object1 := range ps.objects {
		if !object1.static && object1.active {
			for j, object2 := range ps.objects {
				if i != j {
					if object1.BroadPhaseOverlap(object2) {
						if object1.NarrowPhaseOverlap(object2) {

							// activate object
							object2.active = true

							//check contact cache
							inContact := ps.contactCache.Contains(i, j)
							ps.contactCache.Add(i, j)

							//Collision info
							penV := object1.PenetrationVector(object2)

							//position correction
							object1.position = object1.position.Subtract(penV)

							globalContact := object1.ContactPoint(object2)
							// localContact1 := globalContact.Subtract(object1.position)
							// localContact2 := globalContact.Subtract(object2.position)

							//collision normal
							if penV.LengthSquared() > 0 {
								// norm := penV.Normalize()

								//TODO: Add constraint
							}

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

	ps.contactCache.CleanOldContacts()
}
