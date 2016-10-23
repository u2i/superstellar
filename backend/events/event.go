package events

type Event interface {
	handle(c *EventContext)
}
