package backend

import "fmt"

type Vector struct {
	x, y float64
}

func (self *Vector) String() string {
	return fmt.Sprintf("(%f, %f)", self.x, self.y)
}

func (self *Vector) Add(other *Vector) *Vector {
	return &Vector{x: self.x + other.x, y: self.y + other.y}
}

func (self *Vector) Multiply(scalar float64) *Vector {
	return &Vector{x: self.x * scalar, y: self.y * scalar}
}
