package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
	"testing"
)

// BenchmarkCalculateWeightedSum는 가중치 합계 계산 성능을 측정합니다.
func BenchmarkCalculateWeightedSum(b *testing.B) {
	// 테스트 데이터 준비
	scores := map[metadata.MetadataType]shared.ScoreValue{
		metadata.FloorLevel:           shared.ScoreValueFromFloat(85.0),
		metadata.DistanceToStation:    shared.ScoreValueFromFloat(90.0),
		metadata.ElevatorPresence:     shared.ScoreValueFromFloat(100.0),
		metadata.ConstructionYear:     shared.ScoreValueFromFloat(88.0),
		metadata.ConstructionCompany:  shared.ScoreValueFromFloat(82.0),
		metadata.ApartmentSize:        shared.ScoreValueFromFloat(78.0),
		metadata.NearbyAmenities:      shared.ScoreValueFromFloat(85.0),
		metadata.TransportationAccess: shared.ScoreValueFromFloat(92.0),
		metadata.SchoolDistrict:       shared.ScoreValueFromFloat(75.0),
		metadata.CrimeRate:            shared.ScoreValueFromFloat(80.0),
		metadata.GreenSpaceRatio:      shared.ScoreValueFromFloat(70.0),
		metadata.Parking:              shared.ScoreValueFromFloat(85.0),
		metadata.MaintenanceFee:       shared.ScoreValueFromFloat(75.0),
		metadata.HeatingSystem:        shared.ScoreValueFromFloat(78.0),
	}

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
	weights = shared.NormalizeWeights(weights)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CalculateWithStrategy(scores, weights, StrategyWeightedSum)
	}
}

// BenchmarkCalculateGeometricMean은 기하 평균 계산 성능을 측정합니다.
func BenchmarkCalculateGeometricMean(b *testing.B) {
	scores := map[metadata.MetadataType]shared.ScoreValue{
		metadata.FloorLevel:           shared.ScoreValueFromFloat(85.0),
		metadata.DistanceToStation:    shared.ScoreValueFromFloat(90.0),
		metadata.ElevatorPresence:     shared.ScoreValueFromFloat(100.0),
		metadata.ConstructionYear:     shared.ScoreValueFromFloat(88.0),
		metadata.ConstructionCompany:  shared.ScoreValueFromFloat(82.0),
		metadata.ApartmentSize:        shared.ScoreValueFromFloat(78.0),
		metadata.NearbyAmenities:      shared.ScoreValueFromFloat(85.0),
		metadata.TransportationAccess: shared.ScoreValueFromFloat(92.0),
		metadata.SchoolDistrict:       shared.ScoreValueFromFloat(75.0),
		metadata.CrimeRate:            shared.ScoreValueFromFloat(80.0),
		metadata.GreenSpaceRatio:      shared.ScoreValueFromFloat(70.0),
		metadata.Parking:              shared.ScoreValueFromFloat(85.0),
		metadata.MaintenanceFee:       shared.ScoreValueFromFloat(75.0),
		metadata.HeatingSystem:        shared.ScoreValueFromFloat(78.0),
	}

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
	weights = shared.NormalizeWeights(weights)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CalculateWithStrategy(scores, weights, StrategyGeometricMean)
	}
}

// BenchmarkMetadataOperations은 메타데이터 연산 성능을 측정합니다.
func BenchmarkMetadataOperations(b *testing.B) {
	b.Run("GetMetadataByFactorType", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metadata.GetMetadataByFactorType(metadata.FactorInternal)
			metadata.GetMetadataByFactorType(metadata.FactorExternal)
		}
	})

	b.Run("AllMetadataTypes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metadata.AllMetadataTypes()
		}
	})

	b.Run("FactorTypeAccess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = metadata.FloorLevel.FactorType()
			_ = metadata.DistanceToStation.FactorType()
			_ = metadata.SchoolDistrict.FactorType()
		}
	})
}

// BenchmarkWeightNormalization은 가중치 정규화 성능을 측정합니다.
func BenchmarkWeightNormalization(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shared.NormalizeWeights(weights)
	}
}
