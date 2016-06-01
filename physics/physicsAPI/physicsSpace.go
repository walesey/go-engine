package physicsAPI

import vmath "github.com/walesey/go-engine/vectormath"

type PhysicsSpace interface {
	Update(dt float64)
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

type PhysicsSpace2D interface {
	Update(dt float64)
	AddBody(body PhysicsObject2D)
	RemoveBody(body PhysicsObject2D)
	SetGravity(gravity vmath.Vector2)
	GetGravity() vmath.Vector2
}
