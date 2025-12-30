# Metadata Package

메타데이터 정의 및 관리를 위한 패키지입니다.

## 주요 기능

- 메타데이터 타입 정의 (14개 요소)
- 다국어 지원 (영문/한글)
- 인덱스 기반 참조
- 이름 기반 검색

## 메타데이터 목록

| 인덱스 | 영문명 | 한글명 |
|--------|--------|--------|
| 0 | Floor Level | 층수 |
| 1 | Distance to Station | 역까지 거리 |
| 2 | Elevator Presence | 엘리베이터 유무 |
| 3 | Construction Year | 건축년도 |
| 4 | Construction Company | 건설회사 |
| 5 | Apartment Size | 아파트 크기 |
| 6 | Nearby Amenities | 주변 편의시설 |
| 7 | Transportation Access | 교통 접근성 |
| 8 | School District | 학군 |
| 9 | Crime Rate | 범죄율 |
| 10 | Green Space Ratio | 녹지율 |
| 11 | Parking | 주차장 |
| 12 | Maintenance Fee | 관리비 |
| 13 | Heating System | 난방 방식 |

## 사용 예제

```go
import "apart_score/pkg/metadata"

// 기본 조회
fmt.Println(metadata.FloorLevel.String())     // "Floor Level"
fmt.Println(metadata.FloorLevel.KoreanName()) // "층수"
fmt.Println(metadata.FloorLevel.Index())      // 0

// 검색
mt, ok := metadata.GetByEnglishName("Floor Level")
mt, ok := metadata.GetByKoreanName("층수")
mt, ok := metadata.GetByIndex(0)

// 전체 목록
allTypes := metadata.AllMetadataTypes()
for _, mt := range allTypes {
    fmt.Printf("%d: %s\n", mt.Index(), mt.String())
}
```

## 설계 원칙

1. **수정 불가**: 한 번 정의된 메타데이터는 변경 불가
2. **순차 추가**: 새로운 메타데이터는 마지막에만 추가
3. **인덱스 안정성**: iota 기반으로 인덱스 보장
4. **다국어 지원**: 영문/한글 모두 지원
