package analysis

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring/core"
	"apart_score/pkg/shared"
)

type ScoreAnalysis struct {
	Result          *core.ScoreResult
	Strengths       []metadata.MetadataType
	Weaknesses      []metadata.MetadataType
	TopFactors      []ScoreFactor
	ImprovementTips []string
	ComparisonScore float64
}
type ScoreFactor struct {
	Metadata metadata.MetadataType
	Score    shared.ScoreValue
	Weight   shared.Weight
	Impact   shared.ScoreValue
}
