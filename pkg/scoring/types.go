package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
	"time"
)

type ScoreResult struct {
	TotalScore     float64
	WeightedScores [14]float64
	RawScores      [14]shared.ScoreValue
	Weights        [14]shared.Weight
	Method         StrategyType
	Scenario       ScoringScenario
}

type ScoringScenario string

const (
	ScenarioBalanced       ScoringScenario = "balanced"
	ScenarioTransportation ScoringScenario = "transportation"
	ScenarioEducation      ScoringScenario = "education"
	ScenarioCostEffective  ScoringScenario = "cost_effective"
	ScenarioFamilyFriendly ScoringScenario = "family_friendly"
	ScenarioInvestment     ScoringScenario = "investment"
)

type ScoringProfile struct {
	Name        string
	Description string
	Method      StrategyType
	Weights     map[metadata.MetadataType]shared.Weight
	Scenario    ScoringScenario
}

// CalculationStep represents a step in the calculation pipeline
type CalculationStep struct {
	Name        string
	Description string
	Priority    int                    // 실행 우선순위 (낮을수록 먼저 실행)
	Condition   func(ScoreResult) bool // 실행 조건
	Calculator  func(map[metadata.MetadataType]shared.ScoreValue, map[metadata.MetadataType]shared.Weight) float64
}

// CalculationPipeline represents a customizable calculation pipeline
type CalculationPipeline struct {
	Name        string
	Description string
	Steps       []CalculationStep
}

// ScoreInterpretation provides multiple interpretations of a score.
type ScoreInterpretation struct {
	AbsoluteScore   float64   // 절대 점수 (0-100)
	Percentile      float64   // 백분위수 (상위 몇 %)
	Grade           string    // 등급 (A/B/C/D/F)
	Description     string    // 자연어 설명
	ConfidenceLevel float64   // 평가 신뢰도 (%)
	LastUpdated     time.Time // 데이터 최신성
}

// SubjectivityAnalysis provides transparency about the subjective elements in scoring.
type SubjectivityAnalysis struct {
	WeightCustomization    float64  // 가중치 설정의 주관성 정도 (0-100%)
	ScoreInterpretation    float64  // 점수 해석의 주관성 정도 (0-100%)
	ScenarioSelection      float64  // 시나리오 선택의 주관성 정도 (0-100%)
	DataQualityImpact      float64  // 데이터 품질이 결과에 미치는 영향 (0-100%)
	OverallObjectivity     float64  // 전체 평가의 객관성 비율 (0-100%)
	UserBiasFactors        []string // 사용자 편향 요인들
	MethodologicalLimits   []string // 방법론적 제한사항
	RecommendedAdjustments []string // 권장 조정사항
}

// ComprehensiveScore provides multi-dimensional score interpretation.
type ComprehensiveScore struct {
	Score          float64                // 기본 점수
	Interpretation ScoreInterpretation    // 다중 해석
	Subjectivity   SubjectivityAnalysis   // 주관성 분석
	RawComponents  map[string]interface{} // 원본 데이터
	Metadata       map[string]interface{} // 평가 메타데이터
}

// TransparencyDashboard provides complete transparency into the scoring process.
type TransparencyDashboard struct {
	// 점수 분석 섹션
	ScoreBreakdown    ScoreBreakdown    // 점수 구성 요소 세부 분석
	ScoreDistribution ScoreDistribution // 점수 분포 및 통계

	// 투명성 섹션
	AssumptionList     []string            // 평가 가정사항 목록
	MethodologyDetails MethodologyDetails  // 방법론 상세 설명
	UncertaintyFactors []UncertaintyFactor // 불확실성 요인들

	// 대안 분석 섹션
	AlternativeScenarios []AlternativeScenario // 대안 시나리오 비교
	SensitivityAnalysis  SensitivityAnalysis   // 민감도 분석

	// 품질 및 신뢰성 섹션
	DataQualityMetrics DataQualityMetrics // 데이터 품질 메트릭
	BiasIndicators     []BiasIndicator    // 잠재적 편향 지표

	// 사용자 가이드 섹션
	InterpretationGuide InterpretationGuide // 결과 해석 가이드
	RecommendedActions  []RecommendedAction // 권장 조치사항
}

// ScoreBreakdown provides detailed breakdown of how the score was calculated.
type ScoreBreakdown struct {
	TotalScore          float64                   // 최종 점수
	ComponentScores     map[string]ComponentScore // 각 구성 요소의 점수
	WeightContributions map[string]float64        // 각 요소의 가중치 기여도
	StrategyImpact      StrategyImpact            // 전략 선택이 결과에 미친 영향
}

// ComponentScore represents the score contribution of a single component.
type ComponentScore struct {
	RawValue        float64 // 원본 값
	NormalizedValue float64 // 정규화된 값 (0-100)
	Weight          float64 // 적용된 가중치
	Contribution    float64 // 최종 점수에 대한 기여도
	ImpactLevel     string  // 영향도 레벨 (High/Medium/Low)
}

// StrategyImpact shows how different strategies would affect the result.
type StrategyImpact struct {
	UsedStrategy       StrategyType             // 실제 사용된 전략
	AlternativeResults map[StrategyType]float64 // 다른 전략들의 결과
	BestAlternative    StrategyType             // 가장 좋은 대안 전략
	Reasoning          string                   // 전략 선택 근거
}

// ScoreDistribution provides statistical context for the score.
type ScoreDistribution struct {
	ScorePercentile    float64            // 백분위수
	ScoreRange         ScoreRange         // 점수 범위
	ComparativeContext string             // 비교 맥락 설명
	ConfidenceInterval ConfidenceInterval // 신뢰 구간
}

// ScoreRange represents the possible range of scores.
type ScoreRange struct {
	Minimum float64 // 가능한 최소 점수
	Maximum float64 // 가능한 최대 점수
	Average float64 // 평균 점수
	StdDev  float64 // 표준 편차
}

// ConfidenceInterval represents statistical confidence in the score.
type ConfidenceInterval struct {
	LowerBound float64 // 하한
	UpperBound float64 // 상한
	Confidence float64 // 신뢰도 (%)
}

// MethodologyDetails explains the scoring methodology in detail.
type MethodologyDetails struct {
	AlgorithmDescription string             // 알고리즘 설명
	DataSources          []DataSource       // 데이터 출처
	ValidationMethods    []ValidationMethod // 검증 방법
	Assumptions          []Assumption       // 가정사항
	Limitations          []Limitation       // 제한사항
}

// DataSource describes a source of scoring data.
type DataSource struct {
	Name        string    // 출처 이름
	Type        string    // 데이터 타입
	Reliability float64   // 신뢰도 (0-100%)
	LastUpdated time.Time // 최종 업데이트
	Coverage    string    // 적용 범위
}

// ValidationMethod describes how the scoring was validated.
type ValidationMethod struct {
	Method        string    // 검증 방법
	Accuracy      float64   // 정확도 (%)
	SampleSize    int       // 샘플 크기
	DatePerformed time.Time // 수행 일자
}

// Assumption represents an assumption made in the scoring process.
type Assumption struct {
	Description   string // 가정 설명
	Impact        string // 영향도 (High/Medium/Low)
	Justification string // 근거
}

// Limitation represents a limitation of the scoring methodology.
type Limitation struct {
	Description string // 제한사항 설명
	Severity    string // 심각도 (Critical/Important/Minor)
	Mitigation  string // 완화 방안
}

// UncertaintyFactor represents a source of uncertainty in the score.
type UncertaintyFactor struct {
	Factor      string  // 불확실성 요인
	Description string  // 설명
	Impact      float64 // 영향도 (%)
	Probability float64 // 발생 확률 (%)
	Mitigation  string  // 완화 방안
}

// AlternativeScenario represents an alternative scoring scenario.
type AlternativeScenario struct {
	ScenarioName   string  // 시나리오 이름
	Description    string  // 설명
	Score          float64 // 해당 시나리오의 점수
	Difference     float64 // 현재 점수와의 차이
	Reasoning      string  // 적용 근거
	Recommendation string  // 추천 여부
}

// SensitivityAnalysis shows how sensitive the score is to changes in inputs.
type SensitivityAnalysis struct {
	MostSensitiveFactors []SensitivityFactor // 가장 민감한 요소들
	StabilityIndex       float64             // 안정성 지수 (0-100%)
	VariationRange       ScoreRange          // 변동 범위
	RobustnessLevel      string              // 견고성 레벨
}

// SensitivityFactor represents a factor that significantly affects the score.
type SensitivityFactor struct {
	FactorName      string  // 요소 이름
	CurrentValue    float64 // 현재 값
	Sensitivity     float64 // 민감도 (%)
	ImpactDirection string  // 영향 방향 (Positive/Negative)
}

// DataQualityMetrics provides metrics on the quality of input data.
type DataQualityMetrics struct {
	Completeness   float64        // 완전성 (%)
	Accuracy       float64        // 정확성 (%)
	Timeliness     float64        // 적시성 (%)
	Consistency    float64        // 일관성 (%)
	OverallQuality float64        // 종합 품질 (%)
	QualityIssues  []QualityIssue // 품질 문제들
}

// QualityIssue represents a data quality issue.
type QualityIssue struct {
	Issue        string // 문제 설명
	Severity     string // 심각도
	AffectedData string // 영향받는 데이터
	Resolution   string // 해결 방안
}

// BiasIndicator represents a potential bias in the scoring.
type BiasIndicator struct {
	BiasType        string  // 편향 타입
	Description     string  // 설명
	Severity        float64 // 심각도 (%)
	DetectionMethod string  // 발견 방법
	Mitigation      string  // 완화 방안
}

// InterpretationGuide provides guidance on how to interpret the score.
type InterpretationGuide struct {
	ScoreRange           ScoreRange           // 점수 범위 설명
	InterpretationRules  []InterpretationRule // 해석 규칙들
	CommonMisconceptions []string             // 흔한 오해들
	BestPractices        []string             // 모범 사례들
}

// InterpretationRule provides a rule for interpreting scores.
type InterpretationRule struct {
	ScoreRange   ScoreRange // 적용 점수 범위
	Meaning      string     // 의미
	Implications []string   // 함의사항
	Actions      []string   // 권장 행동
}

// RecommendedAction represents a recommended action based on the score.
type RecommendedAction struct {
	Action         string // 권장 행동
	Priority       string // 우선순위 (High/Medium/Low)
	Reasoning      string // 근거
	ExpectedImpact string // 예상 영향
	Timeframe      string // 실행 기간
}
