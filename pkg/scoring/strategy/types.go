package strategy

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring/core"
	"apart_score/pkg/shared"
)

type StrategyType string

const (
	StrategyWeightedSum   StrategyType = "weighted_sum"
	StrategyGeometricMean StrategyType = "geometric_mean"
	StrategyMinMax        StrategyType = "min_max"
	StrategyHarmonicMean  StrategyType = "harmonic_mean"
)

type ApartmentData struct {
	ID       string                                `json:"id"`
	Name     string                                `json:"name"`
	Scores   map[metadata.MetadataType]shared.ScoreValue `json:"scores"`
	Location string                                `json:"location"`
}

type RankingResult struct {
	Apartment   ApartmentData     `json:"apartment"`
	Score       shared.ScoreValue `json:"score"`
	Rank        int               `json:"rank"`
	Percentile  float64           `json:"percentile"`
	Method      core.ScoringMethod `json:"method"`
	Weights     map[metadata.MetadataType]shared.Weight `json:"weights"`
}

type RankingsSummary struct {
	TotalApartments int             `json:"total_apartments"`
	Strategy        StrategyType    `json:"strategy"`
	TopRanked       []RankingResult `json:"top_ranked"`
	ScoreRange      struct {
		Min shared.ScoreValue `json:"min"`
		Max shared.ScoreValue `json:"max"`
		Avg shared.ScoreValue `json:"avg"`
	} `json:"score_range"`
}
