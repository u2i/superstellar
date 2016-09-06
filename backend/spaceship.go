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
