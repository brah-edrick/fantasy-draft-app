package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// =============================================================================
// INTERFACES FOR DEPENDENCY INJECTION
// =============================================================================

// DBExecutor interface for database operations (allows mocking in tests)
type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag, error)
}

// commandTag is a simple interface to abstract pgconn.CommandTag
type commandTag interface {
	RowsAffected() int64
}

// DataGenerator interface for generating synthetic data
type DataGenerator interface {
	GenerateLeague() LeagueFlat
	GenerateRoster(teamID string) FootballTeamRoster
	GenerateCareer(player Player) []PlayerYearlyStatsFootball
}

// =============================================================================
// DEFAULT IMPLEMENTATIONS
// =============================================================================

// DefaultDataGenerator uses the real data generation functions
type DefaultDataGenerator struct {
	uuidGenerator UUIDGenerator
	clock         Clock
	rng           *rand.Rand
}

func NewDefaultDataGenerator() *DefaultDataGenerator {
	return &DefaultDataGenerator{
		uuidGenerator: UUIDGenerator(func() string { return uuid.New().String() }),
		clock:         RealClock{},
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *DefaultDataGenerator) GenerateLeague() LeagueFlat {
	return generateLeagueFlat(g.uuidGenerator, g.clock, g.rng)
}

func (g *DefaultDataGenerator) GenerateRoster(teamID string) FootballTeamRoster {
	return createTeamRoster(teamID)
}

func (g *DefaultDataGenerator) GenerateCareer(player Player) []PlayerYearlyStatsFootball {
	sim := NewCareerSimulator(YearSimulatorConfig{})
	return sim.CreateCareer(player)
}

// =============================================================================
// SEEDER CONFIG AND IMPLEMENTATION
// =============================================================================

// SeederConfig holds all injectable dependencies for the database seeder
type SeederConfig struct {
	// DataGenerator for creating synthetic data (default: DefaultDataGenerator)
	DataGenerator DataGenerator

	// Logger for output (default: log.Printf)
	Logger func(format string, v ...any)

	// Quiet mode suppresses logging
	Quiet bool
}

// DatabaseSeeder handles seeding with injectable dependencies
type DatabaseSeeder struct {
	generator DataGenerator
	logger    func(format string, v ...any)
	quiet     bool
}

// NewDatabaseSeeder creates a seeder with the given config
func NewDatabaseSeeder(cfg SeederConfig) *DatabaseSeeder {
	seeder := &DatabaseSeeder{
		generator: cfg.DataGenerator,
		logger:    cfg.Logger,
		quiet:     cfg.Quiet,
	}

	// Apply defaults
	if seeder.generator == nil {
		seeder.generator = NewDefaultDataGenerator()
	}
	if seeder.logger == nil {
		seeder.logger = log.Printf
	}

	return seeder
}

func (s *DatabaseSeeder) log(format string, v ...any) {
	if !s.quiet {
		s.logger(format, v...)
	}
}

// SeedResult contains the results of a seeding operation
type SeedResult struct {
	ConferencesInserted int
	DivisionsInserted   int
	TeamsInserted       int
	PlayersInserted     int
	YearlyStatsInserted int
}

// Seed performs the database seeding operation
func (s *DatabaseSeeder) Seed(ctx context.Context, tx pgx.Tx) (*SeedResult, error) {
	s.log("🗑️  Purging existing data...")
	if err := purgeDatabase(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to purge database: %w", err)
	}

	s.log("🏈 Generating synthetic data...")
	leagueData := s.generator.GenerateLeague()

	s.log("📝 Inserting conferences...")
	if err := insertConferences(ctx, tx, leagueData.Conferences); err != nil {
		return nil, fmt.Errorf("failed to insert conferences: %w", err)
	}

	s.log("📝 Inserting divisions...")
	if err := insertDivisions(ctx, tx, leagueData.Divisions); err != nil {
		return nil, fmt.Errorf("failed to insert divisions: %w", err)
	}

	s.log("📝 Inserting teams...")
	if err := insertTeams(ctx, tx, leagueData.Teams); err != nil {
		return nil, fmt.Errorf("failed to insert teams: %w", err)
	}

	// Generate rosters and players
	s.log("👥 Generating players and rosters...")
	var allPlayers []Player
	var allCareerStats []PlayerYearlyStatsFootball

	for _, team := range leagueData.Teams {
		roster := s.generator.GenerateRoster(team.ID)
		players := flattenRoster(roster)
		allPlayers = append(allPlayers, players...)

		// Generate career stats for each player
		for _, player := range players {
			career := s.generator.GenerateCareer(player)
			allCareerStats = append(allCareerStats, career...)
		}
	}

	s.log("📝 Inserting %d players...", len(allPlayers))
	if err := insertPlayers(ctx, tx, allPlayers); err != nil {
		return nil, fmt.Errorf("failed to insert players: %w", err)
	}

	s.log("📝 Inserting %d yearly stats records...", len(allCareerStats))
	if err := insertYearlyStats(ctx, tx, allCareerStats); err != nil {
		return nil, fmt.Errorf("failed to insert yearly stats: %w", err)
	}

	result := &SeedResult{
		ConferencesInserted: len(leagueData.Conferences),
		DivisionsInserted:   len(leagueData.Divisions),
		TeamsInserted:       len(leagueData.Teams),
		PlayersInserted:     len(allPlayers),
		YearlyStatsInserted: len(allCareerStats),
	}

	s.log("✅ Database seeded successfully!")
	s.log("   - %d conferences", result.ConferencesInserted)
	s.log("   - %d divisions", result.DivisionsInserted)
	s.log("   - %d teams", result.TeamsInserted)
	s.log("   - %d players", result.PlayersInserted)
	s.log("   - %d yearly stat records", result.YearlyStatsInserted)

	return result, nil
}

// =============================================================================
// DATABASE OPERATIONS (used by both old and new API)
// =============================================================================

// purgeDatabase deletes all data from tables in the correct order (respecting foreign keys)
func purgeDatabase(ctx context.Context, tx pgx.Tx) error {
	// Order matters due to foreign key constraints - delete children first
	tables := []string{
		"fantasy_rosters",
		"fantasy_teams",
		"rankings",
		"ranking_lists",
		"team_depth_charts",
		"yearly_stats",
		"players",
		"pro_teams",
		"divisions",
		"conferences",
		"draft_rooms",
		"users",
	}

	for _, table := range tables {
		_, err := tx.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to purge table %s: %w", table, err)
		}
	}
	return nil
}

func insertConferences(ctx context.Context, tx pgx.Tx, conferences []Conference) error {
	for _, conf := range conferences {
		_, err := tx.Exec(ctx,
			"INSERT INTO conferences (id, name) VALUES ($1, $2)",
			conf.ID, conf.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertDivisions(ctx context.Context, tx pgx.Tx, divisions []Division) error {
	for _, div := range divisions {
		_, err := tx.Exec(ctx,
			"INSERT INTO divisions (id, name, conference_id) VALUES ($1, $2, $3)",
			div.ID, div.Name, div.ConferenceID)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertTeams(ctx context.Context, tx pgx.Tx, teams []Team) error {
	for _, team := range teams {
		_, err := tx.Exec(ctx,
			"INSERT INTO pro_teams (id, city, state, name, abbreviation, division_id) VALUES ($1, $2, $3, $4, $5, $6)",
			team.ID, team.City, team.State, team.Name, team.Abbr, team.DivisionID)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertPlayers(ctx context.Context, tx pgx.Tx, players []Player) error {
	for _, player := range players {
		_, err := tx.Exec(ctx,
			`INSERT INTO players (id, first_name, last_name, position, team_id, height, weight, age, years_of_experience, draft_year, jersey_number, status, skill)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			player.ID, player.FirstName, player.LastName, player.Position, player.TeamID,
			player.Height, player.Weight, player.Age, player.YearsOfExperience, player.DraftYear,
			player.Jersey, player.Status, player.Skill)
		if err != nil {
			return fmt.Errorf("failed to insert player %s %s: %w", player.FirstName, player.LastName, err)
		}
	}
	return nil
}

func insertYearlyStats(ctx context.Context, tx pgx.Tx, stats []PlayerYearlyStatsFootball) error {
	for _, stat := range stats {
		// Marshal the stats to JSON
		statsJSON, err := json.Marshal(stat.Stats)
		if err != nil {
			return fmt.Errorf("failed to marshal stats: %w", err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO yearly_stats (player_id, year, sport_type, stats, games_played)
			 VALUES ($1, $2, 'FOOTBALL', $3, 18)`,
			stat.PlayerID, stat.Year, statsJSON)
		if err != nil {
			return fmt.Errorf("failed to insert yearly stats for player %s year %d: %w", stat.PlayerID, stat.Year, err)
		}
	}
	return nil
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// flattenRoster converts a FootballTeamRoster to a flat slice of Players
func flattenRoster(roster FootballTeamRoster) []Player {
	var players []Player
	players = append(players, roster.QB...)
	players = append(players, roster.RB...)
	players = append(players, roster.WR...)
	players = append(players, roster.TE...)
	players = append(players, roster.PK...)
	return players
}

// =============================================================================
// LEGACY API (backward compatible)
// =============================================================================

// SeedDatabase generates synthetic data and inserts it into the database
// All operations are performed in a single transaction (all-or-nothing)
func SeedDatabase(databaseURL string) error {
	ctx := context.Background()

	// Connect to the database
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(ctx)

	// Start a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Will be ignored if tx.Commit() succeeds

	// Use the new DI-based seeder
	seeder := NewDatabaseSeeder(SeederConfig{})
	_, err = seeder.Seed(ctx, tx)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RunSeed is the main entry point for the seed command
func RunSeed() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://fantasy_user:secret_password@localhost:5432/fantasy_db?sslmode=disable"
	}

	log.Println("🌱 Starting database seed...")
	log.Printf("📡 Connecting to: %s\n", maskPassword(databaseURL))

	if err := SeedDatabase(databaseURL); err != nil {
		log.Fatalf("❌ Seed failed: %v", err)
	}
}

// maskPassword hides the password in the connection string for logging
func maskPassword(url string) string {
	// Simple masking - just show structure without sensitive data
	return "postgres://****:****@****:5432/fantasy_db"
}
