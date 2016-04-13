//for bullet shared lib 2.82 and not double precision
package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

//handle to a physic sdk
type PhysicSDK struct {
	handle C.plPhysicsSdk
}

//Creates a Physics sdk, this is bullet but it could be Physx or ODE
func NewBulletSDK() PhysicSDK {
	return PhysicSDK{C.plNewBulletSdk()}
}

//Free the physics sdk memory, call this when you're done
func (this PhysicSDK) Delete() {
	C.plDeletePhysicsSdk(this.handle)
}
