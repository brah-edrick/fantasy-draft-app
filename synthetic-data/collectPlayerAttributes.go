package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

/**
 * This file is used to aggregate player stats from real data
 * which was publicly available from an API. This module is a set
 * of functions that will process that data and return a set of
 * aggregated stats that can be used to generate synthetic players.
 * I do not own the source data so it is not included in this repository.
 */

type PlayerNameCount struct {
	Name  string
	Count int
}

type JerseyData struct {
	Jersey int
	Count  int
}

type HeightData struct {
	Height int
	Count  int
}

type WeightData struct {
	Weight int
	Count  int
}

type AgeData struct {
	Age   int
	Count int
}

type YearsOfExperienceData struct {
	YearsOfExperience int
	Count             int
}

type PlayerStat struct {
	FirstName         string
	LastName          string
	Height            int
	Weight            int
	Jersey            int
	Age               int
	Position          string
	YearsOfExperience int
}

type AttributeFrequency map[int]int
type NameFrequency map[string]int

type PositionProfile struct {
	Jerseys           AttributeFrequency `json:"jerseys"`
	Heights           AttributeFrequency `json:"heights"`
	Weights           AttributeFrequency `json:"weights"`
	Ages              AttributeFrequency `json:"ages"`
	YearsOfExperience AttributeFrequency `json:"years_of_experience"`
}

func NewPositionProfile() *PositionProfile {
	return &PositionProfile{
		Jerseys:           make(AttributeFrequency),
		Heights:           make(AttributeFrequency),
		Weights:           make(AttributeFrequency),
		Ages:              make(AttributeFrequency),
		YearsOfExperience: make(AttributeFrequency),
	}
}

func importRealData() map[string]interface{} {
	file, err := os.Open("synthetic-data/real-data.json") // you'll need to provide this file in the format specified
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	return data
}

func collectPlayerAttributes(data map[string]interface{}) []PlayerStat {

	athletes, ok := data["athletes"].([]interface{})
	if !ok {
		log.Fatalf("Error: 'athletes' field is not a list")
	}

	stats := make([]PlayerStat, 0)
	for _, p := range athletes {
		player, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		stat, err := normalizePlayerData(player)
		if err != nil {
			// fmt.Printf("Skipping player: %v\n", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats
}

func aggregateAttributesByPosition(stats []PlayerStat) map[string]*PositionProfile {
	aggregatedStats := make(map[string]*PositionProfile)

	for _, stat := range stats {
		prof, ok := aggregatedStats[stat.Position]
		if !ok {
			prof = NewPositionProfile()
			aggregatedStats[stat.Position] = prof
		}

		prof.Jerseys[stat.Jersey]++
		prof.Heights[stat.Height]++
		prof.Weights[stat.Weight]++
		prof.Ages[stat.Age]++
		prof.YearsOfExperience[stat.YearsOfExperience]++

	}
	return aggregatedStats
}

// aggregateFirstNames returns First Name Counts (Global)
func aggregateFirstNames(stats []PlayerStat) NameFrequency {
	aggregated := make(NameFrequency)
	for _, stat := range stats {
		aggregated[stat.FirstName]++
	}
	return aggregated
}

// aggregateLastNames returns Last Name Counts (Global)
func aggregateLastNames(stats []PlayerStat) NameFrequency {
	aggregated := make(NameFrequency)
	for _, stat := range stats {
		aggregated[stat.LastName]++
	}
	return aggregated
}

type AggregatedPlayerStats struct {
	PositionProfile map[string]*PositionProfile `json:"position_profile"`
	FirstNames      NameFrequency               `json:"first_names"`
	LastNames       NameFrequency               `json:"last_names"`
}

func collectAndAggregatePlayerAttributes() AggregatedPlayerStats {
	data := importRealData()
	stats := collectPlayerAttributes(data)
	return AggregatedPlayerStats{
		PositionProfile: aggregateAttributesByPosition(stats),
		FirstNames:      aggregateFirstNames(stats),
		LastNames:       aggregateLastNames(stats),
	}
}

func normalizePlayerData(data map[string]interface{}) (PlayerStat, error) {
	// Assert the types
	emptyStat := PlayerStat{}

	// Position
	positionMap, ok := data["position"].(map[string]interface{})
	if !ok {
		return emptyStat, fmt.Errorf("missing or invalid position map")
	}
	position, ok := positionMap["abbreviation"].(string)
	if !ok {
		return emptyStat, fmt.Errorf("missing or invalid position abbreviation")
	}

	// Status
	statusMap, ok := data["status"].(map[string]interface{})
	if !ok {
		return emptyStat, fmt.Errorf("missing or invalid status map")
	}
	status, ok := statusMap["type"].(string)
	if !ok {
		return emptyStat, fmt.Errorf("missing or invalid status type")
	}
	// Skip free agents because they may not be good enough and will skew our data
	if status == "free-agent" {
		return emptyStat, fmt.Errorf("skip: player is free-agent")
	}

	// Draft information (used to get the years of experience)
	draftMap, ok := data["draft"].(map[string]interface{})
	if !ok {
		return emptyStat, fmt.Errorf("missing or invalid draft map")
	}
	var draftYear int
	if dYearVal, ok := draftMap["year"].(float64); ok {
		draftYear = int(dYearVal)
	} else if dYearVal, ok := draftMap["year"].(int); ok {
		draftYear = dYearVal
	} else {
		return emptyStat, fmt.Errorf("missing or invalid draft year")
	}
	thisYear := time.Now().Year()
	yearsOfExperience := thisYear - draftYear

	// Safely assert other fields
	firstName, _ := data["firstName"].(string)
	lastName, _ := data["lastName"].(string)

	var height, weight, jersey, age int

	if h, ok := data["height"].(float64); ok {
		height = int(h)
	}
	if w, ok := data["weight"].(float64); ok {
		weight = int(w)
	}
	if j, ok := data["jersey"].(string); ok {
		if val, err := strconv.Atoi(j); err == nil {
			jersey = val
		}
	}
	if a, ok := data["age"].(float64); ok {
		age = int(a)
	}

	playerStat := PlayerStat{
		FirstName:         firstName,
		LastName:          lastName,
		Height:            height,
		Weight:            weight,
		Jersey:            jersey,
		Age:               age,
		Position:          position,
		YearsOfExperience: yearsOfExperience,
	}

	return playerStat, nil
}
