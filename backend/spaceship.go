package backend

import "fmt"

// Spaceship struct describes a spaceship.
type Spaceship struct {
	Position *Vector `json:"position"`
	Velocity *Vector `json:"veloctiy"`
	Facing   *Vector `json:"facing"`
}

// NewSpaceship initializes new spaceship facing north with zero velocity.
func NewSpaceship(position *Vector) *Spaceship {
	return &Spaceship{
		Position: position,
		Velocity: ZeroVector(),
		Facing:   NewVector(1.0, 0.0)}
}

// String function returns string representation.
func (s *Spaceship) String() string {
	return fmt.Sprintf("(%v, %v, %v)", s.Position, s.Velocity, s.Facing)
}
