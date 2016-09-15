package physicsAPI

type ConstraintSolver interface {
	SolveGroup(stepTime float32, constraints *[]Constraint)
}

type Constraint interface{}
