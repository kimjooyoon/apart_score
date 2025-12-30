package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
	"fmt"
	"sort"
)

type ScoreAnalysis struct {
	Result          *ScoreResult
	Strengths       []metadata.MetadataType
	Weaknesses      []metadata.MetadataType
	TopFactors      []ScoreFactor
	ImprovementTips []string
	ComparisonScore float64
}
type ScoreFactor struct {
	Metadata metadata.MetadataType
	Score    shared.ScoreValue
	Weight   shared.Weight
	Impact   shared.ScoreValue
}

func AnalyzeScore(result ScoreResult) *ScoreAnalysis {
	analysis := &ScoreAnalysis{
		Result:          &result,
		Strengths:       []metadata.MetadataType{},
		Weaknesses:      []metadata.MetadataType{},
		TopFactors:      []ScoreFactor{},
		ImprovementTips: []string{},
	}
	var factors []ScoreFactor
	for mt, score := range result.RawScores {
		weight := result.Weights[mt]
		// 정수 기반 계산: (score * weight) / WeightScale
		impact := shared.ScoreValue(int64(score) * int64(weight) / shared.WeightScale)
		factors = append(factors, ScoreFactor{
			Metadata: mt,
			Score:    score,
			Weight:   weight,
			Impact:   impact,
		})
		if score >= 80 {
			analysis.Strengths = append(analysis.Strengths, mt)
		} else if score <= 60 {
			analysis.Weaknesses = append(analysis.Weaknesses, mt)
		}
	}
	sort.Slice(factors, func(i, j int) bool {
		return factors[i].Impact > factors[j].Impact
	})
	maxFactors := 5
	if len(factors) < maxFactors {
		maxFactors = len(factors)
	}
	analysis.TopFactors = factors[:maxFactors]
	analysis.ImprovementTips = generateImprovementTips(analysis.Weaknesses)
	averageScore := 75.0
	analysis.ComparisonScore = result.TotalScore.ToFloat() - averageScore
	return analysis
}
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
func RecommendScenario(scores map[metadata.MetadataType]shared.ScoreValue) ScoringScenario {
	type scorePair struct {
		metadata metadata.MetadataType
		score    shared.ScoreValue
	}
	var highScores []scorePair
	for mt, score := range scores {
		if score >= 80 {
			highScores = append(highScores, scorePair{mt, score})
		}
	}
	transportCount := 0
	for _, pair := range highScores {
		if pair.metadata == metadata.DistanceToStation || pair.metadata == metadata.TransportationAccess {
			transportCount++
		}
	}
	if transportCount >= 2 {
		return ScenarioTransportation
	}
	if scores[metadata.SchoolDistrict] >= 85 {
		return ScenarioEducation
	}
	if scores[metadata.MaintenanceFee] >= 80 && scores[metadata.ApartmentSize] >= 75 {
		return ScenarioCostEffective
	}
	return ScenarioBalanced
}
func FormatScoreResult(result ScoreResult) string {
	output := "🏠 아파트 스코어 결과\n"
	output += "━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	output += fmt.Sprintf("총점: %.1f점\n", result.TotalScore.ToFloat())
	output += fmt.Sprintf("방법: %s\n", result.Method)
	output += fmt.Sprintf("시나리오: %s\n", result.Scenario)
	output += "\n📊 상세 점수:\n"
	for _, mt := range metadata.AllMetadataTypes() {
		if rawScore, exists := result.RawScores[mt]; exists {
			weight := result.Weights[mt]
			weighted := result.WeightedScores[mt]
			output += fmt.Sprintf("  %-20s: %.1f점 (가중치: %.1f%%) → %.1f점\n",
				mt.KoreanName(), rawScore.ToFloat(), weight.ToFloat()*100, weighted.ToFloat())
		}
	}
	return output
}
