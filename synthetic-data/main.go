package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func generateLeaugeData() {
	log.Println("Generating League Structure...")
	league := GenerateLeagueFlat()

	// Export to JSON
	outputDir := "synthetic-data/.output"
	// Ensure dir exists
	if _, err := os.Stat("synthetic-data"); err == nil {
		// We are running from root
	} else {
		// running from inside synthetic-data
		outputDir = ".output"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	outputPath := fmt.Sprintf("%s/league.json", outputDir)
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(league); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	log.Printf("SUCCESS: League Data written to %s", outputPath)
}

func exportPlayerStats() {
	log.Println("Exporting Player Stats...")
	playerStats := collectAndAggregatePlayerStats()

	// Export to JSON
	outputDir := "synthetic-data/.output"
	// Ensure dir exists
	if _, err := os.Stat("synthetic-data"); err == nil {
		// We are running from root
	} else {
		// running from inside synthetic-data
		outputDir = ".output"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	outputPath := fmt.Sprintf("%s/player_stats.json", outputDir)
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(playerStats); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	log.Printf("SUCCESS: Player Stats written to %s", outputPath)
}

func main() {
	log.Println("Starting Synthetic Data Generation...")

	generateLeaugeData()
	exportPlayerStats()
	createNewPlayer("PK", "1")
	createNewPlayer("QB", "1")
}
