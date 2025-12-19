package main

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// =============================================================================
// MOCK IMPLEMENTATIONS
// =============================================================================

// MockDataGenerator provides controlled test data
type MockDataGenerator struct {
	LeagueData LeagueFlat
	RosterData FootballTeamRoster
	CareerData []PlayerYearlyStatsFootball
	CallCounts map[string]int
}

func NewMockDataGenerator() *MockDataGenerator {
	return &MockDataGenerator{
		CallCounts: make(map[string]int),
		LeagueData: LeagueFlat{
			Conferences: []Conference{
				{ID: "conf-1", Name: "Test Conference"},
			},
			Divisions: []Division{
				{ID: "div-1", Name: "Test Division", ConferenceID: "conf-1"},
			},
			Teams: []Team{
				{ID: "team-1", City: "Test City", State: "TS", Name: "Testers", Abbr: "TST", DivisionID: "div-1"},
			},
		},
		RosterData: FootballTeamRoster{
			QB: []Player{
				{ID: "player-1", FirstName: "Test", LastName: "QB", Position: "QB", TeamID: "team-1", Height: 75, Weight: 220, Age: 25, YearsOfExperience: 3, DraftYear: 2022, Skill: 0.85, Status: "ACTIVE", Jersey: 12},
			},
			RB: []Player{},
			WR: []Player{},
			TE: []Player{},
			PK: []Player{},
		},
		CareerData: []PlayerYearlyStatsFootball{
			{PlayerID: "player-1", Year: 2024, Stats: FootballYearlyStats{Total: FootballStats{PassingYards: 4000, PassingTDs: 30}}},
		},
	}
}

func (m *MockDataGenerator) GenerateLeague() LeagueFlat {
	m.CallCounts["GenerateLeague"]++
	return m.LeagueData
}

func (m *MockDataGenerator) GenerateRoster(teamID string) FootballTeamRoster {
	m.CallCounts["GenerateRoster"]++
	return m.RosterData
}

func (m *MockDataGenerator) GenerateCareer(player Player) []PlayerYearlyStatsFootball {
	m.CallCounts["GenerateCareer"]++
	return m.CareerData
}

// MockTx implements pgx.Tx for testing
type MockTx struct {
	ExecCalls      []MockExecCall
	ExecErr        error
	ExecErrOnCall  int // Return error on this call number (0 = never)
	currentCall    int
	CommitCalled   bool
	CommitErr      error
	RollbackCalled bool
}

type MockExecCall struct {
	SQL  string
	Args []any
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.currentCall++
	m.ExecCalls = append(m.ExecCalls, MockExecCall{SQL: sql, Args: arguments})

	if m.ExecErrOnCall > 0 && m.currentCall == m.ExecErrOnCall {
		return pgconn.CommandTag{}, m.ExecErr
	}
	if m.ExecErr != nil && m.ExecErrOnCall == 0 {
		return pgconn.CommandTag{}, m.ExecErr
	}
	return pgconn.CommandTag{}, nil
}

func (m *MockTx) Commit(ctx context.Context) error {
	m.CommitCalled = true
	return m.CommitErr
}

func (m *MockTx) Rollback(ctx context.Context) error {
	m.RollbackCalled = true
	return nil
}

// Implement remaining pgx.Tx interface methods (not used in tests)
func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) { return nil, nil }
func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults { return nil }
func (m *MockTx) LargeObjects() pgx.LargeObjects                               { return pgx.LargeObjects{} }
func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row { return nil }
func (m *MockTx) Conn() *pgx.Conn                                               { return nil }

// =============================================================================
// TESTS
// =============================================================================

func TestNewDatabaseSeeder(t *testing.T) {
	t.Run("with default config", func(t *testing.T) {
		seeder := NewDatabaseSeeder(SeederConfig{})

		if seeder.generator == nil {
			t.Error("Expected generator to have a default")
		}
		if seeder.logger == nil {
			t.Error("Expected logger to have a default")
		}
	})

	t.Run("with custom generator", func(t *testing.T) {
		mockGen := NewMockDataGenerator()
		seeder := NewDatabaseSeeder(SeederConfig{
			DataGenerator: mockGen,
		})

		if seeder.generator != mockGen {
			t.Error("Expected custom generator to be used")
		}
	})

	t.Run("with quiet mode", func(t *testing.T) {
		seeder := NewDatabaseSeeder(SeederConfig{
			Quiet: true,
		})

		if !seeder.quiet {
			t.Error("Expected quiet mode to be enabled")
		}
	})
}

func TestDatabaseSeederSeed(t *testing.T) {
	t.Run("successful seed", func(t *testing.T) {
		mockGen := NewMockDataGenerator()
		mockTx := &MockTx{}

		seeder := NewDatabaseSeeder(SeederConfig{
			DataGenerator: mockGen,
			Quiet:         true,
		})

		ctx := context.Background()
		result, err := seeder.Seed(ctx, mockTx)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		// Verify counts
		if result.ConferencesInserted != 1 {
			t.Errorf("Expected 1 conference, got %d", result.ConferencesInserted)
		}
		if result.DivisionsInserted != 1 {
			t.Errorf("Expected 1 division, got %d", result.DivisionsInserted)
		}
		if result.TeamsInserted != 1 {
			t.Errorf("Expected 1 team, got %d", result.TeamsInserted)
		}
		if result.PlayersInserted != 1 {
			t.Errorf("Expected 1 player, got %d", result.PlayersInserted)
		}
		if result.YearlyStatsInserted != 1 {
			t.Errorf("Expected 1 yearly stat, got %d", result.YearlyStatsInserted)
		}

		// Verify generator was called
		if mockGen.CallCounts["GenerateLeague"] != 1 {
			t.Error("Expected GenerateLeague to be called once")
		}
		if mockGen.CallCounts["GenerateRoster"] != 1 {
			t.Error("Expected GenerateRoster to be called once per team")
		}
		if mockGen.CallCounts["GenerateCareer"] != 1 {
			t.Error("Expected GenerateCareer to be called once per player")
		}
	})

	t.Run("purge failure", func(t *testing.T) {
		mockGen := NewMockDataGenerator()
		mockTx := &MockTx{
			ExecErr:       errors.New("purge failed"),
			ExecErrOnCall: 1, // Fail on first exec (purge)
		}

		seeder := NewDatabaseSeeder(SeederConfig{
			DataGenerator: mockGen,
			Quiet:         true,
		})

		ctx := context.Background()
		_, err := seeder.Seed(ctx, mockTx)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, mockTx.ExecErr) && err.Error() != "failed to purge database: failed to purge table fantasy_rosters: purge failed" {
			t.Errorf("Expected purge error, got: %v", err)
		}
	})
}

func TestFlattenRoster(t *testing.T) {
	roster := FootballTeamRoster{
		QB: []Player{{ID: "qb-1"}, {ID: "qb-2"}},
		RB: []Player{{ID: "rb-1"}},
		WR: []Player{{ID: "wr-1"}, {ID: "wr-2"}, {ID: "wr-3"}},
		TE: []Player{{ID: "te-1"}},
		PK: []Player{{ID: "pk-1"}},
	}

	players := flattenRoster(roster)

	if len(players) != 8 {
		t.Errorf("Expected 8 players, got %d", len(players))
	}

	// Check order
	expectedIDs := []string{"qb-1", "qb-2", "rb-1", "wr-1", "wr-2", "wr-3", "te-1", "pk-1"}
	for i, expectedID := range expectedIDs {
		if players[i].ID != expectedID {
			t.Errorf("Expected player %d to have ID %s, got %s", i, expectedID, players[i].ID)
		}
	}
}

func TestFlattenRosterEmpty(t *testing.T) {
	roster := FootballTeamRoster{}
	players := flattenRoster(roster)

	if len(players) != 0 {
		t.Errorf("Expected 0 players for empty roster, got %d", len(players))
	}
}

func TestDefaultDataGenerator(t *testing.T) {
	gen := NewDefaultDataGenerator()

	t.Run("GenerateLeague returns valid data", func(t *testing.T) {
		league := gen.GenerateLeague()

		if len(league.Conferences) == 0 {
			t.Error("Expected conferences")
		}
		if len(league.Divisions) == 0 {
			t.Error("Expected divisions")
		}
		if len(league.Teams) == 0 {
			t.Error("Expected teams")
		}
	})
}

func TestMaskPassword(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"postgres://user:pass@localhost:5432/db", "postgres://****:****@****:5432/fantasy_db"},
		{"anything", "postgres://****:****@****:5432/fantasy_db"},
	}

	for _, tt := range tests {
		result := maskPassword(tt.input)
		if result != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, result)
		}
	}
}

func TestSeederLogging(t *testing.T) {
	var logs []string
	mockLogger := func(format string, v ...any) {
		logs = append(logs, format)
	}

	t.Run("logs when not quiet", func(t *testing.T) {
		logs = nil
		mockGen := NewMockDataGenerator()
		mockTx := &MockTx{}

		seeder := NewDatabaseSeeder(SeederConfig{
			DataGenerator: mockGen,
			Logger:        mockLogger,
			Quiet:         false,
		})

		ctx := context.Background()
		seeder.Seed(ctx, mockTx)

		if len(logs) == 0 {
			t.Error("Expected logs to be written")
		}
	})

	t.Run("no logs when quiet", func(t *testing.T) {
		logs = nil
		mockGen := NewMockDataGenerator()
		mockTx := &MockTx{}

		seeder := NewDatabaseSeeder(SeederConfig{
			DataGenerator: mockGen,
			Logger:        mockLogger,
			Quiet:         true,
		})

		ctx := context.Background()
		seeder.Seed(ctx, mockTx)

		if len(logs) != 0 {
			t.Errorf("Expected no logs in quiet mode, got %d", len(logs))
		}
	})
}

func TestSeedResult(t *testing.T) {
	result := &SeedResult{
		ConferencesInserted: 2,
		DivisionsInserted:   8,
		TeamsInserted:       32,
		PlayersInserted:     544,
		YearlyStatsInserted: 2000,
	}

	if result.ConferencesInserted != 2 {
		t.Errorf("ConferencesInserted mismatch")
	}
	if result.DivisionsInserted != 8 {
		t.Errorf("DivisionsInserted mismatch")
	}
	if result.TeamsInserted != 32 {
		t.Errorf("TeamsInserted mismatch")
	}
	if result.PlayersInserted != 544 {
		t.Errorf("PlayersInserted mismatch")
	}
	if result.YearlyStatsInserted != 2000 {
		t.Errorf("YearlyStatsInserted mismatch")
	}
}
