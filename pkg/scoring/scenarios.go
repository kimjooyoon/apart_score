package scoring

import "apart_score/pkg/metadata"

// NormalizeWeights는 가중치 맵의 합계를 1.0으로 정규화합니다.
func NormalizeWeights(weights map[metadata.MetadataType]Weight) map[metadata.MetadataType]Weight {
	var total Weight
	for _, w := range weights {
		total += w
	}

	if total == 0 {
		return weights
	}

	normalized := make(map[metadata.MetadataType]Weight)
	for mt, w := range weights {
		normalized[mt] = w / total
	}

	return normalized
}

// GetScenarioWeights는 시나리오에 따른 가중치를 반환합니다.
func GetScenarioWeights(scenario ScoringScenario) map[metadata.MetadataType]Weight {
	var weights map[metadata.MetadataType]Weight

	switch scenario {
	case ScenarioTransportation:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:            0.05, // 5%
			metadata.DistanceToStation:     0.25, // 25%
			metadata.ElevatorPresence:      0.05, // 5%
			metadata.ConstructionYear:      0.08, // 8%
			metadata.ConstructionCompany:   0.05, // 5%
			metadata.ApartmentSize:         0.07, // 7%
			metadata.NearbyAmenities:       0.10, // 10%
			metadata.TransportationAccess:  0.20, // 20%
			metadata.SchoolDistrict:        0.05, // 5%
			metadata.CrimeRate:             0.05, // 5%
			metadata.GreenSpaceRatio:       0.02, // 2%
			metadata.Parking:               0.05, // 5%
			metadata.MaintenanceFee:        0.05, // 5%
			metadata.HeatingSystem:         0.03, // 3%
		}

	case ScenarioEducation:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:            0.08, // 8%
			metadata.DistanceToStation:     0.08, // 8%
			metadata.ElevatorPresence:      0.07, // 7%
			metadata.ConstructionYear:      0.10, // 10%
			metadata.ConstructionCompany:   0.08, // 8%
			metadata.ApartmentSize:         0.08, // 8%
			metadata.NearbyAmenities:       0.08, // 8%
			metadata.TransportationAccess:  0.08, // 8%
			metadata.SchoolDistrict:        0.20, // 20%
			metadata.CrimeRate:             0.08, // 8%
			metadata.GreenSpaceRatio:       0.05, // 5%
			metadata.Parking:               0.05, // 5%
			metadata.MaintenanceFee:        0.04, // 4%
			metadata.HeatingSystem:         0.03, // 3%
		}

	case ScenarioCostEffective:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:            0.08, // 8%
			metadata.DistanceToStation:     0.12, // 12%
			metadata.ElevatorPresence:      0.07, // 7%
			metadata.ConstructionYear:      0.08, // 8%
			metadata.ConstructionCompany:   0.05, // 5%
			metadata.ApartmentSize:         0.12, // 12%
			metadata.NearbyAmenities:       0.10, // 10%
			metadata.TransportationAccess:  0.08, // 8%
			metadata.SchoolDistrict:        0.08, // 8%
			metadata.CrimeRate:             0.08, // 8%
			metadata.GreenSpaceRatio:       0.05, // 5%
			metadata.Parking:               0.06, // 6%
			metadata.MaintenanceFee:        0.10, // 10%
			metadata.HeatingSystem:         0.03, // 3%
		}

	case ScenarioFamilyFriendly:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:            0.10, // 10%
			metadata.DistanceToStation:     0.08, // 8%
			metadata.ElevatorPresence:      0.12, // 12%
			metadata.ConstructionYear:      0.10, // 10%
			metadata.ConstructionCompany:   0.08, // 8%
			metadata.ApartmentSize:         0.12, // 12%
			metadata.NearbyAmenities:       0.08, // 8%
			metadata.TransportationAccess:  0.05, // 5%
			metadata.SchoolDistrict:        0.12, // 12%
			metadata.CrimeRate:             0.08, // 8%
			metadata.GreenSpaceRatio:       0.08, // 8%
			metadata.Parking:               0.06, // 6%
			metadata.MaintenanceFee:        0.05, // 5%
			metadata.HeatingSystem:         0.03, // 3%
		}

	case ScenarioInvestment:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:            0.05, // 5%
			metadata.DistanceToStation:     0.15, // 15%
			metadata.ElevatorPresence:      0.05, // 5%
			metadata.ConstructionYear:      0.20, // 20%
			metadata.ConstructionCompany:   0.15, // 15%
			metadata.ApartmentSize:         0.08, // 8%
			metadata.NearbyAmenities:       0.10, // 10%
			metadata.TransportationAccess:  0.12, // 12%
			metadata.SchoolDistrict:        0.05, // 5%
			metadata.CrimeRate:             0.02, // 2%
			metadata.GreenSpaceRatio:       0.05, // 5%
			metadata.Parking:               0.05, // 5%
			metadata.MaintenanceFee:        0.05, // 5%
			metadata.HeatingSystem:         0.03, // 3%
		}

	case ScenarioBalanced:
		fallthrough
	default:
		weights = map[metadata.MetadataType]Weight{
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
	}

	return NormalizeWeights(weights)
}

// GetScenarioDescription은 시나리오의 설명을 반환합니다.
func GetScenarioDescription(scenario ScoringScenario) string {
	switch scenario {
	case ScenarioBalanced:
		return "균형 잡힌 선택"
	case ScenarioTransportation:
		return "교통 중심"
	case ScenarioEducation:
		return "교육 우선"
	case ScenarioCostEffective:
		return "가성비 중시"
	case ScenarioFamilyFriendly:
		return "가족 친화적"
	case ScenarioInvestment:
		return "투자 가치"
	default:
		return "알 수 없는 시나리오"
	}
}

// GetAllScenarios는 사용 가능한 모든 시나리오를 반환합니다.
func GetAllScenarios() []ScoringScenario {
	return []ScoringScenario{
		ScenarioBalanced,
		ScenarioTransportation,
		ScenarioEducation,
		ScenarioCostEffective,
		ScenarioFamilyFriendly,
		ScenarioInvestment,
	}
}