package main

// types.go defines the data structures we use during generation.
// We use the 'main' package here so it can share visibility with our executable script.

// --- Generator Config Types ---

// GeneratorSource holds the raw ingredients we need to cook up a player.
type GeneratorSource struct {
	FirstNames []string
	LastNames  []string
	// We will eventually add "AgeDistribution" maps here
}

// RosterComposition dictates how a team is built.
// Instead of random chance, we force a specific structure.
// Key = Position (e.g., "QB"), Value = Quantity (e.g., 3).
type RosterComposition map[string]int

// NFLRosterComposition is our standard definition for a valid team.
var NFLRosterComposition = RosterComposition{
	"QB":  3,
	"RB":  4,
	"WR":  6,
	"TE":  3,
	"K":   1,
	"DST": 1,
}

// --- Data Model Structs (for JSON Export) ---

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
}
