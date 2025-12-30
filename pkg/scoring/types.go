package scoring

import "apart_score/pkg/metadata"

type ScoreValue float64
type Weight float64
type ScoreResult struct {
	TotalScore     ScoreValue
	WeightedScores map[metadata.MetadataType]ScoreValue
	RawScores      map[metadata.MetadataType]ScoreValue
	Weights        map[metadata.MetadataType]Weight
	Method         ScoringMethod
	Scenario       ScoringScenario
}
type ScoringMethod string

const (
	MethodWeightedSum   ScoringMethod = "weighted_sum"
	MethodGeometricMean ScoringMethod = "geometric_mean"
	MethodMinMax        ScoringMethod = "min_max"
	MethodHarmonicMean  ScoringMethod = "harmonic_mean"
)

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
	Method      ScoringMethod
	Weights     map[metadata.MetadataType]Weight
	Scenario    ScoringScenario
}
