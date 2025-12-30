package scoring

import (
	"apart_score/pkg/metadata"
	"fmt"
	"math"
)

// StrategyType은 계산 전략 유형을 정의합니다.
type StrategyType string

const (
	StrategyWeightedSum   StrategyType = "weighted_sum"
	StrategyGeometricMean StrategyType = "geometric_mean"
	StrategyMinMax        StrategyType = "min_max"
	StrategyHarmonicMean  StrategyType = "harmonic_mean"
)

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
