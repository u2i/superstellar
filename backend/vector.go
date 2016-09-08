package backend

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64	`json:"x"`
	Y float64	`json:"y"`
}

func ZeroVector() *Vector {
	return &Vector{X: 0.0, Y: 0.0}
}

func NewVector(x, y float64) *Vector {
	return &Vector{X: x, Y: y}
}

func (self *Vector) String() string {
	return fmt.Sprintf("(%f, %f)", self.X, self.Y)
}

func (self *Vector) Add(other *Vector) *Vector {
	return &Vector{X: self.X + other.X, Y: self.Y + other.Y}
}

func (self *Vector) Multiply(scalar float64) *Vector {
	return &Vector{X: self.X * scalar, Y: self.Y * scalar}
}

func (self *Vector) Length() float64 {
	return math.Sqrt(self.X * self.X + self.Y * self.Y)
}

func (self *Vector) Normalize() *Vector {
	return &Vector{X: self.X/self.Length(), Y: self.Y/self.Length()}
}
