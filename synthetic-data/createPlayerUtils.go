package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/google/uuid"
)

// Distribution maps a value (T) to its frequency count.
// T must be 'ordered' (int, string, float64) to be sorted for CDF.
type Distribution[T cmp.Ordered] map[T]int

func importPlayerStats() AggregatedPlayerStats {
	file, err := os.Open("synthetic-data/.output/player_stats.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var playerStats AggregatedPlayerStats
	if err := json.NewDecoder(file).Decode(&playerStats); err != nil {
		log.Fatal(err)
	}

	return playerStats
}

func createNewPlayer(position Position, teamId string) Player {
	firstNameGenerator, lastNameGenerator, positionGenerators := createPlayerGeneratorsFromStats()

	firstName := firstNameGenerator()
	lastName := lastNameGenerator()
	positionIndex := slices.IndexFunc(positionGenerators, func(p LabeledPositionGenerators) bool {
		return p.PositionCode == position
	})
	if positionIndex == -1 {
		log.Fatalf("Error: Position %s not found in positionGenerators", position)
	}
	jersey := positionGenerators[positionIndex].Generators.JerseyGenerator()
	height := positionGenerators[positionIndex].Generators.HeightGenerator()
	weight := positionGenerators[positionIndex].Generators.WeightGenerator()
	age := positionGenerators[positionIndex].Generators.AgeGenerator()
	yoe := positionGenerators[positionIndex].Generators.YoeGenerator()

	thisYear := time.Now().Year()
	player := Player{
		ID:                uuid.New().String(),
		DraftYear:         thisYear - yoe,
		FirstName:         firstName,
		LastName:          lastName,
		Position:          string(position),
		Jersey:            jersey,
		Height:            height,
		Weight:            weight,
		Age:               age,
		YearsOfExperience: yoe,
		Status:            "Active",
		Skill:             createRandomSkillFactorWithBellCurve(),
		TeamID:            teamId,
	}

	fmt.Printf("Player created: %+v\n", player)

	return player
}

type Position string

const (
	QB Position = "QB"
	RB Position = "RB"
	WR Position = "WR"
	TE Position = "TE"
	PK Position = "PK"
)

type PositionGenerators struct {
	JerseyGenerator func() int
	HeightGenerator func() int
	WeightGenerator func() int
	AgeGenerator    func() int
	YoeGenerator    func() int
}

type LabeledPositionGenerators struct {
	PositionCode Position
	Generators   PositionGenerators
}

func createPlayerGeneratorsFromStats() (func() string, func() string, []LabeledPositionGenerators) {
	stats := importPlayerStats()
	firstNameGenerator := createGenerateValueFromStat(stats.FirstNames)
	lastNameGenerator := createGenerateValueFromStat(stats.LastNames)
	positionGenerators := createPositionsGeneratorsFromStats(stats)
	return firstNameGenerator, lastNameGenerator, positionGenerators
}

func createPositionsGeneratorsFromStats(stats AggregatedPlayerStats) []LabeledPositionGenerators {
	positionCodes := make([]Position, 0, 5)
	positionCodes = append(positionCodes, QB, RB, WR, TE, PK)
	positionGenerators := make([]LabeledPositionGenerators, 0, 5)
	for _, positionCode := range positionCodes {
		positionMap, ok := stats.PositionProfile[string(positionCode)]
		if !ok {
			log.Fatalf("Error: 'position_profile' field is not a map")
		}
		positionGenerators = append(positionGenerators, LabeledPositionGenerators{
			PositionCode: positionCode,
			Generators:   CreatePositionAttributeGenerators(positionMap),
		})
	}
	return positionGenerators
}

// CreatePositionAttributeGenerators creates generators for all standard position attributes
func CreatePositionAttributeGenerators(profile *PositionProfile) PositionGenerators {
	return PositionGenerators{
		JerseyGenerator: createGenerateValueFromStat(profile.Jerseys),
		HeightGenerator: createGenerateValueFromStat(profile.Heights),
		WeightGenerator: createGenerateValueFromStat(profile.Weights),
		AgeGenerator:    createGenerateValueFromStat(profile.Ages),
		YoeGenerator:    createGenerateValueFromStat(profile.YearsOfExperience),
	}
}

type StatisticToCdf[T cmp.Ordered] struct {
	Values []T
	CDF    []float64
}

// createCdfForStat calculates the Cumulative Distribution Function for a given statistic distribution.
// It returns a struct containing sorted Values and their corresponding CDF probabilities.
// This generic function accepts any map with comparable/ordered keys (int, string, etc.) and int values (counts).
func createCdfForStat[T cmp.Ordered, M ~map[T]int](stats M) *StatisticToCdf[T] {
	// Convert to array
	keys := make([]T, 0, len(stats))
	total := 0
	for k, v := range stats {
		keys = append(keys, k)
		total += v
	}

	// Sort keys to ensure deterministic CDF order
	slices.Sort(keys)

	// Calculate CDF for each value
	cdf := make([]float64, len(keys))
	runningSum := 0
	for i, k := range keys {
		runningSum += stats[k]
		cdf[i] = float64(runningSum) / float64(total)
	}

	// Return CDF entries
	return &StatisticToCdf[T]{
		Values: keys,
		CDF:    cdf,
	}
}

func generateValueFromCdf[T cmp.Ordered](cdf *StatisticToCdf[T]) T {
	randomNumber := rand.Float64()
	for i, cdfValue := range cdf.CDF {
		if randomNumber < cdfValue {
			return cdf.Values[i]
		}
	}
	return cdf.Values[len(cdf.Values)-1]
}

func createGenerateValueFromStat[T cmp.Ordered, M ~map[T]int](stats M) func() T {
	cdf := createCdfForStat(stats)
	return func() T {
		return generateValueFromCdf(cdf)
	}
}

func createRandomSkillFactorWithBellCurve() float64 {
	// Generate random number from bell curve
	// Return value between 0.0 and 1.0
	desiredMean := 0.5
	desiredStdDev := 0.2
	return rand.NormFloat64()*desiredStdDev + desiredMean
}
