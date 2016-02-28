package dynamics

type ConstraintSolver interface {
	SolveGroup(stepTime float64, constraints *[]Constraint)
}
