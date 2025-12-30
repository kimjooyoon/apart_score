package context

import (
	"fmt"
	"time"
)

func NewCompositeContextProvider() *CompositeContextProvider {
	return &CompositeContextProvider{
		regionalProviders: []RegionalContextProvider{NewDefaultRegionalProvider()},
		temporalProviders: []TemporalContextProvider{NewDefaultTemporalProvider()},
	}
}
func (cp *CompositeContextProvider) AddRegionalProvider(provider RegionalContextProvider) {
	cp.regionalProviders = append(cp.regionalProviders, provider)
}
func (cp *CompositeContextProvider) AddTemporalProvider(provider TemporalContextProvider) {
	cp.temporalProviders = append(cp.temporalProviders, provider)
}
func (cp *CompositeContextProvider) GetRegionalContext(location Location) (*RegionalContext, error) {
	var lastErr error
	for _, provider := range cp.regionalProviders {
		if provider.IsRegionSupported(location.RegionType) {
			context, err := provider.GetRegionalContext(location)
			if err == nil {
				return context, nil
			}
			lastErr = err
		}
	}
	if lastErr != nil {
		return nil, fmt.Errorf("지역 컨텍스트를 가져올 수 없습니다: %w", lastErr)
	}
	return nil, fmt.Errorf("지원하는 지역 공급자가 없습니다")
}
func (cp *CompositeContextProvider) GetTemporalContext(timestamp time.Time) (*TemporalContext, error) {
	var lastErr error
	for _, provider := range cp.temporalProviders {
		context, err := provider.GetTemporalContext(timestamp)
		if err == nil {
			return context, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, fmt.Errorf("시간 컨텍스트를 가져올 수 없습니다: %w", lastErr)
	}
	return nil, fmt.Errorf("지원하는 시간 공급자가 없습니다")
}
func (cp *CompositeContextProvider) GetContext(location Location, timestamp time.Time) (*ContextData, error) {
	regional, err := cp.GetRegionalContext(location)
	if err != nil {
		return nil, fmt.Errorf("지역 컨텍스트 조회 실패: %w", err)
	}
	temporal, err := cp.GetTemporalContext(timestamp)
	if err != nil {
		return nil, fmt.Errorf("시간 컨텍스트 조회 실패: %w", err)
	}
	return &ContextData{
		Location: location,
		Regional: regional,
		Temporal: temporal,
	}, nil
}
func (cp *CompositeContextProvider) IsContextValid(context interface{}) bool {
	switch ctx := context.(type) {
	case *RegionalContext:
		return ctx != nil && ctx.RegionType != ""
	case *TemporalContext:
		return ctx != nil && !ctx.Timestamp.IsZero()
	case *ContextData:
		return ctx != nil && ctx.Location.Address != "" &&
			cp.IsContextValid(ctx.Regional) && cp.IsContextValid(ctx.Temporal)
	default:
		return false
	}
}
func (cp *CompositeContextProvider) GetSupportedRegions() []RegionType {
	regions := make(map[RegionType]bool)
	for _, regionType := range []RegionType{
		RegionTypeSeoul, RegionTypeBusan, RegionTypeDaegu, RegionTypeIncheon,
		RegionTypeGwangju, RegionTypeDaejeon, RegionTypeUlsan, RegionTypeSejong,
		RegionTypeGyeonggi, RegionTypeGangwon, RegionTypeChungbuk, RegionTypeChungnam,
		RegionTypeJeonbuk, RegionTypeJeonnam, RegionTypeGyeongbuk, RegionTypeGyeongnam,
		RegionTypeJeju, RegionTypeOther,
	} {
		regions[regionType] = true
	}
	var result []RegionType
	for region := range regions {
		result = append(result, region)
	}
	return result
}
func (cp *CompositeContextProvider) GetRegionalProviders() []RegionalContextProvider {
	return cp.regionalProviders
}
func (cp *CompositeContextProvider) GetTemporalProviders() []TemporalContextProvider {
	return cp.temporalProviders
}
