package state

import (
	"superstellar/backend/constants"
	"superstellar/backend/types"

	"time"
)

type ObjectState struct {
	id                   uint32
	position             *types.Point
	velocity             *types.Vector
	facing               float64
	angularVelocity      float64
	angularVelocityDelta float64
	hp                   uint32

	dirty              bool
	dirtyFramesTimeout uint32
	spawnTimestamp     time.Time
}

func NewObjectState(ID uint32, position *types.Point, velocity *types.Vector, initialHp uint32) *ObjectState {
	return &ObjectState{
		id:                   ID,
		position:             position,
		velocity:             velocity,
		facing:               0.0,
		angularVelocity:      0,
		angularVelocityDelta: 0,
		hp:                   initialHp,

		dirty:              true,
		dirtyFramesTimeout: 0,
		spawnTimestamp:     time.Now(),
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

func (objectState *ObjectState) AngularVelocityDelta() float64 {
	return objectState.angularVelocityDelta
}

func (objectState *ObjectState) SetAngularVelocityDelta(angularVelocityDelta float64) {
	objectState.angularVelocityDelta = angularVelocityDelta
}

func (objectState *ObjectState) Hp() uint32 {
	return objectState.hp
}

func (objectState *ObjectState) SetHp(hp uint32) {
	objectState.hp = hp
}

func (objectState *ObjectState) Dirty() bool {
	return objectState.dirty
}

func (objectState *ObjectState) SpawnTimestamp() time.Time {
	return objectState.spawnTimestamp
}

func (objectState *ObjectState) MarkDirty() {
	objectState.dirty = true
	objectState.dirtyFramesTimeout = constants.DirtyFramesTimeout
}

func (objectState *ObjectState) MarkClean() {
	objectState.dirty = false
}

func (objectState *ObjectState) makeDamage(damage uint32) {
	if objectState.hp < damage {
		objectState.hp = 0
	} else {
		objectState.hp -= damage
	}

	objectState.MarkDirty()
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
