package core

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
)

type Scorer interface {
	Calculate(scores map[metadata.MetadataType]shared.ScoreValue, weights map[metadata.MetadataType]shared.Weight) (ScoreResult, error)
	GetDefaultWeights() map[metadata.MetadataType]shared.Weight
	ValidateWeights(weights map[metadata.MetadataType]shared.Weight) error
}
type DefaultScorer struct {
	method ScoringMethod
}

func NewDefaultScorer(method ScoringMethod) *DefaultScorer {
	return &DefaultScorer{method: method}
}
func (s *DefaultScorer) Calculate(scores map[metadata.MetadataType]shared.ScoreValue, weights map[metadata.MetadataType]shared.Weight) (ScoreResult, error) {
	if err := s.ValidateWeights(weights); err != nil {
		return ScoreResult{}, err
	}
	result := ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]shared.ScoreValue),
		RawScores:      make(map[metadata.MetadataType]shared.ScoreValue),
		Weights:        make(map[metadata.MetadataType]shared.Weight),
		Method:         s.method,
	}
	var totalWeightedSum shared.ScoreValue
	var totalWeight shared.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		rawScore := scores[mt]
		weight := weights[mt]
		// 정수 기반 계산: (score * weight) / WeightScale
		weightedScore := shared.ScoreValue(int64(rawScore) * int64(weight) / shared.WeightScale)
		result.RawScores[mt] = rawScore
		result.Weights[mt] = weight
		result.WeightedScores[mt] = weightedScore
		totalWeightedSum += weightedScore
		totalWeight += weight
	}
	if totalWeight > 0 {
		// 정수 기반 평균 계산: totalWeightedSum / totalWeight (이미 정규화됨)
		result.TotalScore = totalWeightedSum
	}
	return result, nil
}
func (s *DefaultScorer) GetDefaultWeights() map[metadata.MetadataType]shared.Weight {
	weights := map[metadata.MetadataType]shared.Weight{
		metadata.FloorLevel:           shared.WeightFromFloat(0.08),
		metadata.DistanceToStation:    shared.WeightFromFloat(0.15),
		metadata.ElevatorPresence:     shared.WeightFromFloat(0.07),
		metadata.ConstructionYear:     shared.WeightFromFloat(0.10),
		metadata.ConstructionCompany:  shared.WeightFromFloat(0.08),
		metadata.ApartmentSize:        shared.WeightFromFloat(0.08),
		metadata.NearbyAmenities:      shared.WeightFromFloat(0.10),
		metadata.TransportationAccess: shared.WeightFromFloat(0.12),
		metadata.SchoolDistrict:       shared.WeightFromFloat(0.08),
		metadata.CrimeRate:            shared.WeightFromFloat(0.06),
		metadata.GreenSpaceRatio:      shared.WeightFromFloat(0.04),
		metadata.Parking:              shared.WeightFromFloat(0.06),
		metadata.MaintenanceFee:       shared.WeightFromFloat(0.05),
		metadata.HeatingSystem:        shared.WeightFromFloat(0.03),
	}
	return shared.NormalizeWeights(weights)
}
func (s *DefaultScorer) ValidateWeights(weights map[metadata.MetadataType]shared.Weight) error {
	var totalWeight shared.Weight
	for _, mt := range metadata.AllMetadataTypes() {
		weight := weights[mt]
		if weight < 0 || weight > shared.WeightScale {
			return &ValidationError{
				Field:   mt.String(),
				Message: "가중치는 0에서 1000 사이여야 합니다",
			}
		}
		totalWeight += weight
	}
	// 정수 기반 검증: 합계가 WeightScale ±1 범위 내여야 함
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
