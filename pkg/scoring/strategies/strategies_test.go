package strategies

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"testing"
)

func getTestScores() map[metadata.MetadataType]scoring.ScoreValue {
	return map[metadata.MetadataType]scoring.ScoreValue{
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

func getTestWeights() map[metadata.MetadataType]scoring.Weight {
	return map[metadata.MetadataType]scoring.Weight{
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

func TestWeightedSumStrategy_Calculate(t *testing.T) {
	strategy := NewWeightedSumStrategy()
	scores := getTestScores()
	weights := getTestWeights()

	result, err := strategy.Calculate(scores, weights, nil)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.TotalScore <= 0 || result.TotalScore > 100 {
		t.Errorf("Invalid total score: %v", result.TotalScore)
	}

	if result.Method != scoring.MethodWeightedSum {
		t.Errorf("Expected method %v, got %v", scoring.MethodWeightedSum, result.Method)
	}

	// 점수가 유효한 범위인지 확인 (양수)
	if result.TotalScore <= 0 {
		t.Errorf("Expected positive score, got %v", result.TotalScore)
	}
}

func TestGeometricMeanStrategy_Calculate(t *testing.T) {
	strategy := NewGeometricMeanStrategy()
	scores := getTestScores()
	weights := getTestWeights()

	result, err := strategy.Calculate(scores, weights, nil)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.TotalScore <= 0 || result.TotalScore > 100 {
		t.Errorf("Invalid total score: %v", result.TotalScore)
	}

	if result.Method != scoring.MethodGeometricMean {
		t.Errorf("Expected method %v, got %v", scoring.MethodGeometricMean, result.Method)
	}

	// 기하 평균은 산술 평균보다 낮아야 함
	ws := NewWeightedSumStrategy()
	wsResult, _ := ws.Calculate(scores, weights, nil)
	if result.TotalScore >= wsResult.TotalScore {
		t.Error("Geometric mean should be lower than weighted sum for unbalanced scores")
	}
}

func TestMinMaxStrategy_Calculate(t *testing.T) {
	strategy := NewMinMaxStrategy()
	scores := getTestScores()
	weights := getTestWeights()

	result, err := strategy.Calculate(scores, weights, nil)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.Method != scoring.MethodMinMax {
		t.Errorf("Expected method %v, got %v", scoring.MethodMinMax, result.Method)
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

func TestHarmonicMeanStrategy_Calculate(t *testing.T) {
	strategy := NewHarmonicMeanStrategy()
	scores := getTestScores()
	weights := getTestWeights()

	result, err := strategy.Calculate(scores, weights, nil)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.TotalScore <= 0 || result.TotalScore > 100 {
		t.Errorf("Invalid total score: %v", result.TotalScore)
	}

	if result.Method != scoring.MethodHarmonicMean {
		t.Errorf("Expected method %v, got %v", scoring.MethodHarmonicMean, result.Method)
	}

	// 조화 평균은 다른 평균들보다 낮아야 함
	ws := NewWeightedSumStrategy()
	wsResult, _ := ws.Calculate(scores, weights, nil)
	if result.TotalScore >= wsResult.TotalScore {
		t.Error("Harmonic mean should be lower than or equal to weighted sum")
	}
}

func TestStrategyFactory(t *testing.T) {
	factory := NewDefaultStrategyFactory()

	strategies := factory.GetAvailableStrategies()
	expected := []scoring.ScoringMethod{
		scoring.MethodWeightedSum,
		scoring.MethodGeometricMean,
		scoring.MethodMinMax,
		scoring.MethodHarmonicMean,
	}

	if len(strategies) != len(expected) {
		t.Errorf("Expected %d strategies, got %d", len(expected), len(strategies))
	}

	for i, expectedStrategy := range expected {
		if i >= len(strategies) || strategies[i] != expectedStrategy {
			t.Errorf("Expected strategy %v at index %d, got %v", expectedStrategy, i, strategies[i])
		}
	}

	// 각 전략 생성 테스트
	for _, method := range strategies {
		strategy, err := factory.CreateStrategy(method)
		if err != nil {
			t.Errorf("Failed to create strategy %v: %v", method, err)
		}
		if strategy.Name() == "" {
			t.Errorf("Strategy %v has empty name", method)
		}
	}
}

func TestStrategyValidation(t *testing.T) {
	strategy := NewWeightedSumStrategy()

	// 유효한 입력
	validScores := getTestScores()
	validWeights := getTestWeights()

	err := strategy.ValidateInputs(validScores, validWeights)
	if err != nil {
		t.Errorf("Valid inputs should pass validation: %v", err)
	}

	// 잘못된 점수 (음수)
	invalidScores := make(map[metadata.MetadataType]scoring.ScoreValue)
	for k, v := range validScores {
		invalidScores[k] = v
	}
	invalidScores[metadata.FloorLevel] = -10

	err = strategy.ValidateInputs(invalidScores, validWeights)
	if err == nil {
		t.Error("Invalid scores should fail validation")
	}

	// 잘못된 가중치 합계
	invalidWeights := make(map[metadata.MetadataType]scoring.Weight)
	for k, v := range validWeights {
		invalidWeights[k] = v * 2 // 합계가 2가 되도록
	}

	err = strategy.ValidateInputs(validScores, invalidWeights)
	if err == nil {
		t.Error("Invalid weights should fail validation")
	}
}
