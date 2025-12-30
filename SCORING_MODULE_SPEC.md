# 아파트 스코어링 모듈 확장 설계 문서

## 개요

현재 아파트 스코어링 시스템의 확장성을 향상시키기 위한 모듈 설계 문서입니다.
프로젝트 범위는 **실제 데이터 수집/외부 API 연동을 하지 않고**, 이러한 기능들을 지원할 수 있는 스코어링 모듈 구조를 제공하는 것입니다.

## 1. 현재 시스템 분석

### 1.1 현재 구조
```
pkg/scoring/
├── types.go           # 기본 타입 정의
├── engine.go          # 가중치 합계 엔진
├── scenarios.go       # 6가지 시나리오 프리셋
├── profiles.go        # 사용자 프로필 관리
├── analysis.go        # 점수 분석 기능
└── scoring_test.go    # 단위 테스트
```

### 1.2 현재 기능
- ✅ 가중치 합계 기반 점수 계산
- ✅ 6가지 시나리오 프리셋 (균형, 교통, 교육, 가성비, 가족, 투자)
- ✅ 점수 분석 (강점/약점 식별, 개선 제안)
- ✅ 등급 시스템 (A/B/C/D/F)
- ✅ 시나리오 비교 기능

### 1.3 한계점
- ❌ 상대적 평가 기능 부족 (지역 내 순위, 백분위수)
- ❌ 지역별 가중치 자동 조정 불가
- ❌ 시간 기반 점수 변동 반영 불가
- ❌ 외부 데이터 연동 인터페이스 부재
- ❌ 개인화 추천 기능 제한적

## 2. 확장 요구사항

### 2.1 기능 요구사항

#### REQ-001: 상대적 평가 시스템
- 같은 지역/가격대 내 아파트 비교 기능
- 백분위수 기반 순위 제공
- 그룹별 상대 점수 계산

#### REQ-002: 지역별 가중치 자동 조정
- 지역 특성에 따른 가중치 자동 적용 (서울 vs 지방)
- 인구밀도, 교통 상황 등 지역 요인 반영
- 동적 가중치 조정 인터페이스

#### REQ-003: 시간 기반 점수 조정
- 건축년도 기반 노후화 패널티
- 시장 상황 반영 (호황/불황)
- 계절적 요인 고려 (일조량, 날씨 등)

#### REQ-004: 외부 데이터 공급자 인터페이스
- 실거래가 데이터 연동 포인트
- 환경 데이터 (미세먼지, 소음) 공급 인터페이스
- 교통/인프라 데이터 공급 포인트

#### REQ-005: 개인화 추천 확장
- 사용자 행동 패턴 학습 인터페이스
- 협업 필터링을 위한 사용자 유사도 계산
- 실시간 피드백 기반 모델 업데이트 포인트

### 2.2 비기능 요구사항

#### NFR-001: 확장성
- 새로운 계산 전략 쉽게 추가 가능
- 새로운 데이터 공급자 쉽게 통합
- 설정 변경으로 동작 방식 변경 가능

#### NFR-002: 성능
- 점수 계산 시 기존 성능 유지 (100μs 이내)
- 메모리 사용량 증가 최소화
- 동시 요청 처리 가능

#### NFR-003: 안정성
- 인터페이스 변경 시 기존 코드 영향 최소화
- 잘못된 데이터 입력에 대한 방어 로직
- 상세한 오류 메시지 제공

## 3. 아키텍처 설계

### 3.1 전체 아키텍처

```
pkg/scoring/
├── core/              # 핵심 엔진 (기존 유지)
│   ├── types.go       # 기본 타입
│   ├── engine.go      # 기본 스코어러
│   ├── scenarios.go   # 시나리오 프리셋
│   └── analysis.go    # 분석 기능
├── context/           # 컨텍스트 제공자
│   ├── types.go       # 컨텍스트 타입
│   ├── regional.go    # 지역 컨텍스트
│   └── temporal.go    # 시간 컨텍스트
├── providers/         # 데이터 공급자
│   ├── interface.go   # 공급자 인터페이스
│   └── mock.go        # 목업 구현
├── strategies/        # 계산 전략
│   ├── interface.go   # 전략 인터페이스
│   ├── weighted_sum.go # 가중치 합계 (기본)
│   └── geometric_mean.go # 기하 평균
├── relative/          # 상대적 평가
│   ├── interface.go   # 상대 평가 인터페이스
│   ├── percentile.go  # 백분위수 계산
│   └── regional.go    # 지역별 비교
└── personalization/   # 개인화
    ├── interface.go   # 개인화 인터페이스
    └── profile.go     # 사용자 프로필
```

### 3.2 모듈 책임

#### Core 모듈
- **책임**: 기본 스코어링 로직 유지
- **변경**: 최소한의 변경으로 기존 호환성 유지
- **확장**: 새로운 인터페이스 지원

#### Context 모듈
- **책임**: 지역/시간 등의 컨텍스트 정보 제공
- **역할**: 동적 가중치 조정에 필요한 정보 공급

#### Providers 모듈
- **책임**: 외부 데이터 공급 인터페이스 정의
- **역할**: 실제 데이터 수집 로직과 분리

#### Strategies 모듈
- **책임**: 다양한 점수 계산 방법 제공
- **역할**: 계산 로직의 전략 패턴 적용

#### Relative 모듈
- **책임**: 상대적 평가 및 비교 기능
- **역할**: 그룹 내 순위 및 백분위수 계산

#### Personalization 모듈
- **책임**: 사용자별 맞춤 기능
- **역할**: 개인화 추천을 위한 기본 구조

## 4. API 인터페이스 설계

### 4.1 ContextProvider 인터페이스

```go
type ContextProvider interface {
    GetRegionalContext(location Location) (RegionalContext, error)
    GetTemporalContext(timestamp time.Time) (TemporalContext, error)
    IsContextValid(context interface{}) bool
}

type RegionalContext struct {
    RegionType      RegionType      // 수도권/광역시/기타
    PopulationDensity float64       // 인구밀도
    TransportLevel   TransportLevel // 교통 발전도
    CostOfLiving     float64        // 생활비 수준
}

type TemporalContext struct {
    MarketCondition MarketCondition // 호황/불황/정체
    Season          Season         // 계절
    EconomicIndex   float64        // 경제 지표
}
```

### 4.2 DataProvider 인터페이스

```go
type DataProvider interface {
    GetAdditionalScores(apartmentID string, metadataTypes []metadata.MetadataType) (map[metadata.MetadataType]ScoreValue, error)
    GetContextData(location Location) (ContextData, error)
    IsDataAvailable(apartmentID string) bool
    GetDataFreshness(apartmentID string) time.Duration
}

type ContextData struct {
    RealEstatePrices []PricePoint   // 실거래가 히스토리
    EnvironmentalData EnvData       // 환경 데이터
    InfrastructureData InfraData    // 인프라 데이터
}
```

### 4.3 ScoringStrategy 인터페이스

```go
type ScoringStrategy interface {
    Name() string
    Calculate(scores map[metadata.MetadataType]ScoreValue, weights map[metadata.MetadataType]Weight, context Context) (*ScoreResult, error)
    ValidateInputs(scores map[metadata.MetadataType]ScoreValue, weights map[metadata.MetadataType]Weight) error
    GetRequiredContext() []ContextType
}

type Context struct {
    Regional *RegionalContext
    Temporal *TemporalContext
    Additional map[string]interface{}
}
```

### 4.4 RelativeEvaluator 인터페이스

```go
type RelativeEvaluator interface {
    EvaluateRelatively(apartment ApartmentScore, group []ApartmentScore) (*RelativeScore, error)
    CalculatePercentile(score ScoreValue, groupScores []ScoreValue) float64
    FindSimilarApartments(target ApartmentScore, candidates []ApartmentScore, criteria SimilarityCriteria) []ApartmentScore
}

type RelativeScore struct {
    AbsoluteScore  ScoreValue
    PercentileRank float64
    RegionalRank   int
    GroupRank      int
    ScoreDistribution ScoreDistribution
}
```

### 4.5 PersonalizationEngine 인터페이스

```go
type PersonalizationEngine interface {
    LearnFromInteraction(userID string, apartmentID string, interaction UserInteraction) error
    GetPersonalizedWeights(userID string, baseWeights map[metadata.MetadataType]Weight) (map[metadata.MetadataType]Weight, error)
    RecommendApartments(userID string, candidates []ApartmentScore, limit int) ([]Recommendation, error)
    GetUserProfile(userID string) (*UserProfile, error)
}

type UserInteraction struct {
    Action      InteractionType
    Duration    time.Duration
    ScoreViewed bool
    ContactMade bool
    Timestamp   time.Time
}
```

## 5. 구현 계획

### Phase 1: Context Providers (2주)
**목표**: 지역/시간 컨텍스트 제공자 구현
**작업**:
- RegionalContext, TemporalContext 구조체 정의
- 기본 지역 분류 로직 (수도권/광역시/기타)
- 시간 기반 조정 기본 로직
- 기존 엔진에 컨텍스트 통합 포인트 추가

### Phase 2: Data Providers 인터페이스 (2주)
**목표**: 외부 데이터 공급 인터페이스 구축
**작업**:
- DataProvider 인터페이스 정의
- MockProvider 구현 (테스트용)
- 데이터 검증 및 캐싱 로직
- 기존 엔진에 데이터 공급자 통합

### Phase 3: Scoring Strategies 확장 (3주)
**목표**: 다양한 계산 전략 지원
**작업**:
- ScoringStrategy 인터페이스 구현
- GeometricMeanStrategy 추가
- Context-aware 계산 로직
- 전략 선택 및 설정 기능

### Phase 4: Relative Evaluation (3주)
**목표**: 상대적 평가 시스템 구현
**작업**:
- RelativeEvaluator 인터페이스 구현
- 백분위수 계산 로직
- 지역별/가격대별 그룹화 기능
- 상대 점수 시각화 지원

### Phase 5: Personalization 기본 구조 (2주)
**목표**: 개인화 추천 기본 틀 구축
**작업**:
- PersonalizationEngine 인터페이스 정의
- 사용자 프로필 관리 기본 구조
- 상호작용 데이터 수집 인터페이스
- 기본 추천 로직 틀

### 5.1 위험 요소 및 완화 방안

| 위험 요소 | 영향도 | 완화 방안 |
|---------|--------|----------|
| 기존 코드 호환성 | 높음 | 철저한 인터페이스 설계 및 테스트 |
| 성능 저하 | 중간 | 성능 모니터링 및 최적화 |
| 복잡도 증가 | 높음 | 모듈화 및 문서화 강화 |

### 5.2 성공 지표

- ✅ 기존 기능 100% 호환 유지
- ✅ 새로운 기능 인터페이스 명확성
- ✅ 단위 테스트 커버리지 90% 이상
- ✅ 성능 기준 준수 (계산 시간 100μs 이내)
- ✅ 확장성 검증 (새로운 전략 쉽게 추가 가능)

## 6. 결론

이 설계 문서는 현재 아파트 스코어링 시스템을 **확장 가능하고 유연한 모듈 구조**로 발전시키는 청사진입니다.

**핵심 원칙**:
1. **기존 호환성 유지**: 현재 사용자에게 영향 최소화
2. **모듈화**: 각 기능이 독립적으로 확장 가능
3. **인터페이스 중심**: 구현 세부사항과 분리
4. **점진적 개선**: 단계적 구현으로 위험 최소화

이러한 접근을 통해 스코어링 모듈은 **실거래가 연동, 상대적 평가, 개인화 추천** 등 미래 확장 요구사항을 유연하게 수용할 수 있는 기반을 마련하게 됩니다.
