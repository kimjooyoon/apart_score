package metadata

// metadataInfo는 각 메타데이터 타입의 정보를 담는 구조체입니다.
type metadataInfo struct {
	englishName string // 영문명
	koreanName  string // 한글명
	description string // 설명
}

// metadataInfos는 각 MetadataType에 대한 정보를 담는 배열입니다.
// MetadataType iota 순서와 동일하게 정렬되어 있습니다.
var metadataInfos = [MetadataTypeCount]metadataInfo{
	// FloorLevel (0)
	{
		englishName: "Floor Level",
		koreanName:  "층수",
		description: "아파트의 층수 정보 (중간층에 가까울수록 높은 점수)",
	},
	// DistanceToStation (1)
	{
		englishName: "Distance to Station",
		koreanName:  "역까지 거리",
		description: "가장 가까운 역까지의 거리 (가까울수록 높은 점수)",
	},
	// ElevatorPresence (2)
	{
		englishName: "Elevator Presence",
		koreanName:  "엘리베이터 유무",
		description: "엘리베이터 설치 여부 (있으면 높은 점수)",
	},
	// ConstructionYear (3)
	{
		englishName: "Construction Year",
		koreanName:  "건축년도",
		description: "아파트 건축 연도 (최신 건물일수록 높은 점수)",
	},
	// ConstructionCompany (4)
	{
		englishName: "Construction Company",
		koreanName:  "건설회사",
		description: "아파트를 건설한 회사 (신뢰할 수 있는 회사일수록 높은 점수)",
	},
	// ApartmentSize (5)
	{
		englishName: "Apartment Size",
		koreanName:  "아파트 크기",
		description: "아파트의 크기/면적 (적절한 크기일수록 높은 점수)",
	},
	// NearbyAmenities (6)
	{
		englishName: "Nearby Amenities",
		koreanName:  "주변 편의시설",
		description: "주변 편의시설의充実度 (편의시설이 많을수록 높은 점수)",
	},
	// TransportationAccess (7)
	{
		englishName: "Transportation Access",
		koreanName:  "교통 접근성",
		description: "대중교통 접근성 (대중교통 접근성이 좋을수록 높은 점수)",
	},
	// SchoolDistrict (8)
	{
		englishName: "School District",
		koreanName:  "학군",
		description: "주변 학군 수준 (좋은 학군일수록 높은 점수)",
	},
	// CrimeRate (9)
	{
		englishName: "Crime Rate",
		koreanName:  "범죄율",
		description: "주변 지역 범죄율 (낮을수록 높은 점수)",
	},
	// GreenSpaceRatio (10)
	{
		englishName: "Green Space Ratio",
		koreanName:  "녹지율",
		description: "주변 녹지 공간 비율 (높을수록 높은 점수)",
	},
	// Parking (11)
	{
		englishName: "Parking",
		koreanName:  "주차장",
		description: "주차 공간 sufficiency (주차 공간이 충분할수록 높은 점수)",
	},
	// MaintenanceFee (12)
	{
		englishName: "Maintenance Fee",
		koreanName:  "관리비",
		description: "월 관리비 수준 (적절한 수준일수록 높은 점수)",
	},
	// HeatingSystem (13)
	{
		englishName: "Heating System",
		koreanName:  "난방 방식",
		description: "난방 시스템 종류 (효율적인 난방 방식일수록 높은 점수)",
	},
}
