package main

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"apart_score/pkg/shared"
	"fmt"
)

func main() {
	fmt.Println("아파트 스코어링 시스템 시작")
	for i := metadata.MetadataType(0); i < metadata.MetadataTypeCount; i++ {
		fmt.Printf("%d: %s (%s)\n", i.Index(), i.String(), i.KoreanName())
	}
	fmt.Println("\n=== 아파트 스코어링 예제 ===")
	apartmentScores := map[metadata.MetadataType]shared.ScoreValue{
		metadata.FloorLevel:           shared.ScoreValueFromFloat(85.0),
		metadata.DistanceToStation:    shared.ScoreValueFromFloat(95.0),
		metadata.ElevatorPresence:     shared.ScoreValueFromFloat(100.0),
		metadata.ConstructionYear:     shared.ScoreValueFromFloat(90.0),
		metadata.ConstructionCompany:  shared.ScoreValueFromFloat(85.0),
		metadata.ApartmentSize:        shared.ScoreValueFromFloat(75.0),
		metadata.NearbyAmenities:      shared.ScoreValueFromFloat(80.0),
		metadata.TransportationAccess: shared.ScoreValueFromFloat(90.0),
		metadata.SchoolDistrict:       shared.ScoreValueFromFloat(70.0),
		metadata.CrimeRate:            shared.ScoreValueFromFloat(65.0),
		metadata.GreenSpaceRatio:      shared.ScoreValueFromFloat(60.0),
		metadata.Parking:              shared.ScoreValueFromFloat(80.0),
		metadata.MaintenanceFee:       shared.ScoreValueFromFloat(75.0),
		metadata.HeatingSystem:        shared.ScoreValueFromFloat(70.0),
	}
	weights := scoring.GetScenarioWeights(scoring.ScenarioBalanced)
	result, err := scoring.CalculateWithStrategy(apartmentScores, weights, scoring.StrategyWeightedSum)
	if err != nil {
		fmt.Printf("스코어링 실패: %v\n", err)
		return
	}
	fmt.Println("🏠 아파트 스코어 결과")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("총점: %.1f점\n", result.TotalScore)
	fmt.Printf("방법: %s\n", result.Method)
	fmt.Printf("시나리오: %s\n", result.Scenario)
	fmt.Println("\n📊 상세 점수:")
	for _, mt := range shared.FastAllMetadataTypes() {
		idx := int(mt)
		rawScore := result.RawScores[idx]
		weight := result.Weights[idx]
		weighted := result.WeightedScores[idx]
		if rawScore != 0 {
			fmt.Printf("  %-20s: %.1f점 (가중치: %.1f%%) → %.1f점\n",
				mt.KoreanName(), rawScore.ToFloat(), weight.ToFloat()*100, weighted)
		}
	}
	fmt.Println("\n=== 사용자 정의 스코어링 테이블 예제 ===")
	customWeights := map[metadata.MetadataType]shared.Weight{
		metadata.FloorLevel:           shared.WeightFromFloat(0.10),
		metadata.DistanceToStation:    shared.WeightFromFloat(0.30),
		metadata.ElevatorPresence:     shared.WeightFromFloat(0.08),
		metadata.ConstructionYear:     shared.WeightFromFloat(0.05),
		metadata.ConstructionCompany:  shared.WeightFromFloat(0.02),
		metadata.ApartmentSize:        shared.WeightFromFloat(0.02),
		metadata.NearbyAmenities:      shared.WeightFromFloat(0.15),
		metadata.TransportationAccess: shared.WeightFromFloat(0.25),
		metadata.SchoolDistrict:       shared.WeightFromFloat(0.00),
		metadata.CrimeRate:            shared.WeightFromFloat(0.00),
		metadata.GreenSpaceRatio:      shared.WeightFromFloat(0.00),
		metadata.Parking:              shared.WeightFromFloat(0.05),
		metadata.MaintenanceFee:       shared.WeightFromFloat(0.03),
		metadata.HeatingSystem:        shared.WeightFromFloat(0.00),
	}
	customWeights = shared.NormalizeWeights(customWeights)
	customResult, err := scoring.CalculateWithStrategy(apartmentScores, customWeights, scoring.StrategyWeightedSum)
	if err != nil {
		fmt.Printf("사용자 정의 스코어링 실패: %v\n", err)
		return
	}
	fmt.Println("🎯 교통 최우선 스코어링 테이블 결과:")
	fmt.Printf("총점: %.1f점 (기존: %.1f점, 차이: %.1f점)\n",
		customResult.TotalScore, result.TotalScore,
		customResult.TotalScore-result.TotalScore)

	// === 투명성 대시보드 ===
	fmt.Println("\n🔍 투명성 평가 대시보드")
	fmt.Println("═══════════════════════════════════════════════")

	dashboard := scoring.GenerateTransparencyDashboard(result, apartmentScores, weights, scoring.StrategyWeightedSum)
	fmt.Println(scoring.FormatTransparencyDashboard(dashboard))
	fmt.Println("\n=== 메타데이터 팩터 타입 예제 ===")
	fmt.Println("디폴트 팩터 타입 설정:")
	for i := metadata.MetadataType(0); i < metadata.MetadataTypeCount; i++ {
		fmt.Printf("  %s: %s\n", i.KoreanName(), i.FactorType())
	}
	fmt.Println("\n내부 요인 (아파트 자체 속성):")
	internalFactors := metadata.GetMetadataByFactorType(metadata.FactorInternal)
	for _, mt := range internalFactors {
		if mt != 0 {
			fmt.Printf("  - %s\n", mt.KoreanName())
		}
	}
	fmt.Println("\n외부 요인 (주변 환경):")
	externalFactors := metadata.GetMetadataByFactorType(metadata.FactorExternal)
	for _, mt := range externalFactors {
		if mt != 0 {
			fmt.Printf("  - %s\n", mt.KoreanName())
		}
	}
	fmt.Println("\n팩터 타입 변경 예제:")
	fmt.Printf("변경 전 - 층수: %s\n", metadata.FloorLevel.FactorType())
	if err := metadata.SetFactorType(metadata.FloorLevel, metadata.FactorExternal); err != nil {
		fmt.Printf("팩터 타입 변경 실패: %v\n", err)
	} else {
		fmt.Printf("변경 후 - 층수: %s\n", metadata.FloorLevel.FactorType())
		if err := metadata.SetFactorType(metadata.FloorLevel, metadata.FactorInternal); err != nil {
			fmt.Printf("팩터 타입 복원 실패: %v\n", err)
		}
		fmt.Printf("복원 후 - 층수: %s\n", metadata.FloorLevel.FactorType())
	}

	// === 연산 순서 조정 파이프라인 예제 ===
	fmt.Println("\n=== 연산 순서 조정 파이프라인 예제 ===")
	familyPipeline := scoring.CreateFamilyPipeline()
	fmt.Printf("파이프라인: %s\n", familyPipeline.Name)
	fmt.Printf("설명: %s\n", familyPipeline.Description)

	pipelineResult, err := scoring.CalculateWithPipeline(apartmentScores, weights, familyPipeline)
	if err != nil {
		fmt.Printf("파이프라인 계산 실패: %v\n", err)
		return
	}

	fmt.Printf("파이프라인 총점: %.1f점\n", pipelineResult.TotalScore)
	fmt.Println("계산 단계:")
	for i, step := range familyPipeline.Steps {
		fmt.Printf("  %d. %s (%d순위)\n", i+1, step.Name, step.Priority)
		fmt.Printf("     %s\n", step.Description)
	}

	// 기존 방식과 비교
	fmt.Printf("\n비교:\n")
	fmt.Printf("  기존 Weighted Sum: %.1f점\n", result.TotalScore)
	fmt.Printf("  파이프라인 방식: %.1f점\n", pipelineResult.TotalScore)
	fmt.Printf("  차이: %.1f점\n", pipelineResult.TotalScore-result.TotalScore)
}
