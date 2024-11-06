package main

import (
	"cloud-gaming/transport"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/ws", transport.WebsocketHandler)

	port := 8080
	log.Printf("starting server on port %v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
