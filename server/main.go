package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/franxois/ggql-git/server/graph"
	"github.com/rs/cors"
	"github.com/vektah/gqlgen/handler"
)

func main() {
	app := &graph.MyApp{}

	mux := http.NewServeMux()
	mux.Handle("/", handler.Playground("Projects", "/graphql"))
	mux.Handle("/graphql", handler.GraphQL(graph.MakeExecutableSchema(app)))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:8080", "http://localhost:8081"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	fmt.Println("Listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", c.Handler(mux)))
}
