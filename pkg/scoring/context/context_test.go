package context

import (
	"testing"
	"time"
)

func TestDefaultRegionalProvider_GetRegionalContext(t *testing.T) {
	provider := NewDefaultRegionalProvider()

	location := Location{
		Address:   "서울시 강남구 역삼동",
		Region:    "서울특별시",
		SubRegion: "강남구",
	}

	context, err := provider.GetRegionalContext(location)
	if err != nil {
		t.Fatalf("GetRegionalContext failed: %v", err)
	}

	if context.RegionType != RegionTypeSeoul {
		t.Errorf("Expected region type %v, got %v", RegionTypeSeoul, context.RegionType)
	}

	if context.PopulationDensity <= 0 {
		t.Errorf("Invalid population density: %v", context.PopulationDensity)
	}

	if context.CostOfLiving <= 0 {
		t.Errorf("Invalid cost of living: %v", context.CostOfLiving)
	}
}

func TestDefaultTemporalProvider_GetTemporalContext(t *testing.T) {
	provider := NewDefaultTemporalProvider()

	timestamp := time.Date(2024, time.June, 15, 12, 0, 0, 0, time.UTC)

	context, err := provider.GetTemporalContext(timestamp)
	if err != nil {
		t.Fatalf("GetTemporalContext failed: %v", err)
	}

	if context.Year != 2024 {
		t.Errorf("Expected year 2024, got %v", context.Year)
	}

	if context.Month != 6 {
		t.Errorf("Expected month 6, got %v", context.Month)
	}

	if context.Season != SeasonSummer {
		t.Errorf("Expected season %v, got %v", SeasonSummer, context.Season)
	}
}

func TestCompositeContextProvider_GetContext(t *testing.T) {
	provider := NewCompositeContextProvider()

	location := Location{
		Address:    "부산시 해운대구 센텀동",
		Region:     "부산광역시",
		RegionType: RegionTypeBusan,
	}

	timestamp := time.Date(2024, time.March, 20, 10, 30, 0, 0, time.UTC)

	context, err := provider.GetContext(location, timestamp)
	if err != nil {
		t.Fatalf("GetContext failed: %v", err)
	}

	if context.Regional.RegionType != RegionTypeBusan {
		t.Errorf("Expected regional type %v, got %v", RegionTypeBusan, context.Regional.RegionType)
	}

	if context.Temporal.Year != 2024 {
		t.Errorf("Expected temporal year 2024, got %v", context.Temporal.Year)
	}

	if context.Temporal.Season != SeasonSpring {
		t.Errorf("Expected season %v, got %v", SeasonSpring, context.Temporal.Season)
	}
}

func TestDefaultContextAdjuster_AdjustWeights(t *testing.T) {
	adjuster := NewDefaultContextAdjuster()

	baseWeights := map[string]float64{
		"Distance to Station": 0.15,
		"Crime Rate":          0.06,
		"Maintenance Fee":     0.05,
	}

	context := &ContextData{
		Regional: &RegionalContext{
			RegionType:     RegionTypeSeoul,
			TransportLevel: TransportLevelHigh,
			CrimeRate:      160.0, // 높은 범죄율
			CostOfLiving:   110.0, // 높은 생활비
		},
	}

	adjusted, err := adjuster.AdjustWeights(baseWeights, context)
	if err != nil {
		t.Fatalf("AdjustWeights failed: %v", err)
	}

	// 서울(고도화 교통)에서는 역까지 거리 가중치가 증가해야 함
	stationWeight := adjusted["Distance to Station"]
	if stationWeight <= baseWeights["Distance to Station"] {
		t.Errorf("Station weight should increase in Seoul, got %v (base: %v)",
			stationWeight, baseWeights["Distance to Station"])
	}

	// 높은 범죄율에서는 범죄율 가중치가 증가해야 함
	crimeWeight := adjusted["Crime Rate"]
	if crimeWeight <= baseWeights["Crime Rate"] {
		t.Errorf("Crime weight should increase with high crime rate, got %v (base: %v)",
			crimeWeight, baseWeights["Crime Rate"])
	}
}

func TestValidateContext(t *testing.T) {
	// 유효한 컨텍스트
	validContext := &ContextData{
		Location: Location{Address: "서울시 강남구"},
		Regional: &RegionalContext{
			PopulationDensity: 10000,
			CostOfLiving:      100,
		},
		Temporal: &TemporalContext{
			EconomicIndex: 2.5,
			InflationRate: 2.5,
		},
	}

	err := ValidateContext(validContext)
	if err != nil {
		t.Errorf("Valid context should pass validation: %v", err)
	}

	// 유효하지 않은 컨텍스트들
	invalidContexts := []*ContextData{
		nil,
		{Location: Location{}}, // 빈 주소
		{
			Location: Location{Address: "서울"},
			Regional: &RegionalContext{PopulationDensity: -1}, // 음수 인구밀도
		},
		{
			Location: Location{Address: "서울"},
			Temporal: &TemporalContext{EconomicIndex: 100}, // 비정상적 경제지표
		},
	}

	for i, invalidContext := range invalidContexts {
		err := ValidateContext(invalidContext)
		if err == nil {
			t.Errorf("Invalid context %d should fail validation", i)
		}
	}
}
