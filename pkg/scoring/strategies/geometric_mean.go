package strategies

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"apart_score/pkg/scoring/context"
	"fmt"
	"math"
)

type GeometricMeanStrategy struct{}

func NewGeometricMeanStrategy() *GeometricMeanStrategy {
	return &GeometricMeanStrategy{}
}
func (s *GeometricMeanStrategy) Name() string {
	return "Geometric Mean"
}
func (s *GeometricMeanStrategy) Calculate(scores map[metadata.MetadataType]scoring.ScoreValue,
	weights map[metadata.MetadataType]scoring.Weight,
	ctx *context.ContextData) (*scoring.ScoreResult, error) {
	if err := s.ValidateInputs(scores, weights); err != nil {
		return nil, err
	}
	result := &scoring.ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]scoring.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]scoring.ScoreValue),
		Weights:        make(map[metadata.MetadataType]scoring.Weight),
		Method:         scoring.MethodGeometricMean,
	}
	minScore := scoring.ScoreValue(0.1)
	var logSum float64
	var totalWeight scoring.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		if rawScore < minScore {
			rawScore = minScore
		}
		weightedLog := math.Log(float64(rawScore)) * float64(weight)
		result.RawScores[mt] = scores[mt]
		result.Weights[mt] = weight
		result.WeightedScores[mt] = scoring.ScoreValue(math.Exp(weightedLog))
		logSum += weightedLog
		totalWeight += weight
	}
	if totalWeight > 0 {
		result.TotalScore = scoring.ScoreValue(math.Exp(logSum / float64(totalWeight)))
	}
	return result, nil
}
func (s *GeometricMeanStrategy) ValidateInputs(scores map[metadata.MetadataType]scoring.ScoreValue,
	weights map[metadata.MetadataType]scoring.Weight) error {
	for mt, score := range scores {
		if score < 0 {
			return fmt.Errorf("기하 평균 전략은 음수 점수를 허용하지 않습니다 (%s: %.1f)", mt.String(), score)
		}
		if score == 0 {
			return fmt.Errorf("기하 평균 전략은 0점를 허용하지 않습니다 (%s: %.1f)", mt.String(), score)
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
func (s *GeometricMeanStrategy) GetRequiredContext() []context.ContextType {
	return []context.ContextType{}
}
func (s *GeometricMeanStrategy) Description() string {
	return "모든 요소가 균형을 이루어야 하는 경우에 적합한 전략입니다. 하나의 낮은 점수가 전체 점수를 크게 낮춥니다."
}
