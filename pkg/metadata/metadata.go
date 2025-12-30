package metadata

import "fmt"

func (mt MetadataType) Index() int {
	return int(mt)
}
func (mt MetadataType) String() string {
	if mt.IsValid() {
		return metadataInfos[mt].englishName
	}
	return fmt.Sprintf("Unknown(%d)", int(mt))
}
func (mt MetadataType) KoreanName() string {
	if mt.IsValid() {
		return metadataInfos[mt].koreanName
	}
	return fmt.Sprintf("알 수 없음(%d)", int(mt))
}
func (mt MetadataType) Description() string {
	if mt.IsValid() {
		return metadataInfos[mt].description
	}
	return fmt.Sprintf("설명 없음(%d)", int(mt))
}
func (mt MetadataType) IsValid() bool {
	return mt >= 0 && mt < MetadataTypeCount
}
func (mt MetadataType) FactorType() FactorType {
	if mt.IsValid() {
		return metadataInfos[mt].factorType
	}
	return FactorInternal
}
func GetDefaultFactorTypes() map[MetadataType]FactorType {
	result := make(map[MetadataType]FactorType)
	for mt := MetadataType(0); mt < MetadataTypeCount; mt++ {
		result[mt] = mt.FactorType()
	}
	return result
}
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
func GetMetadataByFactorType(factorType FactorType) []MetadataType {
	var result []MetadataType
	for mt := MetadataType(0); mt < MetadataTypeCount; mt++ {
		if mt.FactorType() == factorType {
			result = append(result, mt)
		}
	}
	return result
}
