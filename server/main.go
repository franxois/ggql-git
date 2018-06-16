package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/franxois/ggql-git/server/graph"
	"github.com/vektah/gqlgen/handler"
)

func main() {
	app := &graph.MyApp{}
	http.Handle("/", handler.Playground("Projects", "/graphql"))
	http.Handle("/graphql", handler.GraphQL(graph.MakeExecutableSchema(app)))

	fmt.Println("Listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
