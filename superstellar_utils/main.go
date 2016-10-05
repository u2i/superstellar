package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"superstellar/backend/pb"
	"time"

	"github.com/golang/protobuf/proto"

	"golang.org/x/net/websocket"
)

func main() {
	log.SetFlags(log.Lshortfile)
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) <= 2 {
		log.Fatal("[clients] [sleep time]")
		os.Exit(1)
	}

	clients, err := strconv.ParseInt(os.Args[1], 10, 0)

	if err != nil {
		log.Fatal("Incorrect number of clients")
		os.Exit(1)
	}

	sleepTime, err := time.ParseDuration(os.Args[2])

	if err != nil {
		log.Fatal("Incorrect sleep time")
		os.Exit(1)
	}

	ch := make(chan bool)

	for i := 0; i < int(clients); i++ {
		log.Printf("Running client %d\n", i+1)
		go runClient()
		time.Sleep(sleepTime)
	}

	<-ch
}

func runClient() {
	origin := "http://localhost/"
	url := "ws://localhost:8080/superstellar"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 1024)

	thrust := false
	direction := pb.Direction_NONE

	for {
		if _, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}

		if rand.Int()%100 < 5 {
			thrust = !thrust
		}

		if rand.Int()%100 < 5 {
			switch rand.Int() % 3 {
			case 0:
				direction = pb.Direction_NONE
			case 1:
				direction = pb.Direction_LEFT
			case 2:
				direction = pb.Direction_RIGHT
			}
		}

		userInput := &pb.UserInput{Thrust: thrust, Direction: direction}
		bytes, err := proto.Marshal(userInput)
		if err == nil {
			websocket.Message.Send(ws, bytes)
		}

	}
}
