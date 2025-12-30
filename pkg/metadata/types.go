package metadata

// MetadataType represents different types of apartment metadata for scoring.
type MetadataType int

// FactorType categorizes metadata with clear definitions and examples.
type FactorType string

const (
	// FactorInternal represents building-specific attributes that can be modified or chosen.
	FactorInternal FactorType = "internal"
	// FactorExternal represents environmental attributes that are fixed or external to the building.
	FactorExternal FactorType = "external"
)

// FactorTypeDefinitions provides clear definitions and examples for each factor type.
var FactorTypeDefinitions = map[FactorType]FactorDefinition{
	FactorInternal: {
		Name:        "내부 요인 (건물 속성)",
		Description: "건물 자체의 물리적/구조적 속성으로, 선택이나 개선이 가능한 요소들",
		Examples:    []string{"층수", "엘리베이터", "건축년도", "건설회사", "크기", "주차장", "난방 방식"},
		Characteristics: []string{
			"건물 선택 시 직접 고려 가능",
			"리모델링이나 개선으로 변경 가능",
			"구매/임대 시 협상 대상",
			"건물 관리자가 통제 가능",
		},
	},
	FactorExternal: {
		Name:        "외부 요인 (환경 속성)",
		Description: "건물 주변의 환경적/지역적 속성으로, 건물 외부 요인들",
		Examples:    []string{"역까지 거리", "주변 편의시설", "교통 접근성", "학군", "범죄율", "녹지율"},
		Characteristics: []string{
			"건물 위치에 따라 고정됨",
			"개별 건물로 변경 불가능",
			"지역 개발로 장기적 변화 가능",
			"주변 환경에 의존",
		},
	},
}

// FactorDefinition provides detailed information about a factor type.
type FactorDefinition struct {
	Name            string   // 표시 이름
	Description     string   // 상세 설명
	Examples        []string // 구체적 예시
	Characteristics []string // 주요 특징
}

// Metadata types for apartment scoring
const (
	// FloorLevel represents the floor level of the apartment
	FloorLevel MetadataType = iota
	DistanceToStation
	ElevatorPresence
	ConstructionYear
	ConstructionCompany
	ApartmentSize
	NearbyAmenities
	TransportationAccess
	SchoolDistrict
	CrimeRate
	GreenSpaceRatio
	Parking
	MaintenanceFee
	HeatingSystem
	MetadataTypeCount
)
