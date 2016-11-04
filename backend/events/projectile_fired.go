package events

import "superstellar/backend/state"

type ProjectileFired struct {
	Projectile *state.Projectile
}
