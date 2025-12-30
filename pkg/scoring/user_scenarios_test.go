package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring/strategy"
	"apart_score/pkg/shared"
	"testing"
)

// TestUserScenarios_YoungProfessional은 젊은 직장인의 스코어링 시나리오를 테스트합니다.
func TestUserScenarios_YoungProfessional(t *testing.T) {
	// 젊은 직장인의 선호도: 교통 + 편의시설 + 최신 건물
	youngProfessional := map[metadata.MetadataType]shared.ScoreValue{
		metadata.DistanceToStation:    95.0, // 출퇴근 최우선
		metadata.TransportationAccess: 90.0, // 대중교통 편의
		metadata.NearbyAmenities:      85.0, // 생활 편의
		metadata.ConstructionYear:     88.0, // 최신 건물 선호
		metadata.ElevatorPresence:     100.0, // 필수
		metadata.FloorLevel:           80.0, // 적절한 층수
		metadata.ApartmentSize:        70.0, // 적당한 크기
		metadata.CrimeRate:            75.0, // 안전 중요
		metadata.Parking:              60.0, // 주차는 덜 중요
		metadata.MaintenanceFee:       65.0, // 관리비는 고려
		metadata.ConstructionCompany:  70.0, // 어느 정도 신뢰
		metadata.SchoolDistrict:       50.0, // 교육은 덜 중요
		metadata.GreenSpaceRatio:      55.0, // 환경은 보통
		metadata.HeatingSystem:        75.0, // 난방은 중요
	}

	// 가중치 설정: 교통 40%, 편의시설 30%, 건물 상태 20%, 기타 10%
	weights := map[metadata.MetadataType]shared.Weight{
		metadata.DistanceToStation:    0.20,
		metadata.TransportationAccess: 0.20, // 교통 총합 40%
		metadata.NearbyAmenities:      0.30, // 편의시설 30%
		metadata.ConstructionYear:     0.15,
		metadata.ElevatorPresence:     0.03,
		metadata.HeatingSystem:        0.02, // 건물 상태 총합 20%
		metadata.FloorLevel:           0.03,
		metadata.ApartmentSize:        0.02,
		metadata.CrimeRate:            0.02,
		metadata.Parking:              0.01,
		metadata.MaintenanceFee:       0.01,
		metadata.ConstructionCompany:  0.01,
		metadata.SchoolDistrict:       0.005,
		metadata.GreenSpaceRatio:      0.005, // 기타 10%
	}

	result, err := CalculateWithStrategy(youngProfessional, weights, strategy.StrategyWeightedSum)
	if err != nil {
		t.Fatalf("젊은 직장인 스코어링 실패: %v", err)
	}

	// 교통 관련 점수가 높게 반영되어야 함
	if result.TotalScore < 75.0 {
		t.Errorf("젊은 직장인의 선호도가 충분히 반영되지 않은 낮은 점수: %.1f", result.TotalScore)
	}

	t.Logf("젊은 직장인 스코어링 결과: %.1f점", result.TotalScore)
}

// TestUserScenarios_FamilyOriented은 가족 단위 거주자의 스코어링을 테스트합니다.
func TestUserScenarios_FamilyOriented(t *testing.T) {
	// 가족 단위: 교육 + 안전 + 공간 + 쾌적함
	familyOriented := map[metadata.MetadataType]shared.ScoreValue{
		metadata.SchoolDistrict:       95.0, // 교육 최우선
		metadata.CrimeRate:            90.0, // 안전 중요
		metadata.GreenSpaceRatio:      85.0, // 쾌적함
		metadata.ApartmentSize:        88.0, // 충분한 공간
		metadata.ElevatorPresence:     95.0, // 가족 편의
		metadata.FloorLevel:           75.0, // 적절한 층수
		metadata.DistanceToStation:    70.0, // 교통은 보통
		metadata.TransportationAccess: 65.0,
		metadata.NearbyAmenities:      80.0, // 생활 편의
		metadata.ConstructionYear:     60.0, // 오래된 건물도 괜찮음
		metadata.Parking:              85.0, // 가족 차량 고려
		metadata.MaintenanceFee:       75.0, // 관리비는 적절히
		metadata.ConstructionCompany:  70.0,
		metadata.HeatingSystem:        80.0,
	}

	// 가중치: 교육 30%, 안전/쾌적함 25%, 공간 20%, 편의 15%, 기타 10%
	weights := map[metadata.MetadataType]shared.Weight{
		metadata.SchoolDistrict:       0.30, // 교육 30%
		metadata.CrimeRate:            0.15,
		metadata.GreenSpaceRatio:      0.10, // 안전/쾌적함 25%
		metadata.ApartmentSize:        0.15,
		metadata.ElevatorPresence:     0.03,
		metadata.Parking:              0.02, // 공간/편의 20%
		metadata.FloorLevel:           0.02,
		metadata.NearbyAmenities:      0.10, // 편의시설 10%
		metadata.DistanceToStation:    0.03,
		metadata.TransportationAccess: 0.02,
		metadata.ConstructionYear:     0.02,
		metadata.MaintenanceFee:       0.02,
		metadata.ConstructionCompany:  0.01,
		metadata.HeatingSystem:        0.03, // 기타 10%
	}

	result, err := CalculateWithStrategy(familyOriented, weights, strategy.StrategyWeightedSum)
	if err != nil {
		t.Fatalf("가족 단위 스코어링 실패: %v", err)
	}

	// 교육 점수가 높게 반영되어야 함
	if result.TotalScore < 80.0 {
		t.Errorf("가족 단위 선호도가 충분히 반영되지 않은 점수: %.1f", result.TotalScore)
	}

	t.Logf("가족 단위 스코어링 결과: %.1f점", result.TotalScore)
}