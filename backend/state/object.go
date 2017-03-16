package state

import "superstellar/backend/types"

type Object interface {
	Position() *types.Point
	Velocity() *types.Vector

	SetPosition(*types.Point)
	SetVelocity(*types.Vector)

	Dirty() bool
	MarkDirty()
	MarkClean()

	DetectCollision(other Object) bool
	Collide(other Object)
}
