package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/harrisonmalone/authors-books-app/graph"
	"github.com/harrisonmalone/authors-books-app/graph/database"
	"github.com/harrisonmalone/authors-books-app/graph/generated"
	"github.com/harrisonmalone/authors-books-app/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	if os.Getenv("AUTHORS_BOOKS_ENV") != "prod" {
		err := godotenv.Load()

		if err != nil {
			panic("ENV variables cannot be loaded")
		}
	}

	router := chi.NewRouter()

	
	
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	
	
	
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	router.Use(corsOptions.Handler)

	router.Use(middleware.Authentication())

	database.Initialize()
	
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)

		if oc.OperationName != "IntrospectionQuery" {
			fmt.Printf("%s \n", oc.OperationName)
		}
		return next(ctx)
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	err := http.ListenAndServe("localhost:"+port, router)

	if err != nil {
		panic(err)
	}

}
