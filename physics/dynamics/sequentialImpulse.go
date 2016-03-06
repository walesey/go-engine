package dynamics

import "fmt"

// bias factor for baumgarte stabalization
const biasFactor = 0.2
const iterations = 10

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
	// Initialize constraints
	for _, constraint := range *constraints {
		switch t := constraint.(type) {
		default:
			fmt.Printf("unsupported constraint type: %T\n", t)
		case *ContactConstraint:
			s.initContactConstraint(stepTime, constraint.(*ContactConstraint))
		}
	}

	for i := 0; i < iterations; i = i + 1 {
		for _, constraint := range *constraints {
			switch t := constraint.(type) {
			default:
				fmt.Printf("unsupported constraint type: %T\n", t)
			case *ContactConstraint:
				s.solveContactConstraint(stepTime, constraint.(*ContactConstraint))
			}
		}
	}
}

func (s *SequentialImpulseSolver) contactDeltaVdotN(constraint *ContactConstraint) float64 {
	r1 := constraint.LocalContact1
	r2 := constraint.LocalContact2
	v1 := constraint.Body1.GetVelocity()
	v2 := constraint.Body2.GetVelocity()
	w1 := constraint.Body1.GetAngularVelocityVector()
	w2 := constraint.Body2.GetAngularVelocityVector()
	n := constraint.Normal

	//Delta velocity vector at the contact point
	radialV1 := w1.Cross(r1)
	contactV1 := v1.Add(radialV1)
	radialV2 := w2.Cross(r2)
	contactV2 := v2.Add(radialV2)
	contactDeltaV := contactV1.Subtract(contactV2)
	contactDeltaVdotN := contactDeltaV.Dot(n)

	return contactDeltaVdotN
}

func (s *SequentialImpulseSolver) applyImpulse(constraint *ContactConstraint, impulse float64) {
	r1 := constraint.LocalContact1
	r2 := constraint.LocalContact2
	v1 := constraint.Body1.GetVelocity()
	v2 := constraint.Body2.GetVelocity()
	w1 := constraint.Body1.GetAngularVelocityVector()
	w2 := constraint.Body2.GetAngularVelocityVector()
	n := constraint.Normal
	m1 := constraint.Body1.GetMass()
	m2 := constraint.Body2.GetMass()
	I1_inv := constraint.Body1.inertiaTensor.Inverse()
	I2_inv := constraint.Body2.inertiaTensor.Inverse()

	if !constraint.Body1.IsStatic() {
		constraint.Body1.SetVelocity(v1.Add(n.MultiplyScalar(impulse / m1)))
		constraint.Body1.SetAngularVelocityVector(w1.Add(I1_inv.MultiplyVector(n.MultiplyScalar(-impulse).Cross(r1))))
	}
	if !constraint.Body2.IsStatic() {
		constraint.Body2.SetVelocity(v2.Add(n.MultiplyScalar(-impulse / m2)))
		constraint.Body2.SetAngularVelocityVector(w2.Add(I2_inv.MultiplyVector(n.MultiplyScalar(impulse).Cross(r2))))
	}
}

func (s *SequentialImpulseSolver) initContactConstraint(stepTime float64, constraint *ContactConstraint) {
	m1 := constraint.Body1.GetMass()
	m2 := constraint.Body2.GetMass()
	r1 := constraint.LocalContact1
	r2 := constraint.LocalContact2
	n := constraint.Normal
	I1_inv := constraint.Body1.inertiaTensor.Inverse()
	I2_inv := constraint.Body2.inertiaTensor.Inverse()

	if constraint.Body1.IsStatic() {
		m1 = 99999999999999999999999999.9
		I1_inv = sphereInertia(m1, constraint.Body1.GetRadius()).Inverse()
	}
	if constraint.Body2.IsStatic() {
		m2 = 99999999999999999999999999.9
		I2_inv = sphereInertia(m2, constraint.Body2.GetRadius()).Inverse()
	}

	c1 := r1.Cross(n)
	vec1 := (I1_inv.MultiplyVector(c1)).Cross(r1)
	denom1 := (1.0 / m1) + n.Dot(vec1)
	c2 := r2.Cross(n)
	vec2 := (I2_inv.MultiplyVector(c2)).Cross(r2)
	denom2 := (1.0 / m2) + n.Dot(vec2)
	massMultiplier := 1.0 / (denom1 + denom2)

	contactDeltaVdotN := s.contactDeltaVdotN(constraint)
	constraint.contactDeltaVdotN = contactDeltaVdotN

	constraint.targetV = -constraint.Penetration.Dot(n) / stepTime
	constraint.targetV = constraint.targetV - (contactDeltaVdotN * constraint.Restitution)
	velocityError := constraint.targetV - contactDeltaVdotN

	constraint.impulse = velocityError * massMultiplier
	constraint.contactDeltaVdotN = contactDeltaVdotN
	s.applyImpulse(constraint, constraint.impulse)
}

func (s *SequentialImpulseSolver) solveContactConstraint(stepTime float64, constraint *ContactConstraint) {
	contactDeltaVdotN := s.contactDeltaVdotN(constraint)
	netImpulse := constraint.impulse
	netEffect := contactDeltaVdotN - constraint.contactDeltaVdotN
	if netEffect*netEffect > 0 {
		impulse := (constraint.targetV - contactDeltaVdotN) * (netImpulse / netEffect)
		constraint.impulse = constraint.impulse + impulse
		s.applyImpulse(constraint, impulse)
	}
}
