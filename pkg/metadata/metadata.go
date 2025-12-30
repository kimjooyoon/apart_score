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
