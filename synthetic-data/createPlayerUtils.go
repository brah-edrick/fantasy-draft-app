package main

import (
	"cmp"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"sync"
)

// Distribution maps a value (T) to its frequency count.
// T must be 'ordered' (int, string, float64) to be sorted for CDF.
type Distribution[T cmp.Ordered] map[T]int

func createNewPlayer(position Position, teamId string, generators PlayerGenerators, clock Clock, uuidGenerator UUIDGenerator) Player {
	firstName := generators.FirstNameGenerator()
	lastName := generators.LastNameGenerator()
	positionGenerators := generators.PositionGenerators
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
	thisYear := clock.Now().Year()

	player := Player{
		ID:                uuidGenerator(),
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
		Skill:             generators.SkillGenerator(),
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

type PlayerGenerators struct {
	FirstNameGenerator func() string
	LastNameGenerator  func() string
	PositionGenerators []LabeledPositionGenerators
	SkillGenerator     func() float64
}

type UUIDGenerator func() string

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

func getPlayerGenerators(statsAggregator StatsAggregator, rand *rand.Rand) PlayerGenerators {
	generatorsOnce.Do(func() {
		firstNameGeneratorSingleton, lastNameGeneratorSingleton, positionGeneratorsSingleton = createPlayerGeneratorsFromStats(statsAggregator, rand)
	})
	return PlayerGenerators{
		FirstNameGenerator: firstNameGeneratorSingleton,
		LastNameGenerator:  lastNameGeneratorSingleton,
		PositionGenerators: positionGeneratorsSingleton,
		SkillGenerator:     createRandomSkillFactorWithBellCurve,
	}
}

func createPlayerGeneratorsFromStats(statsAggregator StatsAggregator, rand *rand.Rand) (func() string, func() string, []LabeledPositionGenerators) {
	fmt.Println("Creating player generators from real player stats...")
	fmt.Println("Aggregating player stats...")
	stats := statsAggregator()
	fmt.Println("Creating first name generator...")
	firstNameGenerator := createGenerateValueFromStat(stats.FirstNames, rand)
	fmt.Println("Creating last name generator...")
	lastNameGenerator := createGenerateValueFromStat(stats.LastNames, rand)
	fmt.Println("Creating position generators...")
	positionGenerators := createPositionsGeneratorsFromStats(stats, rand)
	fmt.Println("Player generators created successfully.")
	return firstNameGenerator, lastNameGenerator, positionGenerators
}

func createPositionsGeneratorsFromStats(stats AggregatedPlayerStats, rand *rand.Rand) []LabeledPositionGenerators {
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
			Generators:   CreatePositionAttributeGenerators(positionMap, rand),
		})
	}
	return positionGenerators
}

// CreatePositionAttributeGenerators creates generators for all standard position attributes
func CreatePositionAttributeGenerators(profile *PositionProfile, rand *rand.Rand) PositionGenerators {
	return PositionGenerators{
		JerseyGenerator: createGenerateValueFromStat(profile.Jerseys, rand),
		HeightGenerator: createGenerateValueFromStat(profile.Heights, rand),
		WeightGenerator: createGenerateValueFromStat(profile.Weights, rand),
		AgeGenerator:    createGenerateValueFromStat(profile.Ages, rand),
		YoeGenerator:    createGenerateValueFromStat(profile.YearsOfExperience, rand),
	}
}

type StatisticToCDF[T cmp.Ordered] struct {
	Values []T
	CDF    []float64
}

// createCdfForStat calculates the Cumulative Distribution Function for a given statistic distribution.
// It returns a struct containing sorted Values and their corresponding CDF probabilities.
// This generic function accepts any map with comparable/ordered keys (int, string, etc.) and int values (counts).
func createCDFForStat[T cmp.Ordered, M ~map[T]int](stats M, rand *rand.Rand) *StatisticToCDF[T] {
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

func generateValueFromCDF[T cmp.Ordered](cdf *StatisticToCDF[T], rand *rand.Rand) T {
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

func createGenerateValueFromStat[T cmp.Ordered, M ~map[T]int](stats M, rand *rand.Rand) func() T {
	cdf := createCDFForStat(stats, rand)
	return func() T {
		return generateValueFromCDF(cdf, rand)
	}
}

func createRandomSkillFactorWithBellCurve() float64 {
	// Generate random number from bell curve
	// Return value between 0.0 and 1.0
	desiredMean := 0.5
	desiredStdDev := 0.2
	return rand.NormFloat64()*desiredStdDev + desiredMean
}

// createSkillForDepthPosition generates a skill value based on depth chart position.
// depthPosition is 0-indexed (0 = starter, 1 = backup, etc.)
// This creates a natural falloff down the depth chart while allowing some variance.
func createSkillForDepthPosition(depthPosition int, totalAtPosition int) float64 {
	// Calculate a base skill that decreases with depth
	depthRatio := float64(depthPosition) / float64(max(totalAtPosition-1, 1))

	baseMean := 0.80 - (depthRatio * 0.45)

	variance := 0.08

	skill := rand.NormFloat64()*variance + baseMean

	if skill < 0.15 {
		skill = 0.15
	}
	if skill > 0.95 {
		skill = 0.95
	}

	return skill
}
