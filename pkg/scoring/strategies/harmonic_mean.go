package strategies

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"apart_score/pkg/scoring/context"
	"fmt"
)

type HarmonicMeanStrategy struct{}

func NewHarmonicMeanStrategy() *HarmonicMeanStrategy {
	return &HarmonicMeanStrategy{}
}
func (s *HarmonicMeanStrategy) Name() string {
	return "Harmonic Mean"
}
func (s *HarmonicMeanStrategy) Calculate(scores map[metadata.MetadataType]scoring.ScoreValue,
	weights map[metadata.MetadataType]scoring.Weight,
	ctx *context.ContextData) (*scoring.ScoreResult, error) {
	if err := s.ValidateInputs(scores, weights); err != nil {
		return nil, err
	}
	result := &scoring.ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]scoring.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]scoring.ScoreValue),
		Weights:        make(map[metadata.MetadataType]scoring.Weight),
		Method:         scoring.MethodHarmonicMean,
	}
	minScore := scoring.ScoreValue(0.1)
	var weightedHarmonicSum float64
	var totalWeight scoring.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		if rawScore < minScore {
			rawScore = minScore
		}
		weightedHarmonic := float64(weight) / float64(rawScore)
		result.RawScores[mt] = scores[mt]
		result.Weights[mt] = weight
		result.WeightedScores[mt] = scoring.ScoreValue(float64(weight) * float64(rawScore))
		weightedHarmonicSum += weightedHarmonic
		totalWeight += weight
	}
	if weightedHarmonicSum > 0 && totalWeight > 0 {
		result.TotalScore = scoring.ScoreValue(float64(totalWeight) / weightedHarmonicSum)
	}
	return result, nil
}
func (s *HarmonicMeanStrategy) ValidateInputs(scores map[metadata.MetadataType]scoring.ScoreValue,
	weights map[metadata.MetadataType]scoring.Weight) error {
	for mt, score := range scores {
		if score <= 0 {
			return fmt.Errorf("조화 평균 전략은 0 이하의 점수를 허용하지 않습니다 (%s: %.1f)", mt.String(), score)
		}
	}
	var totalWeight scoring.Weight
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
func (s *HarmonicMeanStrategy) GetRequiredContext() []context.ContextType {
	return []context.ContextType{}
}
func (s *HarmonicMeanStrategy) Description() string {
	return "낮은 점수에 매우 민감하게 반응하는 전략입니다. 모든 요소가 고르게 중요할 때 사용합니다."
}
