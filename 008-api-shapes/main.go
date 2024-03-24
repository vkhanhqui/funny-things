package main

import (
	"api-shapes/pkg/router"
	"api-shapes/transport"
	"api-shapes/transport/rest"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	restAPI := rest.NewUserAPI()
	transport.NewRouter("/v2/users", mux, restAPI)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router.ErrorHandler(mux)))
}
