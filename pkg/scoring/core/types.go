package core

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
)

type ScoreResult struct {
	TotalScore     shared.ScoreValue
	WeightedScores map[metadata.MetadataType]shared.ScoreValue
	RawScores      map[metadata.MetadataType]shared.ScoreValue
	Weights        map[metadata.MetadataType]shared.Weight
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
	Weights     map[metadata.MetadataType]shared.Weight
	Scenario    ScoringScenario
}
