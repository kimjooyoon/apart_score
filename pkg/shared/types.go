// Package shared provides common types and utilities used across multiple packages.
package shared

// ScoreValue represents an apartment score using integer arithmetic for precision.
type ScoreValue int

// Weight represents the weight/importance of a metadata factor using integer arithmetic.
type Weight int

const (
	// ScoreScale is the scaling factor for score values (1000 = 3 decimal places).
	ScoreScale = 1000
	// WeightScale is the scaling factor for weights (1000 = 3 decimal places).
	WeightScale = 1000
)

// ToFloat converts ScoreValue to float64 for external use.
func (s ScoreValue) ToFloat() float64 {
	return float64(s) / ScoreScale
}
func ScoreValueFromFloat(f float64) ScoreValue {
	return ScoreValue(f * ScoreScale)
}
func (w Weight) ToFloat() float64 {
	return float64(w) / WeightScale
}
func WeightFromFloat(f float64) Weight {
	return Weight(f * WeightScale)
}

// ScoreArray represents an array of scores indexed by MetadataType.
type ScoreArray [14]ScoreValue

// WeightArray represents an array of weights indexed by MetadataType.
type WeightArray [14]Weight

// NewScoreArrayFromMap creates a ScoreArray from a map using int as key.
func NewScoreArrayFromMap(scores map[int]ScoreValue) ScoreArray {
	var arr ScoreArray
	for i := 0; i < 14; i++ {
		if score, exists := scores[i]; exists {
			arr[i] = score
		}
	}
	return arr
}

// NewWeightArrayFromMap creates a WeightArray from a map using int as key.
func NewWeightArrayFromMap(weights map[int]Weight) WeightArray {
	var arr WeightArray
	for i := 0; i < 14; i++ {
		if weight, exists := weights[i]; exists {
			arr[i] = weight
		}
	}
	return arr
}

// ToMap converts ScoreArray to a map for compatibility.
func (sa ScoreArray) ToMap() map[int]ScoreValue {
	result := make(map[int]ScoreValue)
	for i, score := range sa {
		if score != 0 {
			result[i] = score
		}
	}
	return result
}

// ToMap converts WeightArray to a map for compatibility.
func (wa WeightArray) ToMap() map[int]Weight {
	result := make(map[int]Weight)
	for i, weight := range wa {
		if weight != 0 {
			result[i] = weight
		}
	}
	return result
}

// DomainGlossary defines standard terminology for the scoring system.
var DomainGlossary = map[string]string{
	"Score":          "0-100 스케일의 상대적 평가 점수 (절대적 우수도를 의미하지 않음)",
	"Weight":         "각 요소의 상대적 중요도 (0.0-1.0, 총합 = 1.0)",
	"Strategy":       "점수 계산을 위한 수학적 방법 (Weighted Sum, Geometric Mean 등)",
	"Scenario":       "특정 사용자 그룹을 위한 가중치 프리셋 (Balanced, Transportation 등)",
	"Percentile":     "전체 비교 대상 중 상위 몇 %인지 (예: 85백분위수 = 상위 15%)",
	"Grade":          "점수를 5단계로 단순화한 등급 (A/B/C/D/F)",
	"Absolute Score": "개별 아파트의 절대적 평가 점수",
	"Relative Score": "다른 아파트들과의 비교를 통한 상대적 평가",
	"Subjectivity":   "사용자의 개인적 선호도가 평가 결과에 미치는 영향 정도",
	"Objectivity":    "데이터 기반의 공통적 평가 기준이 적용되는 정도",
}
