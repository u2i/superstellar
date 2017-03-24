package simulation

import (
	"math"
	"superstellar/backend/constants"
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/types"
	"time"
)

// UpdatePhysics updates world physics for the next simulation step
func UpdatePhysics(space *state.Space, eventDispatcher *events.EventDispatcher, collisionManager *CollisionManager) {
	detectProjectileCollisions(space, eventDispatcher)
	updateSpaceships(space, eventDispatcher)
	updateObjectsStates(space)
	collisionManager.resolveCollisions(space)

	space.PhysicsFrameID++
	eventDispatcher.FirePhysicsReady(&events.PhysicsReady{})

	updateProjectiles(space)
}

func detectProjectileCollisions(space *state.Space, eventDispatcher *events.EventDispatcher) {
	for projectile := range space.Projectiles {
		for clientID, spaceship := range space.Spaceships {
			collisionOccured, collisionPoint := projectile.DetectCollision(spaceship)
			if projectile.ClientID != clientID && collisionOccured {
				spaceship.CollideWithProjectile(projectile)
				applyProjectileImpulse(spaceship, projectile, collisionPoint)
				space.RemoveProjectile(projectile)
				projectileHitMessage := &events.ProjectileHit{Projectile: projectile}
				eventDispatcher.FireProjectileHit(projectileHitMessage)

				if spaceship.HP <= 0 {
					space.RemoveSpaceship(clientID)

					userDiedMessage := &events.UserDied{
						ClientID:      clientID,
						Shooter:       projectile.Spaceship,
						KilledBy:      projectile.ClientID,
						ShotSpaceship: spaceship,
					}
					eventDispatcher.FireUserDied(userDiedMessage)
				}
			}
		}
	}
}

func updateSpaceships(s *state.Space, eventDispatcher *events.EventDispatcher) {
	now := time.Now()

	for _, spaceship := range s.Spaceships {
		// FIRING

		if spaceship.Fire {
			timeSinceLastShot := now.Sub(spaceship.LastShotTime)
			if timeSinceLastShot >= constants.MinFireInterval {
				if spaceship.ShootIfPossible() {
					projectile := state.NewProjectile(s.NextProjectileID(),
						s.PhysicsFrameID, spaceship)

					s.AddProjectile(projectile)
					spaceship.LastShotTime = now

					shotEvent := &events.ProjectileFired{
						Projectile: projectile,
					}
					eventDispatcher.FireProjectileFired(shotEvent)
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

		if spaceship.TargetAngle != nil {
			spaceship.TurnToTarget()
		} else {
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
}

func updateObjectsStates(s *state.Space) {
	for _, object := range s.Objects {

		// POSITION UPDATE

		object.SetPosition(object.Position().Add(object.Velocity()))

		// APPLY ANGULAR VELOCITY

		object.SetAngularVelocity(object.AngularVelocity() + object.AngularVelocityDelta())
		object.SetAngularVelocityDelta(0.0)

		object.SetFacing(object.Facing() - object.AngularVelocity())
		if math.Abs(object.Facing()) > math.Pi {
			object.SetFacing(object.Facing() - math.Copysign(2*math.Pi, object.Facing()))
		}

		// NOTIFY ABOUT NEW FRAME

		object.NotifyAboutNewFrame()
	}
}

func applyProjectileImpulse(spaceship *state.Spaceship, projectile *state.Projectile, collisionPoint *types.Point) {
	impulse := projectile.Velocity.Multiply(constants.ProjectileImpulseStrength)

	momentOfInertia := 0.5 * constants.SpaceshipSize * constants.SpaceshipSize * constants.SpaceshipMass
	r := types.Point{X: collisionPoint.X - spaceship.Position().X, Y: collisionPoint.Y - spaceship.Position().Y}

	torque := (impulse.X*float64(r.Y) - impulse.Y*float64(r.X)) * constants.ProjectileRotationalImpulse

	spaceship.SetVelocity(spaceship.Velocity().Add(impulse.Multiply(1.0 / constants.SpaceshipMass)))
	spaceship.SetAngularVelocity(spaceship.AngularVelocity() - (torque / momentOfInertia))
}

func updateProjectiles(space *state.Space) {
	for projectile := range space.Projectiles {
		projectile.TTL--
		if projectile.TTL > 0 {
			projectile.Position = projectile.Position.Add(projectile.Velocity)
		} else {
			space.RemoveProjectile(projectile)
		}
	}
}
