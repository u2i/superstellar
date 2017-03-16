package simulation

import (
	"container/list"
	"log"
	"math"
	"math/rand"
	"superstellar/backend/constants"
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/types"
	"time"
)

// UpdatePhysics updates world physics for the next simulation step
func UpdatePhysics(space *state.Space, eventDispatcher *events.EventDispatcher) {
	detectProjectileCollisions(space, eventDispatcher)
	updateSpaceships(space, eventDispatcher)
	checkCollisions(space)

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
		var object state.Object
		object = spaceship

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

		// POSITION UPDATE

		object.SetPosition(object.Position().Add(object.Velocity()))


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

		// APPLY ANGULAR VELOCITY

		spaceship.SetAngularVelocity(spaceship.AngularVelocity() + spaceship.AngularVelocityDelta)
		spaceship.AngularVelocityDelta = 0.0

		spaceship.SetFacing(spaceship.Facing() - spaceship.AngularVelocity())
		if math.Abs(spaceship.Facing()) > math.Pi {
			spaceship.SetFacing(spaceship.Facing() - math.Copysign(2*math.Pi, spaceship.Facing()))
		}

		// NOTIFY ABOUT NEW FRAME

		spaceship.NotifyAboutNewFrame()
	}
}

func checkCollisions(s *state.Space) {
	collided := make(map[state.Object]bool)
	oldVelocity := make(map[state.Object]*types.Vector)

	for _, object := range s.Objects {

		collided[object] = true

		for _, otherObject := range s.Spaceships {
			if !collided[otherObject] && object.DetectCollision(otherObject) {
				if _, exists := oldVelocity[object]; !exists {
					oldVelocity[object] = object.Velocity().Multiply(-1.0)
				}

				if _, exists := oldVelocity[otherObject]; !exists {
					oldVelocity[otherObject] = otherObject.Velocity().Multiply(-1.0)
				}

				object.Collide(otherObject)
			}
		}
	}

	queue := list.New()
	collidedThisTurn := make(map[state.Object]bool)
	visited := make(map[state.Object]bool)

	for object := range oldVelocity {
		queue.PushBack(object)
		collidedThisTurn[object] = true
		visited[object] = true
	}

	for e := queue.Front(); e != nil; e = e.Next() {
		spaceship := e.Value.(*state.Spaceship)
		collidedThisTurn[spaceship] = true
		spaceship.SetPosition(spaceship.Position().Add(oldVelocity[spaceship]))

		for _, otherObject := range s.Objects {
			if !collidedThisTurn[otherObject] && spaceship.DetectCollision(otherObject) {
				oldVelocity[otherObject] = otherObject.Velocity().Multiply(-1.0)
				if !visited[otherObject] {
					visited[otherObject] = true
					queue.PushBack(otherObject)
				}

				spaceship.Collide(otherObject)
			}
		}
	}

	// TODO kod przeciwzakrzepowy - wywalic jak zrobimy losowe spawnowanie
	collided2 := make(map[state.Object]bool)

	for _, object := range s.Objects {
		collided2[object] = true
		for _, otherObject := range s.Objects {
			if !collided2[otherObject] && object.DetectCollision(otherObject) {
				log.Printf("COLLISON")
				if val, exists := oldVelocity[object]; exists {
					log.Printf("ov1: %f %f", val.X, val.Y)
				}
				if val, exists := oldVelocity[otherObject]; exists {
					log.Printf("ov2: %f %f", val.X, val.Y)
				}
				log.Printf("v1: %f %f", object.Velocity().X, object.Velocity().Y)
				log.Printf("v2: %f %f", otherObject.Velocity().X, otherObject.Velocity().Y)
				log.Printf("p1: %d %d", object.Position().X, object.Position().Y)
				log.Printf("p2: %d %d", otherObject.Position().X, otherObject.Position().Y)

				randAngle := rand.Float64() * 2 * math.Pi
				randMove := types.NewVector(5000, 0).Rotate(randAngle)
				object.SetPosition(object.Position().Add(randMove))
			}
		}
	}
	// koniec kodu przeciwzakrzepowego
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
