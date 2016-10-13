package backend

import (
	"fmt"
	"superstellar/backend/pb"
	"time"
)

// Direction is a type describing user input on spaceship rotation.
type Direction int

// Constants describing user input on spaceship rotation.
const (
	NONE Direction = iota
	RIGHT
	LEFT
)

const (
	// Acceleration is spaceship's linear acceleration on thruster.
	Acceleration = 20.0

	// AngularVelocity is an angular velocity added on user input.
	AngularVelocity = 0.1

	// MaxSpeed maximum speed of the spacecraft
	MaxSpeed = 2000

	// MinFireInterval is a minimum time between firing.
	MinFireInterval = 500 * time.Millisecond
)

// Spaceship struct describes a spaceship.
type Spaceship struct {
	ID             uint32
	Position       *Point
	Velocity       *Vector
	Facing         *Vector
	InputThrust    bool
	InputDirection Direction
	Fire           bool
	LastShotTime   time.Time
}

// NewSpaceship initializes new spaceship facing north with zero velocity.
func NewSpaceship(id uint32, position *Point) *Spaceship {
	return &Spaceship{
		ID:             id,
		Position:       position,
		Velocity:       ZeroVector(),
		Facing:         NewVector(0.0, 1.0),
		InputThrust:    false,
		InputDirection: NONE,
		Fire:           false,
		LastShotTime:   time.Now(),
	}
}

// String function returns string representation.
func (s *Spaceship) String() string {
	return fmt.Sprintf("(%v, %v, %v)", s.Position, s.Velocity, s.Facing)
}

func (s *Spaceship) getNormalizedFacing() *Vector {
	return s.Facing.Normalize()
}

func (s *Spaceship) updateUserInput(userInput *UserInput) {
	switch userInput.UserInput {
	case pb.UserInput_CENTER:
		s.InputDirection = NONE
	case pb.UserInput_LEFT:
		s.InputDirection = LEFT
	case pb.UserInput_RIGHT:
		s.InputDirection = RIGHT
	case pb.UserInput_THRUST_ON:
		s.InputThrust = true
	case pb.UserInput_THRUST_OFF:
		s.InputThrust = false
	case pb.UserInput_FIRE_START:
		s.Fire = true
	case pb.UserInput_FIRE_STOP:
		s.Fire = false
	}
}

func (s *Spaceship) toProto() *pb.Spaceship {
	return &pb.Spaceship{
		Id:          s.ID,
		Position:    s.Position.toProto(),
		Velocity:    s.Velocity.toProto(),
		Facing:      float32(s.Facing.Radians()),
		InputThrust: s.InputThrust,
	}
}
