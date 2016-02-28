package dynamics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
)

func TestSolveContactConstraint(t *testing.T) {
	constraintSolver := &SequentialImpulseSolver{}

	object1 := &PhysicsObjectImpl{
		position: vmath.Vector3{0, 0, -1},
		velocity: vmath.Vector3{0, 0, 1},
		radius:   1.5,
	}

	object2 := &PhysicsObjectImpl{
		position: vmath.Vector3{0, 0, 1},
		velocity: vmath.Vector3{0, 0, -1},
		radius:   1.5,
	}

	object1.SetMass(1.0)
	object2.SetMass(1.0)

	contactConstraint := &ContactConstraint{
		Body1:         object1,
		Body2:         object2,
		LocalContact1: vmath.Vector3{0, 0, 0.5},
		LocalContact2: vmath.Vector3{0, 0, -0.5},
		Normal:        vmath.Vector3{0, 0, -1},
	}

	constraintSolver.solveContactConstraint(1.0, contactConstraint)
	fmt.Printf("object1.velocity: %v\n", object1.velocity)
	assert.True(t, object1.velocity.ApproxEqual(vmath.Vector3{0, 0, -1}, 0.1))
	fmt.Printf("object2.velocity: %v\n", object2.velocity)
	assert.True(t, object2.velocity.ApproxEqual(vmath.Vector3{0, 0, 1}, 0.1))
}
