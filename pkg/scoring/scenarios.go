package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
)

func GetScenarioWeights(scenario ScoringScenario) map[metadata.MetadataType]shared.Weight {
	var weights map[metadata.MetadataType]shared.Weight
	switch scenario {
	case ScenarioTransportation:
		weights = map[metadata.MetadataType]shared.Weight{
			metadata.FloorLevel:           shared.WeightFromFloat(0.05),
			metadata.DistanceToStation:    shared.WeightFromFloat(0.25),
			metadata.ElevatorPresence:     shared.WeightFromFloat(0.05),
			metadata.ConstructionYear:     shared.WeightFromFloat(0.08),
			metadata.ConstructionCompany:  shared.WeightFromFloat(0.05),
			metadata.ApartmentSize:        shared.WeightFromFloat(0.07),
			metadata.NearbyAmenities:      shared.WeightFromFloat(0.10),
			metadata.TransportationAccess: shared.WeightFromFloat(0.20),
			metadata.SchoolDistrict:       shared.WeightFromFloat(0.05),
			metadata.CrimeRate:            shared.WeightFromFloat(0.05),
			metadata.GreenSpaceRatio:      shared.WeightFromFloat(0.02),
			metadata.Parking:              shared.WeightFromFloat(0.05),
			metadata.MaintenanceFee:       shared.WeightFromFloat(0.05),
			metadata.HeatingSystem:        shared.WeightFromFloat(0.03),
		}
	case ScenarioEducation:
		weights = map[metadata.MetadataType]shared.Weight{
			metadata.FloorLevel:           shared.WeightFromFloat(0.08),
			metadata.DistanceToStation:    shared.WeightFromFloat(0.08),
			metadata.ElevatorPresence:     shared.WeightFromFloat(0.07),
			metadata.ConstructionYear:     shared.WeightFromFloat(0.10),
			metadata.ConstructionCompany:  shared.WeightFromFloat(0.08),
			metadata.ApartmentSize:        shared.WeightFromFloat(0.08),
			metadata.NearbyAmenities:      shared.WeightFromFloat(0.08),
			metadata.TransportationAccess: shared.WeightFromFloat(0.08),
			metadata.SchoolDistrict:       shared.WeightFromFloat(0.20),
			metadata.CrimeRate:            shared.WeightFromFloat(0.08),
			metadata.GreenSpaceRatio:      shared.WeightFromFloat(0.05),
			metadata.Parking:              shared.WeightFromFloat(0.05),
			metadata.MaintenanceFee:       shared.WeightFromFloat(0.04),
			metadata.HeatingSystem:        shared.WeightFromFloat(0.03),
		}
	case ScenarioCostEffective:
		weights = map[metadata.MetadataType]shared.Weight{
			metadata.FloorLevel:           shared.WeightFromFloat(0.08),
			metadata.DistanceToStation:    shared.WeightFromFloat(0.12),
			metadata.ElevatorPresence:     shared.WeightFromFloat(0.07),
			metadata.ConstructionYear:     shared.WeightFromFloat(0.08),
			metadata.ConstructionCompany:  shared.WeightFromFloat(0.05),
			metadata.ApartmentSize:        shared.WeightFromFloat(0.12),
			metadata.NearbyAmenities:      shared.WeightFromFloat(0.10),
			metadata.TransportationAccess: shared.WeightFromFloat(0.08),
			metadata.SchoolDistrict:       shared.WeightFromFloat(0.08),
			metadata.CrimeRate:            shared.WeightFromFloat(0.08),
			metadata.GreenSpaceRatio:      shared.WeightFromFloat(0.05),
			metadata.Parking:              shared.WeightFromFloat(0.06),
			metadata.MaintenanceFee:       shared.WeightFromFloat(0.10),
			metadata.HeatingSystem:        shared.WeightFromFloat(0.03),
		}
	case ScenarioFamilyFriendly:
		weights = map[metadata.MetadataType]shared.Weight{
			metadata.FloorLevel:           shared.WeightFromFloat(0.10),
			metadata.DistanceToStation:    shared.WeightFromFloat(0.08),
			metadata.ElevatorPresence:     shared.WeightFromFloat(0.12),
			metadata.ConstructionYear:     shared.WeightFromFloat(0.10),
			metadata.ConstructionCompany:  shared.WeightFromFloat(0.08),
			metadata.ApartmentSize:        shared.WeightFromFloat(0.12),
			metadata.NearbyAmenities:      shared.WeightFromFloat(0.08),
			metadata.TransportationAccess: shared.WeightFromFloat(0.05),
			metadata.SchoolDistrict:       shared.WeightFromFloat(0.12),
			metadata.CrimeRate:            shared.WeightFromFloat(0.08),
			metadata.GreenSpaceRatio:      shared.WeightFromFloat(0.08),
			metadata.Parking:              shared.WeightFromFloat(0.06),
			metadata.MaintenanceFee:       shared.WeightFromFloat(0.05),
			metadata.HeatingSystem:        shared.WeightFromFloat(0.03),
		}
	case ScenarioInvestment:
		weights = map[metadata.MetadataType]shared.Weight{
			metadata.FloorLevel:           shared.WeightFromFloat(0.05),
			metadata.DistanceToStation:    shared.WeightFromFloat(0.15),
			metadata.ElevatorPresence:     shared.WeightFromFloat(0.05),
			metadata.ConstructionYear:     shared.WeightFromFloat(0.20),
			metadata.ConstructionCompany:  shared.WeightFromFloat(0.15),
			metadata.ApartmentSize:        shared.WeightFromFloat(0.08),
			metadata.NearbyAmenities:      shared.WeightFromFloat(0.10),
			metadata.TransportationAccess: shared.WeightFromFloat(0.12),
			metadata.SchoolDistrict:       shared.WeightFromFloat(0.05),
			metadata.CrimeRate:            shared.WeightFromFloat(0.02),
			metadata.GreenSpaceRatio:      shared.WeightFromFloat(0.05),
			metadata.Parking:              shared.WeightFromFloat(0.05),
			metadata.MaintenanceFee:       shared.WeightFromFloat(0.05),
			metadata.HeatingSystem:        shared.WeightFromFloat(0.03),
		}
	case ScenarioBalanced:
		fallthrough
	default:
		weights = map[metadata.MetadataType]shared.Weight{
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
	}
	return shared.NormalizeWeights(weights)
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
