package scoring

import "apart_score/pkg/metadata"

// Scorer는 점수 계산 인터페이스
type Scorer interface {
	Calculate(scores map[metadata.MetadataType]ScoreValue, weights map[metadata.MetadataType]Weight) (*ScoreResult, error)
	GetDefaultWeights() map[metadata.MetadataType]Weight
	ValidateWeights(weights map[metadata.MetadataType]Weight) error
}

// DefaultScorer는 기본 가중치 합계 스코어러 구현
type DefaultScorer struct {
	method ScoringMethod
}

// NewDefaultScorer는 새로운 기본 스코어러를 생성합니다.
func NewDefaultScorer(method ScoringMethod) *DefaultScorer {
	return &DefaultScorer{method: method}
}

// Calculate는 점수를 계산합니다.
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

	// 각 메타데이터에 대해 가중치 적용 점수 계산
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

	// 최종 점수 계산 (가중치 합계 정규화)
	if totalWeight > 0 {
		result.TotalScore = totalWeightedSum / ScoreValue(totalWeight)
	}

	// 등급 및 백분위수 계산
	result.Grade = calculateGrade(result.TotalScore)

	return result, nil
}

// GetDefaultWeights는 기본 가중치를 반환합니다.
func (s *DefaultScorer) GetDefaultWeights() map[metadata.MetadataType]Weight {
	weights := map[metadata.MetadataType]Weight{
		metadata.FloorLevel:            0.08, // 8%
		metadata.DistanceToStation:     0.15, // 15%
		metadata.ElevatorPresence:      0.07, // 7%
		metadata.ConstructionYear:      0.10, // 10%
		metadata.ConstructionCompany:   0.08, // 8%
		metadata.ApartmentSize:         0.08, // 8%
		metadata.NearbyAmenities:       0.10, // 10%
		metadata.TransportationAccess:  0.12, // 12%
		metadata.SchoolDistrict:        0.08, // 8%
		metadata.CrimeRate:             0.06, // 6%
		metadata.GreenSpaceRatio:       0.04, // 4%
		metadata.Parking:               0.06, // 6%
		metadata.MaintenanceFee:        0.05, // 5%
		metadata.HeatingSystem:         0.03, // 3%
	}
	return NormalizeWeights(weights)
}

// ValidateWeights는 가중치의 유효성을 검증합니다.
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

	// 가중치 합계가 1.0에 가까운지 확인 (약간의 오차 허용)
	if totalWeight < 0.99 || totalWeight > 1.01 {
		return &ValidationError{
			Field:   "total_weight",
			Message: "모든 가중치의 합계는 1.0이어야 합니다",
		}
	}

	return nil
}

// ValidationError는 검증 오류를 나타냅니다.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message + " (필드: " + e.Field + ")"
}

// calculateGrade는 점수에 따라 등급을 계산합니다.
func calculateGrade(score ScoreValue) Grade {
	switch {
	case score >= 90:
		return GradeA
	case score >= 80:
		return GradeB
	case score >= 70:
		return GradeC
	case score >= 60:
		return GradeD
	default:
		return GradeF
	}
}
