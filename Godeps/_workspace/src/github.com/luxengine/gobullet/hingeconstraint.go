package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

type HingeConstraint struct {
	handle C.plHingeConstraint
}

//UNTESTED
func NewHingeConstraint1(rbA, rbB RigidBody, pivotInA, pivotInB, axisInA, axisInB *[3]float32, useReferenceFrameA bool) HingeConstraint {
	return HingeConstraint{C.plNewHingeConstraint1(rbA.handle, rbB.handle, (*C.plReal)(&pivotInA[0]), (*C.plReal)(&pivotInB[0]), (*C.plReal)(&axisInA[0]), (*C.plReal)(&axisInB[0]), C.bool(useReferenceFrameA))}
}

//UNTESTED
func NewHingeConstraint2(rbA RigidBody, pivotInA, axisInA *[3]float32, useReferenceFrameA bool) HingeConstraint {
	return HingeConstraint{C.plNewHingeConstraint2(rbA.handle, (*C.plReal)(&pivotInA[0]), (*C.plReal)(&axisInA[0]), C.bool(useReferenceFrameA))}
}

//UNTESTED
func NewHingeConstraint3(rbA, rbB RigidBody, rbAFrame, rbBFrame Transform, useReferenceFrameA bool) HingeConstraint {
	return HingeConstraint{C.plNewHingeConstraint3(rbA.handle, rbB.handle, rbAFrame.handle, rbBFrame.handle, C.bool(useReferenceFrameA))}
}

//UNTESTED
func NewHingeConstraint4(rbA RigidBody, rbAFrame Transform, useReferenceFrameA bool) HingeConstraint {
	return HingeConstraint{C.plNewHingeConstraint4(rbA.handle, rbAFrame.handle, C.bool(useReferenceFrameA))}
}

//UNTESTED
func (this HingeConstraint) Delete() {
	C.plHingeConstraintDelete(this.handle)
}

//UNTESTED
func (this HingeConstraint) GetFrameOffsetA() Transform {
	return Transform{C.plHingeConstraintGetFrameOffsetA(this.handle)}
}

//UNTESTED
func (this HingeConstraint) GetFrameOffsetB() Transform {
	return Transform{C.plHingeConstraintGetFrameOffsetB(this.handle)}
}

//UNTESTED
func (this HingeConstraint) SetFrames(frameA, frameB Transform) {
	C.plHingeConstraintSetFrames(this.handle, frameA.handle, frameB.handle)
}

//UNTESTED
func (this HingeConstraint) SetAngularOnly(angularOnly bool) {
	C.plHingeConstraintSetAngularOnly(this.handle, C.bool(angularOnly))
}

//UNTESTED
func (this HingeConstraint) EnableAngularMotor(enableMotor bool, targetVelocity, maxMotorImpulse float32) {
	C.plHingeConstraintEnableAngularMotor(this.handle, C.bool(enableMotor), C.plReal(targetVelocity), C.plReal(maxMotorImpulse))
}

//UNTESTED
func (this HingeConstraint) EnableMotor(enableMotor bool) {
	C.plHingeConstraintEnableMotor(this.handle, C.bool(enableMotor))
}

//UNTESTED
func (this HingeConstraint) SetMaxMotorImpulse(maxMotorImpulse float32) {
	C.plHingeConstraintSetMaxMotorImpulse(this.handle, C.plReal(maxMotorImpulse))
}

//UNTESTED
func (this HingeConstraint) SetMotorTargetQuat(qAinB *[4]float32, dt float32) {
	C.plHingeConstraintSetMotorTargetQuat(this.handle, (*C.plReal)(&qAinB[0]), C.plReal(dt))
}

//UNTESTED
func (this HingeConstraint) SetMotorTargetAng(targetAngle, dt float32) {
	C.plHingeConstraintSetMotorTargetAng(this.handle, C.plReal(targetAngle), C.plReal(dt))
}

//UNTESTED
func (this HingeConstraint) SetLimit(low, high, softness, biasFactor, relaxationFactor float32) {
	C.plHingeConstraintSetLimit(this.handle, C.plReal(low), C.plReal(high), C.plReal(softness), C.plReal(biasFactor), C.plReal(relaxationFactor))
}

//UNTESTED
func (this HingeConstraint) SetLimitDefault(low, high float32) {
	C.plHingeConstraintSetLimit(this.handle, C.plReal(low), C.plReal(high), C.plReal(0.9), C.plReal(0.3), C.plReal(1.0))
}

//UNTESTED
func (this HingeConstraint) SetAxis(axisInA *[3]float32) {
	C.plHingeConstraintSetAxis(this.handle, (*C.plReal)(&axisInA[0]))
}

//UNTESTED
func (this HingeConstraint) GetLowerLimit() float32 {
	return float32(C.plHingeConstraintGetLowerLimit(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetUpperLimit() float32 {
	return float32(C.plHingeConstraintGetUpperLimit(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetHingeAngle() float32 {
	return float32(C.plHingeConstraintGetHingeAngle(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetHingeAngle2(transA, transB Transform) float32 {
	return float32(C.plHingeConstraintGetHingeAngle2(this.handle, transA.handle, transB.handle))
}

//UNTESTED
func (this HingeConstraint) GetAFrame() Transform {
	return Transform{C.plHingeConstraintGetAFrame(this.handle)}
}

//UNTESTED
func (this HingeConstraint) GetBFrame() Transform {
	return Transform{C.plHingeConstraintGetBFrame(this.handle)}
}

//UNTESTED
func (this HingeConstraint) GetSolveLimit() int {
	return int(C.plHingeConstraintGetSolveLimit(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetLimitSign() float32 {
	return float32(C.plHingeConstraintGetLimitSign(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetAngularOnly() bool {
	return bool(C.plHingeConstraintGetAngularOnly(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetEnableAngularMotor() bool {
	return bool(C.plHingeConstraintGetEnableAngularMotor(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetMotorTargetVelosity() float32 {
	return float32(C.plHingeConstraintGetMotorTargetVelosity(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetMaxMotorImpulse() float32 {
	return float32(C.plHingeConstraintGetMaxMotorImpulse(this.handle))
}

//UNTESTED
func (this HingeConstraint) GetUseFrameOffset() bool {
	return bool(C.plHingeConstraintGetUseFrameOffset(this.handle))
}

//UNTESTED
func (this HingeConstraint) SetUseFrameOffset(frameOffsetOnOff bool) {
	C.plHingeConstraintSetUseFrameOffset(this.handle, C.bool(frameOffsetOnOff))
}
