package backend

import (
	"fmt"
	"math"
	"superstellar/backend/pb"
)

// IntVector structs holds 2D vector with int coordinates.
type IntVector struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// ZeroIntVector initializes new zero vector.
func ZeroIntVector() *IntVector {
	return &IntVector{X: 0, Y: 0}
}

// NewIntVector initlizes new vector with given parameters.
func NewIntVector(x, y int32) *IntVector {
	return &IntVector{X: x, Y: y}
}

// String returns string representation.
func (v *IntVector) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

// Add adds given Vector and return IntVector.
func (v *IntVector) Add(other *Vector) *IntVector {
	return &IntVector{
		X: v.X + int32(other.X*100),
		Y: v.Y + int32(other.Y*100),
	}
}

func (v *IntVector) toProto() *pb.IntVector {
	return &pb.IntVector{
		X: v.X,
		Y: v.Y,
	}
}

// Length returns length of the vector.
func (v *IntVector) Length() float64 {
	return math.Sqrt(float64(v.X)*float64(v.X) + float64(v.Y)*float64(v.Y))
}

// Normalize returns a new normalized vector.
func (v *IntVector) Normalize() *Vector {
	return &Vector{X: float64(v.X) / v.Length(), Y: float64(v.Y) / v.Length()}
}
