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
	updateProjectiles(space)
}

func detectProjectileCollisions(space *state.Space, eventDispatcher *events.EventDispatcher) {
	for projectile := range space.Projectiles {
		for clientID, spaceship := range space.Spaceships {
			if projectile.ClientID != clientID && projectile.DetectCollision(spaceship) {
				spaceship.CollideWithProjectile(projectile)
				space.RemoveProjectile(projectile)

				if spaceship.HP <= 0 {
					space.RemoveSpaceship(clientID)

					userDiedMessage := &events.UserDied{
						ClientID: clientID,
						KilledBy: projectile.ClientID,
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
		if spaceship.Fire {
			timeSinceLastShot := now.Sub(spaceship.LastShotTime)
			if timeSinceLastShot >= constants.MinFireInterval {
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

		if spaceship.InputThrust {
			deltaVelocity := spaceship.NormalizedFacing().Multiply(constants.SpaceshipAcceleration)
			spaceship.Velocity = spaceship.Velocity.Add(deltaVelocity)
		} else {
			if spaceship.Velocity.Length() != 0 {
				spaceship.Velocity = spaceship.Velocity.Multiply(1 - constants.FrictionCoefficient)

				if spaceship.Velocity.Length() < 1 {
					spaceship.Velocity = types.ZeroVector()
				}
			}
		}

		if spaceship.Position.Add(spaceship.Velocity).Length() > constants.WorldRadius {
			outreachLength := spaceship.Position.Length() - constants.WorldRadius
			gravityAcceleration := -(outreachLength / constants.BoundaryAnnulusWidth) * constants.SpaceshipAcceleration
			deltaVelocity := spaceship.Position.Normalize().Multiply(gravityAcceleration)
			spaceship.Velocity = spaceship.Velocity.Add(deltaVelocity)
		}

		if spaceship.Velocity.Length() > constants.SpaceshipMaxSpeed {
			spaceship.Velocity = spaceship.Velocity.Normalize().Multiply(constants.SpaceshipMaxSpeed)
		}

		spaceship.Position = spaceship.Position.Add(spaceship.Velocity)

		angle := math.Atan2(spaceship.Facing.Y, spaceship.Facing.X)
		switch spaceship.InputDirection {
		case state.LEFT:
			angle += constants.SpaceshipAngularVelocity
		case state.RIGHT:
			angle -= constants.SpaceshipAngularVelocity
		}

		spaceship.Facing = types.NewVector(math.Cos(angle), math.Sin(angle))
	}

	collided := make(map[*state.Spaceship]bool)
	oldVelocity := make(map[*state.Spaceship]*types.Vector)

	for _, spaceship := range s.Spaceships {

		collided[spaceship] = true

		for _, otherSpaceship := range s.Spaceships {
			if !collided[otherSpaceship] && spaceship.DetectCollision(otherSpaceship) {
				if _, exists := oldVelocity[spaceship]; !exists {
					oldVelocity[spaceship] = spaceship.Velocity.Multiply(-1.0)
				}

				if _, exists := oldVelocity[otherSpaceship]; !exists {
					oldVelocity[otherSpaceship] = otherSpaceship.Velocity.Multiply(-1.0)
				}

				spaceship.Collide(otherSpaceship)
			}
		}
	}

	queue := list.New()
	collidedThisTurn := make(map[*state.Spaceship]bool)
	visited := make(map[*state.Spaceship]bool)

	for spaceship := range oldVelocity {
		queue.PushBack(spaceship)
		collidedThisTurn[spaceship] = true
		visited[spaceship] = true
	}

	for e := queue.Front(); e != nil; e = e.Next() {
		spaceship := e.Value.(*state.Spaceship)
		collidedThisTurn[spaceship] = true
		spaceship.Position = spaceship.Position.Add(oldVelocity[spaceship])

		for _, otherSpaceship := range s.Spaceships {
			if !collidedThisTurn[otherSpaceship] && spaceship.DetectCollision(otherSpaceship) {
				oldVelocity[otherSpaceship] = otherSpaceship.Velocity.Multiply(-1.0)
				if !visited[otherSpaceship] {
					visited[otherSpaceship] = true
					queue.PushBack(otherSpaceship)
				}

				spaceship.Collide(otherSpaceship)
			}
		}
	}

	// TODO kod przeciwzakrzepowy - wywalic jak zrobimy losowe spawnowanie
	collided2 := make(map[*state.Spaceship]bool)

	for _, spaceship := range s.Spaceships {
		collided2[spaceship] = true
		for _, otherSpaceship := range s.Spaceships {
			if !collided2[otherSpaceship] && spaceship.DetectCollision(otherSpaceship) {
				log.Printf("COLLISON")
				if val, exists := oldVelocity[spaceship]; exists {
					log.Printf("ov1: %f %f", val.X, val.Y)
				}
				if val, exists := oldVelocity[otherSpaceship]; exists {
					log.Printf("ov2: %f %f", val.X, val.Y)
				}
				log.Printf("v1: %f %f", spaceship.Velocity.X, spaceship.Velocity.Y)
				log.Printf("v2: %f %f", otherSpaceship.Velocity.X, otherSpaceship.Velocity.Y)
				log.Printf("p1: %d %d", spaceship.Position.X, spaceship.Position.Y)
				log.Printf("p2: %d %d", otherSpaceship.Position.X, otherSpaceship.Position.Y)

				randAngle := rand.Float64() * 2 * math.Pi
				randMove := types.NewVector(5000, 0).Rotate(randAngle)
				spaceship.Position = spaceship.Position.Add(randMove)
			}
		}
	}
	// koniec kodu przeciwzakrzepowego

	s.PhysicsFrameID++
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
