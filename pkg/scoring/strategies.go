package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
	"fmt"
	"math"
	"sort"
)

type StrategyType string

const (
	StrategyWeightedSum   StrategyType = "weighted_sum"
	StrategyGeometricMean StrategyType = "geometric_mean"
	StrategyMinMax        StrategyType = "min_max"
	StrategyHarmonicMean  StrategyType = "harmonic_mean"
)

type ApartmentData struct {
	ID       string                                     `json:"id"`
	Name     string                                     `json:"name"`
	Scores   map[metadata.MetadataType]shared.ScoreValue `json:"scores"`
	Location string                                     `json:"location"`
}
type RankingResult struct {
	Apartment  ApartmentData                    `json:"apartment"`
	Score      shared.ScoreValue                `json:"score"`
	Rank       int                              `json:"rank"`
	Percentile float64                          `json:"percentile"`
	Method     ScoringMethod                    `json:"method"`
	Weights    map[metadata.MetadataType]shared.Weight `json:"weights"`
}
type RankingsSummary struct {
	TotalApartments int             `json:"total_apartments"`
	Strategy        StrategyType    `json:"strategy"`
	TopRanked       []RankingResult `json:"top_ranked"`
	ScoreRange      struct {
		Min shared.ScoreValue `json:"min"`
		Max shared.ScoreValue `json:"max"`
		Avg shared.ScoreValue `json:"avg"`
	} `json:"score_range"`
}

func CalculateWithStrategy(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight,
	strategy StrategyType) (ScoreResult, error) {
	if err := validateStrategyInputs(scores, weights); err != nil {
		return ScoreResult{}, err
	}
	switch strategy {
	case StrategyWeightedSum:
		return calculateWeightedSum(scores, weights)
	case StrategyGeometricMean:
		return calculateGeometricMean(scores, weights)
	case StrategyMinMax:
		return calculateMinMax(scores, weights)
	case StrategyHarmonicMean:
		return calculateHarmonicMean(scores, weights)
	default:
		return ScoreResult{}, fmt.Errorf("지원하지 않는 전략: %s", strategy)
	}
}
func calculateWeightedSum(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight) (ScoreResult, error) {
	result := ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]shared.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]shared.ScoreValue),
		Weights:        make(map[metadata.MetadataType]shared.Weight),
		Method:         MethodWeightedSum,
	}
	var totalWeightedSum shared.ScoreValue
	var totalWeight shared.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		weightedScore := rawScore *shared.ScoreValue(weight)
		result.RawScores[mt] = rawScore
		result.Weights[mt] = weight
		result.WeightedScores[mt] = weightedScore
		totalWeightedSum += weightedScore
		totalWeight += weight
	}
	if totalWeight > 0 {
		result.TotalScore = totalWeightedSum /shared.ScoreValue(totalWeight)
	}
	return result, nil
}
func calculateGeometricMean(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight) (ScoreResult, error) {
	result := ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]shared.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]shared.ScoreValue),
		Weights:        make(map[metadata.MetadataType]shared.Weight),
		Method:         MethodGeometricMean,
	}
	minScore := shared.ScoreValueFromFloat(0.1)
	var logSum float64
	var totalWeight shared.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		if rawScore < minScore {
			rawScore = minScore
		}
		weightedLog := math.Log(float64(rawScore)) * float64(weight)
		result.RawScores[mt] = scores[mt]
		result.Weights[mt] = weight
		result.WeightedScores[mt] =shared.ScoreValue(math.Exp(weightedLog))
		logSum += weightedLog
		totalWeight += weight
	}
	if totalWeight > 0 {
		result.TotalScore =shared.ScoreValue(math.Exp(logSum / float64(totalWeight)))
	}
	return result, nil
}
func calculateMinMax(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight) (ScoreResult, error) {
	result := ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]shared.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]shared.ScoreValue),
		Weights:        make(map[metadata.MetadataType]shared.Weight),
		Method:         MethodMinMax,
	}
	minScore := shared.ScoreValueFromFloat(100.0)
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		result.RawScores[mt] = rawScore
		result.Weights[mt] = weight
		result.WeightedScores[mt] = rawScore *shared.ScoreValue(weight)
		if rawScore < minScore {
			minScore = rawScore
		}
	}
	result.TotalScore = minScore
	return result, nil
}
func calculateHarmonicMean(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight) (ScoreResult, error) {
	result := ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]shared.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]shared.ScoreValue),
		Weights:        make(map[metadata.MetadataType]shared.Weight),
		Method:         MethodHarmonicMean,
	}
	minScore := shared.ScoreValueFromFloat(0.1)
	var weightedHarmonicSum float64
	var totalWeight shared.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		if rawScore < minScore {
			rawScore = minScore
		}
		weightedHarmonic := float64(weight) / float64(rawScore)
		result.RawScores[mt] = scores[mt]
		result.Weights[mt] = weight
		result.WeightedScores[mt] =shared.ScoreValue(float64(weight) * float64(rawScore))
		weightedHarmonicSum += weightedHarmonic
		totalWeight += weight
	}
	if weightedHarmonicSum > 0 && totalWeight > 0 {
		result.TotalScore =shared.ScoreValue(float64(totalWeight) / weightedHarmonicSum)
	}
	return result, nil
}
func validateStrategyInputs(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight) error {
	for mt, score := range scores {
		if score < 0 || score > 100*shared.ScoreScale {
			return fmt.Errorf("잘못된 점수 범위 (%s: %.1f)", mt.String(), score.ToFloat())
		}
	}
	var totalWeight shared.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		weight := weights[mt]
		if weight < 0 || weight > shared.WeightScale {
			return fmt.Errorf("잘못된 가중치 범위 (%s: %.3f)", mt.String(), weight.ToFloat())
		}
		totalWeight += weight
	}
	if totalWeight < shared.WeightScale-1 || totalWeight > shared.WeightScale+1 {
		return fmt.Errorf("가중치 합계가 1000이 아닙니다 (현재: %d)", totalWeight)
	}
	return nil
}
func GetAvailableStrategies() []StrategyType {
	return []StrategyType{
		StrategyWeightedSum,
		StrategyGeometricMean,
		StrategyMinMax,
		StrategyHarmonicMean,
	}
}
func GetStrategyDescription(strategy StrategyType) string {
	switch strategy {
	case StrategyWeightedSum:
		return "각 메타데이터의 점수에 가중치를 곱한 후 합계를 계산하는 기본 전략입니다."
	case StrategyGeometricMean:
		return "모든 요소가 균형을 이루어야 하는 경우에 적합한 전략입니다. 하나의 낮은 점수가 전체 점수를 크게 낮춥니다."
	case StrategyMinMax:
		return "모든 요소가 일정 수준 이상이어야 하는 경우에 적합합니다. 가장 낮은 점수가 전체 점수를 결정합니다."
	case StrategyHarmonicMean:
		return "낮은 점수에 매우 민감하게 반응하는 전략입니다. 모든 요소가 고르게 중요할 때 사용합니다."
	default:
		return "알 수 없는 전략입니다."
	}
}
func CalculateRankings(apartments []ApartmentData, weights map[metadata.MetadataType]shared.Weight, strategy StrategyType) (*RankingsSummary, error) {
	if len(apartments) == 0 {
		return nil, fmt.Errorf("순위를 매길 아파트가 없습니다")
	}
	if err := validateStrategyInputs(apartments[0].Scores, weights); err != nil {
		return nil, fmt.Errorf("입력 검증 실패: %w", err)
	}
	var rankings []RankingResult
	var totalScore shared.ScoreValue
	minScore := shared.ScoreValueFromFloat(100.0)
	maxScore := shared.ScoreValue(0.0)
	for _, apt := range apartments {
		result, err := CalculateWithStrategy(apt.Scores, weights, strategy)
		if err != nil {
			return nil, fmt.Errorf("아파트 %s 점수 계산 실패: %w", apt.ID, err)
		}
		ranking := RankingResult{
			Apartment: apt,
			Score:     result.TotalScore,
			Method:    result.Method,
			Weights:   result.Weights,
		}
		rankings = append(rankings, ranking)
		totalScore += result.TotalScore
		if result.TotalScore < minScore {
			minScore = result.TotalScore
		}
		if result.TotalScore > maxScore {
			maxScore = result.TotalScore
		}
	}
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})
	for i := range rankings {
		rankings[i].Rank = i + 1
		if maxScore > minScore {
			rankings[i].Percentile = float64(rankings[i].Score-minScore) / float64(maxScore-minScore) * 100.0
		} else {
			rankings[i].Percentile = 100.0
		}
	}
	summary := &RankingsSummary{
		TotalApartments: len(apartments),
		Strategy:        strategy,
		TopRanked:       rankings,
	}
	summary.ScoreRange.Min = minScore
	summary.ScoreRange.Max = maxScore
	summary.ScoreRange.Avg =shared.ScoreValue(int64(totalScore) / int64(len(apartments)))
	return summary, nil
}
func FormatRankings(summary *RankingsSummary, limit int) string {
	if summary == nil {
		return "순위 데이터가 없습니다."
	}
	output := fmt.Sprintf("🏆 아파트 순위표 (%s 전략)\n", GetStrategyDescription(summary.Strategy))
	output += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	output += fmt.Sprintf("총 아파트 수: %d개\n", summary.TotalApartments)
	output += fmt.Sprintf("점수 범위: %.1f - %.1f (평균: %.1f)\n", summary.ScoreRange.Min.ToFloat(), summary.ScoreRange.Max.ToFloat(), summary.ScoreRange.Avg.ToFloat())
	output += "\n📊 순위 결과:\n"
	displayCount := len(summary.TopRanked)
	if limit > 0 && limit < displayCount {
		displayCount = limit
	}
	for i := 0; i < displayCount; i++ {
		ranking := summary.TopRanked[i]
		rankEmoji := getRankEmoji(i + 1)
		output += fmt.Sprintf("%s %d위: %s (%.1f점, %.1f%%)\n",
			rankEmoji,
			ranking.Rank,
			ranking.Apartment.Name,
			ranking.Score.ToFloat(),
			ranking.Percentile)
	}
	if displayCount < len(summary.TopRanked) {
		output += fmt.Sprintf("\n... 외 %d개 아파트", len(summary.TopRanked)-displayCount)
	}
	return output
}
func getRankEmoji(rank int) string {
	switch rank {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		return "🏠"
	}
}
