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
	"superstellar/backend/ai"
	"superstellar/backend/communication"
	"superstellar/backend/events"
	"superstellar/backend/game"
	"superstellar/backend/monitor"
	"superstellar/backend/persistence"
	"superstellar/backend/simulation"
	"superstellar/backend/state"
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
	userNameRegistry := utils.NewUserNameRegistry()
	idManager := utils.NewIdManager()

	monitor := monitor.NewMonitor(eventDispatcher)

	space := state.NewSpace()
	updater := simulation.NewUpdater(space, monitor, eventDispatcher, idManager)
	eventDispatcher.RegisterUserInputListener(updater)
	eventDispatcher.RegisterTimeTickListener(updater)
	eventDispatcher.RegisterUserJoinedListener(updater)
	eventDispatcher.RegisterUserLeftListener(updater)
	eventDispatcher.RegisterObjectDestroyedListener(updater)
	eventDispatcher.RegisterTargetAngleListener(updater)

	server := communication.NewServer("/superstellar", monitor, eventDispatcher, idManager, userNameRegistry)
	eventDispatcher.RegisterUserLeftListener(server)

	sender := communication.NewSender(server, space, userNameRegistry)
	eventDispatcher.RegisterPhysicsReadyListener(sender)
	eventDispatcher.RegisterProjectileFiredListener(sender)
	eventDispatcher.RegisterProjectileHitListener(sender)
	eventDispatcher.RegisterUserLeftListener(sender)
	eventDispatcher.RegisterUserConnectedListener(sender)
	eventDispatcher.RegisterUserJoinedListener(sender)
	eventDispatcher.RegisterObjectDestroyedListener(sender)

	if _, found := os.LookupEnv("DYNAMODB_ENDPOINT"); found {
		adapter := persistence.NewDynamoDbWriter()

		scoreBoardSerializer := persistence.NewScoreBoardSerializer(userNameRegistry, adapter, idManager, eventDispatcher)
		eventDispatcher.RegisterObjectDestroyedListener(scoreBoardSerializer)

		scoreBoardReader := persistence.NewScoreBoardReader(adapter)
		scoreBoardSender := communication.NewScoreBoardSender(server, scoreBoardReader, userNameRegistry)
		eventDispatcher.RegisterUserJoinedListener(scoreBoardSender)
		eventDispatcher.RegisterScoreSentListener(scoreBoardSender)
	}

	botManager := ai.NewBotManager(eventDispatcher, space, idManager, userNameRegistry)
	eventDispatcher.RegisterTimeTickListener(botManager)
	eventDispatcher.RegisterObjectDestroyedListener(botManager)
	botManager.CreateBots(5)

	if debug {
		fileWriter, err := communication.NewFileWriter(space)
		if err != nil {
			panic(err)
		}

		eventDispatcher.RegisterTimeTickListener(fileWriter)
		go fileWriter.Run()
	}

	monitor.Run()
	go server.Listen()

	go eventDispatcher.RunEventLoop()
	go physicsTicker.Run()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
