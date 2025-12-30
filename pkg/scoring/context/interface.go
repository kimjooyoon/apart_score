package context

import (
	"time"
)

type ContextProvider interface {
	GetRegionalContext(location Location) (*RegionalContext, error)
	GetTemporalContext(timestamp time.Time) (*TemporalContext, error)
	GetContext(location Location, timestamp time.Time) (*ContextData, error)
	IsContextValid(context interface{}) bool
	GetSupportedRegions() []RegionType
}
type ContextAdjuster interface {
	AdjustWeights(baseWeights map[string]float64, context *ContextData) (map[string]float64, error)
	GetAdjustmentFactors(context *ContextData) (map[string]float64, error)
}
type CompositeContextProvider struct {
	regionalProviders []RegionalContextProvider
	temporalProviders []TemporalContextProvider
}
type RegionalContextProvider interface {
	GetRegionalContext(location Location) (*RegionalContext, error)
	IsRegionSupported(regionType RegionType) bool
}
type TemporalContextProvider interface {
	GetTemporalContext(timestamp time.Time) (*TemporalContext, error)
	IsTimeRangeSupported(start, end time.Time) bool
}
