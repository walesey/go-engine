package chipmunkPhysics

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ChipmunkBody struct {
	Body *chipmunk.Body
}

func NewChipmunkBody(mass, i float64) *ChipmunkBody {
	return &ChipmunkBody{
		Body: chipmunk.NewBody(vect.Float(mass), vect.Float(i)),
	}
}

func NewChipmunkBodyStatic() *ChipmunkBody {
	return &ChipmunkBody{
		Body: chipmunk.NewBodyStatic(),
	}
}

func (cBody *ChipmunkBody) KineticEnergy() float64 {
	return float64(cBody.Body.KineticEnergy())
}

func (cBody *ChipmunkBody) SetMass(mass float64) {
	cBody.Body.SetMass(vect.Float(mass))
}

func (cBody *ChipmunkBody) SetMoment(moment float64) {
	cBody.Body.SetMoment(vect.Float(moment))
}

func (cBody *ChipmunkBody) GetMoment() float64 {
	return float64(cBody.Body.Moment())
}

func (cBody *ChipmunkBody) SetAngle(angle float64) {
	cBody.Body.SetAngle(vect.Float(angle))
}

func (cBody *ChipmunkBody) AddAngle(angle float64) {
	cBody.Body.AddAngle(float32(angle))
}

func (cBody *ChipmunkBody) GetMass() float64 {
	return float64(cBody.Body.Mass())
}

func (cBody *ChipmunkBody) SetPosition(pos vmath.Vector2) {
	cBody.Body.SetPosition(convertToVect(pos))
}

func (cBody *ChipmunkBody) AddForce(force vmath.Vector2) {
	cBody.Body.AddForce(float32(force.X), float32(force.Y))
}

func (cBody *ChipmunkBody) SetForce(force vmath.Vector2) {
	cBody.Body.SetForce(float32(force.X), float32(force.Y))
}

func (cBody *ChipmunkBody) AddVelocity(velocity vmath.Vector2) {
	cBody.Body.AddVelocity(float32(velocity.X), float32(velocity.Y))
}

func (cBody *ChipmunkBody) SetVelocity(velocity vmath.Vector2) {
	cBody.Body.SetVelocity(float32(velocity.X), float32(velocity.Y))
}

func (cBody *ChipmunkBody) AddTorque(t float64) {
	cBody.Body.AddTorque(float32(t))
}

func (cBody *ChipmunkBody) GetTorque() float64 {
	return float64(cBody.Body.Torque())
}

func (cBody *ChipmunkBody) GetAngularVelocity() float64 {
	return float64(cBody.Body.AngularVelocity())
}

func (cBody *ChipmunkBody) SetTorque(t float64) {
	cBody.Body.SetTorque(float32(t))
}

func (cBody *ChipmunkBody) AddAngularVelocity(w float64) {
	cBody.Body.AddAngularVelocity(float32(w))
}

func (cBody *ChipmunkBody) SetAngularVelocity(w float64) {
	cBody.Body.SetAngularVelocity(float32(w))
}

func (cBody *ChipmunkBody) GetVelocity() vmath.Vector2 {
	return convertFromVect(cBody.Body.Velocity())
}

func (cBody *ChipmunkBody) GetPosition() vmath.Vector2 {
	return convertFromVect(cBody.Body.Position())
}

func (cBody *ChipmunkBody) GetAngle() float64 {
	return float64(cBody.Body.Angle())
}
