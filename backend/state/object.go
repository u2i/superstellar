package state

import (
	"superstellar/backend/pb"
	"superstellar/backend/types"
)

type Object interface {
	Id() uint32

	Position() *types.Point
	Velocity() *types.Vector
	Facing() float64
	AngularVelocity() float64
	AngularVelocityDelta() float64

	Hp() uint32

	SetPosition(*types.Point)
	SetVelocity(*types.Vector)
	SetFacing(float64)
	SetAngularVelocity(float64)
	SetAngularVelocityDelta(float64)

	Dirty() bool
	MarkDirty()
	MarkClean()

	DetectCollision(other Object) bool
	CollideWithProjectile(*Projectile)
	CollideWith(other Object)
	ObjectDestroyed(other Object)
	DamageValue() uint32

	NotifyAboutNewFrame()
	AddToProtoSpace(*pb.Space)
}
