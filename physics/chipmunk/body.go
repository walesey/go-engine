package chipmunkPhysics

import (
	"github.com/vova616/chipmunk"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ChipmunkBody struct {
	body *chipmunk.Body
}

func NewChipmunkBody() *ChipmunkBody {
	return ChipmunkBody{}
}

func (cBody *ChipmunkBody) ApplyForce(force, position vmath.Vector3) {

}

func (cBody *ChipmunkBody) ApplyTorque(torque vmath.Vector3) {

}

func (cBody *ChipmunkBody) GetPosition() vmath.Vector2 {

}

func (cBody *ChipmunkBody) GetVelocity() vmath.Vector2 {

}

func (cBody *ChipmunkBody) GetRotation() float64 {

}

func (cBody *ChipmunkBody) GetAngularVelocity() float64 {

}

func (cBody *ChipmunkBody) GetMass() float64 {

}

func (cBody *ChipmunkBody) SetPosition(position vmath.Vector2) {

}

func (cBody *ChipmunkBody) SetVelocity(velocity vmath.Vector2) {

}

func (cBody *ChipmunkBody) SetRotation(orientation vmath.Quaternion) {

}

func (cBody *ChipmunkBody) SetAngularVelocity(av float64) {

}
