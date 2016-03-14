package physicsAPI

type ConstraintSolver interface {
	SolveGroup(stepTime float64, constraints *[]Constraint)
}

type Constraint interface{}
