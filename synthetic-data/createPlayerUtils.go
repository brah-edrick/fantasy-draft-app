package main

import (
	"cmp"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Distribution maps a value (T) to its frequency count.
// T must be 'ordered' (int, string, float64) to be sorted for CDF.
type Distribution[T cmp.Ordered] map[T]int

func createNewPlayer(position Position, teamId string) Player {
	firstNameGenerator, lastNameGenerator, positionGenerators := getPlayerGenerators()

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

var (
	firstNameGeneratorSingleton func() string
	lastNameGeneratorSingleton  func() string
	positionGeneratorsSingleton []LabeledPositionGenerators
	generatorsOnce              sync.Once
)

func getPlayerGenerators() (func() string, func() string, []LabeledPositionGenerators) {
	generatorsOnce.Do(func() {
		firstNameGeneratorSingleton, lastNameGeneratorSingleton, positionGeneratorsSingleton = createPlayerGeneratorsFromStats()
	})
	return firstNameGeneratorSingleton, lastNameGeneratorSingleton, positionGeneratorsSingleton
}

func createPlayerGeneratorsFromStats() (func() string, func() string, []LabeledPositionGenerators) {
	fmt.Println("Creating player generators from real player stats...")
	fmt.Println("Aggregating player stats...")
	stats := collectAndAggregatePlayerAttributes()
	fmt.Println("Creating first name generator...")
	firstNameGenerator := createGenerateValueFromStat(stats.FirstNames)
	fmt.Println("Creating last name generator...")
	lastNameGenerator := createGenerateValueFromStat(stats.LastNames)
	fmt.Println("Creating position generators...")
	positionGenerators := createPositionsGeneratorsFromStats(stats)
	fmt.Println("Player generators created successfully.")
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

type StatisticToCDF[T cmp.Ordered] struct {
	Values []T
	CDF    []float64
}

// createCdfForStat calculates the Cumulative Distribution Function for a given statistic distribution.
// It returns a struct containing sorted Values and their corresponding CDF probabilities.
// This generic function accepts any map with comparable/ordered keys (int, string, etc.) and int values (counts).
func createCDFForStat[T cmp.Ordered, M ~map[T]int](stats M) *StatisticToCDF[T] {
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
	return &StatisticToCDF[T]{
		Values: keys,
		CDF:    cdf,
	}
}

func generateValueFromCDF[T cmp.Ordered](cdf *StatisticToCDF[T]) T {
	randomNumber := rand.Float64()
	index := binarySearchUpperBound(cdf, 0, len(cdf.Values)-1, randomNumber)
	return cdf.Values[index]
}

// binarySearchUpperBound returns the index of the first element in the CDF that is greater than or equal to the target value
func binarySearchUpperBound[T cmp.Ordered](cdf *StatisticToCDF[T], left, right int, target float64) int {
	if left == right {
		return left
	}

	midIndex := (left + right) / 2

	if cdf.CDF[midIndex] < target {
		return binarySearchUpperBound(cdf, midIndex+1, right, target)
	} else {
		return binarySearchUpperBound(cdf, left, midIndex, target)
	}
}

func createGenerateValueFromStat[T cmp.Ordered, M ~map[T]int](stats M) func() T {
	cdf := createCDFForStat(stats)
	return func() T {
		return generateValueFromCDF(cdf)
	}
}

func createRandomSkillFactorWithBellCurve() float64 {
	// Generate random number from bell curve
	// Return value between 0.0 and 1.0
	desiredMean := 0.5
	desiredStdDev := 0.2
	return rand.NormFloat64()*desiredStdDev + desiredMean
}
