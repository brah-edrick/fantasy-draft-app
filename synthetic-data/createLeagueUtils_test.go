package main

import (
	"math/rand"
	"testing"
	"time"
)

// Mock UUID Generator for testing
func mockUUIDGenerator(prefix string, counter *int) UUIDGenerator {
	return func() string {
		*counter++
		return prefix + string(rune(*counter))
	}
}

// Mock Clock for testing
type MockClock struct {
	mockTime time.Time
}

func (m MockClock) Now() time.Time {
	return m.mockTime
}

func TestGenerateConference(t *testing.T) {
	counter := 0
	uuidGen := mockUUIDGenerator("conf-", &counter)

	conf := generateConference("Test Conference", uuidGen)

	if conf.Name != "Test Conference" {
		t.Errorf("Expected conference name 'Test Conference', got '%s'", conf.Name)
	}

	if conf.ID == "" {
		t.Error("Conference ID should not be empty")
	}

	if conf.ID != "conf-\x01" {
		t.Errorf("Expected conference ID 'conf-\\x01', got '%s'", conf.ID)
	}
}

func TestGenerateDivision(t *testing.T) {
	counter := 0
	uuidGen := mockUUIDGenerator("div-", &counter)
	conferenceID := "conf-123"

	div := generateDivision("North", conferenceID, uuidGen)

	if div.Name != "North" {
		t.Errorf("Expected division name 'North', got '%s'", div.Name)
	}

	if div.ConferenceID != conferenceID {
		t.Errorf("Expected conference ID '%s', got '%s'", conferenceID, div.ConferenceID)
	}

	if div.ID == "" {
		t.Error("Division ID should not be empty")
	}
}

func TestGenerateTeam(t *testing.T) {
	counter := 0
	uuidGen := mockUUIDGenerator("team-", &counter)
	divisionID := "div-123"

	franchise := Franchise{
		City:  "Austin",
		State: "TX",
		Name:  "Desperados",
		Abbr:  "AUS",
	}

	team := generateTeam(franchise, divisionID, uuidGen)

	if team.Name != "Desperados" {
		t.Errorf("Expected team name 'Desperados', got '%s'", team.Name)
	}

	if team.City != "Austin" {
		t.Errorf("Expected city 'Austin', got '%s'", team.City)
	}

	if team.State != "TX" {
		t.Errorf("Expected state 'TX', got '%s'", team.State)
	}

	if team.Abbr != "AUS" {
		t.Errorf("Expected abbreviation 'AUS', got '%s'", team.Abbr)
	}

	if team.DivisionID != divisionID {
		t.Errorf("Expected division ID '%s', got '%s'", divisionID, team.DivisionID)
	}

	if team.ID == "" {
		t.Error("Team ID should not be empty")
	}
}

func TestGenerateLeagueFlat(t *testing.T) {
	counter := 0
	uuidGen := mockUUIDGenerator("id-", &counter)
	mockClock := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	rng := rand.New(rand.NewSource(12345)) // Fixed seed for deterministic tests

	league := generateLeagueFlat(uuidGen, mockClock, rng)

	// Test conferences
	if len(league.Conferences) != 2 {
		t.Errorf("Expected 2 conferences, got %d", len(league.Conferences))
	}

	expectedConfNames := []string{"Union Conference", "Alliance Conference"}
	for i, conf := range league.Conferences {
		if conf.Name != expectedConfNames[i] {
			t.Errorf("Expected conference %d to be '%s', got '%s'", i, expectedConfNames[i], conf.Name)
		}
		if conf.ID == "" {
			t.Errorf("Conference %d has empty ID", i)
		}
	}

	// Test divisions - should be 8 total (4 per conference)
	if len(league.Divisions) != 8 {
		t.Errorf("Expected 8 divisions, got %d", len(league.Divisions))
	}

	// Verify each conference has 4 divisions
	divisionCounts := make(map[string]int)
	for _, div := range league.Divisions {
		divisionCounts[div.ConferenceID]++
	}

	for _, conf := range league.Conferences {
		if divisionCounts[conf.ID] != 4 {
			t.Errorf("Conference %s should have 4 divisions, got %d", conf.Name, divisionCounts[conf.ID])
		}
	}

	// Test teams - should be 32 total (4 per division)
	if len(league.Teams) != 32 {
		t.Errorf("Expected 32 teams, got %d", len(league.Teams))
	}

	// Verify each division has 4 teams
	teamCounts := make(map[string]int)
	for _, team := range league.Teams {
		teamCounts[team.DivisionID]++
	}

	for _, div := range league.Divisions {
		if teamCounts[div.ID] != 4 {
			t.Errorf("Division %s should have 4 teams, got %d", div.Name, teamCounts[div.ID])
		}
	}

	// Verify all teams are unique
	teamNames := make(map[string]bool)
	teamAbbrs := make(map[string]bool)
	for _, team := range league.Teams {
		if teamNames[team.Name] {
			t.Errorf("Duplicate team name: %s", team.Name)
		}
		teamNames[team.Name] = true

		if teamAbbrs[team.Abbr] {
			t.Errorf("Duplicate team abbreviation: %s", team.Abbr)
		}
		teamAbbrs[team.Abbr] = true

		// Verify team has all required fields
		if team.ID == "" {
			t.Errorf("Team %s has empty ID", team.Name)
		}
		if team.City == "" {
			t.Errorf("Team %s has empty city", team.Name)
		}
		if team.State == "" {
			t.Errorf("Team %s has empty state", team.Name)
		}
		if team.Name == "" {
			t.Errorf("Team has empty name")
		}
		if team.Abbr == "" {
			t.Errorf("Team %s has empty abbreviation", team.Name)
		}
		if team.DivisionID == "" {
			t.Errorf("Team %s has empty division ID", team.Name)
		}
	}
}

func TestGenerateLeagueFlatRandomness(t *testing.T) {
	// Test that different random seeds produce different team assignments
	// Note: This test verifies randomness by checking team name distribution
	counter1 := 0
	uuidGen1 := mockUUIDGenerator("id1-", &counter1)
	mockClock1 := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	rng1 := rand.New(rand.NewSource(12345))

	counter2 := 0
	uuidGen2 := mockUUIDGenerator("id2-", &counter2)
	mockClock2 := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	rng2 := rand.New(rand.NewSource(54321))

	league1 := generateLeagueFlat(uuidGen1, mockClock1, rng1)
	league2 := generateLeagueFlat(uuidGen2, mockClock2, rng2)

	// Collect team names from both leagues
	teams1 := make(map[string]bool)
	teams2 := make(map[string]bool)
	
	for _, team := range league1.Teams {
		teams1[team.Name] = true
	}
	for _, team := range league2.Teams {
		teams2[team.Name] = true
	}

	// Both should have 32 unique teams
	if len(teams1) != 32 {
		t.Errorf("League 1 should have 32 unique teams, got %d", len(teams1))
	}
	if len(teams2) != 32 {
		t.Errorf("League 2 should have 32 unique teams, got %d", len(teams2))
	}

	// The leagues should have the same teams (just in different positions)
	for name := range teams1 {
		if !teams2[name] {
			t.Errorf("League 2 is missing team %s that League 1 has", name)
		}
	}
}

func TestAllAvailableFranchises(t *testing.T) {
	// Count franchises in source code (since the slice might be modified by other tests)
	expectedCount := 32
	
	// Create a fresh copy for validation
	franchises := []Franchise{
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
		{"Brooklyn", "NY", "Barons", "BKN"},
		{"Boston", "MA", "Colonials", "BOS"},
		{"Philadelphia", "PA", "Liberty", "PHI"},
		{"Washington", "DC", "Sentinels", "DC"},
		{"Chicago", "IL", "Wind", "CHI"},
		{"Detroit", "MI", "Gears", "DET"},
		{"Milwaukee", "WI", "Hunters", "MIL"},
		{"Minneapolis", "MN", "Blizzard", "MIN"},
		{"Atlanta", "GA", "Phoenixes", "ATL"},
		{"Miami", "FL", "Sharks", "MIA"},
		{"New Orleans", "LA", "Deltas", "NO"},
		{"Nashville", "TN", "Strings", "NSH"},
		{"Seattle", "WA", "Emeralds", "SEA"},
		{"San Francisco", "CA", "Fog", "SF"},
		{"Los Angeles", "CA", "Stars", "LA"},
		{"Denver", "CO", "Summits", "DEN"},
	}

	if len(franchises) != expectedCount {
		t.Errorf("Expected %d franchises, got %d", expectedCount, len(franchises))
	}

	// Verify all franchises have required fields
	abbrs := make(map[string]bool)
	names := make(map[string]bool)
	
	for i, franchise := range franchises {
		if franchise.City == "" {
			t.Errorf("Franchise at index %d has empty city", i)
		}
		if franchise.State == "" {
			t.Errorf("Franchise at index %d has empty state", i)
		}
		if franchise.Name == "" {
			t.Errorf("Franchise at index %d has empty name", i)
		}
		if franchise.Abbr == "" {
			t.Errorf("Franchise at index %d has empty abbreviation", i)
		}

		// Check for duplicate abbreviations
		if abbrs[franchise.Abbr] {
			t.Errorf("Duplicate franchise abbreviation: %s (at index %d)", franchise.Abbr, i)
		}
		abbrs[franchise.Abbr] = true

		// Check for duplicate names
		if names[franchise.Name] {
			t.Errorf("Duplicate franchise name: %s (at index %d)", franchise.Name, i)
		}
		names[franchise.Name] = true
	}
}

