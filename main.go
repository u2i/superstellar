package main

import (
	"log"
	"net/http"
	"superstellar/backend"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := backend.NewServer("/superstellar")
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
