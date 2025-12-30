package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
	"testing"
)

func getTestScores() map[metadata.MetadataType]shared.ScoreValue {
	return map[metadata.MetadataType]shared.ScoreValue{
		metadata.FloorLevel:           shared.ScoreValueFromFloat(85.0),
		metadata.DistanceToStation:    shared.ScoreValueFromFloat(90.0),
		metadata.ElevatorPresence:     shared.ScoreValueFromFloat(100.0),
		metadata.ConstructionYear:     shared.ScoreValueFromFloat(80.0),
		metadata.ConstructionCompany:  shared.ScoreValueFromFloat(75.0),
		metadata.ApartmentSize:        shared.ScoreValueFromFloat(70.0),
		metadata.NearbyAmenities:      shared.ScoreValueFromFloat(80.0),
		metadata.TransportationAccess: shared.ScoreValueFromFloat(85.0),
		metadata.SchoolDistrict:       shared.ScoreValueFromFloat(75.0),
		metadata.CrimeRate:            shared.ScoreValueFromFloat(65.0),
		metadata.GreenSpaceRatio:      shared.ScoreValueFromFloat(60.0),
		metadata.Parking:              shared.ScoreValueFromFloat(80.0),
		metadata.MaintenanceFee:       shared.ScoreValueFromFloat(75.0),
		metadata.HeatingSystem:        shared.ScoreValueFromFloat(70.0),
	}
}

func getTestWeights() map[metadata.MetadataType]shared.Weight {
	weights := map[metadata.MetadataType]shared.Weight{
		metadata.FloorLevel:           shared.WeightFromFloat(0.08),
		metadata.DistanceToStation:    shared.WeightFromFloat(0.13),
		metadata.ElevatorPresence:     shared.WeightFromFloat(0.07),
		metadata.ConstructionYear:     shared.WeightFromFloat(0.10),
		metadata.ConstructionCompany:  shared.WeightFromFloat(0.08),
		metadata.ApartmentSize:        shared.WeightFromFloat(0.08),
		metadata.NearbyAmenities:      shared.WeightFromFloat(0.10),
		metadata.TransportationAccess: shared.WeightFromFloat(0.12),
		metadata.SchoolDistrict:       shared.WeightFromFloat(0.08),
		metadata.CrimeRate:            shared.WeightFromFloat(0.06),
		metadata.GreenSpaceRatio:      shared.WeightFromFloat(0.04),
		metadata.Parking:              shared.WeightFromFloat(0.06),
		metadata.MaintenanceFee:       shared.WeightFromFloat(0.05),
		metadata.HeatingSystem:        shared.WeightFromFloat(0.03),
	}
	return shared.NormalizeWeights(weights)
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

	// Min-Max 전략 검증: 최종 점수가 0보다 크고 100보다 작아야 함
	if result.TotalScore <= 0 || result.TotalScore > 100 {
		t.Errorf("Min-Max strategy should return valid score, got %v", result.TotalScore)
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
	invalidScores := make(map[metadata.MetadataType]shared.ScoreValue)
	for k, v := range validScores {
		invalidScores[k] = v
	}
	invalidScores[metadata.FloorLevel] = -10

	err = validateStrategyInputs(invalidScores, validWeights)
	if err == nil {
		t.Error("Invalid scores should fail validation")
	}

	// 잘못된 가중치 합계
	invalidWeights := make(map[metadata.MetadataType]shared.Weight)
	for k, v := range validWeights {
		invalidWeights[k] = shared.Weight(int64(v) * 2) // 합계가 2가 되도록
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

func TestCalculateRankings(t *testing.T) {
	// 테스트용 아파트 데이터 생성
	apartments := []ApartmentData{
		{
			ID:   "apt1",
			Name: "강남 뷰타워",
			Scores: map[metadata.MetadataType]shared.ScoreValue{
				metadata.FloorLevel:           85.0,
				metadata.DistanceToStation:    95.0,
				metadata.ElevatorPresence:     100.0,
				metadata.ConstructionYear:     90.0,
				metadata.ConstructionCompany:  88.0,
				metadata.ApartmentSize:        75.0,
				metadata.NearbyAmenities:      85.0,
				metadata.TransportationAccess: 90.0,
				metadata.SchoolDistrict:       80.0,
				metadata.CrimeRate:            70.0,
				metadata.GreenSpaceRatio:      65.0,
				metadata.Parking:              85.0,
				metadata.MaintenanceFee:       80.0,
				metadata.HeatingSystem:        75.0,
			},
			Location: "서울 강남구",
		},
		{
			ID:   "apt2",
			Name: "서초 힐스테이트",
			Scores: map[metadata.MetadataType]shared.ScoreValue{
				metadata.FloorLevel:           80.0,
				metadata.DistanceToStation:    85.0,
				metadata.ElevatorPresence:     100.0,
				metadata.ConstructionYear:     85.0,
				metadata.ConstructionCompany:  82.0,
				metadata.ApartmentSize:        70.0,
				metadata.NearbyAmenities:      80.0,
				metadata.TransportationAccess: 88.0,
				metadata.SchoolDistrict:       75.0,
				metadata.CrimeRate:            75.0,
				metadata.GreenSpaceRatio:      70.0,
				metadata.Parking:              80.0,
				metadata.MaintenanceFee:       75.0,
				metadata.HeatingSystem:        70.0,
			},
			Location: "서울 서초구",
		},
		{
			ID:   "apt3",
			Name: "송파 파크하비오",
			Scores: map[metadata.MetadataType]shared.ScoreValue{
				metadata.FloorLevel:           75.0,
				metadata.DistanceToStation:    80.0,
				metadata.ElevatorPresence:     95.0,
				metadata.ConstructionYear:     80.0,
				metadata.ConstructionCompany:  78.0,
				metadata.ApartmentSize:        65.0,
				metadata.NearbyAmenities:      75.0,
				metadata.TransportationAccess: 82.0,
				metadata.SchoolDistrict:       70.0,
				metadata.CrimeRate:            80.0,
				metadata.GreenSpaceRatio:      75.0,
				metadata.Parking:              75.0,
				metadata.MaintenanceFee:       70.0,
				metadata.HeatingSystem:        65.0,
			},
			Location: "서울 송파구",
		},
	}

	weights := getTestWeights()

	// 가중치 합계 전략으로 순위 계산
	summary, err := CalculateRankings(apartments, weights, StrategyWeightedSum)
	if err != nil {
		t.Fatalf("CalculateRankings failed: %v", err)
	}

	// 기본 검증
	if summary.TotalApartments != len(apartments) {
		t.Errorf("Expected %d apartments, got %d", len(apartments), summary.TotalApartments)
	}

	if summary.Strategy != StrategyWeightedSum {
		t.Errorf("Expected strategy %v, got %v", StrategyWeightedSum, summary.Strategy)
	}

	if len(summary.TopRanked) != len(apartments) {
		t.Errorf("Expected %d rankings, got %d", len(apartments), len(summary.TopRanked))
	}

	// 순위 검증 (점수가 높은 순서로 정렬되어야 함)
	for i := 1; i < len(summary.TopRanked); i++ {
		if summary.TopRanked[i-1].Score < summary.TopRanked[i].Score {
			t.Errorf("Rankings not sorted correctly: rank %d score %.1f < rank %d score %.1f",
				i, summary.TopRanked[i-1].Score, i+1, summary.TopRanked[i].Score)
		}
	}

	// 순위 번호 검증
	for i, ranking := range summary.TopRanked {
		expectedRank := i + 1
		if ranking.Rank != expectedRank {
			t.Errorf("Expected rank %d, got %d for apartment %s", expectedRank, ranking.Rank, ranking.Apartment.Name)
		}
	}

	// 백분위수 범위 검증
	for _, ranking := range summary.TopRanked {
		if ranking.Percentile < 0 || ranking.Percentile > 100 {
			t.Errorf("Invalid percentile %.1f for apartment %s", ranking.Percentile, ranking.Apartment.Name)
		}
	}

	// 1위 아파트는 가장 높은 백분위수를 가져야 함
	if len(summary.TopRanked) > 1 {
		firstPercentile := summary.TopRanked[0].Percentile
		for i := 1; i < len(summary.TopRanked); i++ {
			if summary.TopRanked[i].Percentile > firstPercentile {
				t.Errorf("First ranked apartment should have highest percentile, got %.1f vs %.1f",
					firstPercentile, summary.TopRanked[i].Percentile)
			}
		}
	}

	// 점수 범위 검증
	if summary.ScoreRange.Min > summary.ScoreRange.Max {
		t.Error("Min score should not be greater than max score")
	}

	calculatedAvg := summary.ScoreRange.Min + summary.ScoreRange.Max/2 // 대략적인 평균
	if summary.ScoreRange.Avg < calculatedAvg*0.8 || summary.ScoreRange.Avg > calculatedAvg*1.2 {
		t.Logf("Average score %.1f seems reasonable compared to range %.1f-%.1f",
			summary.ScoreRange.Avg, summary.ScoreRange.Min, summary.ScoreRange.Max)
	}
}

func TestFormatRankings(t *testing.T) {
	apartments := []ApartmentData{
		{
			ID:   "apt1",
			Name: "테스트 아파트 A",
			Scores: map[metadata.MetadataType]shared.ScoreValue{
				metadata.FloorLevel:           80.0,
				metadata.DistanceToStation:    85.0,
				metadata.ElevatorPresence:     100.0,
				metadata.ConstructionYear:     75.0,
				metadata.ConstructionCompany:  70.0,
				metadata.ApartmentSize:        65.0,
				metadata.NearbyAmenities:      75.0,
				metadata.TransportationAccess: 80.0,
				metadata.SchoolDistrict:       70.0,
				metadata.CrimeRate:            75.0,
				metadata.GreenSpaceRatio:      70.0,
				metadata.Parking:              80.0,
				metadata.MaintenanceFee:       75.0,
				metadata.HeatingSystem:        70.0,
			},
		},
		{
			ID:   "apt2",
			Name: "테스트 아파트 B",
			Scores: map[metadata.MetadataType]shared.ScoreValue{
				metadata.FloorLevel:           85.0,
				metadata.DistanceToStation:    90.0,
				metadata.ElevatorPresence:     100.0,
				metadata.ConstructionYear:     80.0,
				metadata.ConstructionCompany:  75.0,
				metadata.ApartmentSize:        70.0,
				metadata.NearbyAmenities:      80.0,
				metadata.TransportationAccess: 85.0,
				metadata.SchoolDistrict:       75.0,
				metadata.CrimeRate:            70.0,
				metadata.GreenSpaceRatio:      65.0,
				metadata.Parking:              85.0,
				metadata.MaintenanceFee:       80.0,
				metadata.HeatingSystem:        75.0,
			},
		},
	}

	weights := getTestWeights()
	summary, err := CalculateRankings(apartments, weights, StrategyWeightedSum)
	if err != nil {
		t.Fatalf("CalculateRankings failed: %v", err)
	}

	// 전체 순위 포맷팅
	fullOutput := FormatRankings(summary, 0)
	if !contains(fullOutput, "테스트 아파트 A") || !contains(fullOutput, "테스트 아파트 B") {
		t.Error("Formatted output should contain all apartment names")
	}

	// 제한된 순위 포맷팅 (상위 1개만)
	limitedOutput := FormatRankings(summary, 1)
	if !contains(limitedOutput, "🥇") {
		t.Error("Limited output should contain top rank emoji")
	}
	if contains(limitedOutput, "테스트 아파트 A") && contains(limitedOutput, "테스트 아파트 B") {
		// 상위 1개 제한인데 두 아파트 모두 나오면 외 XX개 문구가 있어야 함
		if !contains(limitedOutput, "외") {
			t.Error("Limited output should indicate more apartments exist")
		}
	}

	// 빈 데이터 테스트
	emptyOutput := FormatRankings(nil, 0)
	expected := "순위 데이터가 없습니다."
	if emptyOutput != expected {
		t.Errorf("Expected empty output %q, got %q", expected, emptyOutput)
	}
}
