package simulation

import (
	"superstellar/backend/constants"
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/types"
	"time"
)

type ProjectileManager struct {
	space           *state.Space
	eventDispatcher *events.EventDispatcher
}

func NewProjectileManager(space *state.Space, eventDispatcher *events.EventDispatcher) *ProjectileManager {
	return &ProjectileManager{
		space:           space,
		eventDispatcher: eventDispatcher,
	}
}

func (manager *ProjectileManager) updateProjectiles() {
	for projectile := range manager.space.Projectiles {
		projectile.TTL--
		if projectile.TTL > 0 {
			projectile.Position = projectile.Position.Add(projectile.Velocity)
		} else {
			manager.space.RemoveProjectile(projectile)
		}
	}
}

func (manager *ProjectileManager) detectProjectileCollisions() {
	for projectile := range manager.space.Projectiles {
		for clientID, spaceship := range manager.space.Spaceships {
			collisionOccured, collisionPoint := projectile.DetectCollision(spaceship)
			if projectile.ClientID != clientID && collisionOccured {
				spaceship.CollideWithProjectile(projectile)
				projectile.Spaceship.ProjectileHitOtherSpaceship(spaceship)
				manager.applyProjectileImpulse(spaceship, projectile, collisionPoint)
				manager.space.RemoveProjectile(projectile)

				projectileHitMessage := &events.ProjectileHit{Projectile: projectile}
				manager.eventDispatcher.FireProjectileHit(projectileHitMessage)

				if spaceship.HP <= 0 {
					manager.space.RemoveSpaceship(clientID)

					userDiedMessage := &events.UserDied{
						ClientID:      clientID,
						Shooter:       projectile.Spaceship,
						KilledBy:      projectile.ClientID,
						ShotSpaceship: spaceship,
						Timestamp:     time.Now(),
					}
					manager.eventDispatcher.FireUserDied(userDiedMessage)
				}
			}
		}
	}
}

func (manager *ProjectileManager) applyProjectileImpulse(spaceship *state.Spaceship, projectile *state.Projectile, collisionPoint *types.Point) {
	impulse := projectile.Velocity.Multiply(constants.ProjectileImpulseStrength)

	momentOfInertia := 0.5 * constants.SpaceshipSize * constants.SpaceshipSize * constants.SpaceshipMass
	r := types.Point{X: collisionPoint.X - spaceship.Position().X, Y: collisionPoint.Y - spaceship.Position().Y}

	torque := (impulse.X*float64(r.Y) - impulse.Y*float64(r.X)) * constants.ProjectileRotationalImpulse

	spaceship.SetVelocity(spaceship.Velocity().Add(impulse.Multiply(1.0 / constants.SpaceshipMass)))
	spaceship.SetAngularVelocity(spaceship.AngularVelocity() - (torque / momentOfInertia))
}
