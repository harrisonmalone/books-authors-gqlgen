package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/harrisonmalone/authors-books-app/graph"
	"github.com/harrisonmalone/authors-books-app/graph/database"
	"github.com/harrisonmalone/authors-books-app/graph/generated"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	if os.Getenv("AUTHORS_BOOKS_ENV") != "prod" {
		err := godotenv.Load()

		if err != nil {
			panic("ENV variables cannot be loaded")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.Initialize()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
