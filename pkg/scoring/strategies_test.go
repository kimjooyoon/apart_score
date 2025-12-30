package scoring

import (
	"apart_score/pkg/metadata"
	"testing"
)

func getTestScores() map[metadata.MetadataType]ScoreValue {
	return map[metadata.MetadataType]ScoreValue{
		metadata.FloorLevel:           85.0,
		metadata.DistanceToStation:    90.0,
		metadata.ElevatorPresence:     100.0,
		metadata.ConstructionYear:     80.0,
		metadata.ConstructionCompany:  75.0,
		metadata.ApartmentSize:        70.0,
		metadata.NearbyAmenities:      80.0,
		metadata.TransportationAccess: 85.0,
		metadata.SchoolDistrict:       75.0,
		metadata.CrimeRate:            65.0,
		metadata.GreenSpaceRatio:      60.0,
		metadata.Parking:              80.0,
		metadata.MaintenanceFee:       75.0,
		metadata.HeatingSystem:        70.0,
	}
}

func getTestWeights() map[metadata.MetadataType]Weight {
	return map[metadata.MetadataType]Weight{
		metadata.FloorLevel:           0.08,
		metadata.DistanceToStation:    0.14,
		metadata.ElevatorPresence:     0.07,
		metadata.ConstructionYear:     0.09,
		metadata.ConstructionCompany:  0.08,
		metadata.ApartmentSize:        0.08,
		metadata.NearbyAmenities:      0.09,
		metadata.TransportationAccess: 0.11,
		metadata.SchoolDistrict:       0.08,
		metadata.CrimeRate:            0.05,
		metadata.GreenSpaceRatio:      0.03,
		metadata.Parking:              0.05,
		metadata.MaintenanceFee:       0.04,
		metadata.HeatingSystem:        0.02,
	}
}

func TestCalculateWeightedSum(t *testing.T) {
	scores := getTestScores()
	weights := getTestWeights()

	result, err := CalculateWithStrategy(scores, weights, StrategyWeightedSum)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.TotalScore <= 0 || result.TotalScore > 100 {
		t.Errorf("Invalid total score: %v", result.TotalScore)
	}

	if result.Method != MethodWeightedSum {
		t.Errorf("Expected method %v, got %v", MethodWeightedSum, result.Method)
	}

	// 점수가 유효한 범위인지 확인 (양수)
	if result.TotalScore <= 0 {
		t.Errorf("Expected positive score, got %v", result.TotalScore)
	}
}

func TestCalculateGeometricMean(t *testing.T) {
	scores := getTestScores()
	weights := getTestWeights()

	result, err := CalculateWithStrategy(scores, weights, StrategyGeometricMean)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.TotalScore <= 0 || result.TotalScore > 100 {
		t.Errorf("Invalid total score: %v", result.TotalScore)
	}

	if result.Method != MethodGeometricMean {
		t.Errorf("Expected method %v, got %v", MethodGeometricMean, result.Method)
	}

	// 기하 평균은 산술 평균보다 낮아야 함
	wsResult, _ := CalculateWithStrategy(scores, weights, StrategyWeightedSum)
	if result.TotalScore >= wsResult.TotalScore {
		t.Error("Geometric mean should be lower than weighted sum for unbalanced scores")
	}
}

func TestCalculateMinMax(t *testing.T) {
	scores := getTestScores()
	weights := getTestWeights()

	result, err := CalculateWithStrategy(scores, weights, StrategyMinMax)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.Method != MethodMinMax {
		t.Errorf("Expected method %v, got %v", MethodMinMax, result.Method)
	}

	// Min-Max 전략에서는 총점이 최소 점수와 같아야 함
	minScore := 100.0
	for _, score := range scores {
		if float64(score) < minScore {
			minScore = float64(score)
		}
	}

	if float64(result.TotalScore) != minScore {
		t.Errorf("Min-Max strategy should return minimum score %v, got %v", minScore, result.TotalScore)
	}
}

func TestCalculateHarmonicMean(t *testing.T) {
	scores := getTestScores()
	weights := getTestWeights()

	result, err := CalculateWithStrategy(scores, weights, StrategyHarmonicMean)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.TotalScore <= 0 || result.TotalScore > 100 {
		t.Errorf("Invalid total score: %v", result.TotalScore)
	}

	if result.Method != MethodHarmonicMean {
		t.Errorf("Expected method %v, got %v", MethodHarmonicMean, result.Method)
	}

	// 조화 평균은 다른 평균들보다 낮아야 함
	wsResult, _ := CalculateWithStrategy(scores, weights, StrategyWeightedSum)
	if result.TotalScore >= wsResult.TotalScore {
		t.Error("Harmonic mean should be lower than or equal to weighted sum")
	}
}

func TestGetAvailableStrategies(t *testing.T) {
	strategies := GetAvailableStrategies()
	expected := []StrategyType{
		StrategyWeightedSum,
		StrategyGeometricMean,
		StrategyMinMax,
		StrategyHarmonicMean,
	}

	if len(strategies) != len(expected) {
		t.Errorf("Expected %d strategies, got %d", len(expected), len(strategies))
	}

	for i, expectedStrategy := range expected {
		if i >= len(strategies) || strategies[i] != expectedStrategy {
			t.Errorf("Expected strategy %v at index %d, got %v", expectedStrategy, i, strategies[i])
		}
	}
}

func TestStrategyValidation(t *testing.T) {
	// 유효한 입력
	validScores := getTestScores()
	validWeights := getTestWeights()

	err := validateStrategyInputs(validScores, validWeights)
	if err != nil {
		t.Errorf("Valid inputs should pass validation: %v", err)
	}

	// 잘못된 점수 (음수)
	invalidScores := make(map[metadata.MetadataType]ScoreValue)
	for k, v := range validScores {
		invalidScores[k] = v
	}
	invalidScores[metadata.FloorLevel] = -10

	err = validateStrategyInputs(invalidScores, validWeights)
	if err == nil {
		t.Error("Invalid scores should fail validation")
	}

	// 잘못된 가중치 합계
	invalidWeights := make(map[metadata.MetadataType]Weight)
	for k, v := range validWeights {
		invalidWeights[k] = v * 2 // 합계가 2가 되도록
	}

	err = validateStrategyInputs(validScores, invalidWeights)
	if err == nil {
		t.Error("Invalid weights should fail validation")
	}
}

func TestGetStrategyDescription(t *testing.T) {
	tests := []struct {
		strategy  StrategyType
		expected  string
		hasPrefix bool
	}{
		{StrategyWeightedSum, "각 메타데이터의 점수에 가중치를 곱한 후 합계를 계산하는 기본 전략입니다.", false},
		{StrategyGeometricMean, "모든 요소가 균형을 이루어야 하는 경우에 적합한 전략입니다. 하나의 낮은 점수가 전체 점수를 크게 낮춥니다.", false},
		{StrategyMinMax, "모든 요소가 일정 수준 이상이어야 하는 경우에 적합합니다. 가장 낮은 점수가 전체 점수를 결정합니다.", false},
		{StrategyHarmonicMean, "낮은 점수에 매우 민감하게 반응하는 전략입니다. 모든 요소가 고르게 중요할 때 사용합니다.", false},
		{"unknown", "알 수 없는 전략입니다.", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.strategy), func(t *testing.T) {
			result := GetStrategyDescription(tt.strategy)
			if tt.hasPrefix && !contains(result, tt.expected) {
				t.Errorf("Expected description to contain %q, got %q", tt.expected, result)
			} else if !tt.hasPrefix && result != tt.expected {
				t.Errorf("Expected description %q, got %q", tt.expected, result)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
