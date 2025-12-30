package metadata

// GetByIndex returns the metadata type for the given index.
func GetByIndex(index int) (MetadataType, bool) {
	mt := MetadataType(index)
	if mt.IsValid() {
		return mt, true
	}
	return MetadataType(-1), false
}

// GetByEnglishName returns the metadata type for the given English name.
func GetByEnglishName(englishName string) (MetadataType, bool) {
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		if metadataInfos[i].englishName == englishName {
			return i, true
		}
	}
	return MetadataType(-1), false
}

// GetByKoreanName returns the metadata type for the given Korean name.
func GetByKoreanName(koreanName string) (MetadataType, bool) {
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		if metadataInfos[i].koreanName == koreanName {
			return i, true
		}
	}
	return MetadataType(-1), false
}
