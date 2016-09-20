package main

import (
	"log"
	"net/http"
	"superstellar/backend"
	"math/rand"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	rand.Seed(time.Now().UTC().UnixNano())
	// websocket server
	server := backend.NewServer("/superstellar")
	go server.Listen()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
