package relative

import (
	"apart_score/pkg/scoring"
	"testing"
)

func createTestApartments() []ApartmentScore {
	return []ApartmentScore{
		{ID: "apt1", Score: 85.0, Location: "서울 강남"},
		{ID: "apt2", Score: 90.0, Location: "서울 강남"},
		{ID: "apt3", Score: 75.0, Location: "서울 강남"},
		{ID: "apt4", Score: 88.0, Location: "서울 강남"},
		{ID: "apt5", Score: 82.0, Location: "서울 강남"},
		{ID: "apt6", Score: 95.0, Location: "서울 서초"},
		{ID: "apt7", Score: 78.0, Location: "서울 서초"},
	}
}

func TestDefaultRelativeEvaluator_EvaluateRelatively(t *testing.T) {
	evaluator := NewDefaultRelativeEvaluator()
	apartments := createTestApartments()
	target := apartments[0] // apt1, score: 85.0

	result, err := evaluator.EvaluateRelatively(target, apartments)
	if err != nil {
		t.Fatalf("EvaluateRelatively failed: %v", err)
	}

	if result.ApartmentID != target.ID {
		t.Errorf("Expected apartment ID %s, got %s", target.ID, result.ApartmentID)
	}

	if result.AbsoluteScore != target.Score {
		t.Errorf("Expected absolute score %v, got %v", target.Score, result.AbsoluteScore)
	}

	// 백분위수 검증 (85점은 5개 중 3번째이므로 약 40-60% 사이)
	if result.PercentileRank < 0 || result.PercentileRank > 100 {
		t.Errorf("Invalid percentile rank: %v", result.PercentileRank)
	}

	// 순위 검증
	if result.GroupRank <= 0 {
		t.Errorf("Invalid group rank: %v", result.GroupRank)
	}

	// 분포 정보 검증
	if result.ScoreDistribution.Count != len(apartments) {
		t.Errorf("Expected distribution count %d, got %d", len(apartments), result.ScoreDistribution.Count)
	}

	// 비교 정보 검증
	totalComparisons := result.Comparison.BetterThanCount + result.Comparison.WorseThanCount + result.Comparison.SimilarCount
	expectedComparisons := len(apartments) - 1 // 자신 제외
	if totalComparisons != expectedComparisons {
		t.Errorf("Expected %d comparisons, got %d", expectedComparisons, totalComparisons)
	}
}

func TestDefaultRelativeEvaluator_CalculatePercentile(t *testing.T) {
	evaluator := NewDefaultRelativeEvaluator()
	apartments := createTestApartments()

	tests := []struct {
		score     scoring.ScoreValue
		expected  float64
		tolerance float64
	}{
		{95.0, 100.0, 1.0}, // 최고 점수
		{75.0, 0.0, 1.0},   // 최저 점수
		{85.0, 50.0, 20.0}, // 중간 점수 (±20% 오차 허용)
	}

	for _, tt := range tests {
		percentile := evaluator.CalculatePercentile(tt.score, apartments)
		if percentile < tt.expected-tt.tolerance || percentile > tt.expected+tt.tolerance {
			t.Errorf("Score %v: expected percentile around %v (±%v), got %v",
				tt.score, tt.expected, tt.tolerance, percentile)
		}
	}
}

func TestDefaultRelativeEvaluator_FindSimilarApartments(t *testing.T) {
	evaluator := NewDefaultRelativeEvaluator()
	apartments := createTestApartments()
	target := apartments[0] // apt1, score: 85.0

	criteria := SimilarityCriteria{
		ScoreRange: 10.0, // ±10점
		MaxResults: 5,
	}

	similar := evaluator.FindSimilarApartments(target, apartments, criteria)

	// 자신은 제외되어야 함
	for _, apt := range similar {
		if apt.ID == target.ID {
			t.Error("Target apartment should not be in similar results")
		}

		// 점수 범위 검증 (±10점)
		scoreDiff := abs(float64(apt.Score - target.Score))
		if scoreDiff > 10.0 {
			t.Errorf("Similar apartment score %v too different from target %v",
				apt.Score, target.Score)
		}
	}

	// 최대 결과 수 검증
	if len(similar) > criteria.MaxResults {
		t.Errorf("Expected at most %d similar apartments, got %d",
			criteria.MaxResults, len(similar))
	}
}

func TestCalculateDistribution(t *testing.T) {
	evaluator := NewDefaultRelativeEvaluator()
	apartments := createTestApartments()

	distribution := evaluator.calculateDistribution(apartments)

	if distribution.Count != len(apartments) {
		t.Errorf("Expected count %d, got %d", len(apartments), distribution.Count)
	}

	// 평균 계산 검증
	sum := 0.0
	for _, apt := range apartments {
		sum += float64(apt.Score)
	}
	expectedMean := scoring.ScoreValue(sum / float64(len(apartments)))

	if abs(float64(distribution.Mean-expectedMean)) > 0.01 {
		t.Errorf("Expected mean %v, got %v", expectedMean, distribution.Mean)
	}

	// 최소/최대 검증
	if distribution.Min > distribution.Max {
		t.Error("Min should be less than or equal to Max")
	}

	// 사분위수 검증
	if distribution.Q1 > distribution.Q3 {
		t.Error("Q1 should be less than or equal to Q3")
	}

	if distribution.Median < distribution.Q1 || distribution.Median > distribution.Q3 {
		t.Error("Median should be between Q1 and Q3")
	}
}

func TestCalculateRanks(t *testing.T) {
	evaluator := NewDefaultRelativeEvaluator()
	apartments := createTestApartments()
	target := apartments[1] // apt2, score: 90.0 (두 번째로 높은 점수)

	regionalRank, groupRank := evaluator.calculateRanks(target, apartments)

	// 90점은 전체에서 두 번째로 높으므로 순위는 2
	expectedRank := 2

	if regionalRank != expectedRank {
		t.Errorf("Expected regional rank %d, got %d", expectedRank, regionalRank)
	}

	if groupRank != expectedRank {
		t.Errorf("Expected group rank %d, got %d", expectedRank, groupRank)
	}
}

func TestCalculateComparison(t *testing.T) {
	evaluator := NewDefaultRelativeEvaluator()
	apartments := createTestApartments()
	target := apartments[0] // apt1, score: 85.0

	comparison := evaluator.calculateComparison(target, apartments)

	// 85점 기준으로 더 높은 점수: 90, 88, 95 (3개)
	// 더 낮은 점수: 75, 82, 78 (3개)
	// 비슷한 점수 (±5): 88 (1개)

	expectedBetter := 3 // 90, 88, 95
	expectedWorse := 3  // 75, 82, 78

	if comparison.BetterThanCount != expectedBetter {
		t.Errorf("Expected better count %d, got %d", expectedBetter, comparison.BetterThanCount)
	}

	if comparison.WorseThanCount != expectedWorse {
		t.Errorf("Expected worse count %d, got %d", expectedWorse, comparison.WorseThanCount)
	}

	// 순위 백분위수 검증 (0-100 범위)
	if comparison.RankPercentile < 0 || comparison.RankPercentile > 100 {
		t.Errorf("Invalid rank percentile: %v", comparison.RankPercentile)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
