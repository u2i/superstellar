package event_listeners

import "superstellar/backend/events"

type TimeTickListener interface {
	HandleTimeTick(*events.TimeTick)
}
