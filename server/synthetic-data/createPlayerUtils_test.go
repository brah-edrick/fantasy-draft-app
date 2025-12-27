package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestCreateNewPlayer(t *testing.T) {
	counter := 0
	uuidGen := mockUUIDGenerator("player-", &counter)
	mockClock := MockClock{mockTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	teamID := "team-123"

	// Create mock generators
	generators := PlayerGenerators{
		FirstNameGenerator: func() string { return "John" },
		LastNameGenerator:  func() string { return "Doe" },
		PositionGenerators: []LabeledPositionGenerators{
			{
				PositionCode: QB,
				Generators: PositionGenerators{
					JerseyGenerator: func() int { return 12 },
					HeightGenerator: func() int { return 72 },
					WeightGenerator: func() int { return 200 },
					AgeGenerator:    func() int { return 25 },
					YoeGenerator:    func() int { return 3 },
				},
			},
		},
		SkillGenerator: func() float64 { return 0.75 },
	}

	player := createNewPlayer(QB, teamID, generators, mockClock, uuidGen)

	if player.FirstName != "John" {
		t.Errorf("Expected first name 'John', got '%s'", player.FirstName)
	}
	if player.LastName != "Doe" {
		t.Errorf("Expected last name 'Doe', got '%s'", player.LastName)
	}
	if player.Position != "QB" {
		t.Errorf("Expected position 'QB', got '%s'", player.Position)
	}
	if player.Jersey != 12 {
		t.Errorf("Expected jersey 12, got %d", player.Jersey)
	}
	if player.Height != 72 {
		t.Errorf("Expected height 72, got %d", player.Height)
	}
	if player.Weight != 200 {
		t.Errorf("Expected weight 200, got %d", player.Weight)
	}
	if player.Age != 25 {
		t.Errorf("Expected age 25, got %d", player.Age)
	}
	if player.YearsOfExperience != 3 {
		t.Errorf("Expected years of experience 3, got %d", player.YearsOfExperience)
	}
	if player.DraftYear != 2022 { // 2025 - 3
		t.Errorf("Expected draft year 2022, got %d", player.DraftYear)
	}
	if player.Status != "ACTIVE" {
		t.Errorf("Expected status 'ACTIVE', got '%s'", player.Status)
	}
	if player.Skill != 0.75 {
		t.Errorf("Expected skill 0.75, got %f", player.Skill)
	}
	if player.TeamID != teamID {
		t.Errorf("Expected team ID '%s', got '%s'", teamID, player.TeamID)
	}
	if player.ID == "" {
		t.Error("Player ID should not be empty")
	}
}

func TestCreateCDFForStat(t *testing.T) {
	stats := map[int]int{
		1: 10,
		2: 20,
		3: 30,
		4: 40,
	}
	rng := rand.New(rand.NewSource(12345))

	cdf := createCDFForStat(stats, rng)

	// Check that values are sorted
	if len(cdf.Values) != 4 {
		t.Errorf("Expected 4 values, got %d", len(cdf.Values))
	}

	expectedValues := []int{1, 2, 3, 4}
	for i, val := range cdf.Values {
		if val != expectedValues[i] {
			t.Errorf("Expected value at index %d to be %d, got %d", i, expectedValues[i], val)
		}
	}

	// Check CDF values
	if len(cdf.CDF) != 4 {
		t.Errorf("Expected 4 CDF values, got %d", len(cdf.CDF))
	}

	expectedCDF := []float64{0.1, 0.3, 0.6, 1.0}
	for i, val := range cdf.CDF {
		if val != expectedCDF[i] {
			t.Errorf("Expected CDF at index %d to be %f, got %f", i, expectedCDF[i], val)
		}
	}

	// Last CDF value should always be 1.0
	if cdf.CDF[len(cdf.CDF)-1] != 1.0 {
		t.Errorf("Last CDF value should be 1.0, got %f", cdf.CDF[len(cdf.CDF)-1])
	}
}

func TestGenerateValueFromCDF(t *testing.T) {
	stats := map[int]int{
		1: 10,
		2: 20,
		3: 30,
		4: 40,
	}
	rng := rand.New(rand.NewSource(12345))
	cdf := createCDFForStat(stats, rng)

	// Generate multiple values and verify they're within the expected range
	counts := make(map[int]int)
	iterations := 1000

	for range iterations {
		value := generateValueFromCDF(cdf, rng)
		if value < 1 || value > 4 {
			t.Errorf("Generated value %d is out of range [1, 4]", value)
		}
		counts[value]++
	}

	// Verify all values were generated at least once
	for i := 1; i <= 4; i++ {
		if counts[i] == 0 {
			t.Errorf("Value %d was never generated", i)
		}
	}

	// Verify distribution is roughly correct (higher counts for higher frequency values)
	// Value 4 should appear most often (40% probability)
	if counts[4] < counts[1] {
		t.Error("Value 4 should appear more often than value 1")
	}
}

func TestBinarySearchUpperBound(t *testing.T) {
	cdf := &StatisticToCDF[int]{
		Values: []int{1, 2, 3, 4},
		CDF:    []float64{0.1, 0.3, 0.6, 1.0},
	}

	tests := []struct {
		name     string
		target   float64
		expected int
	}{
		{"target below first", 0.05, 0},
		{"target at first", 0.1, 0},
		{"target between first and second", 0.2, 1},
		{"target at second", 0.3, 1},
		{"target between second and third", 0.5, 2},
		{"target at third", 0.6, 2},
		{"target between third and fourth", 0.8, 3},
		{"target at fourth", 1.0, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binarySearchUpperBound(cdf, 0, len(cdf.Values)-1, tt.target)
			if result != tt.expected {
				t.Errorf("Expected index %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCreateGenerateValueFromStat(t *testing.T) {
	stats := map[string]int{
		"John": 50,
		"Jane": 30,
		"Bob":  20,
	}
	rng := rand.New(rand.NewSource(12345))

	generator := createGenerateValueFromStat(stats, rng)

	// Test that generator returns valid names
	counts := make(map[string]int)
	iterations := 1000

	for range iterations {
		name := generator()
		if _, ok := stats[name]; !ok {
			t.Errorf("Generated invalid name: %s", name)
		}
		counts[name]++
	}

	// All names should have been generated at least once
	for name := range stats {
		if counts[name] == 0 {
			t.Errorf("Name '%s' was never generated", name)
		}
	}

	// "John" should appear most often (50% probability)
	if counts["John"] < counts["Bob"] {
		t.Error("John should appear more often than Bob")
	}
}

func TestCreateRandomSkillFactorWithBellCurve(t *testing.T) {
	// Generate multiple skill values and verify they're within expected range
	for range 100 {
		skill := createRandomSkillFactorWithBellCurve()

		// Skill should generally be between 0 and 1, but can technically exceed these bounds
		// with normal distribution. We'll just check it's reasonable.
		if skill < -1.0 || skill > 2.0 {
			t.Errorf("Skill value %f is unreasonably outside expected range", skill)
		}
	}
}

func TestCreateSkillForDepthPosition(t *testing.T) {
	tests := []struct {
		name            string
		depthPosition   int
		totalAtPosition int
		minExpected     float64
		maxExpected     float64
	}{
		{"starter QB", 0, 3, 0.15, 0.95},
		{"backup QB", 1, 3, 0.15, 0.95},
		{"third string QB", 2, 3, 0.15, 0.95},
		{"starter RB", 0, 4, 0.15, 0.95},
		{"fourth string RB", 3, 4, 0.15, 0.95},
		{"single player position", 0, 1, 0.15, 0.95},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run multiple times to account for variance
			for range 10 {
				skill := createSkillForDepthPosition(tt.depthPosition, tt.totalAtPosition)

				if skill < tt.minExpected {
					t.Errorf("Skill %f is below minimum %f", skill, tt.minExpected)
				}
				if skill > tt.maxExpected {
					t.Errorf("Skill %f is above maximum %f", skill, tt.maxExpected)
				}
			}
		})
	}

	// Test that starters generally have higher skill than backups
	starterSkills := make([]float64, 20)
	backupSkills := make([]float64, 20)

	for i := range 20 {
		starterSkills[i] = createSkillForDepthPosition(0, 3)
		backupSkills[i] = createSkillForDepthPosition(2, 3)
	}

	// Calculate averages
	starterAvg := 0.0
	backupAvg := 0.0
	for i := range 20 {
		starterAvg += starterSkills[i]
		backupAvg += backupSkills[i]
	}
	starterAvg /= 20
	backupAvg /= 20

	if starterAvg <= backupAvg {
		t.Errorf("Starters should have higher average skill than backups. Starter avg: %f, Backup avg: %f", starterAvg, backupAvg)
	}

	// Test edge cases to ensure clamping works
	t.Run("skill clamping", func(t *testing.T) {
		// Run many iterations to hopefully hit the clamp boundaries
		minFound := 1.0
		maxFound := 0.0

		for range 10000 {
			skill := createSkillForDepthPosition(0, 3)
			if skill < minFound {
				minFound = skill
			}
			if skill > maxFound {
				maxFound = skill
			}

			// Verify always within bounds
			if skill < 0.15 || skill > 0.95 {
				t.Errorf("Skill %f outside bounds [0.15, 0.95]", skill)
			}
		}

		// We should have values near the boundaries
		if minFound > 0.3 {
			t.Logf("Warning: Minimum skill found (%f) seems high, clamping may not be tested", minFound)
		}
		if maxFound < 0.85 {
			t.Logf("Warning: Maximum skill found (%f) seems low, clamping may not be tested", maxFound)
		}
	})
}

func TestCreatePositionAttributeGenerators(t *testing.T) {
	profile := &PositionProfile{
		Jerseys:           map[int]int{1: 10, 2: 20},
		Heights:           map[int]int{70: 15, 72: 25},
		Weights:           map[int]int{180: 20, 200: 30},
		Ages:              map[int]int{23: 10, 25: 15},
		YearsOfExperience: map[int]int{1: 20, 3: 30},
	}
	rng := rand.New(rand.NewSource(12345))

	generators := CreatePositionAttributeGenerators(profile, rng)

	// Test that generators return valid values
	jersey := generators.JerseyGenerator()
	if jersey != 1 && jersey != 2 {
		t.Errorf("Expected jersey 1 or 2, got %d", jersey)
	}

	height := generators.HeightGenerator()
	if height != 70 && height != 72 {
		t.Errorf("Expected height 70 or 72, got %d", height)
	}

	weight := generators.WeightGenerator()
	if weight != 180 && weight != 200 {
		t.Errorf("Expected weight 180 or 200, got %d", weight)
	}

	age := generators.AgeGenerator()
	if age != 23 && age != 25 {
		t.Errorf("Expected age 23 or 25, got %d", age)
	}

	yoe := generators.YoeGenerator()
	if yoe != 1 && yoe != 3 {
		t.Errorf("Expected years of experience 1 or 3, got %d", yoe)
	}
}

func TestGetPlayerGenerators(t *testing.T) {
	// This function uses a singleton pattern with sync.Once
	// We can test that it returns valid generators
	rng := rand.New(rand.NewSource(12345))

	// Create a mock stats aggregator
	mockAggregator := func() AggregatedPlayerStats {
		return AggregatedPlayerStats{
			PositionProfile: map[string]*PositionProfile{
				"QB": {
					Jerseys:           map[int]int{12: 10},
					Heights:           map[int]int{72: 10},
					Weights:           map[int]int{200: 10},
					Ages:              map[int]int{25: 10},
					YearsOfExperience: map[int]int{3: 10},
				},
				"RB": {
					Jerseys:           map[int]int{28: 10},
					Heights:           map[int]int{70: 10},
					Weights:           map[int]int{210: 10},
					Ages:              map[int]int{24: 10},
					YearsOfExperience: map[int]int{2: 10},
				},
				"WR": {
					Jerseys:           map[int]int{88: 10},
					Heights:           map[int]int{73: 10},
					Weights:           map[int]int{190: 10},
					Ages:              map[int]int{26: 10},
					YearsOfExperience: map[int]int{4: 10},
				},
				"TE": {
					Jerseys:           map[int]int{87: 10},
					Heights:           map[int]int{75: 10},
					Weights:           map[int]int{250: 10},
					Ages:              map[int]int{27: 10},
					YearsOfExperience: map[int]int{5: 10},
				},
				"PK": {
					Jerseys:           map[int]int{4: 10},
					Heights:           map[int]int{68: 10},
					Weights:           map[int]int{175: 10},
					Ages:              map[int]int{28: 10},
					YearsOfExperience: map[int]int{6: 10},
				},
			},
			FirstNames: map[string]int{"John": 10, "Jane": 5},
			LastNames:  map[string]int{"Doe": 10, "Smith": 5},
		}
	}

	generators := getPlayerGenerators(mockAggregator, rng)

	// Test that generators are not nil
	if generators.FirstNameGenerator == nil {
		t.Error("FirstNameGenerator should not be nil")
	}
	if generators.LastNameGenerator == nil {
		t.Error("LastNameGenerator should not be nil")
	}
	if generators.SkillGenerator == nil {
		t.Error("SkillGenerator should not be nil")
	}
	if len(generators.PositionGenerators) != 5 {
		t.Errorf("Expected 5 position generators, got %d", len(generators.PositionGenerators))
	}

	// Test that generators work
	firstName := generators.FirstNameGenerator()
	if firstName != "John" && firstName != "Jane" {
		t.Errorf("Expected first name to be John or Jane, got %s", firstName)
	}

	lastName := generators.LastNameGenerator()
	if lastName != "Doe" && lastName != "Smith" {
		t.Errorf("Expected last name to be Doe or Smith, got %s", lastName)
	}

	skill := generators.SkillGenerator()
	if skill < -1.0 || skill > 2.0 {
		t.Errorf("Skill %f is unreasonably outside expected range", skill)
	}
}

func TestCreatePlayerGeneratorsFromStats(t *testing.T) {
	mockAggregator := func() AggregatedPlayerStats {
		return AggregatedPlayerStats{
			PositionProfile: map[string]*PositionProfile{
				"QB": {
					Jerseys:           map[int]int{12: 10},
					Heights:           map[int]int{72: 10},
					Weights:           map[int]int{200: 10},
					Ages:              map[int]int{25: 10},
					YearsOfExperience: map[int]int{3: 10},
				},
				"RB": {
					Jerseys:           map[int]int{28: 10},
					Heights:           map[int]int{70: 10},
					Weights:           map[int]int{210: 10},
					Ages:              map[int]int{24: 10},
					YearsOfExperience: map[int]int{2: 10},
				},
				"WR": {
					Jerseys:           map[int]int{88: 10},
					Heights:           map[int]int{73: 10},
					Weights:           map[int]int{190: 10},
					Ages:              map[int]int{26: 10},
					YearsOfExperience: map[int]int{4: 10},
				},
				"TE": {
					Jerseys:           map[int]int{87: 10},
					Heights:           map[int]int{75: 10},
					Weights:           map[int]int{250: 10},
					Ages:              map[int]int{27: 10},
					YearsOfExperience: map[int]int{5: 10},
				},
				"PK": {
					Jerseys:           map[int]int{4: 10},
					Heights:           map[int]int{68: 10},
					Weights:           map[int]int{175: 10},
					Ages:              map[int]int{28: 10},
					YearsOfExperience: map[int]int{6: 10},
				},
			},
			FirstNames: map[string]int{"John": 10},
			LastNames:  map[string]int{"Doe": 10},
		}
	}
	rng := rand.New(rand.NewSource(12345))

	firstNameGen, lastNameGen, posGensCopy := createPlayerGeneratorsFromStats(mockAggregator, rng)

	// Verify generators work
	if firstNameGen() != "John" {
		t.Error("First name generator should return John")
	}
	if lastNameGen() != "Doe" {
		t.Error("Last name generator should return Doe")
	}
	if len(posGensCopy) != 5 {
		t.Errorf("Expected 5 position generators, got %d", len(posGensCopy))
	}
}

func TestCreatePositionsGeneratorsFromStats(t *testing.T) {
	stats := AggregatedPlayerStats{
		PositionProfile: map[string]*PositionProfile{
			"QB": {
				Jerseys:           map[int]int{12: 10},
				Heights:           map[int]int{72: 10},
				Weights:           map[int]int{200: 10},
				Ages:              map[int]int{25: 10},
				YearsOfExperience: map[int]int{3: 10},
			},
			"RB": {
				Jerseys:           map[int]int{28: 10},
				Heights:           map[int]int{70: 10},
				Weights:           map[int]int{210: 10},
				Ages:              map[int]int{24: 10},
				YearsOfExperience: map[int]int{2: 10},
			},
			"WR": {
				Jerseys:           map[int]int{88: 10},
				Heights:           map[int]int{73: 10},
				Weights:           map[int]int{190: 10},
				Ages:              map[int]int{26: 10},
				YearsOfExperience: map[int]int{4: 10},
			},
			"TE": {
				Jerseys:           map[int]int{87: 10},
				Heights:           map[int]int{75: 10},
				Weights:           map[int]int{250: 10},
				Ages:              map[int]int{27: 10},
				YearsOfExperience: map[int]int{5: 10},
			},
			"PK": {
				Jerseys:           map[int]int{4: 10},
				Heights:           map[int]int{68: 10},
				Weights:           map[int]int{175: 10},
				Ages:              map[int]int{28: 10},
				YearsOfExperience: map[int]int{6: 10},
			},
		},
		FirstNames: map[string]int{"John": 10},
		LastNames:  map[string]int{"Doe": 10},
	}
	rng := rand.New(rand.NewSource(12345))

	positionGenerators := createPositionsGeneratorsFromStats(stats, rng)

	// Should have 5 positions
	if len(positionGenerators) != 5 {
		t.Errorf("Expected 5 position generators, got %d", len(positionGenerators))
	}

	// Verify all positions are present
	positionCodes := make(map[Position]bool)
	for _, pg := range positionGenerators {
		positionCodes[pg.PositionCode] = true
	}

	expectedPositions := []Position{QB, RB, WR, TE, PK}
	for _, pos := range expectedPositions {
		if !positionCodes[pos] {
			t.Errorf("Position %s not found in generators", pos)
		}
	}
}
