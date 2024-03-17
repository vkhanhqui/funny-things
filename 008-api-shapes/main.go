package main

import (
	"api-shapes/pkg/router"
	trans "api-shapes/transport"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	trans.NewRouter("/users", mux)
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router.ErrorHandler(mux)))
}
