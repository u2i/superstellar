package events

import (
	"log"
	"time"
)

const (
	buffersLength = 10000
	idleDispatcherSleepTime time.Duration = 5 * time.Millisecond
)

type EventDispatcher struct {
	timeTickEvents           chan *TimeTick
	timeTickListeners        []TimeTickListener

	projectileFiredEvents    chan *ProjectileFired
	projectileFiredListeners []ProjectileFiredListener

	userInputEvents    chan *UserInput
	userInputListeners []UserInputListener
}

var instance = &EventDispatcher{
	timeTickEvents: make(chan *TimeTick, buffersLength),
	timeTickListeners: []TimeTickListener{},

	projectileFiredEvents: make(chan *ProjectileFired, buffersLength),
	projectileFiredListeners: []ProjectileFiredListener{},

	userInputEvents: make(chan *UserInput, buffersLength),
	userInputListeners: []UserInputListener{},
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
		case event := <-d.userInputEvents:
			for _, listener := range d.userInputListeners {
				listener.HandleUserInput(event)
			}
		default:
			time.Sleep(idleDispatcherSleepTime)
		}
	}
}

func (d *EventDispatcher) RegisterTimeTickListener(listener TimeTickListener) {
	d.timeTickListeners = append(d.timeTickListeners, listener)
}

func (d *EventDispatcher) FireTimeTick(e *TimeTick) {
	d.timeTickEvents <- e
}

func (d *EventDispatcher) RegisterProjectileFiredListener(listener ProjectileFiredListener) {
	d.projectileFiredListeners = append(d.projectileFiredListeners, listener)
}

func (d *EventDispatcher) FireProjectileFired(e *ProjectileFired) {
	d.projectileFiredEvents <- e
}

func (d *EventDispatcher) RegisterUserInputListener(listener UserInputListener) {
	d.userInputListeners = append(d.userInputListeners, listener)
}

func (d *EventDispatcher) FireUserInput(e *UserInput) {
	d.userInputEvents <- e
}

