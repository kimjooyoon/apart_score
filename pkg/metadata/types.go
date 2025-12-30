package metadata

// MetadataType은 아파트 스코어링에 사용되는 메타데이터 타입을 정의합니다.
// iota를 사용하여 순차적 인덱스를 부여하며, 한 번 정의된 값은 수정되지 않습니다.
type MetadataType int

// 메타데이터 타입 상수 정의 (iota 사용)
// 주의: 이 값들은 수정되지 않으며, 새로운 값만 추가할 수 있습니다.
const (
	FloorLevel            MetadataType = iota // 층수
	DistanceToStation                         // 역까지 거리
	ElevatorPresence                          // 엘리베이터 유무
	ConstructionYear                          // 건축년도
	ConstructionCompany                       // 건설회사
	ApartmentSize                             // 아파트 크기/면적
	NearbyAmenities                          // 주변 편의시설
	TransportationAccess                     // 교통 접근성
	SchoolDistrict                           // 학군
	CrimeRate                                // 범죄율
	GreenSpaceRatio                          // 녹지율
	Parking                                  // 주차장
	MaintenanceFee                           // 관리비
	HeatingSystem                            // 난방 방식

	// MetadataTypeCount는 메타데이터 타입의 총 개수를 나타냅니다.
	// 새로운 메타데이터를 추가할 때는 이 값 위에 추가해야 합니다.
	MetadataTypeCount
)
