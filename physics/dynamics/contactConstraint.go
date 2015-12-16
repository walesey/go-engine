package dynamics

import vmath "github.com/walesey/go-engine/vectormath"

type ContactConstraint struct {
	Normal                       vmath.Vector3
	LocalContact1, LocalContact2 vmath.Vector3
	Object1, Object2             *PhysicsObject
	InContact                    bool
}

func (cc *ContactConstraint) Solve() {

	//velocities
	angularV1 := cc.Object1.AngularVelocityVector()
	angularV2 := cc.Object2.AngularVelocityVector()
	radialV1 := cc.LocalContact1.Cross(angularV1)
	radialV2 := cc.LocalContact2.Cross(angularV2)
	contactV1 := radialV1.Add(cc.Object1.Velocity)
	contactV2 := radialV2.Add(cc.Object2.Velocity)

	if cc.Object2.Static {
		cc.Object2.Mass = 99999999999999999.9
	}

	tensorInv1 := cc.Object1.InertiaTensor().Inverse()
	tensorInv2 := cc.Object2.InertiaTensor().Inverse()

	contactV := contactV1.Subtract(contactV2)
	relativeV := cc.Normal.Dot(contactV)

	velocityImpulse := -relativeV
	vel1 := tensorInv1.Transform(cc.LocalContact1.Cross(cc.Normal))
	vel1 = vel1.Cross(cc.LocalContact1)
	impulseDenom1 := (1.0 / cc.Object1.Mass) + cc.Normal.Dot(vel1)
	vel2 := tensorInv2.Transform(cc.LocalContact2.Cross(cc.Normal))
	vel2 = vel2.Cross(cc.LocalContact2)
	impulseDenom2 := (1.0 / cc.Object2.Mass) + cc.Normal.Dot(vel2)
	impulseDenom := impulseDenom1 + impulseDenom2
	NormalImpulse := velocityImpulse / impulseDenom

	if NormalImpulse > 0.0 {
		NormalImpulse = 0.0
	}

	impulseVector1 := cc.Normal.MultiplyScalar(NormalImpulse)
	impulseVector2 := impulseVector1.MultiplyScalar(-1)
	//project impulse onto the localcontact normal
	localContactNorm1 := cc.LocalContact1.Normalize()
	localContactNorm2 := cc.LocalContact2.Normalize()
	linearImpulse1 := localContactNorm1.MultiplyScalar(impulseVector1.Dot(localContactNorm1))
	linearImpulse2 := localContactNorm2.MultiplyScalar(impulseVector2.Dot(localContactNorm2))

	torqueImpulse1 := impulseVector1.Cross(cc.LocalContact1)
	torqueImpulse2 := impulseVector2.Cross(cc.LocalContact2)

	cc.Object1.Velocity = cc.Object1.Velocity.Add(linearImpulse1.DivideScalar(cc.Object1.Mass))
	newAngularV1 := angularV1.Add(tensorInv1.Transform(torqueImpulse1))
	cc.Object1.SetAngularVelocityVector(newAngularV1)

	if !cc.Object2.Static {
		cc.Object2.Velocity = cc.Object2.Velocity.Add(linearImpulse2.DivideScalar(cc.Object2.Mass))
		newAngularV2 := angularV2.Add(tensorInv2.Transform(torqueImpulse2))
		cc.Object2.SetAngularVelocityVector(newAngularV2)
	}

	if cc.InContact && cc.Object2.Static {
		//restrict velocity to the tangent plane between the two objects
		normalVelocity1 := cc.Normal.MultiplyScalar(cc.Object1.Velocity.Dot(cc.Normal))
		cc.Object1.Velocity = cc.Object1.Velocity.Subtract(normalVelocity1)
		cc.Object1.Velocity = cc.Object1.Velocity.MultiplyScalar(0.8)
	}
}
