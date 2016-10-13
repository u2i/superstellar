package backend

import (
	"fmt"
	"math"
	"superstellar/backend/pb"
)

// Point structs holds 2D vector with int coordinates.
type Point struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// ZeroPoint initializes new zero vector.
func ZeroPoint() *Point {
	return &Point{X: 0, Y: 0}
}

// NewPoint initlizes new vector with given parameters.
func NewPoint(x, y int32) *Point {
	return &Point{X: x, Y: y}
}

// String returns string representation.
func (v *Point) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

// Add adds given Vector and return Point.
func (v *Point) Add(other *Vector) *Point {
	return &Point{
		X: v.X + int32(other.X),
		Y: v.Y + int32(other.Y),
	}
}

func (v *Point) toProto() *pb.Point {
	return &pb.Point{
		X: v.X,
		Y: v.Y,
	}
}

// Length returns length of the vector.
func (v *Point) Length() float64 {
	return math.Sqrt(float64(v.X)*float64(v.X) + float64(v.Y)*float64(v.Y))
}

// Normalize returns a new normalized vector.
func (v *Point) Normalize() *Vector {
	return &Vector{X: float64(v.X) / v.Length(), Y: float64(v.Y) / v.Length()}
}
