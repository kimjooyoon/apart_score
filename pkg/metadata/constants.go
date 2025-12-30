package metadata

type metadataInfo struct {
	englishName string
	koreanName  string
	description string
}

var metadataInfos = [MetadataTypeCount]metadataInfo{
	{
		englishName: "Floor Level",
		koreanName:  "층수",
		description: "아파트의 층수 정보 (중간층에 가까울수록 높은 점수)",
	},
	{
		englishName: "Distance to Station",
		koreanName:  "역까지 거리",
		description: "가장 가까운 역까지의 거리 (가까울수록 높은 점수)",
	},
	{
		englishName: "Elevator Presence",
		koreanName:  "엘리베이터 유무",
		description: "엘리베이터 설치 여부 (있으면 높은 점수)",
	},
	{
		englishName: "Construction Year",
		koreanName:  "건축년도",
		description: "아파트 건축 연도 (최신 건물일수록 높은 점수)",
	},
	{
		englishName: "Construction Company",
		koreanName:  "건설회사",
		description: "아파트를 건설한 회사 (신뢰할 수 있는 회사일수록 높은 점수)",
	},
	{
		englishName: "Apartment Size",
		koreanName:  "아파트 크기",
		description: "아파트의 크기/면적 (적절한 크기일수록 높은 점수)",
	},
	{
		englishName: "Nearby Amenities",
		koreanName:  "주변 편의시설",
		description: "주변 편의시설의充実度 (편의시설이 많을수록 높은 점수)",
	},
	{
		englishName: "Transportation Access",
		koreanName:  "교통 접근성",
		description: "대중교통 접근성 (대중교통 접근성이 좋을수록 높은 점수)",
	},
	{
		englishName: "School District",
		koreanName:  "학군",
		description: "주변 학군 수준 (좋은 학군일수록 높은 점수)",
	},
	{
		englishName: "Crime Rate",
		koreanName:  "범죄율",
		description: "주변 지역 범죄율 (낮을수록 높은 점수)",
	},
	{
		englishName: "Green Space Ratio",
		koreanName:  "녹지율",
		description: "주변 녹지 공간 비율 (높을수록 높은 점수)",
	},
	{
		englishName: "Parking",
		koreanName:  "주차장",
		description: "주차 공간 sufficiency (주차 공간이 충분할수록 높은 점수)",
	},
	{
		englishName: "Maintenance Fee",
		koreanName:  "관리비",
		description: "월 관리비 수준 (적절한 수준일수록 높은 점수)",
	},
	{
		englishName: "Heating System",
		koreanName:  "난방 방식",
		description: "난방 시스템 종류 (효율적인 난방 방식일수록 높은 점수)",
	},
}
