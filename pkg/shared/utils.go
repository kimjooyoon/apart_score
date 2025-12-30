package shared

import "apart_score/pkg/metadata"

func NormalizeWeights(weights map[metadata.MetadataType]Weight) map[metadata.MetadataType]Weight {
	var total Weight
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		return weights
	}
	normalized := make(map[metadata.MetadataType]Weight)
	for mt, w := range weights {
		normalized[mt] = w / total
	}
	return normalized
}
