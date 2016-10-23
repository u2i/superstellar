package events

import (
	"superstellar/backend/space"
)

type ProjectileHitEvent struct {
	Projectile *space.Projectile
	Spaceship  *space.Spaceship
}

func (event *ProjectileHitEvent) handle(context *EventContext) {
	event.Spaceship.CollideWithProjectile(event.Projectile)
	context.Space.RemoveProjectile(event.Projectile)
}
