// Package metadata provides apartment scoring metadata types and utilities.
package metadata

// AllMetadataTypes returns all metadata types as an array.
func AllMetadataTypes() [MetadataTypeCount]MetadataType {
	var types [MetadataTypeCount]MetadataType
	for i := MetadataType(0); i < MetadataTypeCount; i++ {
		types[i] = i
	}
	return types
}
