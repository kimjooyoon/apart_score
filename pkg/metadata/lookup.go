package metadata

// GetByIndex는 인덱스 번호로 메타데이터 타입을 반환합니다.
func GetByIndex(index int) (MetadataType, bool) {
	mt := MetadataType(index)
	if mt.IsValid() {
		return mt, true
	}
	return MetadataType(-1), false
}

// GetByEnglishName은 영문명으로 메타데이터 타입을 반환합니다.
func GetByEnglishName(englishName string) (MetadataType, bool) {
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		if metadataInfos[i].englishName == englishName {
			return i, true
		}
	}
	return MetadataType(-1), false
}

// GetByKoreanName은 한글명으로 메타데이터 타입을 반환합니다.
func GetByKoreanName(koreanName string) (MetadataType, bool) {
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		if metadataInfos[i].koreanName == koreanName {
			return i, true
		}
	}
	return MetadataType(-1), false
}
