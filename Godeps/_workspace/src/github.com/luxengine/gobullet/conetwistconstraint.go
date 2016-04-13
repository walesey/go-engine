package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

type ConeTwistConstraint struct {
	handle C.plConeTwistConstraint
}

//UNTESTED
func NewConeTwistConstraint(rbA RigidBody, rbB RigidBody, rbAFrame Transform, rbBFrame Transform) ConeTwistConstraint {
	return ConeTwistConstraint{C.plNewConeTwistConstraint(rbA.handle, rbB.handle, rbAFrame.handle, rbBFrame.handle)}
}

//UNTESTED
func NewConeTwistConstraint2(rbA RigidBody, rbAFrame Transform) ConeTwistConstraint {
	return ConeTwistConstraint{C.plNewConeTwistConstraint2(rbA.handle, rbAFrame.handle)}
}

//UNTESTED
func (this ConeTwistConstraint) Delete() {
	C.plConeTwistConstraintDelete(this.handle)
}

//UNTESTED
func (this ConeTwistConstraint) SetAngularOnly(angularOnly bool) {
	C.plConeTwistConstraintSetAngularOnly(this.handle, C.bool(angularOnly))
}

//UNTESTED
func (this ConeTwistConstraint) SetLimiti(limitIndex int, limitValue float32) {
	C.plConeTwistConstraintSetLimiti(this.handle, C.int(limitIndex), C.plReal(limitValue))
}

//UNTESTED
//default: =1.f =0.3f =1.0f
func (this ConeTwistConstraint) SetLimits(swingSpan1, swingSpan2, twistSpan, softness, biasFactor, relaxationFactor float32) {
	C.plConeTwistConstraintSetLimits(this.handle, C.plReal(swingSpan1), C.plReal(swingSpan2), C.plReal(twistSpan), C.plReal(softness), C.plReal(biasFactor), C.plReal(relaxationFactor))
}

//UNTESTED
func (this ConeTwistConstraint) GetAFrame() Transform {
	return Transform{C.plConeTwistConstraintGetAFrame(this.handle)}
}

//UNTESTED
func (this ConeTwistConstraint) GetBFrame() Transform {
	return Transform{C.plConeTwistConstraintGetBFrame(this.handle)}
}

//UNTESTED
func (this ConeTwistConstraint) GetSolveTwistLimit() int {
	return int(C.plConeTwistConstraintGetSolveTwistLimit(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) GetSolveSwingLimit() int {
	return int(C.plConeTwistConstraintGetSolveSwingLimit(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) GetTwistLimitSign() float32 {
	return float32(C.plConeTwistConstraintGetTwistLimitSign(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) GetSwingSpan1() float32 {
	return float32(C.plConeTwistConstraintGetSwingSpan1(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) GetSwingSpan2() float32 {
	return float32(C.plConeTwistConstraintGetSwingSpan2(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) GetTwistSpan() float32 {
	return float32(C.plConeTwistConstraintGetTwistSpan(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) GetTwistAngle() float32 {
	return float32(C.plConeTwistConstraintGetTwistAngle(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) IsPastSwingLimit() bool {
	return bool(C.plConeTwistConstraintIsPastSwingLimit(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) SetDamping(damping float32) {
	C.plConeTwistConstraintSetDamping(this.handle, C.plReal(damping))
}

//UNTESTED
func (this ConeTwistConstraint) EnableMotor(b bool) {
	C.plConeTwistConstraintEnableMotor(this.handle, C.bool(b))
}

//UNTESTED
func (this ConeTwistConstraint) SetMaxMotorImpulse(maxMotorImpulse float32) {
	C.plConeTwistConstraintSetMaxMotorImpulse(this.handle, C.plReal(maxMotorImpulse))
}

//UNTESTED
func (this ConeTwistConstraint) SetMaxMotorImpulseNormalized(maxMotorImpulse float32) {
	C.plConeTwistConstraintSetMaxMotorImpulseNormalized(this.handle, C.plReal(maxMotorImpulse))
}

//UNTESTED
func (this ConeTwistConstraint) GetFixThresh() float32 {
	return float32(C.plConeTwistConstraintGetFixThresh(this.handle))
}

//UNTESTED
func (this ConeTwistConstraint) SetFixThresh(fixThresh float32) {
	C.plConeTwistConstraintSetFixThresh(this.handle, C.plReal(fixThresh))
}

//UNTESTED
func (this ConeTwistConstraint) SetMotorTarget(q *[4]float32) {
	C.plConeTwistConstraintSetMotorTarget(this.handle, (*C.plReal)(&q[0]))
}

//UNTESTED
func (this ConeTwistConstraint) SetMotorTargetInConstraintSpace(q *[4]float32) {
	C.plConeTwistConstraintSetMotorTargetInConstraintSpace(this.handle, (*C.plReal)(&q[0]))
}

//UNTESTED
func (this ConeTwistConstraint) GetPointForAngle(fAngleInRadians, fLength float32, dest *[3]float32) {
	C.plConeTwistConstraintGetPointForAngle(this.handle, C.plReal(fAngleInRadians), C.plReal(fLength), (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this ConeTwistConstraint) SetFrames(frameA Transform, frameB Transform) {
	C.plConeTwistConstraintSetFrames(this.handle, frameA.handle, frameB.handle)
}

//UNTESTED
func (this ConeTwistConstraint) GetFrameOffsetA() Transform {
	return Transform{C.plConeTwistConstraintGetFrameOffsetA(this.handle)}
}

//UNTESTED
func (this ConeTwistConstraint) GetFrameOffsetB() Transform {
	return Transform{C.plConeTwistConstraintGetFrameOffsetB(this.handle)}
}
