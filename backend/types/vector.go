package types

import (
	"fmt"
	"math"
	"superstellar/backend/pb"
)

// Vector structs holds 2D vector.
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// ZeroVector initializes new zero vector.
func ZeroVector() *Vector {
	return &Vector{X: 0.0, Y: 0.0}
}

// NewVector initlizes new vector with given parameters.
func NewVector(x, y float64) *Vector {
	return &Vector{X: x, Y: y}
}

// String returns string representation.
func (v *Vector) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}

// Add returns new Vector that is a sum of the two given.
func (v *Vector) Add(other *Vector) *Vector {
	return &Vector{X: v.X + other.X, Y: v.Y + other.Y}
}

// Multiply returns new Vector that is a product of the the vector and
// given scalar.
func (v *Vector) Multiply(scalar float64) *Vector {
	return &Vector{X: v.X * scalar, Y: v.Y * scalar}
}

// Length returns length of the vector.
func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns a new normalized vector.
func (v *Vector) Normalize() *Vector {
	return &Vector{X: v.X / v.Length(), Y: v.Y / v.Length()}
}

// Radians returns vector's angle in radians.
func (v *Vector) Radians() float64 {
	return math.Atan2(-v.Y, v.X)
}

// Rotate returns rotated vector by give angle.
func (v *Vector) Rotate(angle float64) *Vector {
	return &Vector{
		X: math.Cos(angle)*v.X - math.Sin(angle)*v.Y,
		Y: math.Sin(angle)*v.X + math.Cos(angle)*v.Y,
	}
}

// ToProto returns Vector's protobuf representation
func (v *Vector) ToProto() *pb.Vector {
	return &pb.Vector{
		X: float32(v.X),
		Y: float32(v.Y),
	}
}
