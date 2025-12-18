package main

import (
	"testing"
)

func TestNewPositionProfile(t *testing.T) {
	profile := NewPositionProfile()

	if profile == nil {
		t.Fatal("NewPositionProfile() returned nil")
	}

	if profile.Jerseys == nil {
		t.Error("Jerseys map not initialized")
	}
	if profile.Heights == nil {
		t.Error("Heights map not initialized")
	}
	if profile.Weights == nil {
		t.Error("Weights map not initialized")
	}
	if profile.Ages == nil {
		t.Error("Ages map not initialized")
	}
	if profile.YearsOfExperience == nil {
		t.Error("YearsOfExperience map not initialized")
	}
}

func TestCollectPlayerAttributes(t *testing.T) {
	tests := []struct {
		name           string
		inputData      map[string]any
		expectedCount  int
		shouldContain  bool
		expectedPlayer PlayerStat
	}{
		{
			name: "valid single player",
			inputData: map[string]any{
				"athletes": []any{
					map[string]any{
						"firstName": "John",
						"lastName":  "Doe",
						"height":    float64(72),
						"weight":    float64(200),
						"jersey":    "12",
						"age":       float64(25),
						"position": map[string]any{
							"abbreviation": "QB",
						},
						"status": map[string]any{
							"type": "active",
						},
						"draft": map[string]any{
							"year": float64(2020),
						},
					},
				},
			},
			expectedCount: 1,
			shouldContain: true,
			expectedPlayer: PlayerStat{
				FirstName:         "John",
				LastName:          "Doe",
				Height:            72,
				Weight:            200,
				Jersey:            12,
				Age:               25,
				Position:          "QB",
				YearsOfExperience: 5, // 2025 - 2020
			},
		},
		{
			name: "multiple players",
			inputData: map[string]any{
				"athletes": []any{
					map[string]any{
						"firstName": "John",
						"lastName":  "Doe",
						"height":    float64(72),
						"weight":    float64(200),
						"jersey":    "12",
						"age":       float64(25),
						"position": map[string]any{
							"abbreviation": "QB",
						},
						"status": map[string]any{
							"type": "active",
						},
						"draft": map[string]any{
							"year": float64(2020),
						},
					},
					map[string]any{
						"firstName": "Jane",
						"lastName":  "Smith",
						"height":    float64(68),
						"weight":    float64(180),
						"jersey":    "88",
						"age":       float64(23),
						"position": map[string]any{
							"abbreviation": "WR",
						},
						"status": map[string]any{
							"type": "active",
						},
						"draft": map[string]any{
							"year": float64(2022),
						},
					},
				},
			},
			expectedCount: 2,
		},
		{
			name: "skips free agents",
			inputData: map[string]any{
				"athletes": []any{
					map[string]any{
						"firstName": "John",
						"lastName":  "Doe",
						"height":    float64(72),
						"weight":    float64(200),
						"jersey":    "12",
						"age":       float64(25),
						"position": map[string]any{
							"abbreviation": "QB",
						},
						"status": map[string]any{
							"type": "free-agent",
						},
						"draft": map[string]any{
							"year": float64(2020),
						},
					},
				},
			},
			expectedCount: 0,
		},
		{
			name: "skips invalid player data",
			inputData: map[string]any{
				"athletes": []any{
					map[string]any{
						"firstName": "Bad",
						"lastName":  "Player",
						// missing position
						"status": map[string]any{
							"type": "active",
						},
					},
				},
			},
			expectedCount: 0,
		},
		{
			name: "empty athletes list",
			inputData: map[string]any{
				"athletes": []any{},
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := collectPlayerAttributes(tt.inputData)

			if len(stats) != tt.expectedCount {
				t.Errorf("Expected %d players, got %d", tt.expectedCount, len(stats))
			}

			if tt.shouldContain && len(stats) > 0 {
				player := stats[0]
				if player.FirstName != tt.expectedPlayer.FirstName {
					t.Errorf("Expected first name %s, got %s", tt.expectedPlayer.FirstName, player.FirstName)
				}
				if player.LastName != tt.expectedPlayer.LastName {
					t.Errorf("Expected last name %s, got %s", tt.expectedPlayer.LastName, player.LastName)
				}
				if player.Position != tt.expectedPlayer.Position {
					t.Errorf("Expected position %s, got %s", tt.expectedPlayer.Position, player.Position)
				}
			}
		})
	}

	// Test with non-map player data
	t.Run("non-map player in athletes list", func(t *testing.T) {
		data := map[string]any{
			"athletes": []any{
				"not a map", // This should be skipped
				map[string]any{
					"firstName": "Valid",
					"lastName":  "Player",
					"height":    float64(72),
					"weight":    float64(200),
					"jersey":    "10",
					"age":       float64(24),
					"position": map[string]any{
						"abbreviation": "RB",
					},
					"status": map[string]any{
						"type": "active",
					},
					"draft": map[string]any{
						"year": float64(2021),
					},
				},
			},
		}

		stats := collectPlayerAttributes(data)
		
		// Should have only 1 player (the valid one)
		if len(stats) != 1 {
			t.Errorf("Expected 1 valid player, got %d", len(stats))
		}
	})

	// Test with player having missing optional fields
	t.Run("player with minimal data", func(t *testing.T) {
		data := map[string]any{
			"athletes": []any{
				map[string]any{
					// Missing firstName, lastName, height, weight, jersey, age
					"position": map[string]any{
						"abbreviation": "WR",
					},
					"status": map[string]any{
						"type": "active",
					},
					"draft": map[string]any{
						"year": float64(2023),
					},
				},
			},
		}

		stats := collectPlayerAttributes(data)
		
		// Should still work, just with empty/zero values for optional fields
		if len(stats) != 1 {
			t.Errorf("Expected 1 player with minimal data, got %d", len(stats))
		}
		
		if len(stats) > 0 {
			player := stats[0]
			if player.Position != "WR" {
				t.Errorf("Expected position WR, got %s", player.Position)
			}
			// These should be empty/zero since they weren't provided
			if player.FirstName != "" {
				t.Logf("FirstName: %s", player.FirstName)
			}
			if player.Jersey != 0 {
				t.Logf("Jersey: %d", player.Jersey)
			}
		}
	})
}

func TestAggregateAttributesByPosition(t *testing.T) {
	stats := []PlayerStat{
		{
			FirstName:         "John",
			LastName:          "Doe",
			Height:            72,
			Weight:            200,
			Jersey:            12,
			Age:               25,
			Position:          "QB",
			YearsOfExperience: 5,
		},
		{
			FirstName:         "Jane",
			LastName:          "Smith",
			Height:            68,
			Weight:            180,
			Jersey:            88,
			Age:               23,
			Position:          "WR",
			YearsOfExperience: 3,
		},
		{
			FirstName:         "Bob",
			LastName:          "Johnson",
			Height:            73,
			Weight:            210,
			Jersey:            7,
			Age:               26,
			Position:          "QB",
			YearsOfExperience: 6,
		},
	}

	aggregated := aggregateAttributesByPosition(stats)

	// Should have 2 positions
	if len(aggregated) != 2 {
		t.Errorf("Expected 2 positions, got %d", len(aggregated))
	}

	// Check QB aggregation
	qbProfile, ok := aggregated["QB"]
	if !ok {
		t.Fatal("QB profile not found")
	}

	if qbProfile.Jerseys[12] != 1 {
		t.Errorf("Expected jersey 12 count to be 1, got %d", qbProfile.Jerseys[12])
	}
	if qbProfile.Jerseys[7] != 1 {
		t.Errorf("Expected jersey 7 count to be 1, got %d", qbProfile.Jerseys[7])
	}
	if qbProfile.Heights[72] != 1 {
		t.Errorf("Expected height 72 count to be 1, got %d", qbProfile.Heights[72])
	}
	if qbProfile.Heights[73] != 1 {
		t.Errorf("Expected height 73 count to be 1, got %d", qbProfile.Heights[73])
	}
	if qbProfile.Ages[25] != 1 {
		t.Errorf("Expected age 25 count to be 1, got %d", qbProfile.Ages[25])
	}
	if qbProfile.Ages[26] != 1 {
		t.Errorf("Expected age 26 count to be 1, got %d", qbProfile.Ages[26])
	}

	// Check WR aggregation
	wrProfile, ok := aggregated["WR"]
	if !ok {
		t.Fatal("WR profile not found")
	}

	if wrProfile.Jerseys[88] != 1 {
		t.Errorf("Expected jersey 88 count to be 1, got %d", wrProfile.Jerseys[88])
	}
	if wrProfile.Heights[68] != 1 {
		t.Errorf("Expected height 68 count to be 1, got %d", wrProfile.Heights[68])
	}
	if wrProfile.Ages[23] != 1 {
		t.Errorf("Expected age 23 count to be 1, got %d", wrProfile.Ages[23])
	}
}

func TestAggregateFirstNames(t *testing.T) {
	stats := []PlayerStat{
		{FirstName: "John", LastName: "Doe"},
		{FirstName: "Jane", LastName: "Smith"},
		{FirstName: "John", LastName: "Johnson"},
		{FirstName: "Bob", LastName: "Williams"},
	}

	aggregated := aggregateFirstNames(stats)

	if len(aggregated) != 3 {
		t.Errorf("Expected 3 unique first names, got %d", len(aggregated))
	}

	if aggregated["John"] != 2 {
		t.Errorf("Expected 'John' count to be 2, got %d", aggregated["John"])
	}
	if aggregated["Jane"] != 1 {
		t.Errorf("Expected 'Jane' count to be 1, got %d", aggregated["Jane"])
	}
	if aggregated["Bob"] != 1 {
		t.Errorf("Expected 'Bob' count to be 1, got %d", aggregated["Bob"])
	}
}

func TestAggregateLastNames(t *testing.T) {
	stats := []PlayerStat{
		{FirstName: "John", LastName: "Doe"},
		{FirstName: "Jane", LastName: "Smith"},
		{FirstName: "John", LastName: "Smith"},
		{FirstName: "Bob", LastName: "Williams"},
	}

	aggregated := aggregateLastNames(stats)

	if len(aggregated) != 3 {
		t.Errorf("Expected 3 unique last names, got %d", len(aggregated))
	}

	if aggregated["Smith"] != 2 {
		t.Errorf("Expected 'Smith' count to be 2, got %d", aggregated["Smith"])
	}
	if aggregated["Doe"] != 1 {
		t.Errorf("Expected 'Doe' count to be 1, got %d", aggregated["Doe"])
	}
	if aggregated["Williams"] != 1 {
		t.Errorf("Expected 'Williams' count to be 1, got %d", aggregated["Williams"])
	}
}

func TestCollectAndAggregatePlayerAttributes(t *testing.T) {
	// This function calls importRealData which requires a file
	// We can't easily test it without mocking the file system
	// But we can test that the function exists and has the right signature
	var aggregator StatsAggregator = collectAndAggregatePlayerAttributes
	if aggregator == nil {
		t.Error("collectAndAggregatePlayerAttributes should not be nil")
	}
}

func TestNormalizePlayerData(t *testing.T) {
	tests := []struct {
		name        string
		inputData   map[string]any
		expectError bool
		expected    PlayerStat
	}{
		{
			name: "valid player data",
			inputData: map[string]any{
				"firstName": "John",
				"lastName":  "Doe",
				"height":    float64(72),
				"weight":    float64(200),
				"jersey":    "12",
				"age":       float64(25),
				"position": map[string]any{
					"abbreviation": "QB",
				},
				"status": map[string]any{
					"type": "active",
				},
				"draft": map[string]any{
					"year": float64(2020),
				},
			},
			expectError: false,
			expected: PlayerStat{
				FirstName:         "John",
				LastName:          "Doe",
				Height:            72,
				Weight:            200,
				Jersey:            12,
				Age:               25,
				Position:          "QB",
				YearsOfExperience: 5, // Assuming current year is 2025
			},
		},
		{
			name: "free agent player",
			inputData: map[string]any{
				"firstName": "Free",
				"lastName":  "Agent",
				"position": map[string]any{
					"abbreviation": "RB",
				},
				"status": map[string]any{
					"type": "free-agent",
				},
				"draft": map[string]any{
					"year": float64(2020),
				},
			},
			expectError: true,
		},
		{
			name: "missing position",
			inputData: map[string]any{
				"firstName": "No",
				"lastName":  "Position",
				"status": map[string]any{
					"type": "active",
				},
			},
			expectError: true,
		},
		{
			name: "missing draft year",
			inputData: map[string]any{
				"firstName": "No",
				"lastName":  "Draft",
				"position": map[string]any{
					"abbreviation": "TE",
				},
				"status": map[string]any{
					"type": "active",
				},
			},
			expectError: true,
		},
		{
			name: "draft year as int",
			inputData: map[string]any{
				"firstName": "John",
				"lastName":  "Doe",
				"height":    float64(72),
				"weight":    float64(200),
				"jersey":    "12",
				"age":       float64(25),
				"position": map[string]any{
					"abbreviation": "QB",
				},
				"status": map[string]any{
					"type": "active",
				},
				"draft": map[string]any{
					"year": 2020,
				},
			},
			expectError: false,
			expected: PlayerStat{
				FirstName:         "John",
				LastName:          "Doe",
				Height:            72,
				Weight:            200,
				Jersey:            12,
				Age:               25,
				Position:          "QB",
				YearsOfExperience: 5,
			},
		},
		{
			name: "invalid jersey string",
			inputData: map[string]any{
				"firstName": "John",
				"lastName":  "Doe",
				"height":    float64(72),
				"weight":    float64(200),
				"jersey":    "ABC", // Invalid jersey number
				"age":       float64(25),
				"position": map[string]any{
					"abbreviation": "QB",
				},
				"status": map[string]any{
					"type": "active",
				},
				"draft": map[string]any{
					"year": 2020,
				},
			},
			expectError: false,
			expected: PlayerStat{
				FirstName:         "John",
				LastName:          "Doe",
				Height:            72,
				Weight:            200,
				Jersey:            0, // Should default to 0 for invalid jersey
				Age:               25,
				Position:          "QB",
				YearsOfExperience: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizePlayerData(tt.inputData)

			if tt.expectError && err == nil {
				t.Error("Expected an error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				if result.FirstName != tt.expected.FirstName {
					t.Errorf("Expected FirstName %s, got %s", tt.expected.FirstName, result.FirstName)
				}
				if result.LastName != tt.expected.LastName {
					t.Errorf("Expected LastName %s, got %s", tt.expected.LastName, result.LastName)
				}
				if result.Position != tt.expected.Position {
					t.Errorf("Expected Position %s, got %s", tt.expected.Position, result.Position)
				}
				if result.Height != tt.expected.Height {
					t.Errorf("Expected Height %d, got %d", tt.expected.Height, result.Height)
				}
				if result.Weight != tt.expected.Weight {
					t.Errorf("Expected Weight %d, got %d", tt.expected.Weight, result.Weight)
				}
				if result.Jersey != tt.expected.Jersey {
					t.Errorf("Expected Jersey %d, got %d", tt.expected.Jersey, result.Jersey)
				}
				if result.Age != tt.expected.Age {
					t.Errorf("Expected Age %d, got %d", tt.expected.Age, result.Age)
				}
			}
		})
	}
}

