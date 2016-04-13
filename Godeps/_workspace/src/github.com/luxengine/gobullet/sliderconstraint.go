package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

type SliderConstraint struct {
	handle C.plSliderConstraint
}

//UNTESTED
func NewSliderConstraint(b1, b2 RigidBody, frameInA, frameInB Transform, useLinearReferenceFrameA bool) SliderConstraint {
	return SliderConstraint{C.plNewSliderConstraint(b1.handle, b2.handle, frameInA.handle, frameInB.handle, C.bool(useLinearReferenceFrameA))}
}

func NewSliderConstraintWorld(rb RigidBody, tran Transform, useLinearReferenceFrameA bool) SliderConstraint {
	return SliderConstraint{C.plNewSliderConstraintWorld(rb.handle, tran.handle, C.bool(useLinearReferenceFrameA))}
}

//UNTESTED
func (this SliderConstraint) Delete() {
	C.plSliderConstraintDelete(this.handle)
}

//UNTESTED
func (this SliderConstraint) SetLowerLinLimit(lowerLimit float32) {
	C.plSliderConstraintSetLowerLinLimit(this.handle, C.plReal(lowerLimit))
}

//UNTESTED
func (this SliderConstraint) SetUpperLinLimit(upperLimit float32) {
	C.plSliderConstraintSetUpperLinLimit(this.handle, C.plReal(upperLimit))
}

//UNTESTED
func (this SliderConstraint) GetLowerLinLimit() float32 {
	return float32(C.plSliderConstraintGetLowerLinLimit(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetUpperLinLimit() float32 {
	return float32(C.plSliderConstraintGetUpperLinLimit(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetLowerAngLimit() float32 {
	return float32(C.plSliderConstraintGetLowerAngLimit(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetLowerAngLimit(limit float32) {
	C.plSliderConstraintSetLowerAngLimit(this.handle, C.plReal(limit))
}

//UNTESTED
func (this SliderConstraint) GetUpperAngLimit() float32 {
	return float32(C.plSliderConstraintGetUpperAngLimit(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetUpperAngLimit(limit float32) {
	C.plSliderConstraintSetUpperAngLimit(this.handle, C.plReal(limit))
}

//UNTESTED
func (this SliderConstraint) GetUseLinearReferenceFrameA() bool {
	return bool(C.plSliderConstraintGetUseLinearReferenceFrameA(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetSoftnessDirLin() float32 {
	return float32(C.plSliderConstraintGetSoftnessDirLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetRestitutionDirLin() float32 {
	return float32(C.plSliderConstraintGetRestitutionDirLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetDampingDirLin() float32 {
	return float32(C.plSliderConstraintGetDampingDirLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetSoftnessDirAng() float32 {
	return float32(C.plSliderConstraintGetSoftnessDirAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetRestitutionDirAng() float32 {
	return float32(C.plSliderConstraintGetRestitutionDirAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetDampingDirAng() float32 {
	return float32(C.plSliderConstraintGetDampingDirAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetSoftnessLimLin() float32 {
	return float32(C.plSliderConstraintGetSoftnessLimLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetRestitutionLimLin() float32 {
	return float32(C.plSliderConstraintGetRestitutionLimLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetDampingLimLin() float32 {
	return float32(C.plSliderConstraintGetDampingLimLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetSoftnessLimAng() float32 {
	return float32(C.plSliderConstraintGetSoftnessLimAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetRestitutionLimAng() float32 {
	return float32(C.plSliderConstraintGetRestitutionLimAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetDampingLimAng() float32 {
	return float32(C.plSliderConstraintGetDampingLimAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetSoftnessOrthoLin() float32 {
	return float32(C.plSliderConstraintGetSoftnessOrthoLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetRestitutionOrthoLin() float32 {
	return float32(C.plSliderConstraintGetRestitutionOrthoLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetDampingOrthoLin() float32 {
	return float32(C.plSliderConstraintGetDampingOrthoLin(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetSoftnessOrthoAng() float32 {
	return float32(C.plSliderConstraintGetSoftnessOrthoAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetRestitutionOrthoAng() float32 {
	return float32(C.plSliderConstraintGetRestitutionOrthoAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetDampingOrthoAng() float32 {
	return float32(C.plSliderConstraintGetDampingOrthoAng(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetSoftnessDirLin(softness float32) {
	C.plSliderConstraintSetSoftnessDirLin(this.handle, C.plReal(softness))
}

//UNTESTED
func (this SliderConstraint) SetRestitutionDirLin(restitution float32) {
	C.plSliderConstraintSetRestitutionDirLin(this.handle, C.plReal(restitution))
}

//UNTESTED
func (this SliderConstraint) SetDampingDirLin(damping float32) {
	C.plSliderConstraintSetDampingDirLin(this.handle, C.plReal(damping))
}

//UNTESTED
func (this SliderConstraint) SetSoftnessDirAng(softness float32) {
	C.plSliderConstraintSetSoftnessDirAng(this.handle, C.plReal(softness))
}

//UNTESTED
func (this SliderConstraint) SetRestitutionDirAng(restitution float32) {
	C.plSliderConstraintSetRestitutionDirAng(this.handle, C.plReal(restitution))
}

//UNTESTED
func (this SliderConstraint) SetDampingDirAng(damping float32) {
	C.plSliderConstraintSetDampingDirAng(this.handle, C.plReal(damping))
}

//UNTESTED
func (this SliderConstraint) SetSoftnessLimLin(softness float32) {
	C.plSliderConstraintSetSoftnessLimLin(this.handle, C.plReal(softness))
}

//UNTESTED
func (this SliderConstraint) SetRestitutionLimLin(restitution float32) {
	C.plSliderConstraintSetRestitutionLimLin(this.handle, C.plReal(restitution))
}

//UNTESTED
func (this SliderConstraint) SetDampingLimLin(damping float32) {
	C.plSliderConstraintSetDampingLimLin(this.handle, C.plReal(damping))
}

//UNTESTED
func (this SliderConstraint) SetSoftnessLimAng(softness float32) {
	C.plSliderConstraintSetSoftnessLimAng(this.handle, C.plReal(softness))
}

//UNTESTED
func (this SliderConstraint) SetRestitutionLimAng(restitution float32) {
	C.plSliderConstraintSetRestitutionLimAng(this.handle, C.plReal(restitution))
}

//UNTESTED
func (this SliderConstraint) SetDampingLimAng(damping float32) {
	C.plSliderConstraintSetDampingLimAng(this.handle, C.plReal(damping))
}

//UNTESTED
func (this SliderConstraint) SetSoftnessOrthoLin(softness float32) {
	C.plSliderConstraintSetSoftnessOrthoLin(this.handle, C.plReal(softness))
}

//UNTESTED
func (this SliderConstraint) SetRestitutionOrthoLin(restitution float32) {
	C.plSliderConstraintSetRestitutionOrthoLin(this.handle, C.plReal(restitution))
}

//UNTESTED
func (this SliderConstraint) SetDampingOrthoLin(damping float32) {
	C.plSliderConstraintSetDampingOrthoLin(this.handle, C.plReal(damping))
}

//UNTESTED
func (this SliderConstraint) SetSoftnessOrthoAng(softness float32) {
	C.plSliderConstraintSetSoftnessOrthoAng(this.handle, C.plReal(softness))
}

//UNTESTED
func (this SliderConstraint) SetRestitutionOrthoAng(restitution float32) {
	C.plSliderConstraintSetRestitutionOrthoAng(this.handle, C.plReal(restitution))
}

//UNTESTED
func (this SliderConstraint) SetDampingOrthoAng(damping float32) {
	C.plSliderConstraintSetDampingOrthoAng(this.handle, C.plReal(damping))
}

//UNTESTED
func (this SliderConstraint) SetPoweredLinMotor(onoff bool) {
	C.plSliderConstraintSetPoweredLinMotor(this.handle, C.bool(onoff))
}

//UNTESTED
func (this SliderConstraint) GetPoweredLinMotor() bool {
	return bool(C.plSliderConstraintGetPoweredLinMotor(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetTargetLinMotorVelocity(linvel float32) {
	C.plSliderConstraintSetTargetLinMotorVelocity(this.handle, C.plReal(linvel))
}

//UNTESTED
func (this SliderConstraint) GetTargetLinMotorVelocity() float32 {
	return float32(C.plSliderConstraintGetTargetLinMotorVelocity(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetMaxLinMotorForce(maxforce float32) {
	C.plSliderConstraintSetMaxLinMotorForce(this.handle, C.plReal(maxforce))
}

//UNTESTED
func (this SliderConstraint) GetMaxLinMotorForce() float32 {
	return float32(C.plSliderConstraintGetMaxLinMotorForce(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetPoweredAngMotor(onoff bool) {
	C.plSliderConstraintSetPoweredAngMotor(this.handle, C.bool(onoff))
}

//UNTESTED
func (this SliderConstraint) GetPoweredAngMotor() bool {
	return bool(C.plSliderConstraintGetPoweredAngMotor(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetTargetAngMotorVelocity(angvel float32) {
	C.plSliderConstraintSetTargetAngMotorVelocity(this.handle, C.plReal(angvel))
}

//UNTESTED
func (this SliderConstraint) GetTargetAngMotorVelocity() float32 {
	return float32(C.plSliderConstraintGetTargetAngMotorVelocity(this.handle))
}

//UNTESTED
func (this SliderConstraint) SetMaxAngMotorForce(maxforce float32) {
	C.plSliderConstraintSetMaxAngMotorForce(this.handle, C.plReal(maxforce))
}

//UNTESTED
func (this SliderConstraint) GetMaxAngMotorForce() float32 {
	return float32(C.plSliderConstraintGetMaxAngMotorForce(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetLinearPos() float32 {
	return float32(C.plSliderConstraintGetLinearPos(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetAngularPos() float32 {
	return float32(C.plSliderConstraintGetAngularPos(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetSolveLinLimit() bool {
	return bool(C.plSliderConstraintGetSolveLinLimit(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetLinDepth() float32 {
	return float32(C.plSliderConstraintGetLinDepth(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetSolveAngLimit() bool {
	return bool(C.plSliderConstraintGetSolveAngLimit(this.handle))
}

//UNTESTED
func (this SliderConstraint) GetAngDepth() float32 {
	return float32(C.plSliderConstraintGetAngDepth(this.handle))
}
