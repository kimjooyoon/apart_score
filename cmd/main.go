package main

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"fmt"
)

func main() {
	fmt.Println("아파트 스코어링 시스템 시작")
	fmt.Printf("층수 메타데이터: %s (%s)\n", metadata.FloorLevel.String(), metadata.FloorLevel.KoreanName())
	fmt.Printf("역까지 거리 메타데이터: %s (%s)\n", metadata.DistanceToStation.String(), metadata.DistanceToStation.KoreanName())
	fmt.Println("\n=== 모든 메타데이터 목록 ===")
	for i := metadata.MetadataType(0); i < metadata.MetadataTypeCount; i++ {
		fmt.Printf("%d: %s (%s)\n", i.Index(), i.String(), i.KoreanName())
	}
	fmt.Println("\n=== 아파트 스코어링 예제 ===")
	apartmentScores := map[metadata.MetadataType]scoring.ScoreValue{
		metadata.FloorLevel:           85.0,
		metadata.DistanceToStation:    95.0,
		metadata.ElevatorPresence:     100.0,
		metadata.ConstructionYear:     90.0,
		metadata.ConstructionCompany:  85.0,
		metadata.ApartmentSize:        75.0,
		metadata.NearbyAmenities:      80.0,
		metadata.TransportationAccess: 90.0,
		metadata.SchoolDistrict:       70.0,
		metadata.CrimeRate:            65.0,
		metadata.GreenSpaceRatio:      60.0,
		metadata.Parking:              80.0,
		metadata.MaintenanceFee:       75.0,
		metadata.HeatingSystem:        70.0,
	}
	weights := scoring.GetScenarioWeights(scoring.ScenarioBalanced)
	result, err := scoring.CalculateWithStrategy(apartmentScores, weights, scoring.StrategyWeightedSum)
	if err != nil {
		fmt.Printf("스코어링 실패: %v\n", err)
		return
	}
	fmt.Println(scoring.FormatScoreResult(result))
	analysis := scoring.AnalyzeScore(result)
	fmt.Println("\n=== 상세 분석 ===")
	fmt.Printf("강점 (%d개):\n", len(analysis.Strengths))
	for _, mt := range analysis.Strengths {
		fmt.Printf("  - %s\n", mt.KoreanName())
	}
	fmt.Printf("\n약점 (%d개):\n", len(analysis.Weaknesses))
	for _, mt := range analysis.Weaknesses {
		fmt.Printf("  - %s\n", mt.KoreanName())
	}
	fmt.Printf("\n개선 제안:\n")
	for _, tip := range analysis.ImprovementTips {
		fmt.Printf("  - %s\n", tip)
	}
	fmt.Println("\n=== 시나리오 비교 ===")
	scenarios := []scoring.ScoringScenario{
		scoring.ScenarioBalanced,
		scoring.ScenarioTransportation,
		scoring.ScenarioEducation,
		scoring.ScenarioCostEffective,
	}
	for _, scenario := range scenarios {
		weights := scoring.GetScenarioWeights(scenario)
		result, _ := scoring.CalculateWithStrategy(apartmentScores, weights, scoring.StrategyWeightedSum)
		fmt.Printf("%-15s: %.1f점\n",
			scoring.GetScenarioDescription(scenario),
			result.TotalScore)
	}
	recommended := scoring.RecommendScenario(apartmentScores)
	fmt.Printf("\n추천 시나리오: %s\n", scoring.GetScenarioDescription(recommended))
	fmt.Println("\n=== 여러 아파트 순위 비교 ===")
	apartments := []scoring.ApartmentData{
		{
			ID:   "apt001",
			Name: "강남 래미안",
			Scores: map[metadata.MetadataType]scoring.ScoreValue{
				metadata.FloorLevel:           85.0,
				metadata.DistanceToStation:    95.0,
				metadata.ElevatorPresence:     100.0,
				metadata.ConstructionYear:     90.0,
				metadata.ConstructionCompany:  88.0,
				metadata.ApartmentSize:        75.0,
				metadata.NearbyAmenities:      85.0,
				metadata.TransportationAccess: 90.0,
				metadata.SchoolDistrict:       80.0,
				metadata.CrimeRate:            70.0,
				metadata.GreenSpaceRatio:      65.0,
				metadata.Parking:              85.0,
				metadata.MaintenanceFee:       80.0,
				metadata.HeatingSystem:        75.0,
			},
			Location: "서울 강남구",
		},
		{
			ID:   "apt002",
			Name: "서초 아크로텔",
			Scores: map[metadata.MetadataType]scoring.ScoreValue{
				metadata.FloorLevel:           80.0,
				metadata.DistanceToStation:    85.0,
				metadata.ElevatorPresence:     100.0,
				metadata.ConstructionYear:     85.0,
				metadata.ConstructionCompany:  82.0,
				metadata.ApartmentSize:        70.0,
				metadata.NearbyAmenities:      80.0,
				metadata.TransportationAccess: 88.0,
				metadata.SchoolDistrict:       75.0,
				metadata.CrimeRate:            75.0,
				metadata.GreenSpaceRatio:      70.0,
				metadata.Parking:              80.0,
				metadata.MaintenanceFee:       75.0,
				metadata.HeatingSystem:        70.0,
			},
			Location: "서울 서초구",
		},
		{
			ID:   "apt003",
			Name: "송파 헬리오시티",
			Scores: map[metadata.MetadataType]scoring.ScoreValue{
				metadata.FloorLevel:           75.0,
				metadata.DistanceToStation:    80.0,
				metadata.ElevatorPresence:     95.0,
				metadata.ConstructionYear:     80.0,
				metadata.ConstructionCompany:  78.0,
				metadata.ApartmentSize:        65.0,
				metadata.NearbyAmenities:      75.0,
				metadata.TransportationAccess: 82.0,
				metadata.SchoolDistrict:       70.0,
				metadata.CrimeRate:            80.0,
				metadata.GreenSpaceRatio:      75.0,
				metadata.Parking:              75.0,
				metadata.MaintenanceFee:       70.0,
				metadata.HeatingSystem:        65.0,
			},
			Location: "서울 송파구",
		},
		{
			ID:   "apt004",
			Name: "마포 래미안",
			Scores: map[metadata.MetadataType]scoring.ScoreValue{
				metadata.FloorLevel:           70.0,
				metadata.DistanceToStation:    75.0,
				metadata.ElevatorPresence:     90.0,
				metadata.ConstructionYear:     75.0,
				metadata.ConstructionCompany:  72.0,
				metadata.ApartmentSize:        60.0,
				metadata.NearbyAmenities:      70.0,
				metadata.TransportationAccess: 78.0,
				metadata.SchoolDistrict:       65.0,
				metadata.CrimeRate:            85.0,
				metadata.GreenSpaceRatio:      80.0,
				metadata.Parking:              70.0,
				metadata.MaintenanceFee:       65.0,
				metadata.HeatingSystem:        60.0,
			},
			Location: "서울 마포구",
		},
	}
	weights = make(map[metadata.MetadataType]scoring.Weight)
	totalTypes := len(apartmentScores)
	equalWeight := scoring.Weight(1.0 / float64(totalTypes))
	for mt := range apartmentScores {
		weights[mt] = equalWeight
	}
	rankings, err := scoring.CalculateRankings(apartments, weights, scoring.StrategyWeightedSum)
	if err != nil {
		fmt.Printf("순위 계산 실패: %v\n", err)
		return
	}
	fmt.Println(scoring.FormatRankings(rankings, 3))
	fmt.Println("\n=== 메타데이터 팩터 타입 예제 ===")
	fmt.Println("디폴트 팩터 타입 설정:")
	for mt := metadata.MetadataType(0); mt < metadata.MetadataTypeCount; mt++ {
		fmt.Printf("  %s: %s\n", mt.KoreanName(), mt.FactorType())
	}
	fmt.Println("\n내부 요인 (아파트 자체 속성):")
	internalFactors := metadata.GetMetadataByFactorType(metadata.FactorInternal)
	for _, mt := range internalFactors {
		fmt.Printf("  - %s\n", mt.KoreanName())
	}
	fmt.Println("\n외부 요인 (주변 환경):")
	externalFactors := metadata.GetMetadataByFactorType(metadata.FactorExternal)
	for _, mt := range externalFactors {
		fmt.Printf("  - %s\n", mt.KoreanName())
	}
	fmt.Println("\n팩터 타입 변경 예제:")
	fmt.Printf("변경 전 - 층수: %s\n", metadata.FloorLevel.FactorType())
	err = metadata.SetFactorType(metadata.FloorLevel, metadata.FactorExternal)
	if err != nil {
		fmt.Printf("팩터 타입 변경 실패: %v\n", err)
	} else {
		fmt.Printf("변경 후 - 층수: %s\n", metadata.FloorLevel.FactorType())
		_ = metadata.SetFactorType(metadata.FloorLevel, metadata.FactorInternal)
		fmt.Printf("복원 후 - 층수: %s\n", metadata.FloorLevel.FactorType())
	}
}
