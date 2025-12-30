package shared

import "apart_score/pkg/metadata"

// NormalizeWeights는 정수 기반 가중치를 합계가 WeightScale이 되도록 정규화합니다.
func NormalizeWeights(weights map[metadata.MetadataType]Weight) map[metadata.MetadataType]Weight {
	var total Weight
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		return weights
	}

	normalized := make(map[metadata.MetadataType]Weight)

	// 정수 나눗셈으로 정규화 (반올림 고려)
	for mt, w := range weights {
		// 정수 나눗셈: (w * WeightScale + total/2) / total 로 반올림
		scaled := (w * WeightScale + total/2) / total
		normalized[mt] = scaled
	}

	return normalized
}
