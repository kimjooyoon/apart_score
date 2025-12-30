package metadata

import (
	"testing"
)

func TestMetadataType_Index(t *testing.T) {
	tests := []struct {
		name     string
		metadata MetadataType
		expected int
	}{
		{"FloorLevel", FloorLevel, 0},
		{"DistanceToStation", DistanceToStation, 1},
		{"ElevatorPresence", ElevatorPresence, 2},
		{"ConstructionYear", ConstructionYear, 3},
		{"ConstructionCompany", ConstructionCompany, 4},
		{"ApartmentSize", ApartmentSize, 5},
		{"NearbyAmenities", NearbyAmenities, 6},
		{"TransportationAccess", TransportationAccess, 7},
		{"SchoolDistrict", SchoolDistrict, 8},
		{"CrimeRate", CrimeRate, 9},
		{"GreenSpaceRatio", GreenSpaceRatio, 10},
		{"Parking", Parking, 11},
		{"MaintenanceFee", MaintenanceFee, 12},
		{"HeatingSystem", HeatingSystem, 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metadata.Index(); got != tt.expected {
				t.Errorf("Index() = %v, expected %v", got, tt.expected)
			}
		})
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
		name      string
		englishName string
		expected  MetadataType
		expectErr bool
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
		name      string
		koreanName string
		expected  MetadataType
		expectErr bool
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
