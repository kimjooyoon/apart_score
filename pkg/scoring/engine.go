package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
)

// NewDefaultScorer creates a new scorer with the specified method.
func NewDefaultScorer(method StrategyType) *DefaultScorer {
	return &DefaultScorer{method: method}
}

type DefaultScorer struct {
	method StrategyType
}

func (s *DefaultScorer) Calculate(scores map[metadata.MetadataType]shared.ScoreValue, weights map[metadata.MetadataType]shared.Weight) (ScoreResult, error) {
	if err := s.validateWeights(weights); err != nil {
		return ScoreResult{}, err
	}
	result := ScoreResult{
		Method: s.method,
	}
	var totalWeightedSum float64
	var totalWeight shared.Weight
	for _, mt := range shared.FastAllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		weightedScore := shared.ScoreValue(int64(rawScore) * int64(weight) / shared.WeightScale)
		result.RawScores[mt] = rawScore
		result.Weights[mt] = weight
		result.WeightedScores[mt] = weightedScore.ToFloat()
		totalWeightedSum += weightedScore.ToFloat()
		totalWeight += weight
	}
	if totalWeight > 0 {
		result.TotalScore = totalWeightedSum
	}
	return result, nil
}

func (s *DefaultScorer) validateWeights(weights map[metadata.MetadataType]shared.Weight) error {
	var totalWeight shared.Weight
	for _, mt := range shared.FastAllMetadataTypes() {
		weight := weights[mt]
		if weight < 0 || weight > shared.WeightScale {
			return &ValidationError{
				Field:   mt.String(),
				Message: "가중치는 0에서 1000 사이여야 합니다",
			}
		}
		totalWeight += weight
	}
	if totalWeight < shared.WeightScale-1 || totalWeight > shared.WeightScale+1 {
		return &ValidationError{
			Field:   "total_weight",
			Message: "모든 가중치의 합계는 1000이어야 합니다",
		}
	}
	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message + " (필드: " + e.Field + ")"
}
