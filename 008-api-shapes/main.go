package main

import (
	"api-shapes/pkg/router"
	"api-shapes/transport"
	"api-shapes/transport/rest"
	"api-shapes/transport/soap"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	soapUserAPI := soap.NewUserAPI()
	transport.NewRouter("/v1/users", mux, soapUserAPI)

	restUserAPI := rest.NewUserAPI()
	transport.NewRouter("/v2/users", mux, restUserAPI)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router.ErrorHandler(mux)))
}
