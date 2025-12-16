package main

// RosterComposition dictates how a team is built.
// Instead of random chance, we force a specific structure.
// Key = Position (e.g., "QB"), Value = Quantity (e.g., 3).
type RosterComposition map[string]int

// NFLRosterComposition is our standard definition for a valid team.
var NFLRosterComposition = RosterComposition{
	"QB": 3,
	"RB": 4,
	"WR": 6,
	"TE": 3,
	"PK": 1,
}

type FootballTeamRoster struct {
	QB []Player
	RB []Player
	WR []Player
	TE []Player
	PK []Player
}

// --- Data Model Structs ---

type League struct {
	Conferences []Conference `json:"conferences"`
}

type Conference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Division struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ConferenceID string `json:"conference_id"`
}

type Team struct {
	ID         string `json:"id"`
	City       string `json:"city"`
	State      string `json:"state"` // Added per request
	Name       string `json:"name"`
	Abbr       string `json:"abbr"`
	DivisionID string `json:"division_id"`
}

type Player struct {
	ID                string  `json:"id"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Position          string  `json:"position"`
	TeamID            string  `json:"team_id"`
	Height            int     `json:"height"`
	Weight            int     `json:"weight"`
	Age               int     `json:"age"`
	YearsOfExperience int     `json:"years_of_experience"`
	DraftYear         int     `json:"draft_year"`
	Skill             float64 `json:"skill"` // 0.0 - 1.0
	Status            string  `json:"status"`
	Jersey            int     `json:"jersey"`
}

type FootballStats struct {
	PassingAttempts       int
	PassingCompletions    int
	PassingInterceptions  int
	PassingTDs            int
	PassingYards          int
	RushingAttempts       int
	RushingYards          int
	ReceivingYards        int
	RushingTDs            int
	ReceivingReceptions   int
	ReceivingTDs          int
	ReceivingTargets      int
	Fumbles               int
	FumblesLost           int
	FieldGoals            int
	FieldGoalsMade        int
	FieldGoalsMissed      int
	FieldGoalsBlocked     int
	FieldGoalsBlockedMade int
	ExtraPoints           int
	ExtraPointsMade       int
	ExtraPointsMissed     int
}

type FootballYearlyStats struct {
	Total FootballStats
}

type PlayerYearlyStats[T struct{}] struct {
	ID       string `json:"id"`
	PlayerID string `json:"player_id"`
	Year     int    `json:"year"`
	Stats    T      `json:"stats"`
}

type PlayerYearlyStatsFootball struct {
	PlayerID string              `json:"player_id"`
	Year     int                 `json:"year"`
	Stats    FootballYearlyStats `json:"stats"`
}
