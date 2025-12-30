// 정량화 검증을 위한 테스트 계산
package main

import (
	"fmt"
	"math"
)

func main() {
	// 샘플 데이터
	scores := []float64{80.0, 85.0, 90.0}  // 점수들
	weights := []float64{0.3, 0.4, 0.3}    // 가중치 (합계 = 1.0)

	// 1. 가중치 합계 검증
	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}
	fmt.Printf("가중치 합계: %.6f (기대값: 1.0)\n", totalWeight)

	// 2. 가중치 합계 계산
	weightedSum := 0.0
	for i, score := range scores {
		weightedScore := score * weights[i]
		weightedSum += weightedScore
		fmt.Printf("요인 %d: %.1f × %.1f = %.1f\n", i+1, score, weights[i], weightedScore)
	}
	fmt.Printf("최종 점수: %.1f\n", weightedSum)

	// 3. 정밀도 테스트
	if math.Abs(totalWeight - 1.0) < 0.000001 {
		fmt.Println("✅ 가중치 정규화 성공")
	}

	// 4. 범위 검증
	if weightedSum >= 0 && weightedSum <= 100 {
		fmt.Println("✅ 점수 범위 유효")
	}

	// 5. 백분위수 계산 검증
	minScore, maxScore := 70.0, 90.0
	testScore := 80.0
	percentile := (testScore - minScore) / (maxScore - minScore) * 100.0
	fmt.Printf("백분위수 계산: (%.1f - %.1f) / (%.1f - %.1f) × 100 = %.1f%%\n",
		testScore, minScore, maxScore, minScore, percentile)
}
