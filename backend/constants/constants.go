package constants

import (
	"time"
)

const (
	PhysicsFrameDuration = 20 * time.Millisecond

	// BoundaryAnnulusWidth is the width of boundary region (in .01 units), i.e. from WorldRadius till when no more movement is possible
	BoundaryAnnulusWidth = 20000

	// FrictionCoefficient is the coefficient saying how fast a spaceship will slow down when not using acceleration
	FrictionCoefficient = 0.005

	// MinFireInterval is a minimum time between firing.
	MinFireInterval = 250 * time.Millisecond

	// RandomPositionEmptyRadius describes the minimum radius around randomized
	// initial position that needs to be free of any objects.
	RandomPositionEmptyRadius = 5000.0

	// Acceleration is spaceship's linear acceleration on thruster.
	SpaceshipAcceleration = 50.0

	// Maximum angular velocity added on user input.
	SpaceshipMaxAngularSpeed = 0.12

	SpaceshipLinearAngularAcceleration = 0.0001

	SpaceshipNonlinearAngularAcceleration = 2

	SpaceshipAngularFriction = 0.2

	// SpaceshipTurnToAngleP is TurnToTarget propotional gain.
	SpaceshipTurnToAngleP = 0.9

	// MaxSpeed maximum speed of the spacecraft
	SpaceshipMaxSpeed = 600

	// SpaceshipSize is spaceship's radius
	SpaceshipSize = 2200

	// SpaceshipInitialHP spaceship HP
	SpaceshipInitialHP = 500

	// SpaceshipInitialEnergy spaceship Energy
	SpaceshipInitialEnergy = 100

	// Part of killed player's Max HP that killing player receives
	KillRewardRatio = 0.2

	// Energy reward for each
	KillEnergyRewardRatio = 0.1

	// Energy cost of a single shot
	BasicWeaponEnergyCost = 33

	// WorldRadius is the radius of playable world (in .01 units)
	WorldRadius = 100000

	// AutoRepairDelay is time after which spaceship will repair itself
	AutoRepairDelay = 250

	// AutoRepairInterval is the time between repairs
	AutoRepairInterval = 1

	// Time between Energy renewal
	AutoEnergyRepairInterval = 1

	// AutoRepair is the amount of HP player will receive after AutoRepairDelay
	AutoRepairAmount = 2

	// EnergyRepair is the amount of Energy player will receive after AutoRepairDelay
	AutoEnergyRepairAmount = 2;

	// Number of best players in Leaderboard
	LeaderboardLength = 10
)
