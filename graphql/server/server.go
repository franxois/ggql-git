package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	graphql "github.com/franxois/ggql-git/graphql"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler.Playground("GraphQL playground", "/query"))
	mux.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}})))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:8080", "http://localhost:8081"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: !true,
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(mux)))
}
