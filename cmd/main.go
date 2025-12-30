package main

import (
	"fmt"

	"apart_score/pkg/metadata"
)

func main() {
	fmt.Println("아파트 스코어링 시스템 시작")

	// 메타데이터 사용 예시
	fmt.Printf("층수 메타데이터: %s (%s)\n", metadata.FloorLevel.String(), metadata.FloorLevel.KoreanName())
	fmt.Printf("역까지 거리 메타데이터: %s (%s)\n", metadata.DistanceToStation.String(), metadata.DistanceToStation.KoreanName())

	// 모든 메타데이터 출력
	fmt.Println("\n=== 모든 메타데이터 목록 ===")
	for i := metadata.MetadataType(0); i < metadata.MetadataTypeCount; i++ {
		fmt.Printf("%d: %s (%s)\n", i.Index(), i.String(), i.KoreanName())
	}
}
