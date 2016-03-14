package physicsAPI

import vmath "github.com/walesey/go-engine/vectormath"

type PhysicsSpace interface {
	SimulateStep(stepTime float64, subSteps int)
	Delete()
	CreateObject() PhysicsObject
	AddObject(objects ...PhysicsObject)
	RemoveObject(objects ...PhysicsObject)
	SetConstraintSolver(solver ConstraintSolver)
	AddConstraint(constraint Constraint)
	RemoveConstraints(constraint ...Constraint)
	SetGravity(gravity vmath.Vector3)
	GetGravity() vmath.Vector3
}
