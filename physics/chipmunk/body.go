package chipmunkPhysics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

type ChipmunkBody struct {
	Body *chipmunk.Body
}

func NewChipmunkBody(mass, i float32) *ChipmunkBody {
	return &ChipmunkBody{
		Body: chipmunk.NewBody(vect.Float(mass), vect.Float(i)),
	}
}

func NewChipmunkBodyStatic() *ChipmunkBody {
	return &ChipmunkBody{
		Body: chipmunk.NewBodyStatic(),
	}
}

func (cBody *ChipmunkBody) KineticEnergy() float32 {
	return float32(cBody.Body.KineticEnergy())
}

func (cBody *ChipmunkBody) SetMass(mass float32) {
	cBody.Body.SetMass(vect.Float(mass))
}

func (cBody *ChipmunkBody) SetMoment(moment float32) {
	cBody.Body.SetMoment(vect.Float(moment))
}

func (cBody *ChipmunkBody) GetMoment() float32 {
	return cBody.Body.Moment()
}

func (cBody *ChipmunkBody) SetAngle(angle float32) {
	cBody.Body.SetAngle(vect.Float(angle))
}

func (cBody *ChipmunkBody) AddAngle(angle float32) {
	cBody.Body.AddAngle(angle)
}

func (cBody *ChipmunkBody) GetMass() float32 {
	return float32(cBody.Body.Mass())
}

func (cBody *ChipmunkBody) SetPosition(pos mgl32.Vec2) {
	cBody.Body.SetPosition(convertToVect(pos))
}

func (cBody *ChipmunkBody) AddForce(force mgl32.Vec2) {
	cBody.Body.AddForce(force.X(), force.Y())
}

func (cBody *ChipmunkBody) SetForce(force mgl32.Vec2) {
	cBody.Body.SetForce(force.X(), force.Y())
}

func (cBody *ChipmunkBody) AddVelocity(velocity mgl32.Vec2) {
	cBody.Body.AddVelocity(velocity.X(), velocity.Y())
}

func (cBody *ChipmunkBody) SetVelocity(velocity mgl32.Vec2) {
	cBody.Body.SetVelocity(velocity.X(), velocity.Y())
}

func (cBody *ChipmunkBody) AddTorque(t float32) {
	cBody.Body.AddTorque(float32(t))
}

func (cBody *ChipmunkBody) GetTorque() float32 {
	return cBody.Body.Torque()
}

func (cBody *ChipmunkBody) GetAngularVelocity() float32 {
	return cBody.Body.AngularVelocity()
}

func (cBody *ChipmunkBody) SetTorque(t float32) {
	cBody.Body.SetTorque(float32(t))
}

func (cBody *ChipmunkBody) AddAngularVelocity(w float32) {
	cBody.Body.AddAngularVelocity(w)
}

func (cBody *ChipmunkBody) SetAngularVelocity(w float32) {
	cBody.Body.SetAngularVelocity(w)
}

func (cBody *ChipmunkBody) GetVelocity() mgl32.Vec2 {
	return convertFromVect(cBody.Body.Velocity())
}

func (cBody *ChipmunkBody) GetPosition() mgl32.Vec2 {
	return convertFromVect(cBody.Body.Position())
}

func (cBody *ChipmunkBody) GetAngle() float32 {
	return float32(cBody.Body.Angle())
}
