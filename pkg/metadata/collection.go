package metadata

func AllMetadataTypes() []MetadataType {
	types := make([]MetadataType, MetadataTypeCount)
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		types[i] = i
	}
	return types
}
