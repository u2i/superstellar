package event_dispatcher

import (
	"superstellar/backend/events"
	"superstellar/backend/event_listeners"
	"log"
	"time"
)

const (
	buffersLength = 10000
	idleDispatcherSleepTime time.Duration = 5 * time.Millisecond
)

type EventDispatcher struct {
	timeTickEvents           chan *events.TimeTick
	timeTickListeners        []event_listeners.TimeTickListener

	projectileFiredEvents    chan *events.ProjectileFired
	projectileFiredListeners []event_listeners.ProjectileFiredListener
}

var instance = &EventDispatcher{
	timeTickEvents: make(chan *events.TimeTick, buffersLength),
	timeTickListeners: []event_listeners.TimeTickListener{},

	projectileFiredEvents: make(chan *events.ProjectileFired, buffersLength),
	projectileFiredListeners: []event_listeners.ProjectileFiredListener{},
}

func Instance() *EventDispatcher {
	return instance
}

func (d *EventDispatcher) RunEventLoop() {
	log.Printf("Enabling event dispatcher...")

	for {
		select {
		case event := <-d.timeTickEvents:
			for _, listener := range d.timeTickListeners {
				listener.HandleTimeTick(event)
			}
		case event := <-d.projectileFiredEvents:
			for _, listener := range d.projectileFiredListeners {
				listener.HandleProjectileFired(event)
			}
		default:
			time.Sleep(idleDispatcherSleepTime)
		}
	}
}

func (d *EventDispatcher) RegisterTimeTickListener(listener event_listeners.TimeTickListener) {
	d.timeTickListeners = append(d.timeTickListeners, listener)
}

func (d *EventDispatcher) FireTimeTick(e *events.TimeTick) {
	d.timeTickEvents <- e
}

func (d *EventDispatcher) RegisterProjectileFiredListener(listener event_listeners.ProjectileFiredListener) {
	d.projectileFiredListeners = append(d.projectileFiredListeners, listener)
}

func (d *EventDispatcher) FireProjectileFired(e *events.ProjectileFired) {
	d.projectileFiredEvents <- e
}
