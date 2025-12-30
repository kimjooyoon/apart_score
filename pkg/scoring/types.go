package scoring

import "apart_score/pkg/metadata"

// ScoreValue는 개별 메타데이터의 점수 값 (0-100)
type ScoreValue float64

// Weight는 메타데이터의 가중치 (0.0-1.0)
type Weight float64

// ScoreResult는 최종 계산된 점수 결과
type ScoreResult struct {
	TotalScore       ScoreValue               // 총점 (0-100)
	WeightedScores   map[metadata.MetadataType]ScoreValue // 각 메타데이터의 가중치 적용 점수
	RawScores        map[metadata.MetadataType]ScoreValue // 각 메타데이터의 원본 점수
	Weights          map[metadata.MetadataType]Weight     // 적용된 가중치
	Method           ScoringMethod             // 사용된 계산 방법
	Scenario         ScoringScenario           // 적용된 시나리오
	Grade            Grade                     // 등급 (A, B, C, D)
	Percentile       float64                   // 백분위수 (0-100)
}

// ScoringMethod는 점수 계산 방법
type ScoringMethod string

const (
	MethodWeightedSum   ScoringMethod = "weighted_sum"   // 가중치 합계
	MethodGeometricMean ScoringMethod = "geometric_mean" // 기하 평균
	MethodMinMax        ScoringMethod = "min_max"        // 최소값 우선
	MethodHarmonicMean  ScoringMethod = "harmonic_mean" // 조화 평균
)

// ScoringScenario는 미리 정의된 가중치 시나리오
type ScoringScenario string

const (
	ScenarioBalanced      ScoringScenario = "balanced"       // 균형 잡힌
	ScenarioTransportation ScoringScenario = "transportation" // 교통 중심
	ScenarioEducation      ScoringScenario = "education"      // 교육 중심
	ScenarioCostEffective  ScoringScenario = "cost_effective" // 가성비 중심
	ScenarioFamilyFriendly ScoringScenario = "family_friendly" // 가족 친화적
	ScenarioInvestment     ScoringScenario = "investment"     // 투자 가치 중심
)

// Grade는 점수 등급
type Grade string

const (
	GradeA Grade = "A" // 우수 (90-100)
	GradeB Grade = "B" // 양호 (80-89)
	GradeC Grade = "C" // 보통 (70-79)
	GradeD Grade = "D" // 미흡 (60-69)
	GradeF Grade = "F" // 불량 (0-59)
)

// ScoringProfile은 사용자의 스코어링 설정
type ScoringProfile struct {
	Name        string                              // 프로필 이름
	Description string                              // 설명
	Method      ScoringMethod                       // 계산 방법
	Weights     map[metadata.MetadataType]Weight    // 커스텀 가중치
	Scenario    ScoringScenario                     // 기본 시나리오
}
