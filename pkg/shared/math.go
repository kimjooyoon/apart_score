package shared

// MulDiv performs safe integer multiplication and division to prevent overflow.
func MulDiv(a, b, divisor ScoreValue) ScoreValue {
	product := int64(a) * int64(b)
	return ScoreValue(product / int64(divisor))
}
func MulDivWeight(score ScoreValue, weight Weight) ScoreValue {
	return MulDiv(score, ScoreValue(weight), ScoreValue(WeightScale))
}
func FastAverage(totalWeightedSum float64, totalWeight Weight) float64 {
	return totalWeightedSum
}
