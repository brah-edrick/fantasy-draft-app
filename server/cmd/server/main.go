package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"fantasy-draft/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Get DB URL from environment (set in docker-compose)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// 2. Create a connection pool (thread-safe for concurrent GraphQL resolvers)
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	fmt.Println("✅ Successfully connected to Postgres with connection pool!")

	// 3. Create the GraphQL resolver with connection pool
	resolver := graph.NewResolver(pool)

	// 4. Create the GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// 5. Register Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the Fantasy Draft API! Visit /playground for the GraphQL Playground.")
	})

	// GraphQL Playground - interactive query interface
	http.Handle("/playground", playground.Handler("Fantasy Draft GraphQL", "/graphql"))

	// GraphQL endpoint
	http.Handle("/graphql", srv)

	// 6. Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s...\n", port)
	fmt.Printf("📊 GraphQL Playground: http://localhost:%s/playground\n", port)
	fmt.Printf("🔗 GraphQL Endpoint: http://localhost:%s/graphql\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
