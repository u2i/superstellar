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
		for id, object := range manager.space.Objects {
			collisionOccured, collisionPoint := projectile.DetectCollision(object)
			if projectile.ClientID != id && collisionOccured {
				object.CollideWithProjectile(projectile)
				projectile.Spaceship.ProjectileHitOtherObject(object)
				manager.applyProjectileImpulse(object, projectile, collisionPoint)
				manager.space.RemoveProjectile(projectile)

				projectileHitMessage := &events.ProjectileHit{Projectile: projectile}
				manager.eventDispatcher.FireProjectileHit(projectileHitMessage)

				if object.Hp() <= 0 {
					manager.space.RemoveObject(id)

					objectDestroyedMessage := &events.ObjectDestroyed{
						DestroyedObject: object,
						DestroyedBy:     projectile.Spaceship,
						Timestamp:       time.Now(),
					}
					manager.eventDispatcher.FireObjectDestroyed(objectDestroyedMessage)
				}
			}
		}
	}
}

func (manager *ProjectileManager) applyProjectileImpulse(object state.Object, projectile *state.Projectile, collisionPoint *types.Point) {
	impulse := projectile.Velocity.Multiply(constants.ProjectileImpulseStrength)

	momentOfInertia := 0.5 * constants.SpaceshipSize * constants.SpaceshipSize * constants.SpaceshipMass
	r := types.Point{X: collisionPoint.X - object.Position().X, Y: collisionPoint.Y - object.Position().Y}

	torque := (impulse.X*float64(r.Y) - impulse.Y*float64(r.X)) * constants.ProjectileRotationalImpulse

	object.SetVelocity(object.Velocity().Add(impulse.Multiply(1.0 / constants.SpaceshipMass)))
	object.SetAngularVelocity(object.AngularVelocity() - (torque / momentOfInertia))
}
