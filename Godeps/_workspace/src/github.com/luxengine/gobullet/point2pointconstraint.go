package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

type Point2PointConstraint struct {
	handle C.plPoint2PointConstraint
}

//UNTESTED
func NewPoint2PointConstraint(rbA RigidBody, rbB RigidBody, pivotInA *[3]float32, pivotInB *[3]float32) Point2PointConstraint {
	return Point2PointConstraint{C.plNewPoint2PointConstraint(rbA.handle, rbB.handle, (*C.plReal)(&pivotInA[0]), (*C.plReal)(&pivotInB[0]))}
}

//UNTESTED
func NewPoint2PointConstraint2(rbA RigidBody, pivotInA *[3]float32) Point2PointConstraint {
	return Point2PointConstraint{C.plNewPoint2PointConstraint2(rbA.handle, (*C.plReal)(&pivotInA[0]))}
}

//UNTESTED
func (this Point2PointConstraint) Delete() {
	C.plPoint2PointConstraintDelete(this.handle)
}

//UNTESTED
func (this Point2PointConstraint) SetPivotA(pivotA *[3]float32) {
	C.plPoint2PointConstraintSetPivotA(this.handle, (*C.plReal)(&pivotA[0]))
}

//UNTESTED
func (this Point2PointConstraint) SetPivotB(pivotB *[3]float32) {
	C.plPoint2PointConstraintSetPivotB(this.handle, (*C.plReal)(&pivotB[0]))
}

//UNTESTED
func (this Point2PointConstraint) GetPivotInA(dest *[3]float32) {
	C.plPoint2PointConstraintGetPivotInA(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this Point2PointConstraint) GetPivotInB(dest *[3]float32) {
	C.plPoint2PointConstraintGetPivotInB(this.handle, (*C.plReal)(&dest[0]))
}
