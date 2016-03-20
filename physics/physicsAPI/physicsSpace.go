package physicsAPI

import vmath "github.com/walesey/go-engine/vectormath"

type PhysicsSpace interface {
	SimulateStep(stepTime float64, subSteps int)
	Delete()
	AddObject(objects ...PhysicsObject)
	RemoveObject(objects ...PhysicsObject)
	AddCharacterController(characterController CharacterController)
	SetConstraintSolver(solver ConstraintSolver)
	AddConstraint(constraint Constraint)
	RemoveConstraints(constraint ...Constraint)
	SetGravity(gravity vmath.Vector3)
	GetGravity() vmath.Vector3
}
