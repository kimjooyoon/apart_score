package metadata

import "fmt"

// Index returns the integer index of the metadata type.
func (mt MetadataType) Index() int {
	return int(mt)
}

// String returns the English name of the metadata type.
func (mt MetadataType) String() string {
	if mt.IsValid() {
		return metadataInfos[mt].englishName
	}
	return fmt.Sprintf("Unknown(%d)", int(mt))
}

// KoreanName returns the Korean name of the metadata type.
func (mt MetadataType) KoreanName() string {
	if mt.IsValid() {
		return metadataInfos[mt].koreanName
	}
	return fmt.Sprintf("알 수 없음(%d)", int(mt))
}

// Description returns the description of the metadata type.
func (mt MetadataType) Description() string {
	if mt.IsValid() {
		return metadataInfos[mt].description
	}
	return fmt.Sprintf("설명 없음(%d)", int(mt))
}

// IsValid checks if the metadata type is valid.
func (mt MetadataType) IsValid() bool {
	return mt >= 0 && mt < MetadataTypeCount
}

// FactorType returns the factor type (internal/external) of the metadata.
func (mt MetadataType) FactorType() FactorType {
	if mt.IsValid() {
		return metadataInfos[mt].factorType
	}
	return FactorInternal
}

// GetDefaultFactorTypes returns the default factor types for all metadata.
func GetDefaultFactorTypes() map[MetadataType]FactorType {
	result := make(map[MetadataType]FactorType)
	for mt := MetadataType(0); mt < MetadataTypeCount; mt++ {
		result[mt] = mt.FactorType()
	}
	return result
}

// SetFactorType sets the factor type for a specific metadata type.
func SetFactorType(mt MetadataType, factorType FactorType) error {
	if !mt.IsValid() {
		return fmt.Errorf("유효하지 않은 메타데이터 타입: %d", mt)
	}
	if factorType != FactorInternal && factorType != FactorExternal {
		return fmt.Errorf("유효하지 않은 팩터 타입: %s", factorType)
	}
	metadataInfos[mt].factorType = factorType
	return nil
}

// GetMetadataByFactorType returns all metadata types of the specified factor type.
func GetMetadataByFactorType(factorType FactorType) [MetadataTypeCount]MetadataType {
	var result [MetadataTypeCount]MetadataType
	index := 0
	for mt := MetadataType(0); mt < MetadataTypeCount; mt++ {
		if mt.FactorType() == factorType && index < int(MetadataTypeCount) {
			result[index] = mt
			index++
		}
	}
	return result
}
