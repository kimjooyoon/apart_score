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

// ScenarioDefinitions provides clear definitions and characteristics for each scenario.
var ScenarioDefinitions = map[ScoringScenario]ScenarioDefinition{
	ScenarioBalanced: {
		Name:        "균형 잡힌 선택",
		TargetUser:  "일반적인 아파트 구매자",
		Description: "모든 요소를 골고루 고려하는 균형 잡힌 평가",
		KeyWeights: map[string]float64{
			"층수": 7.3, "역까지 거리": 13.6, "엘리베이터": 6.4, "건축년도": 9.1,
			"건설회사": 7.3, "크기": 7.3, "편의시설": 9.1, "교통": 10.9,
			"학군": 7.3, "범죄율": 5.5, "녹지율": 3.6, "주차장": 5.5,
			"관리비": 4.5, "난방": 2.7,
		},
		UseCase:      "처음 아파트 구매, 일반적인 투자, 임대용 부동산",
		Strengths:    []string{"안정적", "편향 없음", "대부분 상황 적합"},
		Limitations:  []string{"개인 우선순위 반영 부족", "특별한 요구사항 무시"},
		Alternatives: []ScoringScenario{ScenarioFamilyFriendly, ScenarioCostEffective},
	},
	ScenarioTransportation: {
		Name:        "교통 중심",
		TargetUser:  "직장인, 자가용 없는 세대",
		Description: "교통 접근성을 최우선으로 하는 평가",
		KeyWeights: map[string]float64{
			"역까지 거리": 25.0, "교통 접근성": 20.0, "편의시설": 15.0, "층수": 10.0,
			"엘리베이터": 8.0, "학군": 5.0, "크기": 7.0, "건축년도": 8.0,
			"건설회사": 5.0, "범죄율": 5.0, "녹지율": 2.0, "주차장": 5.0,
			"관리비": 5.0, "난방": 3.0,
		},
		UseCase:      "직장 근처 거주, 대중교통 의존, 통근 시간 최소화",
		Strengths:    []string{"교통 편의성 극대화", "시간 절약", "이동 비용 감소"},
		Limitations:  []string{"주거 환경 무시 가능", "교통 편의성 외 요인 저평가"},
		Alternatives: []ScoringScenario{ScenarioBalanced, ScenarioEducation},
	},
	ScenarioEducation: {
		Name:        "교육 우선",
		TargetUser:  "자녀 교육 우선 가정",
		Description: "학군과 교육 환경을 최우선으로 하는 평가",
		KeyWeights: map[string]float64{
			"학군": 30.0, "범죄율": 15.0, "녹지율": 10.0, "편의시설": 10.0,
			"층수": 8.0, "크기": 8.0, "엘리베이터": 7.0, "건축년도": 8.0,
			"건설회사": 8.0, "역까지 거리": 8.0, "교통": 5.0, "주차장": 5.0,
			"관리비": 4.0, "난방": 2.0,
		},
		UseCase:      "자녀 교육 우선, 안전한 환경 선호, 장기 거주 계획",
		Strengths:    []string{"교육 환경 최적화", "안전성 확보", "주거 안정성"},
		Limitations:  []string{"경제성 무시 가능", "교통 불편 감수"},
		Alternatives: []ScoringScenario{ScenarioFamilyFriendly, ScenarioBalanced},
	},
	ScenarioCostEffective: {
		Name:        "가성비 중시",
		TargetUser:  "예산 제한이 있는 구매자",
		Description: "가격 대비 성능을 중시하는 평가",
		KeyWeights: map[string]float64{
			"관리비": 15.0, "크기": 15.0, "건축년도": 12.0, "건설회사": 10.0,
			"층수": 8.0, "엘리베이터": 7.0, "역까지 거리": 8.0, "교통": 7.0,
			"편의시설": 7.0, "학군": 5.0, "범죄율": 5.0, "녹지율": 3.0,
			"주차장": 5.0, "난방": 3.0,
		},
		UseCase:      "예산 제한, 유지비 최소화, 실용성 우선",
		Strengths:    []string{"비용 효율성", "경제적 타당성", "실용적 선택"},
		Limitations:  []string{"품질 저하 가능", "위치 제한", "편의성 희생"},
		Alternatives: []ScoringScenario{ScenarioBalanced, ScenarioInvestment},
	},
	ScenarioFamilyFriendly: {
		Name:        "가족 친화적",
		TargetUser:  "자녀 있는 가정",
		Description: "가족 구성원을 고려한 종합적 평가",
		KeyWeights: map[string]float64{
			"크기": 15.0, "학군": 12.0, "엘리베이터": 12.0, "건축년도": 10.0,
			"층수": 10.0, "건설회사": 8.0, "범죄율": 8.0, "녹지율": 8.0,
			"편의시설": 8.0, "역까지 거리": 8.0, "교통": 5.0, "주차장": 6.0,
			"관리비": 5.0, "난방": 3.0,
		},
		UseCase:      "자녀 양육, 가족 편의성, 안전과 공간 중시",
		Strengths:    []string{"가족 편의성", "공간 확보", "안전성 균형"},
		Limitations:  []string{"가격 상승", "교통 불편 가능"},
		Alternatives: []ScoringScenario{ScenarioEducation, ScenarioBalanced},
	},
	ScenarioInvestment: {
		Name:        "투자 가치",
		TargetUser:  "부동산 투자자",
		Description: "장기적 가치 상승 잠재력을 고려한 평가",
		KeyWeights: map[string]float64{
			"건축년도": 20.0, "건설회사": 15.0, "역까지 거리": 12.0, "교통": 10.0,
			"편의시설": 10.0, "학군": 8.0, "범죄율": 7.0, "크기": 8.0,
			"층수": 5.0, "엘리베이터": 5.0, "녹지율": 5.0, "주차장": 5.0,
			"관리비": 5.0, "난방": 2.0,
		},
		UseCase:      "임대 수익, 가치 상승, 장기 투자",
		Strengths:    []string{"투자 수익성", "가치 상승 잠재력", "시장성"},
		Limitations:  []string{"주거 편의성 무시", "주관적 가치 판단"},
		Alternatives: []ScoringScenario{ScenarioBalanced, ScenarioCostEffective},
	},
}

// ScenarioDefinition provides comprehensive information about a scoring scenario.
type ScenarioDefinition struct {
	Name         string             // 시나리오 이름
	TargetUser   string             // 대상 사용자
	Description  string             // 상세 설명
	KeyWeights   map[string]float64 // 주요 가중치 (백분율)
	UseCase      string             // 사용 사례
	Strengths    []string           // 장점들
	Limitations  []string           // 제한사항
	Alternatives []ScoringScenario  // 대안 시나리오들
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
