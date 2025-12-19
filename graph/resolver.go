package graph

import "github.com/jackc/pgx/v5/pgxpool"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

// Resolver is the root resolver that holds dependencies
type Resolver struct {
	// DB is a connection pool for database queries (thread-safe for concurrent resolvers)
	DB *pgxpool.Pool
}

// NewResolver creates a new resolver with all dependencies
func NewResolver(db *pgxpool.Pool) *Resolver {
	return &Resolver{
		DB: db,
	}
}
