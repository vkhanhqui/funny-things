package main

import (
	"cloud-gaming/service"
	"cloud-gaming/transport"
	"fmt"
	"log"
	"net/http"
)

func main() {
	wk := service.NewWorker()
	wk.Run()

	svc := service.NewService()
	handler := transport.NewHandler(*svc)

	http.Handle("/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/ws", handler.Websocket)

	port := 8080
	log.Printf("starting server on port %v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
