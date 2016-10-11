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

	if len(os.Args) <= 3 {
		log.Fatal("[host] [clients] [sleep time]")
		os.Exit(1)
	}

	host := os.Args[1]

	clients, err := strconv.ParseInt(os.Args[2], 10, 0)

	if err != nil {
		log.Fatal("Incorrect number of clients")
		os.Exit(1)
	}

	sleepTime, err := time.ParseDuration(os.Args[3])

	if err != nil {
		log.Fatal("Incorrect sleep time")
		os.Exit(1)
	}

	ch := make(chan bool)

	for i := 0; i < int(clients); i++ {
		log.Printf("Running client %d\n", i+1)
		go runClient(host)
		time.Sleep(sleepTime)
	}

	<-ch
}

func runClient(host string) {
	origin := "http://" + host
	url := "ws://" + host + ":8080/superstellar"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 1024)

	thrust := false
	fire := false
	userInput := pb.UserInput_CENTER

	for {
		if _, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}

		if rand.Int()%100 < 5 {
			thrust = !thrust
			if thrust {
				userInput = pb.UserInput_THRUST_ON
			} else {
				userInput = pb.UserInput_THRUST_OFF
			}
		} else if rand.Int()%100 < 5 {
			switch rand.Int() % 3 {
			case 0:
				userInput = pb.UserInput_CENTER
			case 1:
				userInput = pb.UserInput_LEFT
			case 2:
				userInput = pb.UserInput_RIGHT
			}
		} else if rand.Int()%100 < 5 {
			fire = !fire
			if fire {
				userInput = pb.UserInput_FIRE_START
			} else {
				userInput = pb.UserInput_FIRE_STOP
			}
		}

		send(ws, userInput)
	}
}

func send(ws *websocket.Conn, userInput pb.UserInput) {
	userMessage := &pb.UserMessage{UserInput: userInput}
	bytes, err := proto.Marshal(userMessage)
	if err == nil {
		websocket.Message.Send(ws, bytes)
	}
}
