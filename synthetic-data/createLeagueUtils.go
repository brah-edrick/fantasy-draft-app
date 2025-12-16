package main

import (
	"math/rand"

	"github.com/google/uuid"
)

// Franchise represents a team identity
type Franchise struct {
	City  string
	State string
	Name  string
	Abbr  string
}

// A curated list of 32 synthetic teams
var allAvailableFranchises = []Franchise{
	{"Austin", "TX", "Desperados", "AUS"},
	{"Portland", "OR", "Lumberjacks", "POR"},
	{"Salt Lake", "UT", "Peaks", "SLC"},
	{"Orlando", "FL", "Orbit", "ORL"},
	{"San Diego", "CA", "Destroyers", "SD"},
	{"Columbus", "OH", "Aviators", "COL"},
	{"Sacramento", "CA", "Miners", "SAC"},
	{"San Antonio", "TX", "Marshals", "SA"},
	{"Memphis", "TN", "Pharaohs", "MEM"},
	{"Oklahoma City", "OK", "Twisters", "OKC"},
	{"Las Vegas", "NV", "High Rollers", "LV"},
	{"Raleigh", "NC", "Capitals", "RAL"},
	{"Birmingham", "AL", "Vulcans", "BHM"},
	{"Louisville", "KY", "Jockeys", "LOU"},
	{"Virginia Beach", "VA", "Neptunes", "VB"},
	{"Omaha", "NE", "Mammoths", "OMA"},
	// East Coast / Metro
	{"Brooklyn", "NY", "Barons", "BKN"},
	{"Boston", "MA", "Colonials", "BOS"},
	{"Philadelphia", "PA", "Liberty", "PHI"},
	{"Washington", "DC", "Sentinels", "DC"},
	// Midwest / North
	{"Chicago", "IL", "Wind", "CHI"},
	{"Detroit", "MI", "Gears", "DET"},
	{"Milwaukee", "WI", "Hunters", "MIL"},
	{"Minneapolis", "MN", "Blizzard", "MIN"},
	// South
	{"Atlanta", "GA", "Phoenixes", "ATL"},
	{"Miami", "FL", "Sharks", "MIA"},
	{"New Orleans", "LA", "Deltas", "NO"},
	{"Nashville", "TN", "Strings", "NSH"},
	// West
	{"Seattle", "WA", "Emeralds", "SEA"},
	{"San Francisco", "CA", "Fog", "SF"},
	{"Los Angeles", "CA", "Stars", "LA"},
	{"Denver", "CO", "Summits", "DEN"},
}

func generateConference(name string) Conference {
	return Conference{
		ID:   uuid.NewString(),
		Name: name,
	}
}

func generateDivision(name string, conferenceID string) Division {
	return Division{
		ID:           uuid.NewString(),
		Name:         name,
		ConferenceID: conferenceID,
	}
}

func generateTeam(franchise Franchise, divisionID string) Team {
	return Team{
		ID:         uuid.NewString(),
		Name:       franchise.Name,
		DivisionID: divisionID,
		City:       franchise.City,
		State:      franchise.State,
		Abbr:       franchise.Abbr,
	}
}

type LeagueFlat struct {
	Conferences []Conference `json:"conferences"`
	Divisions   []Division   `json:"divisions"`
	Teams       []Team       `json:"teams"`
}

func generateLeagueFlat() LeagueFlat {
	returnValue := LeagueFlat{}

	// Generate Conferences
	confNames := []string{"Union Conference", "Alliance Conference"}
	generatedConferences := make([]Conference, len(confNames))
	for confIndex, confName := range confNames {
		generatedConferences[confIndex] = generateConference(confName)
	}
	returnValue.Conferences = generatedConferences

	// Generate Divisions
	divisionNames := []string{"North", "South", "East", "West"}
	generatedDivisions := make([]Division, len(divisionNames)*len(generatedConferences))
	for confIndex, generatedConference := range generatedConferences {
		for divIndex, divName := range divisionNames {
			generatedDivisions[confIndex*len(divisionNames)+divIndex] = generateDivision(divName, generatedConference.ID)
		}
	}
	returnValue.Divisions = generatedDivisions

	// Generate Teams
	availableFranchises := allAvailableFranchises
	generatedTeams := make([]Team, len(availableFranchises))
	for divisionIndex, generatedDivision := range generatedDivisions {
		// each division has 4 teams
		divisionSize := 4
		for teamIndex := range divisionSize {
			randomIndex := rand.Intn(len(availableFranchises))
			generatedTeams[divisionIndex*divisionSize+teamIndex] = generateTeam(availableFranchises[randomIndex], generatedDivision.ID)
			// remove the franchise from the list
			availableFranchises = append(availableFranchises[:randomIndex], availableFranchises[randomIndex+1:]...)
		}

	}
	returnValue.Teams = generatedTeams

	return returnValue
}
