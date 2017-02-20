package main

//go:generate go run backend/code_generation/generate_event_dispatcher.go

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

import (
	_ "net/http/pprof"
	"os"
	"superstellar/backend/communication"
	"superstellar/backend/events"
	"superstellar/backend/game"
	"superstellar/backend/monitor"
	"superstellar/backend/simulation"
	"superstellar/backend/state"
	"superstellar/backend/ai"
	"superstellar/backend/utils"
)

func main() {
	log.SetFlags(log.Lshortfile)

	debug := false
	if len(os.Args) == 2 && os.Args[1] == "-d" {
		debug = true
	}

	rand.Seed(time.Now().UTC().UnixNano())

	eventDispatcher := events.NewEventDispatcher()
	physicsTicker := game.NewPhysicsTicker(eventDispatcher)

	monitor := monitor.NewMonitor(eventDispatcher)

	space := state.NewSpace()
	updater := simulation.NewUpdater(space, monitor, eventDispatcher)
	eventDispatcher.RegisterUserInputListener(updater)
	eventDispatcher.RegisterTimeTickListener(updater)
	eventDispatcher.RegisterUserJoinedListener(updater)
	eventDispatcher.RegisterUserLeftListener(updater)
	eventDispatcher.RegisterUserDiedListener(updater)
	eventDispatcher.RegisterTargetAngleListener(updater)

	clientIdSequencer := utils.NewIdSequencer()

	srv := communication.NewServer("/superstellar", monitor, eventDispatcher, clientIdSequencer)
	eventDispatcher.RegisterUserLeftListener(srv)

	sender := communication.NewSender(srv, space)
	eventDispatcher.RegisterPhysicsReadyListener(sender)
	eventDispatcher.RegisterProjectileFiredListener(sender)
	eventDispatcher.RegisterProjectileHitListener(sender)
	eventDispatcher.RegisterUserLeftListener(sender)
	eventDispatcher.RegisterUserJoinedListener(sender)
	eventDispatcher.RegisterUserDiedListener(sender)

	botManager := ai.NewBotManager(space, clientIdSequencer)
	botManager.CreateNewBot()
	botManager.CreateNewBot()
	botManager.CreateNewBot()
	botManager.CreateNewBot()
	eventDispatcher.RegisterTimeTickListener(botManager)

	if debug {
		fileWriter, err := communication.NewFileWriter(space)
		if err != nil {
			panic(err)
		}

		eventDispatcher.RegisterTimeTickListener(fileWriter)
		go fileWriter.Run()
	}

	monitor.Run()
	go srv.Listen()

	go eventDispatcher.RunEventLoop()
	go physicsTicker.Run()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
