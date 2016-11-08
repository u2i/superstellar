package main

//go:generate go run backend/code_generation/generate_event_dispatcher.go

import (
	"log"
	"math/rand"
	"net/http"
	"superstellar/backend/server"
	"time"
)

import (
	_ "net/http/pprof"
	"superstellar/backend/events"
	"superstellar/backend/game"
	"superstellar/backend/simulation"
	"superstellar/backend/state"
	"superstellar/backend/monitor"
)

func main() {
	log.SetFlags(log.Lshortfile)

	rand.Seed(time.Now().UTC().UnixNano())

	eventDispatcher := events.NewEventDispatcher()
	physicsTicker := game.NewPhysicsTicker(eventDispatcher)

	monitor := monitor.NewMonitor()

	space := state.NewSpace()
	updater := simulation.NewUpdater(space, monitor, eventDispatcher)
	eventDispatcher.RegisterUserInputListener(updater)
	eventDispatcher.RegisterTimeTickListener(updater)
	eventDispatcher.RegisterUserJoinedListener(updater)
	eventDispatcher.RegisterUserLeftListener(updater)

	srv := server.NewServer("/superstellar", monitor, eventDispatcher)

	sender := server.NewSender(srv, space)
	eventDispatcher.RegisterTimeTickListener(sender)
	eventDispatcher.RegisterProjectileFiredListener(sender)
	eventDispatcher.RegisterUserLeftListener(sender)

	monitor.Run()
	go srv.Listen()
	go eventDispatcher.RunEventLoop()
	go physicsTicker.Run()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
