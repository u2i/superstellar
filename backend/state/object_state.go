package state

import (
	"superstellar/backend/types"
	"superstellar/backend/constants"
	"math"
)

type ObjectState struct {
	id                   uint32
	position             *types.Point
	velocity             *types.Vector
	facing               float64
	angularVelocity      float64
	dirty                bool
	dirtyFramesTimeout   uint32
}

func NewObjectState(ID uint32, position *types.Point, velocity *types.Vector) *ObjectState {
	return &ObjectState{
		id:                   ID,
		position:             position,
		velocity:             velocity,
		facing:               0.0,
		angularVelocity:      0,
		dirty:                true,
		dirtyFramesTimeout:   0,
	}
}

func (objectState *ObjectState) Id() uint32 {
	return objectState.id
}

func (objectState *ObjectState) Position() *types.Point {
	return objectState.position
}

func (objectState *ObjectState) SetPosition(position *types.Point) {
	objectState.position = position
}

func (objectState *ObjectState) Velocity() *types.Vector {
	return objectState.velocity
}

func (objectState *ObjectState) SetVelocity(velocity *types.Vector) {
	objectState.velocity = velocity
}

func (objectState *ObjectState) Facing() float64 {
	return objectState.facing
}

func (objectState *ObjectState) SetFacing(facing float64) {
	objectState.facing = facing
}

func (objectState *ObjectState) AngularVelocity() float64 {
	return objectState.angularVelocity
}

func (objectState *ObjectState) SetAngularVelocity(angularVelocity float64) {
	objectState.angularVelocity = angularVelocity
}

func (objectState *ObjectState) Dirty() bool {
	return objectState.dirty
}

func (objectState *ObjectState) MarkDirty() {
	objectState.dirty = true
	objectState.dirtyFramesTimeout = constants.DirtyFramesTimeout
}

func (objectState *ObjectState) MarkClean() {
	objectState.dirty = false
}

func (objectState *ObjectState) HandleDirtyTimeout() {
	if objectState.dirtyFramesTimeout == 0 {
		objectState.MarkDirty()
	} else {
		objectState.dirtyFramesTimeout--
	}
}

func (objectState *ObjectState) NotifyAboutNewFrame() {
	objectState.HandleDirtyTimeout()
}

// DetectCollision returns true if receiver object collides with other object.
func (object *ObjectState) DetectCollision(other Object) bool {
	v := types.Point{X: object.Position().X - other.Position().X, Y: object.Position().Y - other.Position().Y}
	dist := v.Length()

	return dist < 2*constants.SpaceshipSize
}

// Collide transforms colliding ships' parameters.
func (object *ObjectState) Collide(other Object) {
	v := types.Point{
		X: object.Position().X - other.Position().X,
		Y: object.Position().Y - other.Position().Y,
	}

	transformAngle := -math.Atan2(float64(v.Y), float64(v.X))
	newV1 := object.Velocity().Rotate(transformAngle)
	newV2 := other.Velocity().Rotate(transformAngle)

	switchedV1 := types.Vector{X: newV2.X, Y: newV1.Y}
	switchedV2 := types.Vector{X: newV1.X, Y: newV2.Y}

	object.SetVelocity(switchedV1.Rotate(-transformAngle))
	other.SetVelocity(switchedV2.Rotate(-transformAngle))

	object.MarkDirty()
	other.MarkDirty()
}