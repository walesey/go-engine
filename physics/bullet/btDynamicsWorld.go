package bullet

import (
	"github.com/luxengine/gobullet"
	"github.com/walesey/go-engine/physics/physicsAPI"
	vmath "github.com/walesey/go-engine/vectormath"
)

type BtDynamicsWorld struct {
	world gobullet.DynamicsWorld
}

func NewBtDynamicsWorld(sdk gobullet.PhysicSDK) physicsAPI.PhysicsSpace {
	return &BtDynamicsWorld{sdk.CreateDynamicsWorld()}
}

func (w *BtDynamicsWorld) SimulateStep(stepTime float64, subSteps int) {
	stepDt := stepTime / float64(subSteps)
	for iteration := 0; iteration < subSteps; iteration = iteration + 1 {
		w.world.StepSimulation(float32(stepDt))
	}
}

func (w *BtDynamicsWorld) Delete() {
	w.world.Delete()
}

func (w *BtDynamicsWorld) CreateObject() physicsAPI.PhysicsObject { return nil }

// AddObject Add objects to the world
func (w *BtDynamicsWorld) AddObject(objects ...physicsAPI.PhysicsObject) {
	for _, object := range objects {
		btRigidBody, ok := object.(*BtRigidBody)
		if ok {
			w.world.AddRigidBodyWithGroup(btRigidBody.rigidBody, 1, 1)
		}
	}
}

func (w *BtDynamicsWorld) RemoveObject(objects ...physicsAPI.PhysicsObject) {
	for _, object := range objects {
		btRigidBody, ok := object.(*BtRigidBody)
		if ok {
			w.world.RemoveRigidBody(btRigidBody.rigidBody)
		}
	}
}

func (w *BtDynamicsWorld) AddCharacterController(characterController physicsAPI.CharacterController) {
	btCharacterController, ok := characterController.(*BtCharacterController)
	if ok {
		w.world.AddCollisionObject(btCharacterController.ghost, 1, 1)
		w.world.AddAction(btCharacterController.kcc)
	}
}

func (w *BtDynamicsWorld) SetConstraintSolver(solver physicsAPI.ConstraintSolver) {}
func (w *BtDynamicsWorld) AddConstraint(constraint physicsAPI.Constraint)         {}
func (w *BtDynamicsWorld) RemoveConstraints(constraint ...physicsAPI.Constraint)  {}

func (w *BtDynamicsWorld) SetGravity(gravity vmath.Vector3) {
	setVector(w.world.SetGravity, gravity)
}

func (w *BtDynamicsWorld) GetGravity() vmath.Vector3 {
	return getVector(w.world.GetGravity)
}
