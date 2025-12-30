package scoring

import (
	"apart_score/pkg/metadata"
	"fmt"
	"math"
	"sort"
)

// StrategyType은 계산 전략 유형을 정의합니다.
type StrategyType string

const (
	StrategyWeightedSum   StrategyType = "weighted_sum"
	StrategyGeometricMean StrategyType = "geometric_mean"
	StrategyMinMax        StrategyType = "min_max"
	StrategyHarmonicMean  StrategyType = "harmonic_mean"
)

// ApartmentData는 아파트 정보를 담는 구조체입니다.
type ApartmentData struct {
	ID       string                                `json:"id"`
	Name     string                                `json:"name"`     // 아파트 이름
	Scores   map[metadata.MetadataType]ScoreValue `json:"scores"`   // 각 메타데이터 점수
	Location string                                `json:"location"` // 위치 정보
}

// RankingResult는 순위 결과를 담는 구조체입니다.
type RankingResult struct {
	Apartment   ApartmentData `json:"apartment"`
	Score       ScoreValue    `json:"score"`
	Rank        int           `json:"rank"`
	Percentile  float64       `json:"percentile"` // 백분위수 (0-100)
	Method      ScoringMethod `json:"method"`
	Weights     map[metadata.MetadataType]Weight `json:"weights"`
}

// RankingsSummary는 순위 요약 정보를 담는 구조체입니다.
type RankingsSummary struct {
	TotalApartments int         `json:"total_apartments"`
	Strategy        StrategyType `json:"strategy"`
	TopRanked       []RankingResult `json:"top_ranked"`
	ScoreRange      struct {
		Min ScoreValue `json:"min"`
		Max ScoreValue `json:"max"`
		Avg ScoreValue `json:"avg"`
	} `json:"score_range"`
}

// CalculateWithStrategy는 지정된 전략으로 점수를 계산합니다.
func CalculateWithStrategy(scores map[metadata.MetadataType]ScoreValue,
	weights map[metadata.MetadataType]Weight,
	strategy StrategyType) (*ScoreResult, error) {

	if err := validateStrategyInputs(scores, weights); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("지원하지 않는 전략: %s", strategy)
	}
}

// calculateWeightedSum은 가중치 합계 전략으로 계산합니다.
func calculateWeightedSum(scores map[metadata.MetadataType]ScoreValue,
	weights map[metadata.MetadataType]Weight) (*ScoreResult, error) {

	result := &ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]ScoreValue),
		RawScores:      make(map[metadata.MetadataType]ScoreValue),
		Weights:        make(map[metadata.MetadataType]Weight),
		Method:         MethodWeightedSum,
	}

	var totalWeightedSum ScoreValue
	var totalWeight Weight

	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		weightedScore := rawScore * ScoreValue(weight)

		result.RawScores[mt] = rawScore
		result.Weights[mt] = weight
		result.WeightedScores[mt] = weightedScore

		totalWeightedSum += weightedScore
		totalWeight += weight
	}

	if totalWeight > 0 {
		result.TotalScore = totalWeightedSum / ScoreValue(totalWeight)
	}

	return result, nil
}

// calculateGeometricMean은 기하 평균 전략으로 계산합니다.
func calculateGeometricMean(scores map[metadata.MetadataType]ScoreValue,
	weights map[metadata.MetadataType]Weight) (*ScoreResult, error) {

	result := &ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]ScoreValue),
		RawScores:      make(map[metadata.MetadataType]ScoreValue),
		Weights:        make(map[metadata.MetadataType]Weight),
		Method:         MethodGeometricMean,
	}

	minScore := ScoreValue(0.1)
	var logSum float64
	var totalWeight Weight

	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]

		if rawScore < minScore {
			rawScore = minScore
		}

		weightedLog := math.Log(float64(rawScore)) * float64(weight)

		result.RawScores[mt] = scores[mt]
		result.Weights[mt] = weight
		result.WeightedScores[mt] = ScoreValue(math.Exp(weightedLog))

		logSum += weightedLog
		totalWeight += weight
	}

	if totalWeight > 0 {
		result.TotalScore = ScoreValue(math.Exp(logSum / float64(totalWeight)))
	}

	return result, nil
}

// calculateMinMax는 최소값 우선 전략으로 계산합니다.
func calculateMinMax(scores map[metadata.MetadataType]ScoreValue,
	weights map[metadata.MetadataType]Weight) (*ScoreResult, error) {

	result := &ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]ScoreValue),
		RawScores:      make(map[metadata.MetadataType]ScoreValue),
		Weights:        make(map[metadata.MetadataType]Weight),
		Method:         MethodMinMax,
	}

	minScore := ScoreValue(100.0)

	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]

		result.RawScores[mt] = rawScore
		result.Weights[mt] = weight
		result.WeightedScores[mt] = rawScore * ScoreValue(weight)

		if rawScore < minScore {
			minScore = rawScore
		}
	}

	result.TotalScore = minScore
	return result, nil
}

// calculateHarmonicMean은 조화 평균 전략으로 계산합니다.
func calculateHarmonicMean(scores map[metadata.MetadataType]ScoreValue,
	weights map[metadata.MetadataType]Weight) (*ScoreResult, error) {

	result := &ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]ScoreValue),
		RawScores:      make(map[metadata.MetadataType]ScoreValue),
		Weights:        make(map[metadata.MetadataType]Weight),
		Method:         MethodHarmonicMean,
	}

	minScore := ScoreValue(0.1)
	var weightedHarmonicSum float64
	var totalWeight Weight

	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]

		if rawScore < minScore {
			rawScore = minScore
		}

		weightedHarmonic := float64(weight) / float64(rawScore)

		result.RawScores[mt] = scores[mt]
		result.Weights[mt] = weight
		result.WeightedScores[mt] = ScoreValue(float64(weight) * float64(rawScore))

		weightedHarmonicSum += weightedHarmonic
		totalWeight += weight
	}

	if weightedHarmonicSum > 0 && totalWeight > 0 {
		result.TotalScore = ScoreValue(float64(totalWeight) / weightedHarmonicSum)
	}

	return result, nil
}

// validateStrategyInputs는 전략 계산을 위한 입력값을 검증합니다.
func validateStrategyInputs(scores map[metadata.MetadataType]ScoreValue,
	weights map[metadata.MetadataType]Weight) error {

	for mt, score := range scores {
		if score < 0 || score > 100 {
			return fmt.Errorf("잘못된 점수 범위 (%s: %.1f)", mt.String(), score)
		}
	}

	var totalWeight Weight
	for _, mt := range metadata.AllMetadataTypes() {
		weight := weights[mt]
		if weight < 0 || weight > 1 {
			return fmt.Errorf("잘못된 가중치 범위 (%s: %.3f)", mt.String(), weight)
		}
		totalWeight += weight
	}

	if totalWeight < 0.99 || totalWeight > 1.01 {
		return fmt.Errorf("가중치 합계가 1.0이 아닙니다 (현재: %.3f)", totalWeight)
	}

	return nil
}

// GetAvailableStrategies는 사용 가능한 전략 목록을 반환합니다.
func GetAvailableStrategies() []StrategyType {
	return []StrategyType{
		StrategyWeightedSum,
		StrategyGeometricMean,
		StrategyMinMax,
		StrategyHarmonicMean,
	}
}

// GetStrategyDescription은 전략에 대한 설명을 반환합니다.
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

// CalculateRankings는 여러 아파트의 점수를 계산하고 순위를 매깁니다.
func CalculateRankings(apartments []ApartmentData, weights map[metadata.MetadataType]Weight, strategy StrategyType) (*RankingsSummary, error) {
	if len(apartments) == 0 {
		return nil, fmt.Errorf("순위를 매길 아파트가 없습니다")
	}

	if err := validateStrategyInputs(apartments[0].Scores, weights); err != nil {
		return nil, fmt.Errorf("입력 검증 실패: %w", err)
	}

	// 각 아파트의 점수 계산
	var rankings []RankingResult
	var totalScore ScoreValue
	minScore := ScoreValue(100.0)
	maxScore := ScoreValue(0.0)

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

	// 점수 기준 내림차순 정렬 (높은 점수가 1위)
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	// 순위 부여 및 백분위수 계산
	for i := range rankings {
		rankings[i].Rank = i + 1

		// 백분위수 계산 (높은 점수일수록 높은 백분위수)
		if maxScore > minScore {
			rankings[i].Percentile = float64(rankings[i].Score-minScore) / float64(maxScore-minScore) * 100.0
		} else {
			rankings[i].Percentile = 100.0 // 모든 점수가 같을 경우
		}
	}

	// 요약 정보 생성
	summary := &RankingsSummary{
		TotalApartments: len(apartments),
		Strategy:        strategy,
		TopRanked:       rankings,
	}
	summary.ScoreRange.Min = minScore
	summary.ScoreRange.Max = maxScore
	summary.ScoreRange.Avg = totalScore / ScoreValue(len(apartments))

	return summary, nil
}

// FormatRankings는 순위 결과를 읽기 쉽게 포맷팅합니다.
func FormatRankings(summary *RankingsSummary, limit int) string {
	if summary == nil {
		return "순위 데이터가 없습니다."
	}

	output := fmt.Sprintf("🏆 아파트 순위표 (%s 전략)\n", GetStrategyDescription(summary.Strategy))
	output += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	output += fmt.Sprintf("총 아파트 수: %d개\n", summary.TotalApartments)
	output += fmt.Sprintf("점수 범위: %.1f - %.1f (평균: %.1f)\n", summary.ScoreRange.Min, summary.ScoreRange.Max, summary.ScoreRange.Avg)
	output += fmt.Sprintf("\n📊 순위 결과:\n")

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
			ranking.Score,
			ranking.Percentile)
	}

	if displayCount < len(summary.TopRanked) {
		output += fmt.Sprintf("\n... 외 %d개 아파트", len(summary.TopRanked)-displayCount)
	}

	return output
}

// getRankEmoji는 순위에 따른 이모지를 반환합니다.
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
