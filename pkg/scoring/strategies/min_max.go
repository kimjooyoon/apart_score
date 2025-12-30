package strategies

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"apart_score/pkg/scoring/context"
)

type MinMaxStrategy struct{}

func NewMinMaxStrategy() *MinMaxStrategy {
	return &MinMaxStrategy{}
}
func (s *MinMaxStrategy) Name() string {
	return "Min-Max (Minimum Priority)"
}
func (s *MinMaxStrategy) Calculate(scores map[metadata.MetadataType]scoring.ScoreValue,
	weights map[metadata.MetadataType]scoring.Weight,
	ctx *context.ContextData) (*scoring.ScoreResult, error) {
	if err := s.ValidateInputs(scores, weights); err != nil {
		return nil, err
	}
	result := &scoring.ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]scoring.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]scoring.ScoreValue),
		Weights:        make(map[metadata.MetadataType]scoring.Weight),
		Method:         scoring.MethodMinMax,
	}
	minScore := scoring.ScoreValue(100.0)
	maxScore := scoring.ScoreValue(0.0)
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		result.RawScores[mt] = rawScore
		result.Weights[mt] = weight
		result.WeightedScores[mt] = rawScore * scoring.ScoreValue(weight)
		if rawScore < minScore {
			minScore = rawScore
		}
		if rawScore > maxScore {
			maxScore = rawScore
		}
	}
	result.TotalScore = minScore
	return result, nil
}
func (s *MinMaxStrategy) ValidateInputs(scores map[metadata.MetadataType]scoring.ScoreValue,
	weights map[metadata.MetadataType]scoring.Weight) error {
	ws := &WeightedSumStrategy{}
	return ws.ValidateInputs(scores, weights)
}
func (s *MinMaxStrategy) GetRequiredContext() []context.ContextType {
	return []context.ContextType{}
}
func (s *MinMaxStrategy) Description() string {
	return "모든 요소가 일정 수준 이상이어야 하는 경우에 적합합니다. 가장 낮은 점수가 전체 점수를 결정합니다."
}
