package main

import (
	"fmt"
	"math/rand"
)

func createPlayerYear(player Player, year int) PlayerYearlyStatsFootball {
	playerYear := PlayerYearlyStatsFootball{
		PlayerID: player.ID,
		Year:     year,
		Stats:    walkThroughPlayerYear(player, 18, year),
	}
	// pretty print the player year
	fmt.Printf("Player Year: %+v\n", playerYear)
	return playerYear
}

func walkThroughPlayerYear(player Player, gamesInYear int, year int) FootballYearlyStats {
	playerYearsOfExperience := player.DraftYear - year
	_ = FootballYearlyStats{Total: FootballStats{}} // placeholder
	isInjured := false
	injuryGameCount := 0
	yearlyStats := FootballStats{
		PassingAttempts:      0,
		PassingCompletions:   0,
		PassingInterceptions: 0,
		PassingTDs:           0,
		PassingYards:         0,
		RushingAttempts:      0,
		RushingYards:         0,
		ReceivingYards:       0,
		RushingTDs:           0,
		ReceivingReceptions:  0,
		ReceivingTDs:         0,
		ReceivingTargets:     0,
		Fumbles:              0,
		FumblesLost:          0,
	}
	for range gamesInYear {
		if isInjured {
			injuryGameCount--
			if injuryGameCount <= 0 {
				isInjured = false
			}
			continue
		}
		wasInjured, injuryGamesAffected := rollForInjury(player.Age, player.Position)

		if wasInjured {
			isInjured = true
			injuryGameCount = injuryGamesAffected
		}

		gameStats := multiplyYearlyStatsByPlayerSkill(player, playerYearsOfExperience, generatePlayerGameStats(player, playerYearsOfExperience))

		yearlyStats.PassingAttempts += gameStats.PassingAttempts
		yearlyStats.PassingCompletions += gameStats.PassingCompletions
		yearlyStats.PassingInterceptions += gameStats.PassingInterceptions
		yearlyStats.PassingTDs += gameStats.PassingTDs
		yearlyStats.PassingYards += gameStats.PassingYards
		yearlyStats.RushingAttempts += gameStats.RushingAttempts
		yearlyStats.RushingYards += gameStats.RushingYards
		yearlyStats.ReceivingYards += gameStats.ReceivingYards
		yearlyStats.RushingTDs += gameStats.RushingTDs
		yearlyStats.ReceivingReceptions += gameStats.ReceivingReceptions
		yearlyStats.ReceivingTDs += gameStats.ReceivingTDs
		yearlyStats.ReceivingTargets += gameStats.ReceivingTargets
		yearlyStats.Fumbles += gameStats.Fumbles
		yearlyStats.FumblesLost += gameStats.FumblesLost
	}
	return FootballYearlyStats{Total: yearlyStats}
}

func rollForInjury(playerAge int, playerPosition string) (bool, int) {
	injuryRate := 0.0
	if playerAge < 25 {
		injuryRate = 0.04
	} else if playerAge < 30 {
		injuryRate = 0.06
	} else if playerAge < 35 {
		injuryRate = 0.10
	} else {
		injuryRate = 0.12
	}

	if playerPosition == "QB" {
		injuryRate = injuryRate * 0.5
	} else if playerPosition == "RB" {
		injuryRate = injuryRate * 1
	} else if playerPosition == "WR" {
		injuryRate = injuryRate * 0.8
	} else if playerPosition == "TE" {
		injuryRate = injuryRate * 0.8
	} else if playerPosition == "PK" {
		injuryRate = injuryRate * 0.25
	}

	wasInjured := rand.Float64() < injuryRate

	injuryGameCount := 0
	if wasInjured {
		injuryGameCount = normalIntInRange(1, 20)
	}

	return wasInjured, injuryGameCount
}

func generatePlayerGameStats(player Player, yearsOfExperience int) FootballStats {
	switch player.Position {
	case "QB":
		return QuarterBackGameStatsGenerator().generate(player, yearsOfExperience)
	case "RB":
		return RunningBackGameStatsGenerator().generate(player, yearsOfExperience)
	case "WR":
		return WideReceiverGameStatsGenerator().generate(player, yearsOfExperience)
	case "TE":
		return TightEndGameStatsGenerator().generate(player, yearsOfExperience)
	case "PK":
		return KickerGameStatsGenerator().generate(player, yearsOfExperience)
	default:
		return FootballStats{}
	}
}

type PlayerGameStatsGenerator interface {
	generate(player Player, yearsOfExperience int) FootballStats
}

// quarterBackGenerator is a concrete type that implements PlayerGameStatsGenerator
type quarterBackGenerator struct{}

func (q quarterBackGenerator) generate(player Player, yearsOfExperience int) FootballStats {
	passingTouchdowns := normalIntInRange(0, 4)
	passingInterceptions := normalIntInRange(0, 2)
	passingAttempts := normalIntInRange(25, 45)
	passingCompletions := normalIntInRange(15, 32)
	passingAverage := normalIntInRange(8, 14)
	passingYards := passingCompletions * passingAverage
	rushingAttempts := normalIntInRange(1, 6)
	rushingYards := normalIntInRange(5, 35)
	rushingTDs := normalIntInRange(0, 1)
	fumbles := normalIntInRange(0, 1)
	fumblesLost := normalIntInRange(0, fumbles)

	return FootballStats{
		PassingAttempts:       passingAttempts,
		PassingCompletions:    passingCompletions,
		PassingInterceptions:  passingInterceptions,
		PassingTDs:            passingTouchdowns,
		PassingYards:          passingYards,
		RushingAttempts:       rushingAttempts,
		RushingYards:          rushingYards,
		RushingTDs:            rushingTDs,
		ReceivingReceptions:   0,
		ReceivingTDs:          0,
		ReceivingTargets:      0,
		Fumbles:               fumbles,
		FumblesLost:           fumblesLost,
		FieldGoals:            0,
		FieldGoalsMade:        0,
		FieldGoalsMissed:      0,
		FieldGoalsBlocked:     0,
		FieldGoalsBlockedMade: 0,
		ExtraPoints:           0,
		ExtraPointsMade:       0,
		ExtraPointsMissed:     0,
	}
}

func QuarterBackGameStatsGenerator() PlayerGameStatsGenerator {
	return quarterBackGenerator{}
}

type runningBackGenerator struct{}

func (r runningBackGenerator) generate(player Player, yearsOfExperience int) FootballStats {
	rushingAttempts := normalIntInRange(12, 25)
	rushingAverage := normalIntInRange(4, 6)
	rushingYards := rushingAttempts * rushingAverage
	rushingTDs := normalIntInRange(0, 2)
	fumbles := normalIntInRange(0, 1)
	fumblesLost := normalIntInRange(0, fumbles)
	receivingReceptions := normalIntInRange(2, 6)
	receivingTargets := normalIntInRange(3, 8)
	receivingAverage := normalIntInRange(6, 12)
	receivingYards := receivingReceptions * receivingAverage
	receivingTDs := normalIntInRange(0, 1)

	return FootballStats{
		PassingAttempts:       0,
		PassingCompletions:    0,
		PassingInterceptions:  0,
		PassingTDs:            0,
		PassingYards:          0,
		RushingAttempts:       rushingAttempts,
		RushingYards:          rushingYards,
		RushingTDs:            rushingTDs,
		ReceivingReceptions:   receivingReceptions,
		ReceivingTDs:          receivingTDs,
		ReceivingTargets:      receivingTargets,
		ReceivingYards:        receivingYards,
		Fumbles:               fumbles,
		FumblesLost:           fumblesLost,
		FieldGoals:            0,
		FieldGoalsMade:        0,
		FieldGoalsMissed:      0,
		FieldGoalsBlocked:     0,
		FieldGoalsBlockedMade: 0,
		ExtraPoints:           0,
		ExtraPointsMade:       0,
		ExtraPointsMissed:     0,
	}
}

func RunningBackGameStatsGenerator() PlayerGameStatsGenerator {
	return runningBackGenerator{}
}

type wideReceiverGenerator struct{}

func (w wideReceiverGenerator) generate(player Player, yearsOfExperience int) FootballStats {
	receivingReceptions := normalIntInRange(4, 10)
	receivingTargets := normalIntInRange(6, 14)
	receivingAverage := normalIntInRange(12, 18)
	receivingYards := receivingReceptions * receivingAverage
	rushingAttempts := normalIntInRange(0, 2)
	rushingAverage := normalIntInRange(5, 14)
	rushingYards := rushingAttempts * rushingAverage
	rushingTDs := normalIntInRange(0, 1)
	receivingTDs := normalIntInRange(0, 2)
	fumbles := normalIntInRange(0, 1)
	fumblesLost := normalIntInRange(0, fumbles)

	return FootballStats{
		PassingAttempts:       0,
		PassingCompletions:    0,
		PassingInterceptions:  0,
		PassingTDs:            0,
		PassingYards:          0,
		RushingAttempts:       rushingAttempts,
		RushingYards:          rushingYards,
		RushingTDs:            rushingTDs,
		ReceivingReceptions:   receivingReceptions,
		ReceivingTDs:          receivingTDs,
		ReceivingTargets:      receivingTargets,
		ReceivingYards:        receivingYards,
		Fumbles:               fumbles,
		FumblesLost:           fumblesLost,
		FieldGoals:            0,
		FieldGoalsMade:        0,
		FieldGoalsMissed:      0,
		FieldGoalsBlocked:     0,
		FieldGoalsBlockedMade: 0,
		ExtraPoints:           0,
		ExtraPointsMade:       0,
		ExtraPointsMissed:     0,
	}
}

func WideReceiverGameStatsGenerator() PlayerGameStatsGenerator {
	return wideReceiverGenerator{}
}

type tightEndGenerator struct{}

func (te tightEndGenerator) generate(player Player, yearsOfExperience int) FootballStats {
	receivingReceptions := normalIntInRange(3, 8)
	receivingTargets := normalIntInRange(5, 11)
	receivingAverage := normalIntInRange(10, 14)
	receivingYards := receivingReceptions * receivingAverage
	rushingAttempts := normalIntInRange(0, 1)
	rushingAverage := normalIntInRange(4, 10)
	rushingYards := rushingAttempts * rushingAverage
	rushingTDs := normalIntInRange(0, 1)
	receivingTDs := normalIntInRange(0, 1)
	fumbles := normalIntInRange(0, 1)
	fumblesLost := normalIntInRange(0, fumbles)

	return FootballStats{
		PassingAttempts:       0,
		PassingCompletions:    0,
		PassingInterceptions:  0,
		PassingTDs:            0,
		PassingYards:          0,
		RushingAttempts:       rushingAttempts,
		RushingYards:          rushingYards,
		RushingTDs:            rushingTDs,
		ReceivingReceptions:   receivingReceptions,
		ReceivingTDs:          receivingTDs,
		ReceivingTargets:      receivingTargets,
		ReceivingYards:        receivingYards,
		Fumbles:               fumbles,
		FumblesLost:           fumblesLost,
		FieldGoals:            0,
		FieldGoalsMade:        0,
		FieldGoalsMissed:      0,
		FieldGoalsBlocked:     0,
		FieldGoalsBlockedMade: 0,
		ExtraPoints:           0,
		ExtraPointsMade:       0,
		ExtraPointsMissed:     0,
	}
}

func TightEndGameStatsGenerator() PlayerGameStatsGenerator {
	return tightEndGenerator{}
}

type kickerGenerator struct{}

func (k kickerGenerator) generate(player Player, yearsOfExperience int) FootballStats {
	fieldGoals := normalIntInRange(0, 50)
	fieldGoalsMade := normalIntInRange(0, fieldGoals)
	fieldGoalsMissed := fieldGoals - fieldGoalsMade
	fieldGoalsBlocked := normalIntInRange(0, 5)
	fieldGoalsBlockedMade := normalIntInRange(0, fieldGoalsBlocked)
	extraPoints := normalIntInRange(0, 2)
	extraPointsMade := normalIntInRange(0, extraPoints)
	extraPointsMissed := extraPoints - extraPointsMade

	return FootballStats{
		PassingAttempts:       0,
		PassingCompletions:    0,
		PassingInterceptions:  0,
		PassingTDs:            0,
		PassingYards:          0,
		RushingAttempts:       0,
		RushingYards:          0,
		RushingTDs:            0,
		ReceivingReceptions:   0,
		ReceivingTDs:          0,
		ReceivingTargets:      0,
		ReceivingYards:        0,
		Fumbles:               0,
		FumblesLost:           0,
		FieldGoals:            fieldGoals,
		FieldGoalsMade:        fieldGoalsMade,
		FieldGoalsMissed:      fieldGoalsMissed,
		FieldGoalsBlocked:     fieldGoalsBlocked,
		FieldGoalsBlockedMade: fieldGoalsBlockedMade,
		ExtraPoints:           extraPoints,
		ExtraPointsMade:       extraPointsMade,
		ExtraPointsMissed:     extraPointsMissed,
	}
}

func KickerGameStatsGenerator() PlayerGameStatsGenerator {
	return kickerGenerator{}
}

func multiplyStatByPlayerSkill(player Player, yearsOfExperience int, stat int) int {
	return int(float64(stat) * (1 + float64(yearsOfExperience)/100) * player.Skill)
}

func multiplyYearlyStatsByPlayerSkill(player Player, yearsofExperience int, stats FootballStats) FootballStats {
	adjustedStats := FootballStats{
		PassingAttempts:       multiplyStatByPlayerSkill(player, yearsofExperience, stats.PassingAttempts),
		PassingCompletions:    multiplyStatByPlayerSkill(player, yearsofExperience, stats.PassingCompletions),
		PassingInterceptions:  multiplyStatByPlayerSkill(player, yearsofExperience, stats.PassingInterceptions),
		PassingTDs:            multiplyStatByPlayerSkill(player, yearsofExperience, stats.PassingTDs),
		PassingYards:          multiplyStatByPlayerSkill(player, yearsofExperience, stats.PassingYards),
		RushingAttempts:       multiplyStatByPlayerSkill(player, yearsofExperience, stats.RushingAttempts),
		RushingYards:          multiplyStatByPlayerSkill(player, yearsofExperience, stats.RushingYards),
		RushingTDs:            multiplyStatByPlayerSkill(player, yearsofExperience, stats.RushingTDs),
		ReceivingReceptions:   multiplyStatByPlayerSkill(player, yearsofExperience, stats.ReceivingReceptions),
		ReceivingTDs:          multiplyStatByPlayerSkill(player, yearsofExperience, stats.ReceivingTDs),
		ReceivingTargets:      multiplyStatByPlayerSkill(player, yearsofExperience, stats.ReceivingTargets),
		ReceivingYards:        multiplyStatByPlayerSkill(player, yearsofExperience, stats.ReceivingYards),
		Fumbles:               multiplyStatByPlayerSkill(player, yearsofExperience, stats.Fumbles),
		FumblesLost:           multiplyStatByPlayerSkill(player, yearsofExperience, stats.FumblesLost),
		FieldGoals:            multiplyStatByPlayerSkill(player, yearsofExperience, stats.FieldGoals),
		FieldGoalsMade:        multiplyStatByPlayerSkill(player, yearsofExperience, stats.FieldGoalsMade),
		FieldGoalsMissed:      multiplyStatByPlayerSkill(player, yearsofExperience, stats.FieldGoalsMissed),
		FieldGoalsBlocked:     multiplyStatByPlayerSkill(player, yearsofExperience, stats.FieldGoalsBlocked),
		FieldGoalsBlockedMade: multiplyStatByPlayerSkill(player, yearsofExperience, stats.FieldGoalsBlockedMade),
		ExtraPoints:           multiplyStatByPlayerSkill(player, yearsofExperience, stats.ExtraPoints),
		ExtraPointsMade:       multiplyStatByPlayerSkill(player, yearsofExperience, stats.ExtraPointsMade),
		ExtraPointsMissed:     multiplyStatByPlayerSkill(player, yearsofExperience, stats.ExtraPointsMissed),
	}
	return adjustedStats
}

// stats utils

func normalInRange(low, high float64) float64 {
	mean := (low + high) / 2
	// Use 3 standard deviations to cover the range (99.7% of values)
	stdDev := (high - low) / 6

	result := rand.NormFloat64()*stdDev + mean

	// Clamp to bounds for the rare outliers beyond 3 sigma
	if result < low {
		result = low
	}
	if result > high {
		result = high
	}

	return result
}

func normalIntInRange(low, high int) int {
	return int(normalInRange(float64(low), float64(high)+0.5))
}
