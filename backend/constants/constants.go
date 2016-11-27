package constants

import (
	"time"
)

const (
	PhysicsFrameDuration = 20 * time.Millisecond

	UpdateSendInterval = 200 * time.Millisecond

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
	SpaceshipInitialEnergy = 300

	// Part of killed player's Max HP that killing player receives
	KillRewardRatio = 0.2

	// Energy reward for each
	KillEnergyRewardRatio = 0.2

	// Energy cost of a single shot
	BasicWeaponEnergyCost = 100

	// WorldRadius is the radius of playable world (in .01 units)
	WorldRadius = 100000

	// AutoRepairDelay is time after which spaceship will repair itself
	AutoRepairDelay = 250

	// AutoRepairInterval is the time between repairs
	AutoRepairInterval = 1

	// AutoRepair is the amount of HP player will receive after AutoRepairDelay
	AutoRepairAmount = 2

	// Time between Energy renewal
	AutoEnergyRechargeInterval = 1

	// EnergyRecharge is the amount of Energy player will receive after AutoRepairDelay
	AutoEnergyRechargeAmount = 3

	// Number of best players in Leaderboard
	LeaderboardLength = 10

	// MinimumUsernameLength is minimal the number of characters in username
	MinimumUsernameLength = 3

	// MaximumUsernameLength is maximal the number of characters in username
	MaximumUsernameLength = 25

	// ProjectileSpeed describes projectile speed. Captain Obvious
	ProjectileSpeed = 4000

	// DefaultTTL describes the default number of frames the projectile lives
	ProjectileDefaultTTL = 15

	// ProjectileDamage is the damage that is inflicted on a user upon colliding with projectile
	ProjectileDamage = 50
)
