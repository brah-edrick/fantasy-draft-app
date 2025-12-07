package main

import (
	"math/rand"

	"github.com/google/uuid"
)

// Franchise represents a pre-defined cool team identity
type Franchise struct {
	City   string
	State  string
	Name   string
	Abbr   string
}

// A curated list of 32 non-sucky synthetic teams
var availableFranchises = []Franchise{
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
	{"Atlanta", "GA", "Phoenix", "ATL"},
	{"Miami", "FL", "Sharks", "MIA"},
	{"New Orleans", "LA", "Deltas", "NO"},
	{"Nashville", "TN", "Strings", "NSH"},
	// West
	{"Seattle", "WA", "Emeralds", "SEA"},
	{"San Francisco", "CA", "Fog", "SF"},
	{"Los Angeles", "CA", "Stars", "LA"},
	{"Denver", "CO", "Summits", "DEN"},
}


func GenerateConference(name string) Conference {
	return Conference{
		ID:   uuid.NewString(),
		Name: name,
	}
}

func GenerateDivision(name string, conferenceID string) Division {
	return Division{
		ID:   uuid.NewString(),
		Name: name,
		ConferenceID: conferenceID,
	}
}

func GenerateTeam(franchise Franchise, divisionID string) Team {
	return Team{
		ID:   uuid.NewString(),
		Name: franchise.Name,
		DivisionID: divisionID,
		City: franchise.City,
		State: franchise.State,
		Abbr: franchise.Abbr,
	}
}

type LeagueFlat struct {
	Conferences []Conference `json:"conferences"`
	Divisions   []Division   `json:"divisions"`
	Teams       []Team       `json:"teams"`
}

func GenerateLeagueFlat() LeagueFlat {
	returnValue := LeagueFlat{}

	// Generate Conferences
	confNames := []string{"Union Conference", "Alliance Conference"};
	generatedConferences := make([]Conference, len(confNames));
	for confIndex, confName := range confNames {
		generatedConferences[confIndex] = GenerateConference(confName)
	}
	returnValue.Conferences = generatedConferences;

	// Generate Divisions
	divisionNames := []string{"North", "South", "East", "West"}
	generatedDivisions := make([]Division, len(divisionNames)*len(generatedConferences));
	for confIndex, generatedConference := range generatedConferences {
		for divIndex, divName := range divisionNames {
			generatedDivisions[confIndex*len(divisionNames)+divIndex] = GenerateDivision(divName, generatedConference.ID)
		}
	}
	returnValue.Divisions = generatedDivisions;


	// Generate Teams
	copyOfAvailableFranchises := availableFranchises
	generatedTeams := make([]Team, len(copyOfAvailableFranchises));
	for divisionIndex, generatedDivision := range generatedDivisions {
		// each division has 4 teams
		divisionSize := 4
		for teamIndex := 0; teamIndex < divisionSize; teamIndex++ {
			randomIndex := rand.Intn(len(copyOfAvailableFranchises))
			generatedTeams[divisionIndex*divisionSize+teamIndex] = GenerateTeam(copyOfAvailableFranchises[randomIndex], generatedDivision.ID)
			// remove the franchise from the list
			copyOfAvailableFranchises = append(copyOfAvailableFranchises[:randomIndex], copyOfAvailableFranchises[randomIndex+1:]...)
		}

	}
	returnValue.Teams = generatedTeams;

	return returnValue
}
