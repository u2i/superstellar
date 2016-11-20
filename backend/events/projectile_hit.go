package events

import "superstellar/backend/state"

type ProjectileHit struct {
	Projectile *state.Projectile
}
