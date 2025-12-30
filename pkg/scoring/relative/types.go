package relative

import "apart_score/pkg/scoring"

type ApartmentScore struct {
	ID       string                 `json:"id"`
	Score    scoring.ScoreValue     `json:"score"`
	Location string                 `json:"location"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
type RelativeScore struct {
	ApartmentID       string             `json:"apartment_id"`
	AbsoluteScore     scoring.ScoreValue `json:"absolute_score"`
	PercentileRank    float64            `json:"percentile_rank"`
	RegionalRank      int                `json:"regional_rank"`
	GroupRank         int                `json:"group_rank"`
	ScoreDistribution ScoreDistribution  `json:"score_distribution"`
	Comparison        ScoreComparison    `json:"comparison"`
}
type ScoreDistribution struct {
	Mean   scoring.ScoreValue `json:"mean"`
	Median scoring.ScoreValue `json:"median"`
	StdDev scoring.ScoreValue `json:"std_dev"`
	Min    scoring.ScoreValue `json:"min"`
	Max    scoring.ScoreValue `json:"max"`
	Q1     scoring.ScoreValue `json:"q1"`
	Q3     scoring.ScoreValue `json:"q3"`
	Count  int                `json:"count"`
}
type ScoreComparison struct {
	BetterThanCount int     `json:"better_than_count"`
	WorseThanCount  int     `json:"worse_than_count"`
	SimilarCount    int     `json:"similar_count"`
	RankPercentile  float64 `json:"rank_percentile"`
}
type SimilarityCriteria struct {
	LocationWeight    float64 `json:"location_weight"`
	ScoreRange        float64 `json:"score_range"`
	PriceRange        float64 `json:"price_range"`
	MaxResults        int     `json:"max_results"`
	IncludeHigherRank bool    `json:"include_higher_rank"`
}
type EvaluationGroup struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Criteria    GroupCriteria    `json:"criteria"`
	Apartments  []ApartmentScore `json:"apartments"`
}
type GroupCriteria struct {
	LocationPattern string `json:"location_pattern"`
	ScoreRange      struct {
		Min scoring.ScoreValue `json:"min"`
		Max scoring.ScoreValue `json:"max"`
	} `json:"score_range"`
	PriceRange struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"price_range"`
	MaxGroupSize        int     `json:"max_group_size"`
	SimilarityThreshold float64 `json:"similarity_threshold"`
}
type EvaluationResult struct {
	TargetApartment   RelativeScore     `json:"target_apartment"`
	Group             EvaluationGroup   `json:"group"`
	SimilarApartments []ApartmentScore  `json:"similar_apartments"`
	Summary           EvaluationSummary `json:"summary"`
}
type EvaluationSummary struct {
	TotalEvaluated    int     `json:"total_evaluated"`
	GroupSize         int     `json:"group_size"`
	AverageScore      float64 `json:"average_score"`
	ScoreStdDeviation float64 `json:"score_std_deviation"`
	GeneratedAt       string  `json:"generated_at"`
}
