package scoring

import "apart_score/pkg/metadata"

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
func GetScenarioWeights(scenario ScoringScenario) map[metadata.MetadataType]Weight {
	var weights map[metadata.MetadataType]Weight
	switch scenario {
	case ScenarioTransportation:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:           0.05,
			metadata.DistanceToStation:    0.25,
			metadata.ElevatorPresence:     0.05,
			metadata.ConstructionYear:     0.08,
			metadata.ConstructionCompany:  0.05,
			metadata.ApartmentSize:        0.07,
			metadata.NearbyAmenities:      0.10,
			metadata.TransportationAccess: 0.20,
			metadata.SchoolDistrict:       0.05,
			metadata.CrimeRate:            0.05,
			metadata.GreenSpaceRatio:      0.02,
			metadata.Parking:              0.05,
			metadata.MaintenanceFee:       0.05,
			metadata.HeatingSystem:        0.03,
		}
	case ScenarioEducation:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:           0.08,
			metadata.DistanceToStation:    0.08,
			metadata.ElevatorPresence:     0.07,
			metadata.ConstructionYear:     0.10,
			metadata.ConstructionCompany:  0.08,
			metadata.ApartmentSize:        0.08,
			metadata.NearbyAmenities:      0.08,
			metadata.TransportationAccess: 0.08,
			metadata.SchoolDistrict:       0.20,
			metadata.CrimeRate:            0.08,
			metadata.GreenSpaceRatio:      0.05,
			metadata.Parking:              0.05,
			metadata.MaintenanceFee:       0.04,
			metadata.HeatingSystem:        0.03,
		}
	case ScenarioCostEffective:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:           0.08,
			metadata.DistanceToStation:    0.12,
			metadata.ElevatorPresence:     0.07,
			metadata.ConstructionYear:     0.08,
			metadata.ConstructionCompany:  0.05,
			metadata.ApartmentSize:        0.12,
			metadata.NearbyAmenities:      0.10,
			metadata.TransportationAccess: 0.08,
			metadata.SchoolDistrict:       0.08,
			metadata.CrimeRate:            0.08,
			metadata.GreenSpaceRatio:      0.05,
			metadata.Parking:              0.06,
			metadata.MaintenanceFee:       0.10,
			metadata.HeatingSystem:        0.03,
		}
	case ScenarioFamilyFriendly:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:           0.10,
			metadata.DistanceToStation:    0.08,
			metadata.ElevatorPresence:     0.12,
			metadata.ConstructionYear:     0.10,
			metadata.ConstructionCompany:  0.08,
			metadata.ApartmentSize:        0.12,
			metadata.NearbyAmenities:      0.08,
			metadata.TransportationAccess: 0.05,
			metadata.SchoolDistrict:       0.12,
			metadata.CrimeRate:            0.08,
			metadata.GreenSpaceRatio:      0.08,
			metadata.Parking:              0.06,
			metadata.MaintenanceFee:       0.05,
			metadata.HeatingSystem:        0.03,
		}
	case ScenarioInvestment:
		weights = map[metadata.MetadataType]Weight{
			metadata.FloorLevel:           0.05,
			metadata.DistanceToStation:    0.15,
			metadata.ElevatorPresence:     0.05,
			metadata.ConstructionYear:     0.20,
			metadata.ConstructionCompany:  0.15,
			metadata.ApartmentSize:        0.08,
			metadata.NearbyAmenities:      0.10,
			metadata.TransportationAccess: 0.12,
			metadata.SchoolDistrict:       0.05,
			metadata.CrimeRate:            0.02,
			metadata.GreenSpaceRatio:      0.05,
			metadata.Parking:              0.05,
			metadata.MaintenanceFee:       0.05,
			metadata.HeatingSystem:        0.03,
		}
	case ScenarioBalanced:
		fallthrough
	default:
		weights = map[metadata.MetadataType]Weight{
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
	}
	return NormalizeWeights(weights)
}
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
