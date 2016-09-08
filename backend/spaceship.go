package backend

import "fmt"

type Spaceship struct {
	Position *Vector	`json:"position"`
	Velocity *Vector	`json:"veloctiy"`
	Facing   *Vector	`json:"facing"`
}

func (self *Spaceship) String() string {
	return fmt.Sprintf("(%v, %v, %v)", self.Position, self.Velocity, self.Facing)
}

func NewSpaceship(position *Vector) *Spaceship {
	return &Spaceship{
		Position: position,
		Velocity: ZeroVector(),
		Facing:   NewVector(1.0, 0.0)}
}
