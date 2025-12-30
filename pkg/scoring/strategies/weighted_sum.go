package strategies

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"apart_score/pkg/scoring/context"
	"fmt"
)

type WeightedSumStrategy struct{}

func NewWeightedSumStrategy() *WeightedSumStrategy {
	return &WeightedSumStrategy{}
}
func (s *WeightedSumStrategy) Name() string {
	return "Weighted Sum"
}
func (s *WeightedSumStrategy) Calculate(scores map[metadata.MetadataType]scoring.ScoreValue,
	weights map[metadata.MetadataType]scoring.Weight,
	ctx *context.ContextData) (*scoring.ScoreResult, error) {
	if err := s.ValidateInputs(scores, weights); err != nil {
		return nil, err
	}
	result := &scoring.ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]scoring.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]scoring.ScoreValue),
		Weights:        make(map[metadata.MetadataType]scoring.Weight),
		Method:         scoring.MethodWeightedSum,
	}
	var totalWeightedSum scoring.ScoreValue
	var totalWeight scoring.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		weightedScore := rawScore * scoring.ScoreValue(weight)
		result.RawScores[mt] = rawScore
		result.Weights[mt] = weight
		result.WeightedScores[mt] = weightedScore
		totalWeightedSum += weightedScore
		totalWeight += weight
	}
	if totalWeight > 0 {
		result.TotalScore = totalWeightedSum / scoring.ScoreValue(totalWeight)
	}
	return result, nil
}
func (s *WeightedSumStrategy) ValidateInputs(scores map[metadata.MetadataType]scoring.ScoreValue,
	weights map[metadata.MetadataType]scoring.Weight) error {
	for mt, score := range scores {
		if score < 0 || score > 100 {
			return fmt.Errorf("잘못된 점수 범위 (%s: %.1f)", mt.String(), score)
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
func (s *WeightedSumStrategy) GetRequiredContext() []context.ContextType {
	return []context.ContextType{}
}
func (s *WeightedSumStrategy) Description() string {
	return "각 메타데이터의 점수에 가중치를 곱한 후 합계를 계산하는 기본 전략입니다."
}
