package types

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

// NewPoint initializes new vector with given parameters.
func NewPoint(x, y int32) *Point {
	return &Point{X: x, Y: y}
}

// NewPointFromPolar initializes new vector form given polar coordinates.
func NewPointFromPolar(angle float64, radius uint32) *Point {
	return &Point{
		X: int32(float64(radius) * math.Cos(angle)),
		Y: int32(float64(radius) * math.Sin(angle)),
	}
}

// String returns string representation.
func (point *Point) String() string {
	return fmt.Sprintf("(%d, %d)", point.X, point.Y)
}

// Add adds given Vector and return Point.
func (point *Point) Add(other *Vector) *Point {
	return &Point{
		X: point.X + int32(other.X),
		Y: point.Y + int32(other.Y),
	}
}

// Distance returns the distance between two points
func (point *Point) Distance(other *Point) float64 {
	dx := point.X - other.X
	dy := point.Y - other.Y

	return math.Sqrt(float64(dx*dx + dy*dy))
}

// ToProto returns protobuf representation
func (point *Point) ToProto() *pb.Point {
	return &pb.Point{
		X: point.X,
		Y: point.Y,
	}
}

// Length returns length of the vector.
func (point *Point) Length() float64 {
	return math.Sqrt(float64(point.X)*float64(point.X) + float64(point.Y)*float64(point.Y))
}

// Normalize returns a new normalized vector.
func (point *Point) Normalize() *Vector {
	length := point.Length()
	return &Vector{X: float64(point.X) / length, Y: float64(point.Y) / length}
}
