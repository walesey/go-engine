package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

type Constraint interface {
	GetBreakingImpulseThreshold() float32
	SetBreakingImpulseThreshold(threshold float32)
	IsEnabled() bool
	SetEnabled(enabled bool)
	GetRigidBodyA() RigidBody
	GetRigidBodyB() RigidBody
	GetUserConstraintType() int
	SetUserConstraintType(userConstraintType int)
	SetUserConstraintId(uid int)
	GetUserConstraintId() int
	SetUserConstraintPtr(ptr *interface{})
	GetUserConstraintPtr() *interface{}
	GetUid() int
	NeedsFeedback() bool
	EnableFeedback(needsFeedback bool)
	GetAppliedImpulse() float32
	//SetJointFeedback(jointFeedback btJointFeedback*)
	//GetJointFeedback() btJointFeedback
	//GetConstraintType() btTypedConstraintType
}
