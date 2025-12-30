package main

import (
	"fmt"

	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
)

func main() {
	fmt.Println("아파트 스코어링 시스템 시작")

	// 메타데이터 사용 예시
	fmt.Printf("층수 메타데이터: %s (%s)\n", metadata.FloorLevel.String(), metadata.FloorLevel.KoreanName())
	fmt.Printf("역까지 거리 메타데이터: %s (%s)\n", metadata.DistanceToStation.String(), metadata.DistanceToStation.KoreanName())

	// 모든 메타데이터 출력
	fmt.Println("\n=== 모든 메타데이터 목록 ===")
	for i := metadata.MetadataType(0); i < metadata.MetadataTypeCount; i++ {
		fmt.Printf("%d: %s (%s)\n", i.Index(), i.String(), i.KoreanName())
	}

	// 스코어링 예제
	fmt.Println("\n=== 아파트 스코어링 예제 ===")

	// 용인시 A아파트 점수 데이터 (예시)
	apartmentScores := map[metadata.MetadataType]scoring.ScoreValue{
		metadata.FloorLevel:            85.0,  // 5층 (중간층)
		metadata.DistanceToStation:     95.0,  // 역까지 5분
		metadata.ElevatorPresence:      100.0, // 엘리베이터 있음
		metadata.ConstructionYear:      90.0,  // 2020년 건축
		metadata.ConstructionCompany:   85.0,  // 신뢰할 수 있는 회사
		metadata.ApartmentSize:         75.0,  // 적절한 크기
		metadata.NearbyAmenities:       80.0,  // 편의시설 보통
		metadata.TransportationAccess:  90.0,  // 대중교통 접근성 좋음
		metadata.SchoolDistrict:        70.0,  // 학군 보통
		metadata.CrimeRate:             65.0,  // 범죄율 약간 높음
		metadata.GreenSpaceRatio:       60.0,  // 녹지율 낮음
		metadata.Parking:               80.0,  // 주차장 충분
		metadata.MaintenanceFee:        75.0,  // 관리비 적절
		metadata.HeatingSystem:         70.0,  // 난방 시스템 보통
	}

	// 균형 잡힌 시나리오로 점수 계산
	result, err := scoring.QuickScore(apartmentScores, scoring.ScenarioBalanced)
	if err != nil {
		fmt.Printf("스코어링 실패: %v\n", err)
		return
	}

	// 결과 출력
	fmt.Println(scoring.FormatScoreResult(result))

	// 분석 결과
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

	// 시나리오 비교
	fmt.Println("\n=== 시나리오 비교 ===")
	scenarios := []scoring.ScoringScenario{
		scoring.ScenarioBalanced,
		scoring.ScenarioTransportation,
		scoring.ScenarioEducation,
		scoring.ScenarioCostEffective,
	}

	for _, scenario := range scenarios {
		result, _ := scoring.QuickScore(apartmentScores, scenario)
		fmt.Printf("%-15s: %.1f점 (%s)\n",
			scoring.GetScenarioDescription(scenario),
			result.TotalScore,
			result.Grade)
	}

	// 추천 시나리오
	recommended := scoring.RecommendScenario(apartmentScores)
	fmt.Printf("\n추천 시나리오: %s\n", scoring.GetScenarioDescription(recommended))
}
