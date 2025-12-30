package shared

// 내부적으로는 1000배 스케일링된 정수를 사용 (정밀도 유지)
type ScoreValue int     // 0-100000 (0.0-100.0점)
type Weight int         // 0-1000 (0.0-1.0 가중치)

// 정수 ↔ 실수 변환 상수
const (
	ScoreScale = 1000  // 점수 스케일링 (소수점 3자리 정밀도)
	WeightScale = 1000 // 가중치 스케일링
)

// ToFloat converts ScoreValue to float64 for external use
func (s ScoreValue) ToFloat() float64 {
	return float64(s) / ScoreScale
}

// FromFloat converts float64 to ScoreValue for internal use
func ScoreValueFromFloat(f float64) ScoreValue {
	return ScoreValue(f * ScoreScale)
}

// ToFloat converts Weight to float64 for external use
func (w Weight) ToFloat() float64 {
	return float64(w) / WeightScale
}

// FromFloat converts float64 to Weight for internal use
func WeightFromFloat(f float64) Weight {
	return Weight(f * WeightScale)
}
