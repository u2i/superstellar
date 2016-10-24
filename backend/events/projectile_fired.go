package events

import "superstellar/backend/space"

type ProjectileFired struct {
	Projectile *space.Projectile
}
