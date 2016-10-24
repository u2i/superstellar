package event_listeners

import "superstellar/backend/events"

type ProjectileFiredListener interface {
	HandleProjectileFired(*events.ProjectileFired)
}
