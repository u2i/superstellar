package main

import (
	"log"
	"math/rand"
	"net/http"
	"superstellar/backend/server"
	"time"
)

import (
	_ "net/http/pprof"
	"superstellar/backend/event_dispatcher"
	"superstellar/backend/game"
)

func main() {
	log.SetFlags(log.Lshortfile)

	rand.Seed(time.Now().UTC().UnixNano())

	eventDispatcher := event_dispatcher.Instance()
	physicsTicker := game.NewPhysicsTicker(eventDispatcher)

	server := server.NewServer("/superstellar", eventDispatcher)
	go server.Listen()

	go eventDispatcher.RunEventLoop()

	go physicsTicker.Run()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
