package main

import (
	"api-shapes/pkg/router"
	trans "api-shapes/transport"
	"api-shapes/transport/rest"
	"api-shapes/transport/soap"
	"log"
	"net/http"

	"api-shapes/transport/graphql"
	"api-shapes/transport/graphql/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	graphqltrans "github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	mux := http.NewServeMux()

	restUserAPI := rest.NewUserAPI()
	trans.NewRouter("/rest/users", mux, restUserAPI)

	soapUserAPI := soap.NewUserAPI()
	trans.NewRouter("/soap/users", mux, soapUserAPI)

	es := graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(graphql.NewResolver()),
	})
	srv := handler.New(es)
	srv.AddTransport(graphqltrans.POST{})

	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	mux.Handle("/graphql", srv)

	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router.ErrorHandler(mux)))
}
