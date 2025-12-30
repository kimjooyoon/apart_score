package metadata

import (
	"testing"
)

func TestMetadataType_OrderStability(t *testing.T) {
	// enum 순서 변경을 감지하기 위한 테스트
	// 이름이 순서대로 매칭되었는지 검증 (크기는 검증하지 않음)

	orderMap := map[string]int{
		"Floor Level":           0,
		"Distance to Station":   1,
		"Elevator Presence":     2,
		"Construction Year":     3,
		"Construction Company":  4,
		"Apartment Size":        5,
		"Nearby Amenities":      6,
		"Transportation Access": 7,
		"School District":       8,
		"Crime Rate":            9,
		"Green Space Ratio":     10,
		"Parking":               11,
		"Maintenance Fee":       12,
		"Heating System":        13,
	}

	// 대표적인 값들로 순서 검증 (모두 테스트하지 않음)
	keyChecks := []MetadataType{FloorLevel, ConstructionYear, SchoolDistrict, HeatingSystem}

	for _, mt := range keyChecks {
		expectedIndex, exists := orderMap[mt.String()]
		if !exists {
			t.Errorf("정의되지 않은 메타데이터: %s", mt.String())
			continue
		}

		if got := mt.Index(); got != expectedIndex {
			t.Errorf("순서 변경 감지: %s의 인덱스가 %d에서 %d로 변경됨 (iota 기반)", mt.String(), expectedIndex, got)
		}
	}

	// iota 기반 순차 할당 검증
	if FloorLevel.Index() != 0 {
		t.Error("FloorLevel이 0번이 아님 - iota 순서 변경됨")
	}

	if HeatingSystem.Index() != 13 {
		t.Error("HeatingSystem이 13번이 아님 - iota 순서 변경됨")
	}
}

func TestFactorType(t *testing.T) {
	tests := []struct {
		name        string
		metadata    MetadataType
		expected    FactorType
		description string
	}{
		{"FloorLevel", FloorLevel, FactorInternal, "층수는 아파트 내부 요인"},
		{"DistanceToStation", DistanceToStation, FactorExternal, "역거리는 외부 환경 요인"},
		{"ElevatorPresence", ElevatorPresence, FactorInternal, "엘리베이터는 아파트 내부 요인"},
		{"ConstructionYear", ConstructionYear, FactorInternal, "건축년도는 아파트 내부 요인"},
		{"NearbyAmenities", NearbyAmenities, FactorExternal, "주변 편의시설은 외부 요인"},
		{"SchoolDistrict", SchoolDistrict, FactorExternal, "학군은 외부 환경 요인"},
		{"Parking", Parking, FactorInternal, "주차장은 아파트 내부 요인"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metadata.FactorType(); got != tt.expected {
				t.Errorf("%s: FactorType() = %v, expected %v (%s)", tt.name, got, tt.expected, tt.description)
			}
		})
	}
}

func TestGetDefaultFactorTypes(t *testing.T) {
	defaultTypes := GetDefaultFactorTypes()

	if len(defaultTypes) != int(MetadataTypeCount) {
		t.Errorf("Expected %d factor types, got %d", MetadataTypeCount, len(defaultTypes))
	}

	// 특정 값들 검증
	if defaultTypes[FloorLevel] != FactorInternal {
		t.Error("FloorLevel should be internal factor")
	}

	if defaultTypes[DistanceToStation] != FactorExternal {
		t.Error("DistanceToStation should be external factor")
	}

	if defaultTypes[SchoolDistrict] != FactorExternal {
		t.Error("SchoolDistrict should be external factor")
	}
}

func TestSetFactorType(t *testing.T) {
	// 원래 값 저장
	originalType := FloorLevel.FactorType()

	// 팩터 타입 변경
	err := SetFactorType(FloorLevel, FactorExternal)
	if err != nil {
		t.Fatalf("SetFactorType failed: %v", err)
	}

	// 변경된 값 검증
	if FloorLevel.FactorType() != FactorExternal {
		t.Error("FactorType should be changed to external")
	}

	// 다시 원래 값으로 복원
	err = SetFactorType(FloorLevel, originalType)
	if err != nil {
		t.Fatalf("SetFactorType failed to restore: %v", err)
	}

	if FloorLevel.FactorType() != originalType {
		t.Error("FactorType should be restored to original")
	}

	// 유효하지 않은 메타데이터 타입 테스트
	err = SetFactorType(MetadataType(-1), FactorInternal)
	if err == nil {
		t.Error("SetFactorType should fail for invalid metadata type")
	}

	// 유효하지 않은 팩터 타입 테스트
	err = SetFactorType(FloorLevel, FactorType("invalid"))
	if err == nil {
		t.Error("SetFactorType should fail for invalid factor type")
	}
}

func TestGetMetadataByFactorType(t *testing.T) {
	internalTypes := GetMetadataByFactorType(FactorInternal)
	externalTypes := GetMetadataByFactorType(FactorExternal)

	// 총 개수 검증
	if len(internalTypes)+len(externalTypes) != int(MetadataTypeCount) {
		t.Errorf("Total metadata types should be %d, got %d internal + %d external",
			MetadataTypeCount, len(internalTypes), len(externalTypes))
	}

	// 내부 요인에 층수가 포함되는지 검증
	found := false
	for _, mt := range internalTypes {
		if mt == FloorLevel {
			found = true
			break
		}
	}
	if !found {
		t.Error("FloorLevel should be in internal factors")
	}

	// 외부 요인에 학군이 포함되는지 검증
	found = false
	for _, mt := range externalTypes {
		if mt == SchoolDistrict {
			found = true
			break
		}
	}
	if !found {
		t.Error("SchoolDistrict should be in external factors")
	}
}

func TestMetadataType_String(t *testing.T) {
	tests := []struct {
		name     string
		metadata MetadataType
		expected string
	}{
		{"FloorLevel", FloorLevel, "Floor Level"},
		{"DistanceToStation", DistanceToStation, "Distance to Station"},
		{"ElevatorPresence", ElevatorPresence, "Elevator Presence"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metadata.String(); got != tt.expected {
				t.Errorf("String() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestMetadataType_KoreanName(t *testing.T) {
	tests := []struct {
		name     string
		metadata MetadataType
		expected string
	}{
		{"FloorLevel", FloorLevel, "층수"},
		{"DistanceToStation", DistanceToStation, "역까지 거리"},
		{"ElevatorPresence", ElevatorPresence, "엘리베이터 유무"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metadata.KoreanName(); got != tt.expected {
				t.Errorf("KoreanName() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestGetByIndex(t *testing.T) {
	tests := []struct {
		name      string
		index     int
		expected  MetadataType
		expectErr bool
	}{
		{"Valid FloorLevel", 0, FloorLevel, false},
		{"Valid DistanceToStation", 1, DistanceToStation, false},
		{"Invalid negative", -1, MetadataType(-1), true},
		{"Invalid out of range", 100, MetadataType(-1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetByIndex(tt.index)
			if tt.expectErr {
				if ok {
					t.Errorf("GetByIndex(%d) expected error, but got %v", tt.index, got)
				}
			} else {
				if !ok || got != tt.expected {
					t.Errorf("GetByIndex(%d) = %v, %v, expected %v, true", tt.index, got, ok, tt.expected)
				}
			}
		})
	}
}

func TestGetByEnglishName(t *testing.T) {
	tests := []struct {
		name        string
		englishName string
		expected    MetadataType
		expectErr   bool
	}{
		{"Valid Floor Level", "Floor Level", FloorLevel, false},
		{"Valid Distance to Station", "Distance to Station", DistanceToStation, false},
		{"Invalid name", "Invalid Name", MetadataType(-1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetByEnglishName(tt.englishName)
			if tt.expectErr {
				if ok {
					t.Errorf("GetByEnglishName(%s) expected error, but got %v", tt.englishName, got)
				}
			} else {
				if !ok || got != tt.expected {
					t.Errorf("GetByEnglishName(%s) = %v, %v, expected %v, true", tt.englishName, got, ok, tt.expected)
				}
			}
		})
	}
}

func TestGetByKoreanName(t *testing.T) {
	tests := []struct {
		name       string
		koreanName string
		expected   MetadataType
		expectErr  bool
	}{
		{"Valid 층수", "층수", FloorLevel, false},
		{"Valid 역까지 거리", "역까지 거리", DistanceToStation, false},
		{"Invalid name", "잘못된 이름", MetadataType(-1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetByKoreanName(tt.koreanName)
			if tt.expectErr {
				if ok {
					t.Errorf("GetByKoreanName(%s) expected error, but got %v", tt.koreanName, got)
				}
			} else {
				if !ok || got != tt.expected {
					t.Errorf("GetByKoreanName(%s) = %v, %v, expected %v, true", tt.koreanName, got, ok, tt.expected)
				}
			}
		})
	}
}

func TestAllMetadataTypes(t *testing.T) {
	allTypes := AllMetadataTypes()

	if len(allTypes) != int(MetadataTypeCount) {
		t.Errorf("AllMetadataTypes() length = %d, expected %d", len(allTypes), MetadataTypeCount)
	}

	// 각 타입이 올바른 순서로 있는지 확인
	for i, mt := range allTypes {
		if mt.Index() != i {
			t.Errorf("AllMetadataTypes()[%d] = %v, expected index %d", i, mt.Index(), i)
		}
	}
}
