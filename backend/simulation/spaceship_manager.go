package simulation

import (
	"math"
	"superstellar/backend/constants"
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/types"
	"time"
)

type SpaceshipManager struct {
	space           *state.Space
	eventDispatcher *events.EventDispatcher
}

func NewSpaceshipManager(space *state.Space, eventDispatcher *events.EventDispatcher) *SpaceshipManager {
	return &SpaceshipManager{
		space:           space,
		eventDispatcher: eventDispatcher,
	}
}

func (manager *SpaceshipManager) updateSpaceships() {
	now := time.Now()

	for _, spaceship := range manager.space.Spaceships {
		// FIRING

		if spaceship.IsFiring() {
			timeSinceLastShot := now.Sub(spaceship.LastShotTime)
			if timeSinceLastShot >= constants.MinFireInterval {
				if spaceship.ShootIfPossible() {
					projectile := state.NewProjectile(manager.space.NextProjectileID(),
						manager.space.PhysicsFrameID, spaceship)

					manager.space.AddProjectile(projectile)
					spaceship.LastShotTime = now

					shotEvent := &events.ProjectileFired{
						Projectile: projectile,
					}
					manager.eventDispatcher.FireProjectileFired(shotEvent)
				}
			}
		}

		// THRUST

		if spaceship.InputThrust || spaceship.InputBoost {
			deltaVelocity := types.NewVector(math.Cos(spaceship.Facing()), -math.Sin(spaceship.Facing())).Multiply(constants.SpaceshipAcceleration)
			spaceship.SetVelocity(spaceship.Velocity().Add(deltaVelocity))
		} else {
			if spaceship.Velocity().Length() != 0 {
				spaceship.SetVelocity(spaceship.Velocity().Multiply(1 - constants.FrictionCoefficient))

				if spaceship.Velocity().Length() < 1 {
					spaceship.SetVelocity(types.ZeroVector())
				}
			}
		}

		// ANNULUS

		if spaceship.Position().Add(spaceship.Velocity()).Length() > constants.WorldRadius {
			outreachLength := spaceship.Position().Length() - constants.WorldRadius
			gravityAcceleration := -(outreachLength / constants.BoundaryAnnulusWidth) * constants.SpaceshipAcceleration
			deltaVelocity := spaceship.Position().Normalize().Multiply(gravityAcceleration)
			spaceship.SetVelocity(spaceship.Velocity().Add(deltaVelocity))
		}

		// BOOST

		boostActive := false

		if spaceship.InputBoost {
			boostActive = spaceship.BoostIfPossible()
		}

		spaceship.InputBoost = boostActive

		maxVelocity := float64(constants.SpaceshipMaxSpeed)
		if boostActive {
			maxVelocity *= constants.SpaceshipBoostFactor
		}

		// VELOCITY LIMITING

		if spaceship.Velocity().Length() > maxVelocity {
			spaceship.SetVelocity(spaceship.Velocity().Normalize().Multiply(maxVelocity))
			// TODO: add easing when returning to base velocity
		}

		// TURNING

		switch spaceship.InputDirection {
		case state.LEFT:
			spaceship.LeftTurn()
		case state.RIGHT:
			spaceship.RightTurn()
		default:
			spaceship.ApplyAngularFriction()
		}
	}
}
