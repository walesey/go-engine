package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

import (
	"reflect"
)

/************************************
************Dynamics world***********
************************************/
type DynamicsWorld struct {
	handle C.plDynamicsWorld
}

//Create a new dynamics world to hold all the rigid bodies, perform simulation and callback
func (this PhysicSDK) CreateDynamicsWorld() DynamicsWorld {
	return DynamicsWorld{C.plCreateDynamicsWorld(this.handle)}
}

//Free the Dynamics world memory
func (this DynamicsWorld) Delete() {
	C.plDeleteDynamicsWorld(this.handle)
}

//Add a rigid body to this dynamics world, will perform collision on the next StepSimulation
func (this DynamicsWorld) AddRigidBody(b RigidBody) {
	C.plAddRigidBody(this.handle, b.handle)
}

//Remove body from this dynamics world
func (this DynamicsWorld) RemoveRigidBody(b RigidBody) {
	C.plRemoveRigidBody(this.handle, b.handle)
}

//timestep: time since last step in SECONDS (so for 60 fps: 0.16666...)
func (this DynamicsWorld) StepSimulation(timeStep float32) {
	C.plStepSimulation(this.handle, C.plReal(timeStep))
}

//timestep: time since last step in SECONDS (so for 60 fps: 0.16666...)
func (this DynamicsWorld) StepSimulationSubStep(timeStep float32, substeps int) {
	C.plStepSimulationSubStep(this.handle, C.plReal(timeStep), C.int(substeps))
}

var sliderconstrainttype = reflect.TypeOf(SliderConstraint{})
var generic6dofconstrainttype = reflect.TypeOf(Generic6DofConstraint{})
var hingeconstraint = reflect.TypeOf(HingeConstraint{})
var point2pointconstraint = reflect.TypeOf(Point2PointConstraint{})
var conetwistconstraint = reflect.TypeOf(ConeTwistConstraint{})

func (this DynamicsWorld) AddConstraint(constraint interface{}, disablecollisionbetweenlinked bool) {
	switch reflect.TypeOf(constraint) {
	case sliderconstrainttype:
		C.plDynamicsWorldAddConstraintSlider(this.handle, constraint.(SliderConstraint).handle, C.bool(disablecollisionbetweenlinked))
	case generic6dofconstrainttype:
		C.plDynamicsWorldAddConstraintG6dof(this.handle, constraint.(Generic6DofConstraint).handle, C.bool(disablecollisionbetweenlinked))
	case hingeconstraint:
		C.plDynamicsWorldAddConstraintHinge(this.handle, constraint.(HingeConstraint).handle, C.bool(disablecollisionbetweenlinked))
	case point2pointconstraint:
		C.plDynamicsWorldAddConstraintP2P(this.handle, constraint.(Point2PointConstraint).handle, C.bool(disablecollisionbetweenlinked))
	case conetwistconstraint:
		C.plDynamicsWorldAddConstraintConeTwist(this.handle, constraint.(ConeTwistConstraint).handle, C.bool(disablecollisionbetweenlinked))
	}
}

//UNTESTED
func (this DynamicsWorld) RayTestClosest(start, end, dest *[3]float32, frac *float32) RigidBody {
	return RigidBody{C.plDynamicsWorldRayTestClosest(this.handle, (*C.plReal)(&start[0]), (*C.plReal)(&end[0]), (*C.plReal)(&dest[0]), (*C.plReal)(frac))}
}

//UNTESTED
func (this DynamicsWorld) AddCollisionObject(ghost PairCachingGhostObject, group, mask int16) {
	C.plDynamicsWorldAddCollisionObject(this.handle, ghost.handle, C.short(group), C.short(mask))
}

//UNTESTED
func (this DynamicsWorld) AddAction(kcc KinematicCharacterController) {
	C.plDynamicsWorldAddAction(this.handle, kcc.handle)
}

//UNTESTED
func (this DynamicsWorld) GetGravity(vec *[3]float32) {
	C.plDynamicsWorldGetGravity(this.handle, (*C.plReal)(&vec[0]))
}

//UNTESTED
func (this DynamicsWorld) SetGravity(vec *[3]float32) {
	C.plDynamicsWorldSetGravity(this.handle, (*C.plReal)(&vec[0]))
}

//UNTESTED
func (this DynamicsWorld) RegisterGImpactAlgorithm() {
	C.plDynamicsWorldRegisterGImpactAlgorithm(this.handle)
}

//UNTESTED
func (this DynamicsWorld) AddRigidBodyWithGroup(rhandle RigidBody, group int16, mask int16) {
	C.plDynamicsWorldAddRigidBodyWithGroup(this.handle, rhandle.handle, C.short(group), C.short(mask))
}
