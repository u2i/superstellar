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
	"superstellar/backend/state"
	"superstellar/backend/simulation"
)

func main() {
	log.SetFlags(log.Lshortfile)

	rand.Seed(time.Now().UTC().UnixNano())

	eventDispatcher := events.NewEventDispatcher()
	physicsTicker := game.NewPhysicsTicker(eventDispatcher)

	space := state.NewSpace()
	updater := simulation.NewUpdater(space, eventDispatcher)
	eventDispatcher.RegisterUserInputListener(updater)
	eventDispatcher.RegisterTimeTickListener(updater)

	server := server.NewServer("/superstellar", eventDispatcher, space)
	go server.Listen()

	go eventDispatcher.RunEventLoop()

	go physicsTicker.Run()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
