package game

import (
	"superstellar/backend/events"
	"superstellar/backend/event_dispatcher"
	"time"
	"superstellar/backend/constants"
)

type PhysicsTicker struct {
	currentFrameId  uint32
	eventDispatcher *event_dispatcher.EventDispatcher
}

func NewPhysicsTicker(eventDispatcher *event_dispatcher.EventDispatcher) *PhysicsTicker {
	return &PhysicsTicker{
		currentFrameId: 1,
		eventDispatcher: eventDispatcher,
	}
}

func (ticker *PhysicsTicker) Run() {
	var i uint32
	i = 0
	for range time.Tick(constants.PhysicsFrameDuration) {
		event := &events.TimeTick{
			FrameId: i,
		}
		ticker.eventDispatcher.FireTimeTick(event)
		i++
	}
}
