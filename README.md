# 아파트 스코어링 시스템 (Apart Score)

[![Go Version](https://img.shields.io/badge/go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)

**외부에서 수집된 아파트 정보를 기반으로 DDD(Domain-Driven Design) 기반의 스코어링 시스템입니다.**
사용자가 설정한 메타데이터 기준에 따라 아파트 점수를 계산하고 평가합니다.

## 🚀 빠른 시작

### 설치 및 실행

```bash
# 프로젝트 클론 (또는 다운로드)
git clone <repository-url>
cd apart_score

# 빌드 및 실행
go build -o apart_score ./cmd
./apart_score
```

### 기본 사용법

```go
package main

import (
    "fmt"
    "apart_score/pkg/metadata"
    "apart_score/pkg/scoring"
)

func main() {
    // 아파트 점수 데이터
    scores := map[metadata.MetadataType]scoring.ScoreValue{
        metadata.FloorLevel:         85.0,  // 층수 점수
        metadata.DistanceToStation:  90.0,  // 역까지 거리
        metadata.ElevatorPresence:   100.0, // 엘리베이터
        // ... 다른 메타데이터들
    }

    // 균형 잡힌 시나리오로 점수 계산
    result, err := scoring.QuickScore(scores, scoring.ScenarioBalanced)
    if err != nil {
        panic(err)
    }

    fmt.Printf("총점: %.1f점 (등급: %s)\n", result.TotalScore, result.Grade)
    fmt.Println(scoring.FormatScoreResult(result))
}
```

## 📊 메타데이터 스코어링

아파트의 점수는 14개 메타데이터 요소들을 기반으로 계산됩니다:

### 메타데이터 요소들

| 요소 | 영문명 | 설명 | 가중치 범위 |
|-----|--------|------|------------|
| 층수 | Floor Level | 중간층에 가까울수록 높은 점수 | 7.2% |
| 역까지 거리 | Distance to Station | 가까울수록 높은 점수 | 13.5% |
| 엘리베이터 유무 | Elevator Presence | 있으면 높은 점수 | 6.3% |
| 건축년도 | Construction Year | 최신 건물일수록 높은 점수 | 9.0% |
| 건설회사 | Construction Company | 신뢰할 수 있는 회사일수록 높은 점수 | 7.2% |
| 아파트 크기 | Apartment Size | 적절한 크기일수록 높은 점수 | 7.2% |
| 주변 편의시설 | Nearby Amenities | 편의시설이 많을수록 높은 점수 | 9.0% |
| 교통 접근성 | Transportation Access | 대중교통 접근성이 좋을수록 높은 점수 | 10.8% |
| 학군 | School District | 좋은 학군일수록 높은 점수 | 7.2% |
| 범죄율 | Crime Rate | 낮을수록 높은 점수 | 5.4% |
| 녹지율 | Green Space Ratio | 높을수록 높은 점수 | 3.6% |
| 주차장 | Parking | 주차 공간이 충분할수록 높은 점수 | 5.4% |
| 관리비 | Maintenance Fee | 적절한 수준일수록 높은 점수 | 4.5% |
| 난방 방식 | Heating System | 효율적인 난방 방식일수록 높은 점수 | 2.7% |

### 스코어링 예시

```go
// 용인시 A아파트 평가
scores := map[metadata.MetadataType]scoring.ScoreValue{
    metadata.FloorLevel:         85.0,  // 5층 (중간층) - 양호
    metadata.DistanceToStation:  95.0,  // 역까지 5분 - 우수
    metadata.ElevatorPresence:   100.0, // 엘리베이터 있음 - 우수
    metadata.ConstructionYear:   90.0,  // 2020년 건축 - 우수
    // ... 다른 요소들
}

// 결과: 총점 82.9점 (등급: B)
```

## 🏗️ 아키텍처

### 전체 구조

```
apart_score/
├── pkg/
│   ├── metadata/           # 메타데이터 정의 및 관리
│   └── scoring/            # 스코어링 엔진
│       ├── core/           # 기본 엔진 (가중치 합계)
│       ├── context/        # 지역/시간 컨텍스트 (확장)
│       ├── providers/      # 외부 데이터 공급 (확장)
│       ├── strategies/     # 계산 전략 (확장)
│       ├── relative/       # 상대적 평가 (확장)
│       └── personalization/ # 개인화 (확장)
├── cmd/
│   └── main.go             # 데모 애플리케이션
├── SCORING_MODULE_SPEC.md  # 확장 설계 문서
├── Makefile                # 빌드 자동화
├── remove_comments.sh      # 주석 제거 스크립트
└── README.md
```

### 아키텍처 원칙

- **DDD 기반**: 도메인 중심 설계
- **모듈화**: 각 기능이 독립적으로 확장 가능
- **인터페이스 중심**: 구현 세부사항과 분리
- **확장성**: 새로운 계산 전략/데이터 공급자 쉽게 추가

## 📚 API 문서

### 메타데이터 패키지

```go
// 메타데이터 타입 조회
metadata.FloorLevel.String()           // "Floor Level"
metadata.FloorLevel.KoreanName()       // "층수"
metadata.FloorLevel.Index()            // 0

// 메타데이터 검색
mt, ok := metadata.GetByEnglishName("Floor Level")
mt, ok := metadata.GetByKoreanName("층수")
mt, ok := metadata.GetByIndex(0)

// 모든 메타데이터 조회
allTypes := metadata.AllMetadataTypes()
```

### 스코어링 패키지

```go
// 빠른 점수 계산
result, err := scoring.QuickScore(scores, scoring.ScenarioBalanced)

// 상세 점수 계산
manager, _ := scoring.NewScoringManager(profile)
result, err := manager.CalculateScore(scores)

// 시나리오별 가중치 조회
weights := scoring.GetScenarioWeights(scoring.ScenarioTransportation)

// 점수 분석
analysis := scoring.AnalyzeScore(result)
fmt.Println("강점:", analysis.Strengths)
fmt.Println("약점:", analysis.Weaknesses)
fmt.Println("개선 제안:", analysis.ImprovementTips)

// 시나리오 비교
for _, scenario := range scoring.GetAllScenarios() {
    result, _ := scoring.QuickScore(scores, scenario)
    fmt.Printf("%s: %.1f점\n", scoring.GetScenarioDescription(scenario), result.TotalScore)
}
```

## 🎯 사용 시나리오

### 1. 기본 평가
```go
// 간단한 아파트 평가
result, _ := scoring.QuickScore(scores, scoring.ScenarioBalanced)
fmt.Println(scoring.FormatScoreResult(result))
```

### 2. 시나리오 비교
```go
// 여러 시나리오로 비교 평가
scenarios := []scoring.ScoringScenario{
    scoring.ScenarioBalanced,
    scoring.ScoringTransportation,
    scoring.ScoringEducation,
}

for _, scenario := range scenarios {
    result, _ := scoring.QuickScore(scores, scenario)
    // 비교 분석
}
```

### 3. 상세 분석
```go
// 강점/약점 분석 및 개선 제안
analysis := scoring.AnalyzeScore(result)
for _, tip := range analysis.ImprovementTips {
    fmt.Println("개선 제안:", tip)
}
```

## 🔧 개발 도구

### 빌드 명령어

```bash
# 일반 빌드
make build

# 테스트 실행
make test

# 주석 제거 (배포용)
make clean-comments

# 백업 복원
make restore-backups

# 전체 정리
make clean-all
```

### 주석 제거 스크립트

```bash
# 단일 파일
./remove_comments.sh pkg/scoring/types.go

# Makefile 사용
make clean-comments-single FILE=pkg/scoring/types.go
```

## 📋 메타데이터 설계 원칙

1. **수정 불가**: 한 번 정의된 메타데이터는 수정되지 않음 (오직 추가만 가능)
2. **인덱스 기반**: iota를 사용한 순차적 인덱스 번호 부여
3. **다국어 지원**: 영문명과 한글명 모두 지원
4. **외부 참조**: 다른 패키지에서 쉽게 참조 가능

## 🚀 확장 계획

현재 시스템은 **핵심 기능에 집중**하며, 다음과 같은 확장을 지원하는 구조를 제공합니다:

- **상대적 평가**: 지역 내 순위, 백분위수 계산
- **지역별 가중치**: 수도권/지방에 따른 자동 조정
- **외부 데이터 연동**: 실거래가, 환경 데이터 공급자 인터페이스
- **개인화 추천**: 사용자 맞춤 가중치 및 추천

자세한 확장 설계는 [`SCORING_MODULE_SPEC.md`](./SCORING_MODULE_SPEC.md)를 참고하세요.

## 🤝 기여하기

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참고하세요.

## 📞 연락처

프로젝트 관리자 - [Your Name](mailto:your.email@example.com)

프로젝트 링크: [https://github.com/your-username/apart_score](https://github.com/your-username/apart_score)
