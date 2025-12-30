package metadata

// AllMetadataTypes는 모든 유효한 메타데이터 타입의 슬라이스를 반환합니다.
func AllMetadataTypes() []MetadataType {
	types := make([]MetadataType, MetadataTypeCount)
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		types[i] = i
	}
	return types
}
