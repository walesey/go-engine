package physics

import (
	"github.com/walesey/go-engine/vectormath"
)

type Collider interface {
	Overlap(other Collider) bool
}

type BoundingBox struct {
	bounds vectormath.Vector3
}

type PhysicsObject struct {
	Position, Velocity           vectormath.Vector3
	Orientation, AngularVelocity vectormath.Quaternion
	broadPhase, narrowPhase      Collider
}

type PhysicsSpace struct {
	objects []*PhysicsObject
	StepDt  float64
}

func NewPhysicsSpace() *PhysicsSpace {
	return &PhysicsSpace{
		StepDt:  0.018,
		objects: make([]*PhysicsObject, 0, 500),
	}
}

func NewPhysicsObject() *PhysicsObject {
	return &PhysicsObject{
		Position:        vectormath.Vector3{0, 0, 0},
		Velocity:        vectormath.Vector3{0, 0, 0},
		Orientation:     vectormath.IdentityQuaternion(),
		AngularVelocity: vectormath.Quaternion{1, 0, 0, 0},
	}
}

func (ps *PhysicsSpace) PhysicsStep() {
	for _, object := range ps.objects {
		object.doStep(ps.StepDt)
	}
}

func (ps *PhysicsSpace) Add(objects ...*PhysicsObject) {
	ps.objects = append(ps.objects, objects...)
}

func (ps *PhysicsSpace) Remove(objects ...*PhysicsObject) {
	//find the address in the slice
	for _, remove := range objects {
		for index, object := range ps.objects {
			if object == remove {
				if index+1 == len(ps.objects) {
					ps.objects = ps.objects[:index]
				} else {
					ps.objects = append(ps.objects[:index], ps.objects[index+1:]...)
				}
				break
			}
		}
	}
}

func (obj *PhysicsObject) doStep(dt float64) {
	//apply position increment
	obj.Position = obj.Position.Add(obj.Velocity.MultiplyScalar(dt))

	//apply orientation increment
	axis := vectormath.Vector3{obj.AngularVelocity.X, obj.AngularVelocity.Y, obj.AngularVelocity.Z}
	obj.Orientation = vectormath.AngleAxis(dt*obj.AngularVelocity.W, axis).Multiply(obj.Orientation)
}
