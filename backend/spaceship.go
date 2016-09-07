package backend

import "fmt"

type Spaceship struct {
	position Vector
	velocity Vector
	facing   Vector
}

func (self *Spaceship) String() string {
	return fmt.Sprintf("(%v, %v, %v)", self.position, self.velocity, self.facing)
}

func NewSpaceship(position *Vector) *Spaceship {
	return &Spaceship{
		position: position,
		velocity: ZeroVector(),
		facing:   NewVector(1.0, 0.0)}
}
