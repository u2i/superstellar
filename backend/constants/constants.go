package constants

import "time"

const (
	PhysicsFrameDuration = 20 * time.Millisecond

	// BoundaryAnnulusWidth is the width of boundary region (in .01 units), i.e. from WorldRadius till when no more movement is possible
	BoundaryAnnulusWidth = 20000

	// FrictionCoefficient is the coefficient saying how fast a spaceship will slow down when not using acceleration
	FrictionCoefficient = 0.02

	// MinFireInterval is a minimum time between firing.
	MinFireInterval = 500 * time.Millisecond

	// RandomPositionEmptyRadius describes the minimum radius around randomized
	// initial position that needs to be free of any objects.
	RandomPositionEmptyRadius = 5000.0

	// Acceleration is spaceship's linear acceleration on thruster.
	SpaceshipAcceleration = 30.0

	// AngularVelocity is an angular velocity added on user input.
	SpaceshipAngularVelocity = 0.11

	// MaxSpeed maximum speed of the spacecraft
	SpaceshipMaxSpeed = 1500

	// SpaceshipSize is spaceship's radius
	SpaceshipSize = 2000

	// SpaceshipInitialHP spaceship HP
	SpaceshipInitialHP = 500

	// WorldRadius is the radius of playable world (in .01 units)
	WorldRadius = 100000
)
