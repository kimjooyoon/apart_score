package context

import (
	"apart_score/pkg/metadata"
	"fmt"
)

type DefaultContextAdjuster struct{}

func NewDefaultContextAdjuster() *DefaultContextAdjuster {
	return &DefaultContextAdjuster{}
}
func (ca *DefaultContextAdjuster) AdjustWeights(baseWeights map[string]float64, context *ContextData) (map[string]float64, error) {
	if context == nil {
		return baseWeights, nil
	}
	adjusted := make(map[string]float64)
	for key, weight := range baseWeights {
		adjusted[key] = weight
	}
	if context.Regional != nil {
		ca.adjustForRegionalContext(adjusted, context.Regional)
	}
	if context.Temporal != nil {
		ca.adjustForTemporalContext(adjusted, context.Temporal)
	}
	return adjusted, nil
}
func (ca *DefaultContextAdjuster) adjustForRegionalContext(weights map[string]float64, regional *RegionalContext) {
	if regional.TransportLevel == TransportLevelHigh {
		if weight, exists := weights[metadata.DistanceToStation.String()]; exists {
			weights[metadata.DistanceToStation.String()] = weight * 1.2
		}
		if weight, exists := weights[metadata.TransportationAccess.String()]; exists {
			weights[metadata.TransportationAccess.String()] = weight * 1.1
		}
	}
	costMultiplier := regional.CostOfLiving / 100.0
	if weight, exists := weights[metadata.MaintenanceFee.String()]; exists {
		weights[metadata.MaintenanceFee.String()] = weight * costMultiplier
	}
	if regional.CrimeRate > 150 {
		if weight, exists := weights[metadata.CrimeRate.String()]; exists {
			weights[metadata.CrimeRate.String()] = weight * 1.3
		}
	}
	if regional.GreenSpaceRatio < 30 {
		if weight, exists := weights[metadata.GreenSpaceRatio.String()]; exists {
			weights[metadata.GreenSpaceRatio.String()] = weight * 1.2
		}
	}
	if regional.SchoolDistrictRank > 80 {
		if weight, exists := weights[metadata.SchoolDistrict.String()]; exists {
			weights[metadata.SchoolDistrict.String()] = weight * 1.15
		}
	}
}
func (ca *DefaultContextAdjuster) adjustForTemporalContext(weights map[string]float64, temporal *TemporalContext) {
	switch temporal.MarketCondition {
	case MarketConditionBoom:
		if weight, exists := weights[metadata.ConstructionYear.String()]; exists {
			weights[metadata.ConstructionYear.String()] = weight * 1.1
		}
		if weight, exists := weights[metadata.ConstructionCompany.String()]; exists {
			weights[metadata.ConstructionCompany.String()] = weight * 1.1
		}
	case MarketConditionSlump:
		if weight, exists := weights[metadata.MaintenanceFee.String()]; exists {
			weights[metadata.MaintenanceFee.String()] = weight * 1.2
		}
	case MarketConditionRecession:
		if weight, exists := weights[metadata.FloorLevel.String()]; exists {
			weights[metadata.FloorLevel.String()] = weight * 1.1
		}
		if weight, exists := weights[metadata.ElevatorPresence.String()]; exists {
			weights[metadata.ElevatorPresence.String()] = weight * 1.1
		}
	}
	switch temporal.Season {
	case SeasonWinter:
		if weight, exists := weights[metadata.HeatingSystem.String()]; exists {
			weights[metadata.HeatingSystem.String()] = weight * 1.3
		}
	case SeasonSummer:
		if weight, exists := weights[metadata.FloorLevel.String()]; exists {
			weights[metadata.FloorLevel.String()] = weight * 1.1
		}
	}
	if temporal.InflationRate > 3.0 {
		if weight, exists := weights[metadata.MaintenanceFee.String()]; exists {
			weights[metadata.MaintenanceFee.String()] = weight * (1.0 + (temporal.InflationRate-3.0)/10.0)
		}
	}
}
func (ca *DefaultContextAdjuster) GetAdjustmentFactors(context *ContextData) (map[string]float64, error) {
	baseWeights := ca.getDefaultWeights()
	adjustedWeights, err := ca.AdjustWeights(baseWeights, context)
	if err != nil {
		return nil, err
	}
	factors := make(map[string]float64)
	for key, baseWeight := range baseWeights {
		if adjustedWeight, exists := adjustedWeights[key]; exists && baseWeight > 0 {
			factors[key] = adjustedWeight / baseWeight
		}
	}
	return factors, nil
}
func (ca *DefaultContextAdjuster) getDefaultWeights() map[string]float64 {
	return map[string]float64{
		metadata.FloorLevel.String():           0.08,
		metadata.DistanceToStation.String():    0.15,
		metadata.ElevatorPresence.String():     0.07,
		metadata.ConstructionYear.String():     0.10,
		metadata.ConstructionCompany.String():  0.08,
		metadata.ApartmentSize.String():        0.08,
		metadata.NearbyAmenities.String():      0.10,
		metadata.TransportationAccess.String(): 0.12,
		metadata.SchoolDistrict.String():       0.08,
		metadata.CrimeRate.String():            0.06,
		metadata.GreenSpaceRatio.String():      0.04,
		metadata.Parking.String():              0.06,
		metadata.MaintenanceFee.String():       0.05,
		metadata.HeatingSystem.String():        0.03,
	}
}
func ValidateContext(context *ContextData) error {
	if context == nil {
		return fmt.Errorf("컨텍스트가 nil입니다")
	}
	if context.Location.Address == "" {
		return fmt.Errorf("주소가 비어있습니다")
	}
	if context.Regional != nil {
		if context.Regional.PopulationDensity < 0 {
			return fmt.Errorf("인구밀도가 음수입니다")
		}
		if context.Regional.CostOfLiving <= 0 {
			return fmt.Errorf("생활비 지수가 0 이하입니다")
		}
	}
	if context.Temporal != nil {
		if context.Temporal.EconomicIndex < -50 || context.Temporal.EconomicIndex > 50 {
			return fmt.Errorf("경제 지수가 비정상적입니다")
		}
		if context.Temporal.InflationRate < 0 || context.Temporal.InflationRate > 20 {
			return fmt.Errorf("물가 상승률이 비정상적입니다")
		}
	}
	return nil
}
