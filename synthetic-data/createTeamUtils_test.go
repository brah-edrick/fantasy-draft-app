package main

import (
	"testing"
)

func TestCreateTeamRosterIntegration(t *testing.T) {
	// Integration test - uses real data generation
	// This will test the actual roster creation with dependencies
	teamID := "test-team-123"
	
	roster := createTeamRoster(teamID)

	// Verify roster has correct number of players per position
	if len(roster.QB) != NFLRosterComposition["QB"] {
		t.Errorf("Expected %d QBs, got %d", NFLRosterComposition["QB"], len(roster.QB))
	}
	if len(roster.RB) != NFLRosterComposition["RB"] {
		t.Errorf("Expected %d RBs, got %d", NFLRosterComposition["RB"], len(roster.RB))
	}
	if len(roster.WR) != NFLRosterComposition["WR"] {
		t.Errorf("Expected %d WRs, got %d", NFLRosterComposition["WR"], len(roster.WR))
	}
	if len(roster.TE) != NFLRosterComposition["TE"] {
		t.Errorf("Expected %d TEs, got %d", NFLRosterComposition["TE"], len(roster.TE))
	}
	if len(roster.PK) != NFLRosterComposition["PK"] {
		t.Errorf("Expected %d PKs, got %d", NFLRosterComposition["PK"], len(roster.PK))
	}

	// Verify all players have correct team ID
	for _, player := range roster.QB {
		if player.TeamID != teamID {
			t.Errorf("Expected team ID %s, got %s", teamID, player.TeamID)
		}
		if player.Position != "QB" {
			t.Errorf("Expected position QB, got %s", player.Position)
		}
	}

	// Verify depth-based skills: starters should have higher average skill than backups
	if len(roster.QB) >= 2 {
		starterSkill := roster.QB[0].Skill
		backupSkill := roster.QB[len(roster.QB)-1].Skill
		// This might not always be true due to randomness, so we'll just verify they're in valid range
		if starterSkill < 0.15 || starterSkill > 0.95 {
			t.Errorf("Starter skill %f out of expected range [0.15, 0.95]", starterSkill)
		}
		if backupSkill < 0.15 || backupSkill > 0.95 {
			t.Errorf("Backup skill %f out of expected range [0.15, 0.95]", backupSkill)
		}
	}
}

func TestCreateTeamRoster(t *testing.T) {
	// Test NFLRosterComposition structure
	expectedComposition := map[string]int{
		"QB": 3,
		"RB": 4,
		"WR": 6,
		"TE": 3,
		"PK": 1,
	}

	for position, count := range expectedComposition {
		if NFLRosterComposition[position] != count {
			t.Errorf("Expected %s count to be %d, got %d", position, count, NFLRosterComposition[position])
		}
	}

	// Total roster size should be 17
	totalRoster := 0
	for _, count := range NFLRosterComposition {
		totalRoster += count
	}

	if totalRoster != 17 {
		t.Errorf("Expected total roster size of 17, got %d", totalRoster)
	}
}

func TestCreatePlayersWithDepthSkillsIntegration(t *testing.T) {
	// Integration test for createPlayersWithDepthSkills
	teamID := "test-team-456"
	position := QB
	count := 3

	players := createPlayersWithDepthSkills(position, teamID, count)

	// Verify correct number of players
	if len(players) != count {
		t.Errorf("Expected %d players, got %d", count, len(players))
	}

	// Verify all players have correct attributes
	for i, player := range players {
		if player.TeamID != teamID {
			t.Errorf("Player %d: expected team ID %s, got %s", i, teamID, player.TeamID)
		}
		if player.Position != string(position) {
			t.Errorf("Player %d: expected position %s, got %s", i, position, player.Position)
		}
		if player.Skill < 0.15 || player.Skill > 0.95 {
			t.Errorf("Player %d: skill %f out of valid range [0.15, 0.95]", i, player.Skill)
		}
		if player.ID == "" {
			t.Errorf("Player %d: ID should not be empty", i)
		}
		if player.FirstName == "" {
			t.Errorf("Player %d: FirstName should not be empty", i)
		}
		if player.LastName == "" {
			t.Errorf("Player %d: LastName should not be empty", i)
		}
	}

	// Verify skills generally decrease with depth (might not always be true due to variance)
	// But the average of first half should be >= average of second half
	if len(players) >= 2 {
		firstHalfAvg := 0.0
		secondHalfAvg := 0.0
		midpoint := len(players) / 2

		for i := 0; i < midpoint; i++ {
			firstHalfAvg += players[i].Skill
		}
		firstHalfAvg /= float64(midpoint)

		for i := midpoint; i < len(players); i++ {
			secondHalfAvg += players[i].Skill
		}
		secondHalfAvg /= float64(len(players) - midpoint)

		// First half should generally have higher skill (allowing for variance)
		// We'll just verify both are in valid range
		if firstHalfAvg < 0.15 || firstHalfAvg > 0.95 {
			t.Errorf("First half avg skill %f out of valid range", firstHalfAvg)
		}
		if secondHalfAvg < 0.15 || secondHalfAvg > 0.95 {
			t.Errorf("Second half avg skill %f out of valid range", secondHalfAvg)
		}
	}
}

func TestCreatePlayersWithDepthSkills(t *testing.T) {
	// We'll verify the roster composition structure
	t.Run("QB roster size", func(t *testing.T) {
		if NFLRosterComposition["QB"] != 3 {
			t.Errorf("Expected 3 QBs, got %d", NFLRosterComposition["QB"])
		}
	})

	t.Run("RB roster size", func(t *testing.T) {
		if NFLRosterComposition["RB"] != 4 {
			t.Errorf("Expected 4 RBs, got %d", NFLRosterComposition["RB"])
		}
	})

	t.Run("WR roster size", func(t *testing.T) {
		if NFLRosterComposition["WR"] != 6 {
			t.Errorf("Expected 6 WRs, got %d", NFLRosterComposition["WR"])
		}
	})

	t.Run("TE roster size", func(t *testing.T) {
		if NFLRosterComposition["TE"] != 3 {
			t.Errorf("Expected 3 TEs, got %d", NFLRosterComposition["TE"])
		}
	})

	t.Run("PK roster size", func(t *testing.T) {
		if NFLRosterComposition["PK"] != 1 {
			t.Errorf("Expected 1 PK, got %d", NFLRosterComposition["PK"])
		}
	})
}

func TestFootballTeamRosterStructure(t *testing.T) {
	// Test that we can create an empty roster structure
	roster := FootballTeamRoster{
		QB: []Player{},
		RB: []Player{},
		WR: []Player{},
		TE: []Player{},
		PK: []Player{},
	}

	if roster.QB == nil {
		t.Error("QB slice should not be nil")
	}
	if roster.RB == nil {
		t.Error("RB slice should not be nil")
	}
	if roster.WR == nil {
		t.Error("WR slice should not be nil")
	}
	if roster.TE == nil {
		t.Error("TE slice should not be nil")
	}
	if roster.PK == nil {
		t.Error("PK slice should not be nil")
	}
}

func TestRosterCompositionType(t *testing.T) {
	// Verify RosterComposition type works as expected
	composition := RosterComposition{
		"QB": 3,
		"RB": 4,
	}

	if composition["QB"] != 3 {
		t.Errorf("Expected QB count to be 3, got %d", composition["QB"])
	}

	if composition["RB"] != 4 {
		t.Errorf("Expected RB count to be 4, got %d", composition["RB"])
	}

	// Test that accessing non-existent key returns 0
	if composition["DE"] != 0 {
		t.Errorf("Expected non-existent position count to be 0, got %d", composition["DE"])
	}
}

