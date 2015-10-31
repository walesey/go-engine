package physics

import (
	//"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	contactCache := NewContactCache()
	obj1 := NewPhysicsObject()
	obj2 := NewPhysicsObject()
	pair := PhysicsPair{obj1, obj2}
	contactCache.Add(pair)
}
