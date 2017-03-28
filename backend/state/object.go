package state

import (
	"superstellar/backend/pb"
	"superstellar/backend/types"
)

type Object interface {
	Position() *types.Point
	Velocity() *types.Vector
	Facing() float64
	AngularVelocity() float64
	AngularVelocityDelta() float64

	SetPosition(*types.Point)
	SetVelocity(*types.Vector)
	SetFacing(float64)
	SetAngularVelocity(float64)
	SetAngularVelocityDelta(float64)

	Dirty() bool
	MarkDirty()
	MarkClean()

	DetectCollision(other Object) bool

	NotifyAboutNewFrame()
	AddToProtoSpace(*pb.Space)

	CollideWith(other Object)
	DamageValue() uint32
}
