package game

import (
	"superstellar/backend/constants"
	"superstellar/backend/events"
	"time"
)

type PhysicsTicker struct {
	currentFrameId  uint32
	eventDispatcher *events.EventDispatcher
}

func NewPhysicsTicker(eventDispatcher *events.EventDispatcher) *PhysicsTicker {
	return &PhysicsTicker{
		currentFrameId:  1,
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
