package events

import "superstellar/backend/space"

const QueueLength = 1000

type EventProcessor struct {
	eventQueue   chan Event
	eventContext *EventContext
}

func NewEventProcessor(space *space.Space) *EventProcessor {
	context := &EventContext{
		Space: space,
	}

	return &EventProcessor{
		eventQueue: make(chan Event, QueueLength),
		eventContext: context,
	}
}

func (processor *EventProcessor) ProcessEvents() {
	for {
		select {
		case event := <-processor.eventQueue:
			event.handle(processor.eventContext)
		default:
			return
		}
	}
}

func (processor *EventProcessor) AddEvent(event Event) {
	processor.eventQueue <- event
}
