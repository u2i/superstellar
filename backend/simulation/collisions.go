package simulation

import "superstellar/backend/state"

type Collision interface {
	collide(state.Object, state.Object)
}
