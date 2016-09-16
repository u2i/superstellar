package backend

import (
	"fmt"
	"superstellar/backend/proto"
)

// Direction is a type describing user input on spaceship rotation.
type Direction int

// Constants describing user input on spaceship rotation.
const (
	LEFT Direction = iota
	RIGHT
	NONE
)

const (
	// Acceleration is spaceship's linear acceleration on thruster.
	Acceleration = 0.01

	// AngularVelocity is an angular velocity added on user input.
	AngularVelocity = 0.1
)

// Spaceship struct describes a spaceship.
type Spaceship struct {
	ID             uint32
	Position       *IntVector `json:"position"`
	Velocity       *Vector    `json:"velocity"`
	Facing         *Vector    `json:"facing"`
	InputThrust    bool       `json:"thrust"`
	InputDirection Direction  `json:"direction"`
}

// NewSpaceship initializes new spaceship facing north with zero velocity.
func NewSpaceship(id uint32, position *IntVector) *Spaceship {
	return &Spaceship{
		ID:             id,
		Position:       position,
		Velocity:       ZeroVector(),
		Facing:         NewVector(0.0, 1.0),
		InputThrust:    false,
		InputDirection: NONE,
	}
}

// String function returns string representation.
func (s *Spaceship) String() string {
	return fmt.Sprintf("(%v, %v, %v)", s.Position, s.Velocity, s.Facing)
}

func (s *Spaceship) getNormalizedFacing() *Vector {
	return s.Facing.Normalize()
}

func (s *Spaceship) toProto() *proto.Spaceship {
	return &proto.Spaceship{
		Id:          s.ID,
		Position:    s.Position.toProto(),
		Velocity:    s.Velocity.toProto(),
		Facing:      float32(s.Facing.Radians()),
		InputThrust: s.InputThrust,
	}
}
