package strategies

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"apart_score/pkg/scoring/context"
)

type ScoringStrategy interface {
	Name() string
	Calculate(scores map[metadata.MetadataType]scoring.ScoreValue,
		weights map[metadata.MetadataType]scoring.Weight,
		ctx *context.ContextData) (*scoring.ScoreResult, error)
	ValidateInputs(scores map[metadata.MetadataType]scoring.ScoreValue,
		weights map[metadata.MetadataType]scoring.Weight) error
	GetRequiredContext() []context.ContextType
	Description() string
}
type StrategyFactory interface {
	CreateStrategy(strategyType scoring.ScoringMethod) (ScoringStrategy, error)
	GetAvailableStrategies() []scoring.ScoringMethod
}
type DefaultStrategyFactory struct{}

func NewDefaultStrategyFactory() *DefaultStrategyFactory {
	return &DefaultStrategyFactory{}
}
func (f *DefaultStrategyFactory) CreateStrategy(strategyType scoring.ScoringMethod) (ScoringStrategy, error) {
	switch strategyType {
	case scoring.MethodWeightedSum:
		return NewWeightedSumStrategy(), nil
	case scoring.MethodGeometricMean:
		return NewGeometricMeanStrategy(), nil
	case scoring.MethodMinMax:
		return NewMinMaxStrategy(), nil
	case scoring.MethodHarmonicMean:
		return NewHarmonicMeanStrategy(), nil
	default:
		return nil, &UnsupportedStrategyError{Strategy: strategyType}
	}
}
func (f *DefaultStrategyFactory) GetAvailableStrategies() []scoring.ScoringMethod {
	return []scoring.ScoringMethod{
		scoring.MethodWeightedSum,
		scoring.MethodGeometricMean,
		scoring.MethodMinMax,
		scoring.MethodHarmonicMean,
	}
}

type UnsupportedStrategyError struct {
	Strategy scoring.ScoringMethod
}

func (e *UnsupportedStrategyError) Error() string {
	return "지원하지 않는 전략: " + string(e.Strategy)
}
