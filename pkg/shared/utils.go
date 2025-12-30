package shared

import "apart_score/pkg/metadata"

// NormalizeWeights normalizes weights to sum to WeightScale using integer arithmetic.
func NormalizeWeights(weights map[metadata.MetadataType]Weight) map[metadata.MetadataType]Weight {
	total := Weight(0)
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		return weights
	}
	normalized := make(map[metadata.MetadataType]Weight)
	for mt, w := range weights {
		scaled := (w*WeightScale + total/2) / total
		normalized[mt] = scaled
	}
	return normalized
}
