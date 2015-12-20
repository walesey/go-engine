package dynamics

type Constraint interface {
	Solve(dt float64)
}
