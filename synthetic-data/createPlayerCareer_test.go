package main

import (
	"testing"
	"time"
)

func TestRollForInjury(t *testing.T) {
	tests := []struct {
		name           string
		playerAge      int
		playerPosition string
		iterations     int
	}{
		{"young QB", 23, "QB", 1000},
		{"old QB", 36, "QB", 1000},
		{"young RB", 22, "RB", 1000},
		{"old RB", 34, "RB", 1000},
		{"middle age WR", 27, "WR", 1000},
		{"PK", 30, "PK", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injuredCount := 0
			totalGames := 0

			for range tt.iterations {
				injured, games := rollForInjury(tt.playerAge, tt.playerPosition)
				if injured {
					injuredCount++
					totalGames += games

					// Games missed should be between 1 and 20
					if games < 1 || games > 20 {
						t.Errorf("Games missed %d is out of range [1, 20]", games)
					}
				} else {
					// If not injured, games should be 0
					if games != 0 {
						t.Errorf("Expected 0 games missed when not injured, got %d", games)
					}
				}
			}

			injuryRate := float64(injuredCount) / float64(tt.iterations)

			// Injury rate should be reasonable (between 0% and 20%)
			if injuryRate < 0.0 || injuryRate > 0.20 {
				t.Errorf("Injury rate %f is outside reasonable range [0.0, 0.20]", injuryRate)
			}

			// Older players should generally have higher injury rates than younger players
			// RBs should have higher injury rates than QBs and PKs
			t.Logf("%s: Injury rate: %.2f%%, Avg games missed when injured: %.1f",
				tt.name, injuryRate*100, float64(totalGames)/float64(max(injuredCount, 1)))
		})
	}
}

func TestRollForInjuryPositionRates(t *testing.T) {
	// Test that position affects injury rates correctly
	iterations := 10000
	age := 28 // Middle age for consistent comparison

	positions := []string{"QB", "RB", "WR", "TE", "PK"}
	rates := make(map[string]float64)

	for _, pos := range positions {
		injuredCount := 0
		for range iterations {
			injured, _ := rollForInjury(age, pos)
			if injured {
				injuredCount++
			}
		}
		rates[pos] = float64(injuredCount) / float64(iterations)
	}

	// QB and PK should have lower injury rates than RB
	if rates["QB"] >= rates["RB"] {
		t.Errorf("QB injury rate (%f) should be lower than RB (%f)", rates["QB"], rates["RB"])
	}

	if rates["PK"] >= rates["RB"] {
		t.Errorf("PK injury rate (%f) should be lower than RB (%f)", rates["PK"], rates["RB"])
	}

	// PK should have the lowest injury rate
	if rates["PK"] >= rates["QB"] {
		t.Errorf("PK injury rate (%f) should be lowest, but is higher than QB (%f)", rates["PK"], rates["QB"])
	}
}

func TestNewCareerSimulator(t *testing.T) {
	t.Run("with all configs provided", func(t *testing.T) {
		mockClock := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
		mockInjuryRoller := func(age int, position string) (bool, int) {
			return false, 0
		}
		mockStatsGenerator := func(player Player, yoe int) FootballStats {
			return FootballStats{PassingYards: 100}
		}
		mockStatMultiplier := func(player Player, yoe int, stats FootballStats) FootballStats {
			return stats
		}

		cfg := YearSimulatorConfig{
			Clock:          mockClock,
			GamesPerSeason: 16,
			InjuryRoller:   mockInjuryRoller,
			StatsGenerator: mockStatsGenerator,
			StatMultiplier: mockStatMultiplier,
		}

		sim := NewCareerSimulator(cfg)

		if sim.clock != mockClock {
			t.Error("Clock was not set correctly")
		}
		if sim.gamesPerSeason != 16 {
			t.Errorf("Expected 16 games per season, got %d", sim.gamesPerSeason)
		}
	})

	t.Run("with default configs", func(t *testing.T) {
		cfg := YearSimulatorConfig{}
		sim := NewCareerSimulator(cfg)

		// Should use defaults
		if sim.clock == nil {
			t.Error("Clock should default to RealClock")
		}
		if sim.gamesPerSeason != 18 {
			t.Errorf("Expected default 18 games per season, got %d", sim.gamesPerSeason)
		}
		if sim.injuryRoller == nil {
			t.Error("InjuryRoller should have a default")
		}
		if sim.statsGenerator == nil {
			t.Error("StatsGenerator should have a default")
		}
		if sim.statMultiplier == nil {
			t.Error("StatMultiplier should have a default")
		}
	})
}

func TestCareerSimulatorCreateCareer(t *testing.T) {
	t.Run("rookie player", func(t *testing.T) {
		mockClock := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
		cfg := YearSimulatorConfig{
			Clock:          mockClock,
			GamesPerSeason: 18,
			InjuryRoller:   func(age int, position string) (bool, int) { return false, 0 },
			StatsGenerator: func(player Player, yoe int) FootballStats { return FootballStats{} },
			StatMultiplier: func(player Player, yoe int, stats FootballStats) FootballStats { return stats },
		}

		sim := NewCareerSimulator(cfg)

		player := Player{
			ID:                "player-1",
			DraftYear:         2025,
			Age:               22,
			Position:          "QB",
			YearsOfExperience: 0,
			Skill:             0.8,
		}

		career := sim.CreateCareer(player)

		// Rookie should have exactly 1 year
		if len(career) != 1 {
			t.Errorf("Expected 1 year for rookie, got %d", len(career))
		}

		if career[0].Year != 2025 {
			t.Errorf("Expected rookie year to be 2025, got %d", career[0].Year)
		}

		if career[0].PlayerID != player.ID {
			t.Errorf("Expected player ID %s, got %s", player.ID, career[0].PlayerID)
		}
	})

	t.Run("veteran player", func(t *testing.T) {
		mockClock := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
		cfg := YearSimulatorConfig{
			Clock:          mockClock,
			GamesPerSeason: 18,
			InjuryRoller:   func(age int, position string) (bool, int) { return false, 0 },
			StatsGenerator: func(player Player, yoe int) FootballStats { return FootballStats{} },
			StatMultiplier: func(player Player, yoe int, stats FootballStats) FootballStats { return stats },
		}

		sim := NewCareerSimulator(cfg)

		player := Player{
			ID:                "player-1",
			DraftYear:         2020,
			Age:               27,
			Position:          "QB",
			YearsOfExperience: 5,
			Skill:             0.8,
		}

		career := sim.CreateCareer(player)

		// Veteran drafted in 2020, current year 2025 = 5 years
		if len(career) != 5 {
			t.Errorf("Expected 5 years for veteran, got %d", len(career))
		}

		// Check that years are sequential
		for i, year := range career {
			expectedYear := 2020 + i
			if year.Year != expectedYear {
				t.Errorf("Expected year %d at index %d, got %d", expectedYear, i, year.Year)
			}
		}
	})
}

func TestCareerSimulatorSimulateYear(t *testing.T) {
	t.Run("full healthy season", func(t *testing.T) {
		mockClock := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
		gamesPlayed := 0
		cfg := YearSimulatorConfig{
			Clock:          mockClock,
			GamesPerSeason: 18,
			InjuryRoller:   func(age int, position string) (bool, int) { return false, 0 },
			StatsGenerator: func(player Player, yoe int) FootballStats {
				gamesPlayed++
				return FootballStats{PassingYards: 250, PassingTDs: 2}
			},
			StatMultiplier: func(player Player, yoe int, stats FootballStats) FootballStats {
				return stats
			},
		}

		sim := NewCareerSimulator(cfg)

		player := Player{
			ID:                "player-1",
			DraftYear:         2020,
			Age:               27,
			Position:          "QB",
			YearsOfExperience: 5,
			Skill:             0.8,
		}

		yearStats := sim.SimulateYear(player, 2025)

		// Should have played all 18 games
		if gamesPlayed != 18 {
			t.Errorf("Expected 18 games played, got %d", gamesPlayed)
		}

		// Should have accumulated stats
		expectedYards := 250 * 18
		if yearStats.Total.PassingYards != expectedYards {
			t.Errorf("Expected %d passing yards, got %d", expectedYards, yearStats.Total.PassingYards)
		}

		expectedTDs := 2 * 18
		if yearStats.Total.PassingTDs != expectedTDs {
			t.Errorf("Expected %d passing TDs, got %d", expectedTDs, yearStats.Total.PassingTDs)
		}
	})

	t.Run("season with injury", func(t *testing.T) {
		mockClock := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
		gamesPlayed := 0
		injuryAtGame := 10

		cfg := YearSimulatorConfig{
			Clock:          mockClock,
			GamesPerSeason: 18,
			InjuryRoller: func(age int, position string) (bool, int) {
				// Injury after 10 games, miss 5 games
				if gamesPlayed == injuryAtGame {
					return true, 5
				}
				return false, 0
			},
			StatsGenerator: func(player Player, yoe int) FootballStats {
				gamesPlayed++
				return FootballStats{PassingYards: 250}
			},
			StatMultiplier: func(player Player, yoe int, stats FootballStats) FootballStats {
				return stats
			},
		}

		sim := NewCareerSimulator(cfg)

		player := Player{
			ID:                "player-1",
			DraftYear:         2020,
			Age:               27,
			Position:          "QB",
			YearsOfExperience: 5,
			Skill:             0.8,
		}

		yearStats := sim.SimulateYear(player, 2025)

		// Should have played only 13 games (18 - 5 missed)
		if gamesPlayed != 13 {
			t.Errorf("Expected 13 games played (18 - 5 missed), got %d", gamesPlayed)
		}

		expectedYards := 250 * 13
		if yearStats.Total.PassingYards != expectedYards {
			t.Errorf("Expected %d passing yards, got %d", expectedYards, yearStats.Total.PassingYards)
		}
	})
}

func TestCreatePlayerCareer(t *testing.T) {
	// Test the wrapper function
	player := Player{
		ID:                "player-1",
		DraftYear:         2024,
		Age:               22,
		Position:          "QB",
		FirstName:         "Test",
		LastName:          "Player",
		YearsOfExperience: 1,
		Skill:             0.8,
	}

	// This will use the real simulator
	career := createPlayerCareer(player)

	// Should have at least 1 year
	if len(career) == 0 {
		t.Error("Career should have at least 1 year")
	}

	// Verify player ID matches
	for _, year := range career {
		if year.PlayerID != player.ID {
			t.Errorf("Expected player ID %s, got %s", player.ID, year.PlayerID)
		}
	}
}

func TestRealClockNow(t *testing.T) {
	clock := RealClock{}
	now := clock.Now()
	
	// Just verify it returns a valid time
	if now.IsZero() {
		t.Error("RealClock.Now() should not return zero time")
	}
}

func TestGeneratePlayerGameStats(t *testing.T) {
	player := Player{
		ID:                "player-1",
		DraftYear:         2020,
		Age:               27,
		Position:          "QB",
		YearsOfExperience: 5,
		Skill:             0.8,
	}

	tests := []struct {
		position     string
		checkStat    func(FootballStats) bool
		description  string
	}{
		{
			position: "QB",
			checkStat: func(stats FootballStats) bool {
				return stats.PassingAttempts > 0 || stats.PassingYards > 0
			},
			description: "QB should have passing stats",
		},
		{
			position: "RB",
			checkStat: func(stats FootballStats) bool {
				return stats.RushingAttempts > 0 || stats.RushingYards > 0
			},
			description: "RB should have rushing stats",
		},
		{
			position: "WR",
			checkStat: func(stats FootballStats) bool {
				return stats.ReceivingReceptions > 0 || stats.ReceivingYards > 0
			},
			description: "WR should have receiving stats",
		},
		{
			position: "TE",
			checkStat: func(stats FootballStats) bool {
				return stats.ReceivingReceptions >= 0 // TE might have 0 receptions in a game
			},
			description: "TE should have valid receiving stats",
		},
		{
			position: "PK",
			checkStat: func(stats FootballStats) bool {
				return stats.FieldGoals >= 0 // PK might have 0 field goals in a game
			},
			description: "PK should have kicking stats",
		},
		{
			position: "UNKNOWN",
			checkStat: func(stats FootballStats) bool {
				// Unknown position should return empty stats
				return stats.PassingAttempts == 0 && stats.RushingAttempts == 0 && 
					stats.ReceivingReceptions == 0 && stats.FieldGoals == 0
			},
			description: "Unknown position should return empty stats",
		},
	}

	for _, tt := range tests {
		t.Run(tt.position, func(t *testing.T) {
			player.Position = tt.position
			stats := generatePlayerGameStats(player, 5)

			// Run the check multiple times to ensure consistency
			for range 10 {
				stats = generatePlayerGameStats(player, 5)
				if !tt.checkStat(stats) {
					t.Errorf("%s failed: stats = %+v", tt.description, stats)
				}
			}
		})
	}
}

func TestMultiplyStatByPlayerSkill(t *testing.T) {
	player := Player{
		Skill: 0.8,
	}

	tests := []struct {
		name              string
		yearsOfExperience int
		stat              int
		description       string
	}{
		{"rookie", 0, 100, "rookie with base stat 100"},
		{"veteran", 5, 100, "5-year veteran with base stat 100"},
		{"star", 10, 200, "10-year star with base stat 200"},
		{"zero stat", 0, 0, "zero stat should remain zero"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := multiplyStatByPlayerSkill(player, tt.yearsOfExperience, tt.stat)

			if tt.stat == 0 {
				if result != 0 {
					t.Errorf("Zero stat should remain zero, got %d", result)
				}
				return
			}

			// Result should be less than or equal to original stat * (1 + yoe/100) * skill
			maxExpected := int(float64(tt.stat) * (1.0 + float64(tt.yearsOfExperience)/100.0))

			if result > maxExpected {
				t.Errorf("Result %d is higher than expected maximum %d", result, maxExpected)
			}

			// Result should generally be less than original for 0.8 skill
			// (unless years of experience compensates)
			if tt.yearsOfExperience == 0 && result > tt.stat {
				t.Errorf("Rookie with 0.8 skill should have stat <= original. Got %d from %d", result, tt.stat)
			}
		})
	}

	// Test with very low skill player
	t.Run("low skill player", func(t *testing.T) {
		lowSkillPlayer := Player{Skill: 0.2}
		result := multiplyStatByPlayerSkill(lowSkillPlayer, 0, 100)
		
		// Should be significantly reduced
		if result >= 100 {
			t.Errorf("Low skill player should have reduced stats, got %d from 100", result)
		}
	})

	// Test with high skill player
	t.Run("high skill player", func(t *testing.T) {
		highSkillPlayer := Player{Skill: 0.95}
		result := multiplyStatByPlayerSkill(highSkillPlayer, 0, 100)
		
		// Should be close to original or slightly less
		if result < 50 {
			t.Errorf("High skill player should have stats close to original, got %d from 100", result)
		}
	})
}

func TestMultiplyYearlyStatsByPlayerSkill(t *testing.T) {
	player := Player{
		Skill: 0.9,
	}

	stats := FootballStats{
		PassingAttempts:    100,
		PassingCompletions: 65,
		PassingYards:       1000,
		PassingTDs:         10,
		RushingAttempts:    20,
		RushingYards:       100,
	}

	adjusted := multiplyYearlyStatsByPlayerSkill(player, 5, stats)

	// All stats should be adjusted
	if adjusted.PassingAttempts == stats.PassingAttempts {
		t.Error("PassingAttempts should be adjusted")
	}
	if adjusted.PassingCompletions == stats.PassingCompletions {
		t.Error("PassingCompletions should be adjusted")
	}
	if adjusted.PassingYards == stats.PassingYards {
		t.Error("PassingYards should be adjusted")
	}
	if adjusted.PassingTDs == stats.PassingTDs {
		t.Error("PassingTDs should be adjusted")
	}
}

func TestNormalInRange(t *testing.T) {
	tests := []struct {
		name string
		low  float64
		high float64
	}{
		{"small range", 1.0, 5.0},
		{"large range", 0.0, 100.0},
		{"negative range", -10.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for range 100 {
				result := normalInRange(tt.low, tt.high)

				// Result should be clamped within bounds
				if result < tt.low {
					t.Errorf("Result %f is below low bound %f", result, tt.low)
				}
				if result > tt.high {
					t.Errorf("Result %f is above high bound %f", result, tt.high)
				}
			}
		})
	}
}

func TestNormalIntInRange(t *testing.T) {
	tests := []struct {
		name string
		low  int
		high int
	}{
		{"small range", 1, 5},
		{"large range", 0, 100},
		{"single value", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for range 100 {
				result := normalIntInRange(tt.low, tt.high)

				// Result should be within bounds
				if result < tt.low {
					t.Errorf("Result %d is below low bound %d", result, tt.low)
				}
				if result > tt.high {
					t.Errorf("Result %d is above high bound %d", result, tt.high)
				}
			}
		})
	}
}

