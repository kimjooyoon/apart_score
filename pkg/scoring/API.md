# Scoring Package API Reference

이 문서는 아파트 스코어링 시스템의 pkg/scoring 패키지 API를 설명합니다.

## 목차

- [기본 개념](#기본-개념)
- [메인 API](#메인-api)
- [타입 정의](#타입-정의)
- [사용 예제](#사용-예제)

## 기본 개념

### ScoreValue
각 메타데이터의 점수 값 (0-100 사이의 float64)

### Weight
메타데이터의 가중치 (0.0-1.0 사이의 float64, 합계는 1.0)

### ScoringScenario
미리 정의된 가중치 시나리오

## 메인 API

### QuickScore

빠른 점수 계산을 위한 헬퍼 함수

```go
func QuickScore(scores map[metadata.MetadataType]ScoreValue, scenario ScoringScenario) (*ScoreResult, error)
```

**파라미터:**
- `scores`: 메타데이터별 점수 맵
- `scenario`: 적용할 시나리오

**반환값:**
- `*ScoreResult`: 계산된 점수 결과
- `error`: 계산 중 발생한 오류

**예제:**
```go
scores := map[metadata.MetadataType]ScoreValue{
    metadata.FloorLevel: 85.0,
    metadata.DistanceToStation: 90.0,
    // ...
}

result, err := scoring.QuickScore(scores, scoring.ScenarioBalanced)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("총점: %.1f점\n", result.TotalScore)
```

### NewScoringManager

상세한 점수 계산을 위한 매니저 생성

```go
func NewScoringManager(profile ScoringProfile) (*ScoringManager, error)
```

**파라미터:**
- `profile`: 스코어링 프로필 설정

**반환값:**
- `*ScoringManager`: 스코어링 매니저 인스턴스
- `error`: 초기화 중 발생한 오류

### ScoringManager.CalculateScore

매니저를 통한 점수 계산

```go
func (sm *ScoringManager) CalculateScore(scores map[metadata.MetadataType]ScoreValue) (*ScoreResult, error)
```

**파라미터:**
- `scores`: 메타데이터별 점수 맵

**반환값:**
- `*ScoreResult`: 계산된 점수 결과
- `error`: 계산 중 발생한 오류

## 타입 정의

### ScoreResult

점수 계산 결과를 담는 구조체

```go
type ScoreResult struct {
    TotalScore       ScoreValue                           // 총점 (0-100)
    WeightedScores   map[metadata.MetadataType]ScoreValue // 각 메타데이터의 가중치 적용 점수
    RawScores        map[metadata.MetadataType]ScoreValue // 각 메타데이터의 원본 점수
    Weights          map[metadata.MetadataType]Weight     // 적용된 가중치
    Method           ScoringMethod                         // 사용된 계산 방법
    Scenario         ScoringScenario                       // 적용된 시나리오
    Grade            Grade                                 // 등급 (A, B, C, D, F)
    Percentile       float64                               // 백분위수 (0-100)
}
```

### ScoringScenario

사용 가능한 시나리오들

```go
const (
    ScenarioBalanced      ScoringScenario = "balanced"       // 균형 잡힌
    ScenarioTransportation ScoringScenario = "transportation" // 교통 중심
    ScenarioEducation      ScoringScenario = "education"      // 교육 중심
    ScenarioCostEffective  ScoringScenario = "cost_effective" // 가성비 중심
    ScenarioFamilyFriendly ScoringScenario = "family_friendly" // 가족 친화적
    ScenarioInvestment     ScoringScenario = "investment"     // 투자 가치 중심
)
```

### Grade

점수 등급

```go
const (
    GradeA Grade = "A" // 우수 (90-100)
    GradeB Grade = "B" // 양호 (80-89)
    GradeC Grade = "C" // 보통 (70-79)
    GradeD Grade = "D" // 미흡 (60-69)
    GradeF Grade = "F" // 불량 (0-59)
)
```

## 유틸리티 함수들

### GetScenarioWeights

시나리오별 가중치 조회

```go
func GetScenarioWeights(scenario ScoringScenario) map[metadata.MetadataType]Weight
```

### GetScenarioDescription

시나리오 설명 조회

```go
func GetScenarioDescription(scenario ScoringScenario) string
```

### GetAllScenarios

사용 가능한 모든 시나리오 조회

```go
func GetAllScenarios() []ScoringScenario
```

### AnalyzeScore

점수 상세 분석

```go
func AnalyzeScore(result *ScoreResult) *ScoreAnalysis
```

### FormatScoreResult

점수 결과를 읽기 쉽게 포맷팅

```go
func FormatScoreResult(result *ScoreResult) string
```

### RecommendScenario

점수 기반 시나리오 추천

```go
func RecommendScenario(scores map[metadata.MetadataType]ScoreValue) ScoringScenario
```

## 사용 예제

### 기본 사용법

```go
package main

import (
    "fmt"
    "apart_score/pkg/metadata"
    "apart_score/pkg/scoring"
)

func main() {
    // 아파트 점수 데이터 준비
    scores := map[metadata.MetadataType]scoring.ScoreValue{
        metadata.FloorLevel:         85.0,
        metadata.DistanceToStation:  90.0,
        metadata.ElevatorPresence:   100.0,
        metadata.ConstructionYear:   90.0,
        metadata.ConstructionCompany: 85.0,
        metadata.ApartmentSize:      75.0,
        metadata.NearbyAmenities:    80.0,
        metadata.TransportationAccess: 90.0,
        metadata.SchoolDistrict:     70.0,
        metadata.CrimeRate:          65.0,
        metadata.GreenSpaceRatio:    60.0,
        metadata.Parking:            80.0,
        metadata.MaintenanceFee:     75.0,
        metadata.HeatingSystem:      70.0,
    }

    // 빠른 점수 계산
    result, err := scoring.QuickScore(scores, scoring.ScenarioBalanced)
    if err != nil {
        fmt.Printf("점수 계산 실패: %v\n", err)
        return
    }

    // 결과 출력
    fmt.Println(scoring.FormatScoreResult(result))

    // 상세 분석
    analysis := scoring.AnalyzeScore(result)
    fmt.Printf("강점 개수: %d\n", len(analysis.Strengths))
    fmt.Printf("약점 개수: %d\n", len(analysis.Weaknesses))
    fmt.Printf("개선 제안 개수: %d\n", len(analysis.ImprovementTips))
}
```

### 시나리오 비교

```go
// 여러 시나리오로 비교
scenarios := scoring.GetAllScenarios()

fmt.Println("=== 시나리오별 점수 비교 ===")
for _, scenario := range scenarios {
    result, _ := scoring.QuickScore(scores, scenario)
    desc := scoring.GetScenarioDescription(scenario)
    fmt.Printf("%-15s: %.1f점 (%s)\n", desc, result.TotalScore, result.Grade)
}

// 추천 시나리오
recommended := scoring.RecommendScenario(scores)
fmt.Printf("\n추천 시나리오: %s\n", scoring.GetScenarioDescription(recommended))
```

### 커스텀 프로필 사용

```go
// 커스텀 가중치로 프로필 생성
customWeights := map[metadata.MetadataType]scoring.Weight{
    metadata.DistanceToStation: 0.3,  // 교통을 매우 중요하게
    metadata.SchoolDistrict:    0.2,  // 교육을 중요하게
    // ... 다른 가중치들
}

profile, err := scoring.CreateCustomProfile(
    "커스텀 프로필",
    "교통과 교육을 중시하는 프로필",
    customWeights,
)
if err != nil {
    fmt.Printf("프로필 생성 실패: %v\n", err)
    return
}

// 커스텀 프로필로 계산
manager, _ := scoring.NewScoringManager(*profile)
result, _ := manager.CalculateScore(scores)
fmt.Println(scoring.FormatScoreResult(result))
```

## 에러 처리

패키지의 모든 함수들은 적절한 에러를 반환합니다:

```go
result, err := scoring.QuickScore(scores, scenario)
if err != nil {
    switch e := err.(type) {
    case *scoring.ValidationError:
        fmt.Printf("검증 오류 (%s): %s\n", e.Field, e.Message)
    default:
        fmt.Printf("알 수 없는 오류: %v\n", err)
    }
    return
}
```

## 성능 고려사항

- 점수 계산은 메모리 할당을 최소화하도록 최적화됨
- 일반적인 사용에서는 100μs 이내에 계산 완료
- 대량 데이터 처리 시 배치 처리를 권장
