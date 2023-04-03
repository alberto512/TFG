package main

import (
	"log"
	"net/http"
	"os"
	"tfg/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/rs/cors"

	"tfg/internal/middleware"
	"tfg/internal/mongo"

	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	// Only for development
	/*
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	*/

	// Start mongo connection
	mongo.Start()
	defer mongo.Close()

	// Get port from env variable or default port
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Start router, add middleware, add default route and start server
	router := chi.NewRouter()

	router.Use(cors.Default().Handler)

	router.Use(middleware.Middleware())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", srv)

	log.Printf("Server running in %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
