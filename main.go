package main

import (
	"log"
	"math/rand"
	"net/http"
	"superstellar/backend/server"
	"time"
)

import _ "net/http/pprof"

func main() {
	log.SetFlags(log.Lshortfile)

	rand.Seed(time.Now().UTC().UnixNano())
	// websocket server
	server := server.NewServer("/superstellar")
	go server.Listen()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
