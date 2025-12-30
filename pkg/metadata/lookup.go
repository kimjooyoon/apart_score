package metadata

func GetByIndex(index int) (MetadataType, bool) {
	mt := MetadataType(index)
	if mt.IsValid() {
		return mt, true
	}
	return MetadataType(-1), false
}
func GetByEnglishName(englishName string) (MetadataType, bool) {
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		if metadataInfos[i].englishName == englishName {
			return i, true
		}
	}
	return MetadataType(-1), false
}
func GetByKoreanName(koreanName string) (MetadataType, bool) {
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		if metadataInfos[i].koreanName == koreanName {
			return i, true
		}
	}
	return MetadataType(-1), false
}
