package backend

import "fmt"

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
	Acceleration = 0.1

	// AngularVelocity is an angular velocity added on user input.
	AngularVelocity = 0.23
)

// Spaceship struct describes a spaceship.
type Spaceship struct {
	Position       *Vector   `json:"position"`
	Velocity       *Vector   `json:"veloctiy"`
	Facing         *Vector   `json:"facing"`
	InputThrust    bool      `json:"thrust"`
	InputDirection Direction `json:"direction"`
}

// NewSpaceship initializes new spaceship facing north with zero velocity.
func NewSpaceship(position *Vector) *Spaceship {
	return &Spaceship{
		Position:       position,
		Velocity:       ZeroVector(),
		Facing:         NewVector(1.0, 0.0),
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
