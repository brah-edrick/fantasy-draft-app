package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("Starting Synthetic Data Generation...")
	fmt.Println("Generating League Data...")
	uuidGenerator := UUIDGenerator(func() string { return uuid.New().String() })
	clock := RealClock{}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	leagueData := generateLeagueFlat(uuidGenerator, clock, rng)
	fmt.Println("Creating Sample Team Roster...")
	rosters := make(map[string]FootballTeamRoster)
	for _, team := range leagueData.Teams {
		fmt.Printf("\nCreating Team Roster for Team: %s %s | %s\n", team.City, team.Name, team.ID)
		rosters[team.ID] = createTeamRoster(team.ID)
	}

	for _, roster := range rosters {
		for _, player := range roster.QB {
			createPlayerCareer(player)
		}
		for _, player := range roster.RB {
			createPlayerCareer(player)
		}
		for _, player := range roster.WR {
			createPlayerCareer(player)
		}
		for _, player := range roster.TE {
			createPlayerCareer(player)
		}
		for _, player := range roster.PK {
			createPlayerCareer(player)
		}
	}
	fmt.Println("\nSynthetic Data Generation Complete.")
}
