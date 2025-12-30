// Package shared provides common types and utilities shared across packages.
package shared

import "apart_score/pkg/metadata"

// CachedMetadataTypes holds pre-computed array of all metadata types for performance.
var CachedMetadataTypes = func() [metadata.MetadataTypeCount]metadata.MetadataType {
	var types [metadata.MetadataTypeCount]metadata.MetadataType
	for i := metadata.MetadataType(0); i < metadata.MetadataTypeCount; i++ {
		types[i] = i
	}
	return types
}()

// FastAllMetadataTypes returns the cached array of all metadata types for performance.
func FastAllMetadataTypes() [metadata.MetadataTypeCount]metadata.MetadataType {
	return CachedMetadataTypes
}
