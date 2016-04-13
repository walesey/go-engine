package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

/********************************
***********Constraints***********
********************************/
type Generic6DofConstraint struct {
	handle C.plGeneric6DofConstraint
}

//UNTESTED
func NewGeneric6DofConstraint(b1, b2 RigidBody, t1, t2 Transform, useLinearReferenceFrameA bool) Generic6DofConstraint {
	return Generic6DofConstraint{C.plNewGeneric6DofConstraint(b1.handle, b2.handle, t1.handle, t2.handle, C.bool(useLinearReferenceFrameA))}
}

//UNTESTED
func (this Generic6DofConstraint) Delete() {
	C.plGeneric6DofDelete(this.handle)
}

//UNTESTED
func (this Generic6DofConstraint) SetLinearLowerLimit(limit *[3]float32) {
	C.plGeneric6DofConstraintSetLinearLowerLimit(this.handle, (*C.plReal)(&limit[0]))
}

//UNTESTED
func (this Generic6DofConstraint) GetLinearLowerLimit(dest *[3]float32) {
	C.plGeneric6DofConstraintGetLinearLowerLimit(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this Generic6DofConstraint) SetLinearUpperLimit(limit *[3]float32) {
	C.plGeneric6DofConstraintSetLinearUpperLimit(this.handle, (*C.plReal)(&limit[0]))
}

//UNTESTED
func (this Generic6DofConstraint) GetLinearUpperLimit(dest *[3]float32) {
	C.plGeneric6DofConstraintGetLinearUpperLimit(this.handle, (*C.plReal)(&dest[0]))
}
