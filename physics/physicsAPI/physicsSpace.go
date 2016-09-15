package physicsAPI

import "github.com/go-gl/mathgl/mgl32"

type PhysicsSpace interface {
	Update(dt float64)
	SimulateStep(stepTime float32, subSteps int)
	Delete()
	AddObject(objects ...PhysicsObject)
	RemoveObject(objects ...PhysicsObject)
	AddCharacterController(characterController CharacterController)
	SetConstraintSolver(solver ConstraintSolver)
	AddConstraint(constraint Constraint)
	RemoveConstraints(constraint ...Constraint)
	SetGravity(gravity mgl32.Vec3)
	GetGravity() mgl32.Vec3
}

type PhysicsSpace2D interface {
	Update(dt float64)
	AddBody(body PhysicsObject2D)
	RemoveBody(body PhysicsObject2D)
	SetGravity(gravity mgl32.Vec2)
	GetGravity() mgl32.Vec2
}
