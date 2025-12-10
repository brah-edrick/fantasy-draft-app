package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting Synthetic Data Generation...")
	fmt.Println("Generating League Data...")
	leagueData := GenerateLeagueFlat()
	fmt.Println("Creating Sample Team Roster...")
	for _, team := range leagueData.Teams {
		fmt.Printf("\nCreating Team Roster for Team: %s %s | %s\n", team.City, team.Name, team.ID)
		createTeamRoster(team.ID)
	}
	fmt.Println("\nSynthetic Data Generation Complete.")
}
