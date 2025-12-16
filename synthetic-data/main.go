package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting Synthetic Data Generation...")
	fmt.Println("Generating League Data...")
	leagueData := generateLeagueFlat()
	fmt.Println("Creating Sample Team Roster...")
	rosters := make(map[string]FootballTeamRoster)
	for _, team := range leagueData.Teams {
		fmt.Printf("\nCreating Team Roster for Team: %s %s | %s\n", team.City, team.Name, team.ID)
		rosters[team.ID] = createTeamRoster(team.ID)
	}

	for _, roster := range rosters {
		fmt.Println("\nGenerating Player Yearly Stats...")
		for _, player := range roster.QB {
			createPlayerYear(player, 2025)
		}
		for _, player := range roster.RB {
			createPlayerYear(player, 2025)
		}
		for _, player := range roster.WR {
			createPlayerYear(player, 2025)
		}
		for _, player := range roster.TE {
			createPlayerYear(player, 2025)
		}
		for _, player := range roster.PK {
			createPlayerYear(player, 2025)
		}
	}
	fmt.Println("\nSynthetic Data Generation Complete.")
}
