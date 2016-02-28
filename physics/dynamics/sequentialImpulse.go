package dynamics

import (
	"fmt"

	vmath "github.com/walesey/go-engine/vectormath"
)

// bias factor for baumgarte stabalization
const biasFactor = 1.0

// SequentialImpulseSolver - solver that uses iterative impulses to solve groups of constraints simultaniously
type SequentialImpulseSolver struct {
	stepTime float64
}

// NewSequentialImpulseSolver - create a new instance of the SequentialImpulseSolver
func NewSequentialImpulseSolver() ConstraintSolver {
	return &SequentialImpulseSolver{}
}

//SolveGroup Solve a group of constraints simultaniously
func (s *SequentialImpulseSolver) SolveGroup(stepTime float64, constraints *[]Constraint) {
	for _, constraint := range *constraints {
		switch t := constraint.(type) {
		default:
			fmt.Printf("unsupported constraint type: %T\n", t)
		case *ContactConstraint:
			s.solveContactConstraint(stepTime, constraint.(*ContactConstraint))
		}
	}
}

func (s *SequentialImpulseSolver) solveContactConstraint(stepTime float64, constraint *ContactConstraint) {
	m1 := constraint.Body1.GetMass()
	m2 := constraint.Body2.GetMass()
	r1 := constraint.LocalContact1
	r2 := constraint.LocalContact2
	x1 := constraint.Body1.GetPosition()
	x2 := constraint.Body2.GetPosition()
	v1 := constraint.Body1.GetVelocity()
	v2 := constraint.Body2.GetVelocity()
	w1 := constraint.Body1.GetAngularVelocityVector()
	w2 := constraint.Body2.GetAngularVelocityVector()
	n := constraint.Normal
	I1_inv := constraint.Body1.inertiaTensor.Inverse()
	I2_inv := constraint.Body2.inertiaTensor.Inverse()

	if constraint.Body1.IsStatic() {
		m1 = 999999999999999999.9
		I1_inv = sphereInertia(m1, constraint.Body1.GetRadius()).Inverse()
	}
	if constraint.Body2.IsStatic() {
		m2 = 999999999999999999.9
		I2_inv = sphereInertia(m2, constraint.Body2.GetRadius()).Inverse()
	}

	// Solve contact constraint variables
	K := (1.0 / m1) + (1.0 / m2) + vmath.RowMat3ColumnProduct(I1_inv, r1.Cross(n), r1.Cross(n)) + vmath.RowMat3ColumnProduct(I2_inv, r2.Cross(n), r2.Cross(n))
	J_v1 := n.MultiplyScalar(-1)
	J_w1 := r1.Cross(n).MultiplyScalar(-1)
	J_v2 := n
	J_w2 := r2.Cross(n)

	// baumgarte stabalization
	b := x2.Add(r2).Subtract(x1).Subtract(r1).Dot(n) * (biasFactor / stepTime)

	// find lagrange Multiplier
	lm := -K * (b + J_v1.Dot(v1) + J_w1.Dot(w1) + J_v2.Dot(v2) + J_w2.Dot(w2))

	// Find force and torque impulses
	f1 := J_v1.MultiplyScalar(lm)
	t1 := J_w1.MultiplyScalar(lm)
	f2 := J_v2.MultiplyScalar(lm)
	t2 := J_w2.MultiplyScalar(lm)

	fmt.Printf("b: %v\n", b)
	fmt.Printf("lm: %v\n", lm)
	fmt.Printf("K: %v\n", K)
	fmt.Printf("f1: %v\n", f1)
	fmt.Printf("t1: %v\n", t1)
	fmt.Printf("f2: %v\n", f2)
	fmt.Printf("t2: %v\n", t2)

	// apply impulses
	constraint.Body1.SetVelocity(v1.Add(f1.DivideScalar(m1)))
	constraint.Body2.SetVelocity(v2.Add(f2.DivideScalar(m2)))
	constraint.Body1.SetAngularVelocityVector(w1.Add(I1_inv.MultiplyVector(t1)))
	constraint.Body2.SetAngularVelocityVector(w2.Add(I2_inv.MultiplyVector(t2)))
}
