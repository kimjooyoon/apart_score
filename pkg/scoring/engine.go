package scoring

import "apart_score/pkg/metadata"

type Scorer interface {
	Calculate(scores map[metadata.MetadataType]ScoreValue, weights map[metadata.MetadataType]Weight) (*ScoreResult, error)
	GetDefaultWeights() map[metadata.MetadataType]Weight
	ValidateWeights(weights map[metadata.MetadataType]Weight) error
}
type DefaultScorer struct {
	method ScoringMethod
}

func NewDefaultScorer(method ScoringMethod) *DefaultScorer {
	return &DefaultScorer{method: method}
}
func (s *DefaultScorer) Calculate(scores map[metadata.MetadataType]ScoreValue, weights map[metadata.MetadataType]Weight) (*ScoreResult, error) {
	if err := s.ValidateWeights(weights); err != nil {
		return nil, err
	}
	result := &ScoreResult{
		WeightedScores: make(map[metadata.MetadataType]ScoreValue),
		RawScores:      make(map[metadata.MetadataType]ScoreValue),
		Weights:        make(map[metadata.MetadataType]Weight),
		Method:         s.method,
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
func (s *DefaultScorer) GetDefaultWeights() map[metadata.MetadataType]Weight {
	weights := map[metadata.MetadataType]Weight{
		metadata.FloorLevel:           0.08,
		metadata.DistanceToStation:    0.15,
		metadata.ElevatorPresence:     0.07,
		metadata.ConstructionYear:     0.10,
		metadata.ConstructionCompany:  0.08,
		metadata.ApartmentSize:        0.08,
		metadata.NearbyAmenities:      0.10,
		metadata.TransportationAccess: 0.12,
		metadata.SchoolDistrict:       0.08,
		metadata.CrimeRate:            0.06,
		metadata.GreenSpaceRatio:      0.04,
		metadata.Parking:              0.06,
		metadata.MaintenanceFee:       0.05,
		metadata.HeatingSystem:        0.03,
	}
	return NormalizeWeights(weights)
}
func (s *DefaultScorer) ValidateWeights(weights map[metadata.MetadataType]Weight) error {
	var totalWeight Weight
	for _, mt := range metadata.AllMetadataTypes() {
		weight := weights[mt]
		if weight < 0 || weight > 1 {
			return &ValidationError{
				Field:   mt.String(),
				Message: "가중치는 0.0에서 1.0 사이여야 합니다",
			}
		}
		totalWeight += weight
	}
	if totalWeight < 0.99 || totalWeight > 1.01 {
		return &ValidationError{
			Field:   "total_weight",
			Message: "모든 가중치의 합계는 1.0이어야 합니다",
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
