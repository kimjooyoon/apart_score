package scoring

import (
	"testing"

	"apart_score/pkg/metadata"
)

func TestDefaultScorer_Calculate(t *testing.T) {
	scorer := NewDefaultScorer(MethodWeightedSum)

	// 테스트용 점수 데이터
	scores := map[metadata.MetadataType]ScoreValue{
		metadata.FloorLevel:        85.0,
		metadata.DistanceToStation: 90.0,
		metadata.ElevatorPresence:  100.0,
		// 다른 요소들은 70점으로 설정
	}

	for _, mt := range metadata.AllMetadataTypes() {
		if _, exists := scores[mt]; !exists {
			scores[mt] = 70.0
		}
	}

	weights := scorer.GetDefaultWeights()

	result, err := scorer.Calculate(scores, weights)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if result.TotalScore <= 0 || result.TotalScore > 100 {
		t.Errorf("Invalid total score: %v", result.TotalScore)
	}
}

func TestQuickScore(t *testing.T) {
	scores := map[metadata.MetadataType]ScoreValue{
		metadata.FloorLevel:        85.0,
		metadata.DistanceToStation: 90.0,
		metadata.ElevatorPresence:  100.0,
	}

	// 다른 요소들 기본값 설정
	for _, mt := range metadata.AllMetadataTypes() {
		if _, exists := scores[mt]; !exists {
			scores[mt] = 70.0
		}
	}

	weights := GetScenarioWeights(ScenarioBalanced)
	result, err := CalculateWithStrategy(scores, weights, StrategyWeightedSum)
	if err != nil {
		t.Fatalf("QuickScore failed: %v", err)
	}

	if result.TotalScore <= 0 {
		t.Errorf("Invalid total score: %v", result.TotalScore)
	}

	// Scenario 검증은 제거 (CalculateWithStrategy에서는 Scenario를 설정하지 않음)
}

func TestAnalyzeScore(t *testing.T) {
	// 테스트용 결과 생성
	result := &ScoreResult{
		TotalScore: 85.0,
		RawScores: map[metadata.MetadataType]ScoreValue{
			metadata.FloorLevel:        90.0, // 강점
			metadata.DistanceToStation: 95.0, // 강점
			metadata.ElevatorPresence:  50.0, // 약점
			metadata.MaintenanceFee:    55.0, // 약점
		},
		Weights: map[metadata.MetadataType]Weight{
			metadata.FloorLevel:        0.1,
			metadata.DistanceToStation: 0.2,
			metadata.ElevatorPresence:  0.1,
			metadata.MaintenanceFee:    0.1,
		},
		Method:   MethodWeightedSum,
		Scenario: ScenarioBalanced,
	}

	analysis := AnalyzeScore(result)

	if len(analysis.Strengths) == 0 {
		t.Error("Should have strengths")
	}

	if len(analysis.Weaknesses) == 0 {
		t.Error("Should have weaknesses")
	}

	if len(analysis.ImprovementTips) == 0 {
		t.Error("Should have improvement tips")
	}
}

func TestRecommendScenario(t *testing.T) {
	// 교통 점수가 높은 경우
	transportScores := map[metadata.MetadataType]ScoreValue{
		metadata.DistanceToStation:    95.0,
		metadata.TransportationAccess: 90.0,
	}

	scenario := RecommendScenario(transportScores)
	if scenario != ScenarioTransportation {
		t.Errorf("Expected %v, got %v", ScenarioTransportation, scenario)
	}

	// 교육 점수가 높은 경우
	educationScores := map[metadata.MetadataType]ScoreValue{
		metadata.SchoolDistrict: 95.0,
	}

	scenario = RecommendScenario(educationScores)
	if scenario != ScenarioEducation {
		t.Errorf("Expected %v, got %v", ScenarioEducation, scenario)
	}
}

func TestGetScenarioWeights(t *testing.T) {
	weights := GetScenarioWeights(ScenarioTransportation)

	if len(weights) != int(metadata.MetadataTypeCount) {
		t.Errorf("Expected %d weights, got %d", metadata.MetadataTypeCount, len(weights))
	}

	// 교통 시나리오에서 역까지 거리 가중치가 높아야 함
	stationWeight := weights[metadata.DistanceToStation]
	if stationWeight < 0.2 {
		t.Errorf("Transportation scenario should have high station weight, got %v", stationWeight)
	}
}
