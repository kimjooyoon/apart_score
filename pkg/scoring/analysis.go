package scoring

import (
	"fmt"
	"sort"

	"apart_score/pkg/metadata"
)

// ScoreAnalysis는 점수 분석 결과를 담는 구조체
type ScoreAnalysis struct {
	Result           *ScoreResult
	Strengths        []metadata.MetadataType // 강점 요소들
	Weaknesses       []metadata.MetadataType // 약점 요소들
	TopFactors       []ScoreFactor           // 상위 영향 요인들
	ImprovementTips  []string                // 개선 제안
	ComparisonScore  float64                 // 평균 대비 점수
}

// ScoreFactor는 점수에 영향을 미치는 요소
type ScoreFactor struct {
	Metadata metadata.MetadataType
	Score    ScoreValue
	Weight   Weight
	Impact   ScoreValue // 영향력 (score * weight)
}

// AnalyzeScore는 점수를 분석합니다.
func AnalyzeScore(result *ScoreResult) *ScoreAnalysis {
	analysis := &ScoreAnalysis{
		Result:          result,
		Strengths:       []metadata.MetadataType{},
		Weaknesses:      []metadata.MetadataType{},
		TopFactors:      []ScoreFactor{},
		ImprovementTips: []string{},
	}

	// 각 요소의 영향력 계산
	var factors []ScoreFactor
	for mt, score := range result.RawScores {
		weight := result.Weights[mt]
		impact := score * ScoreValue(weight)

		factors = append(factors, ScoreFactor{
			Metadata: mt,
			Score:    score,
			Weight:   weight,
			Impact:   impact,
		})

		// 강점과 약점 분류 (80점 이상/이하 기준)
		if score >= 80 {
			analysis.Strengths = append(analysis.Strengths, mt)
		} else if score <= 60 {
			analysis.Weaknesses = append(analysis.Weaknesses, mt)
		}
	}

	// 영향력이 큰 순서대로 정렬
	sort.Slice(factors, func(i, j int) bool {
		return factors[i].Impact > factors[j].Impact
	})

	// 상위 5개 영향 요인 선택
	maxFactors := 5
	if len(factors) < maxFactors {
		maxFactors = len(factors)
	}
	analysis.TopFactors = factors[:maxFactors]

	// 개선 제안 생성
	analysis.ImprovementTips = generateImprovementTips(analysis.Weaknesses)

	// 평균 대비 점수 계산 (임의의 평균값 사용)
	averageScore := 75.0
	analysis.ComparisonScore = float64(result.TotalScore) - averageScore

	return analysis
}

// generateImprovementTips는 개선 제안을 생성합니다.
func generateImprovementTips(weaknesses []metadata.MetadataType) []string {
	tips := []string{}

	for _, mt := range weaknesses {
		switch mt {
		case metadata.FloorLevel:
			tips = append(tips, "중간층에 가까운 아파트를 고려해보세요")
		case metadata.DistanceToStation:
			tips = append(tips, "역과 가까운 아파트를 선택하세요")
		case metadata.ElevatorPresence:
			tips = append(tips, "엘리베이터가 있는 건물을 우선 고려하세요")
		case metadata.ConstructionYear:
			tips = append(tips, "최근에 지어진 아파트를 선택하세요")
		case metadata.ConstructionCompany:
			tips = append(tips, "신뢰할 수 있는 건설사의 아파트를 고려하세요")
		case metadata.ApartmentSize:
			tips = append(tips, "적절한 크기의 아파트를 선택하세요")
		case metadata.SchoolDistrict:
			tips = append(tips, "좋은 학군이 있는 지역을 고려하세요")
		case metadata.CrimeRate:
			tips = append(tips, "범죄율이 낮은 안전한 지역을 선택하세요")
		case metadata.MaintenanceFee:
			tips = append(tips, "관리비가 적절한 수준인 아파트를 선택하세요")
		}
	}

	return tips
}

// CompareScores는 두 점수를 비교합니다.
func CompareScores(result1, result2 *ScoreResult) string {
	diff := result1.TotalScore - result2.TotalScore

	if diff > 10 {
		return fmt.Sprintf("첫 번째 옵션이 %.1f점 더 높습니다", diff)
	} else if diff < -10 {
		return fmt.Sprintf("두 번째 옵션이 %.1f점 더 높습니다", -diff)
	} else {
		return fmt.Sprintf("두 옵션의 점수가 비슷합니다 (차이: %.1f점)", diff)
	}
}

// RecommendScenario는 주어진 점수에 기반하여 추천 시나리오를 제안합니다.
func RecommendScenario(scores map[metadata.MetadataType]ScoreValue) ScoringScenario {
	// 가장 높은 점수를 받은 요소들을 분석
	type scorePair struct {
		metadata metadata.MetadataType
		score    ScoreValue
	}

	var highScores []scorePair
	for mt, score := range scores {
		if score >= 80 {
			highScores = append(highScores, scorePair{mt, score})
		}
	}

	// 교통 관련 요소가 많으면 교통 중심 추천
	transportCount := 0
	for _, pair := range highScores {
		if pair.metadata == metadata.DistanceToStation || pair.metadata == metadata.TransportationAccess {
			transportCount++
		}
	}
	if transportCount >= 2 {
		return ScenarioTransportation
	}

	// 교육 관련 요소가 높으면 교육 중심 추천
	if scores[metadata.SchoolDistrict] >= 85 {
		return ScenarioEducation
	}

	// 관리비가 낮고 크기가 적당하면 가성비 중심 추천
	if scores[metadata.MaintenanceFee] >= 80 && scores[metadata.ApartmentSize] >= 75 {
		return ScenarioCostEffective
	}

	// 기본적으로 균형 잡힌 추천
	return ScenarioBalanced
}

// FormatScoreResult는 점수 결과를 읽기 쉽게 포맷팅합니다.
func FormatScoreResult(result *ScoreResult) string {
	output := fmt.Sprintf("🏠 아파트 스코어 결과\n")
	output += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	output += fmt.Sprintf("총점: %.1f점 (등급: %s)\n", result.TotalScore, result.Grade)
	output += fmt.Sprintf("방법: %s\n", result.Method)
	output += fmt.Sprintf("시나리오: %s\n", result.Scenario)
	output += fmt.Sprintf("\n📊 상세 점수:\n")

	for _, mt := range metadata.AllMetadataTypes() {
		if rawScore, exists := result.RawScores[mt]; exists {
			weight := result.Weights[mt]
			weighted := result.WeightedScores[mt]
			output += fmt.Sprintf("  %-20s: %.1f점 (가중치: %.1f%%) → %.1f점\n",
				mt.KoreanName(), rawScore, float64(weight)*100, weighted)
		}
	}

	return output
}
