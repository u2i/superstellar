package communication

import (
	"superstellar/backend/events"
	"time"
	"superstellar/backend/constants"
)

type Ticker struct {
	eventDispatcher *events.EventDispatcher
}

func NewTicker(eventDispatcher *events.EventDispatcher) *Ticker {
	return &Ticker{
		eventDispatcher: eventDispatcher,
	}
}

func (ticker *Ticker) Run() {
	var i uint32
	i = 0
	for range time.Tick(constants.UpdateSendInterval) {
		event := &events.CommunicationTimeTick{}
		ticker.eventDispatcher.FireCommunicationTimeTick(event)
		i++
	}
}
