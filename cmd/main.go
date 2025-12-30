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
	result, err := scoring.QuickScore(apartmentScores, scoring.ScenarioBalanced)
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
		result, _ := scoring.QuickScore(apartmentScores, scenario)
		fmt.Printf("%-15s: %.1f점\n",
			scoring.GetScenarioDescription(scenario),
			result.TotalScore)
	}
	recommended := scoring.RecommendScenario(apartmentScores)
	fmt.Printf("\n추천 시나리오: %s\n", scoring.GetScenarioDescription(recommended))
}
