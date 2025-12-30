package metadata

// AllMetadataTypes는 모든 메타데이터 타입을 배열로 반환합니다.
func AllMetadataTypes() [MetadataTypeCount]MetadataType {
	var types [MetadataTypeCount]MetadataType
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		types[i] = i
	}
	return types
}
