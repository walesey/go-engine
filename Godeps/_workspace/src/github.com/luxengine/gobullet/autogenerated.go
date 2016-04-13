
package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

import (
	"unsafe"
)

//UNTESTED
func (this Point2PointConstraint) GetBreakingImpulseThreshold() float32 {
	return float32(C.plPoint2PointConstraintGetBreakingImpulseThreshold(this.handle))
}

//UNTESTED
func (this Point2PointConstraint) SetBreakingImpulseThreshold(threshold float32) {
	C.plPoint2PointConstraintSetBreakingImpulseThreshold(this.handle, C.plReal(threshold))
}

//UNTESTED
func (this Point2PointConstraint) IsEnabled() bool {
	return bool(C.plPoint2PointConstraintIsEnabled(this.handle))
}

//UNTESTED
func (this Point2PointConstraint) SetEnabled(enabled bool) {
	C.plPoint2PointConstraintSetEnabled(this.handle, C.bool(enabled))
}

//UNTESTED
func (this Point2PointConstraint) GetRigidBodyA() RigidBody {
	return RigidBody{C.plPoint2PointConstraintGetRigidBodyA(this.handle)}
}

//UNTESTED
func (this Point2PointConstraint) GetRigidBodyB() RigidBody {
	return RigidBody{C.plPoint2PointConstraintGetRigidBodyB(this.handle)}
}

//UNTESTED
func (this Point2PointConstraint) GetUserConstraintType() int {
	return int(C.plPoint2PointConstraintGetUserConstraintType(this.handle))
}

//UNTESTED
func (this Point2PointConstraint) SetUserConstraintType(userConstraintType int) {
	C.plPoint2PointConstraintSetUserConstraintType(this.handle, C.int(userConstraintType))
}

//UNTESTED
func (this Point2PointConstraint) SetUserConstraintId(uid int) {
	C.plPoint2PointConstraintSetUserConstraintId(this.handle, C.int(uid))
}

//UNTESTED
func (this Point2PointConstraint) GetUserConstraintId() int {
	return int(C.plPoint2PointConstraintGetUserConstraintId(this.handle))
}

//UNTESTED
func (this Point2PointConstraint) SetUserConstraintPtr(ptr *interface{}) {
	C.plPoint2PointConstraintSetUserConstraintPtr(this.handle, unsafe.Pointer(ptr))
}

//UNTESTED
func (this Point2PointConstraint) GetUserConstraintPtr() *interface{}{
	return (*interface{})(C.plPoint2PointConstraintGetUserConstraintPtr(this.handle))
}

//UNTESTED
func (this Point2PointConstraint) GetUid() int {
	return int(C.plPoint2PointConstraintGetUid(this.handle))
}

//UNTESTED
func (this Point2PointConstraint) NeedsFeedback() bool {
	return bool(C.plPoint2PointConstraintNeedsFeedback(this.handle))
}

//UNTESTED
func (this Point2PointConstraint) EnableFeedback(needsFeedback bool) {
	C.plPoint2PointConstraintEnableFeedback(this.handle, C.bool(needsFeedback))
}

//UNTESTED
func (this Point2PointConstraint) GetAppliedImpulse() float32 {
	return float32(C.plPoint2PointConstraintGetAppliedImpulse(this.handle))
}

/*
//UNTESTED
func (this Point2PointConstraint) SetJointFeedback(jointFeedback btJointFeedback*) {
	C.plPoint2PointConstraintSetJointFeedback(this.handle)
}

//UNTESTED
func (this Point2PointConstraint) GetJointFeedback() btJointFeedback {
	return btJointFeedback(C.plPoint2PointConstraintGetJointFeedback(this.handle))
}

//UNTESTED
func (this Point2PointConstraint) GetConstraintType() btTypedConstraintType {
	return btTypedConstraintType(C.plPoint2PointConstraintGetConstraintType(this.handle))
}
*/



//UNTESTED
func (this HingeConstraint) GetBreakingImpulseThreshold() float32 {
	return float32(C.plHingeConstraintGetBreakingImpulseThreshold(this.handle))
}

//UNTESTED
func (this HingeConstraint) SetBreakingImpulseThreshold(threshold float32) {
	C.plHingeConstraintSetBreakingImpulseThreshold(this.handle, C.plReal(threshold))
}

//UNTESTED
func (this HingeConstraint) IsEnabled() bool {
	return bool(C.plHingeConstraintIsEnabled(this.handle))
}

//UNTESTED
func (this HingeConstraint) SetEnabled(enabled bool) {
	C.plHingeConstraintSetEnabled(this.handle, C.bool(enabled))
}

//UNTESTED
func (this HingeConstraint) GetRigidBodyA() RigidBody {
	return RigidBody{C.plHingeConstraintGetRigidBodyA(this.handle)}
}

//UNTESTED
func (this HingeConstraint) GetRigidBodyB() RigidBody {
	return RigidBody{C.plHingeConstraintGetRigidBodyB(this.handle)}
}

//UNTESTED
func (this HingeConstraint) GetUserConstraintType() int {
	return int(C.plHingeConstraintGetUserConstraintType(this.handle))
}

//UNTESTED
func (this HingeConstraint) SetUserConstraintType(userConstraintType int) {
	C.plHingeConstraintSetUserConstraintType(this.handle, C.int(userConstraintType))
}

//UNTESTED
func (this HingeConstraint) SetUserConstraintId(uid int) {
	C.plHingeConstraintSetUserConstraintId(this.handle, C.int(uid))
}

//UNTESTED
func (this HingeConstraint) GetUserConstraintId() int {
	return int(C.plHingeConstraintGetUserConstraintId(this.handle))
}

//UNTESTED
func (this HingeConstraint) SetUserConstraintPtr(ptr *interface{}) {
	C.plHingeConstraintSetUserConstraintPtr(this.handle, unsafe.Pointer(ptr))
}

//UNTESTED
func (this HingeConstraint) GetUserConstraintPtr() *interface{}{
	return (*interface{})(C.plHingeConstraintGetUserConstraintPtr(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetUid() int {
	return int(C.plHingeConstraintGetUid(this.handle))
}

//UNTESTED
func (this HingeConstraint) NeedsFeedback() bool {
	return bool(C.plHingeConstraintNeedsFeedback(this.handle))
}

//UNTESTED
func (this HingeConstraint) EnableFeedback(needsFeedback bool) {
	C.plHingeConstraintEnableFeedback(this.handle, C.bool(needsFeedback))
}

//UNTESTED
func (this HingeConstraint) GetAppliedImpulse() float32 {
	return float32(C.plHingeConstraintGetAppliedImpulse(this.handle))
}

/*
//UNTESTED
func (this HingeConstraint) SetJointFeedback(jointFeedback btJointFeedback*) {
	C.plHingeConstraintSetJointFeedback(this.handle)
}

//UNTESTED
func (this HingeConstraint) GetJointFeedback() btJointFeedback {
	return btJointFeedback(C.plHingeConstraintGetJointFeedback(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetConstraintType() btTypedConstraintType {
	return btTypedConstraintType(C.plHingeConstraintGetConstraintType(this.handle))
}
*/



//UNTESTED
func (this ConeTwistConstraint) GetBreakingImpulseThreshold() float32 {
	return float32(C.plConeTwistConstraintGetBreakingImpulseThreshold(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) SetBreakingImpulseThreshold(threshold float32) {
	C.plConeTwistConstraintSetBreakingImpulseThreshold(this.handle, C.plReal(threshold))
}

//UNTESTED
func (this ConeTwistConstraint) IsEnabled() bool {
	return bool(C.plConeTwistConstraintIsEnabled(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) SetEnabled(enabled bool) {
	C.plConeTwistConstraintSetEnabled(this.handle, C.bool(enabled))
}

//UNTESTED
func (this ConeTwistConstraint) GetRigidBodyA() RigidBody {
	return RigidBody{C.plConeTwistConstraintGetRigidBodyA(this.handle)}
}

//UNTESTED
func (this ConeTwistConstraint) GetRigidBodyB() RigidBody {
	return RigidBody{C.plConeTwistConstraintGetRigidBodyB(this.handle)}
}

//UNTESTED
func (this ConeTwistConstraint) GetUserConstraintType() int {
	return int(C.plConeTwistConstraintGetUserConstraintType(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) SetUserConstraintType(userConstraintType int) {
	C.plConeTwistConstraintSetUserConstraintType(this.handle, C.int(userConstraintType))
}

//UNTESTED
func (this ConeTwistConstraint) SetUserConstraintId(uid int) {
	C.plConeTwistConstraintSetUserConstraintId(this.handle, C.int(uid))
}

//UNTESTED
func (this ConeTwistConstraint) GetUserConstraintId() int {
	return int(C.plConeTwistConstraintGetUserConstraintId(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) SetUserConstraintPtr(ptr *interface{}) {
	C.plConeTwistConstraintSetUserConstraintPtr(this.handle, unsafe.Pointer(ptr))
}

//UNTESTED
func (this ConeTwistConstraint) GetUserConstraintPtr() *interface{}{
	return (*interface{})(C.plConeTwistConstraintGetUserConstraintPtr(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) GetUid() int {
	return int(C.plConeTwistConstraintGetUid(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) NeedsFeedback() bool {
	return bool(C.plConeTwistConstraintNeedsFeedback(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) EnableFeedback(needsFeedback bool) {
	C.plConeTwistConstraintEnableFeedback(this.handle, C.bool(needsFeedback))
}

//UNTESTED
func (this ConeTwistConstraint) GetAppliedImpulse() float32 {
	return float32(C.plConeTwistConstraintGetAppliedImpulse(this.handle))
}

/*
//UNTESTED
func (this ConeTwistConstraint) SetJointFeedback(jointFeedback btJointFeedback*) {
	C.plConeTwistConstraintSetJointFeedback(this.handle)
}

//UNTESTED
func (this ConeTwistConstraint) GetJointFeedback() btJointFeedback {
	return btJointFeedback(C.plConeTwistConstraintGetJointFeedback(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) GetConstraintType() btTypedConstraintType {
	return btTypedConstraintType(C.plConeTwistConstraintGetConstraintType(this.handle))
}
*/



//UNTESTED
func (this Generic6DofConstraint) GetBreakingImpulseThreshold() float32 {
	return float32(C.plGeneric6DofConstraintGetBreakingImpulseThreshold(this.handle))
}

//UNTESTED
func (this Generic6DofConstraint) SetBreakingImpulseThreshold(threshold float32) {
	C.plGeneric6DofConstraintSetBreakingImpulseThreshold(this.handle, C.plReal(threshold))
}

//UNTESTED
func (this Generic6DofConstraint) IsEnabled() bool {
	return bool(C.plGeneric6DofConstraintIsEnabled(this.handle))
}

//UNTESTED
func (this Generic6DofConstraint) SetEnabled(enabled bool) {
	C.plGeneric6DofConstraintSetEnabled(this.handle, C.bool(enabled))
}

//UNTESTED
func (this Generic6DofConstraint) GetRigidBodyA() RigidBody {
	return RigidBody{C.plGeneric6DofConstraintGetRigidBodyA(this.handle)}
}

//UNTESTED
func (this Generic6DofConstraint) GetRigidBodyB() RigidBody {
	return RigidBody{C.plGeneric6DofConstraintGetRigidBodyB(this.handle)}
}

//UNTESTED
func (this Generic6DofConstraint) GetUserConstraintType() int {
	return int(C.plGeneric6DofConstraintGetUserConstraintType(this.handle))
}

//UNTESTED
func (this Generic6DofConstraint) SetUserConstraintType(userConstraintType int) {
	C.plGeneric6DofConstraintSetUserConstraintType(this.handle, C.int(userConstraintType))
}

//UNTESTED
func (this Generic6DofConstraint) SetUserConstraintId(uid int) {
	C.plGeneric6DofConstraintSetUserConstraintId(this.handle, C.int(uid))
}

//UNTESTED
func (this Generic6DofConstraint) GetUserConstraintId() int {
	return int(C.plGeneric6DofConstraintGetUserConstraintId(this.handle))
}

//UNTESTED
func (this Generic6DofConstraint) SetUserConstraintPtr(ptr *interface{}) {
	C.plGeneric6DofConstraintSetUserConstraintPtr(this.handle, unsafe.Pointer(ptr))
}

//UNTESTED
func (this Generic6DofConstraint) GetUserConstraintPtr() *interface{}{
	return (*interface{})(C.plGeneric6DofConstraintGetUserConstraintPtr(this.handle))
}

//UNTESTED
func (this Generic6DofConstraint) GetUid() int {
	return int(C.plGeneric6DofConstraintGetUid(this.handle))
}

//UNTESTED
func (this Generic6DofConstraint) NeedsFeedback() bool {
	return bool(C.plGeneric6DofConstraintNeedsFeedback(this.handle))
}

//UNTESTED
func (this Generic6DofConstraint) EnableFeedback(needsFeedback bool) {
	C.plGeneric6DofConstraintEnableFeedback(this.handle, C.bool(needsFeedback))
}

//UNTESTED
func (this Generic6DofConstraint) GetAppliedImpulse() float32 {
	return float32(C.plGeneric6DofConstraintGetAppliedImpulse(this.handle))
}

/*
//UNTESTED
func (this Generic6DofConstraint) SetJointFeedback(jointFeedback btJointFeedback*) {
	C.plGeneric6DofConstraintSetJointFeedback(this.handle)
}

//UNTESTED
func (this Generic6DofConstraint) GetJointFeedback() btJointFeedback {
	return btJointFeedback(C.plGeneric6DofConstraintGetJointFeedback(this.handle))
}

//UNTESTED
func (this Generic6DofConstraint) GetConstraintType() btTypedConstraintType {
	return btTypedConstraintType(C.plGeneric6DofConstraintGetConstraintType(this.handle))
}
*/



//UNTESTED
func (this SliderConstraint) GetBreakingImpulseThreshold() float32 {
	return float32(C.plSliderConstraintGetBreakingImpulseThreshold(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetBreakingImpulseThreshold(threshold float32) {
	C.plSliderConstraintSetBreakingImpulseThreshold(this.handle, C.plReal(threshold))
}

//UNTESTED
func (this SliderConstraint) IsEnabled() bool {
	return bool(C.plSliderConstraintIsEnabled(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetEnabled(enabled bool) {
	C.plSliderConstraintSetEnabled(this.handle, C.bool(enabled))
}

//UNTESTED
func (this SliderConstraint) GetRigidBodyA() RigidBody {
	return RigidBody{C.plSliderConstraintGetRigidBodyA(this.handle)}
}

//UNTESTED
func (this SliderConstraint) GetRigidBodyB() RigidBody {
	return RigidBody{C.plSliderConstraintGetRigidBodyB(this.handle)}
}

//UNTESTED
func (this SliderConstraint) GetUserConstraintType() int {
	return int(C.plSliderConstraintGetUserConstraintType(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetUserConstraintType(userConstraintType int) {
	C.plSliderConstraintSetUserConstraintType(this.handle, C.int(userConstraintType))
}

//UNTESTED
func (this SliderConstraint) SetUserConstraintId(uid int) {
	C.plSliderConstraintSetUserConstraintId(this.handle, C.int(uid))
}

//UNTESTED
func (this SliderConstraint) GetUserConstraintId() int {
	return int(C.plSliderConstraintGetUserConstraintId(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetUserConstraintPtr(ptr *interface{}) {
	C.plSliderConstraintSetUserConstraintPtr(this.handle, unsafe.Pointer(ptr))
}

//UNTESTED
func (this SliderConstraint) GetUserConstraintPtr() *interface{}{
	return (*interface{})(C.plSliderConstraintGetUserConstraintPtr(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetUid() int {
	return int(C.plSliderConstraintGetUid(this.handle))
}

//UNTESTED
func (this SliderConstraint) NeedsFeedback() bool {
	return bool(C.plSliderConstraintNeedsFeedback(this.handle))
}

//UNTESTED
func (this SliderConstraint) EnableFeedback(needsFeedback bool) {
	C.plSliderConstraintEnableFeedback(this.handle, C.bool(needsFeedback))
}

//UNTESTED
func (this SliderConstraint) GetAppliedImpulse() float32 {
	return float32(C.plSliderConstraintGetAppliedImpulse(this.handle))
}

/*
//UNTESTED
func (this SliderConstraint) SetJointFeedback(jointFeedback btJointFeedback*) {
	C.plSliderConstraintSetJointFeedback(this.handle)
}

//UNTESTED
func (this SliderConstraint) GetJointFeedback() btJointFeedback {
	return btJointFeedback(C.plSliderConstraintGetJointFeedback(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetConstraintType() btTypedConstraintType {
	return btTypedConstraintType(C.plSliderConstraintGetConstraintType(this.handle))
}
*/


