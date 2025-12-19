package main

import (
	"fmt"
	"os"
)

func main() {
	// Check for command-line arguments
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		RunSeed()
		return
	}

	// Show usage if no valid command
	fmt.Println("Usage: go run . <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  seed    Purge and seed the database with synthetic data")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("  REAL_DATA_FILE=/path/to/real-data.json \\")
	fmt.Println("  DATABASE_URL=\"postgres://user:pass@localhost:5432/db\" \\")
	fmt.Println("  go run . seed")
}
