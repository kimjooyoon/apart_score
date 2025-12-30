package context

import (
	"fmt"
	"strings"
)

type DefaultRegionalProvider struct {
	regionData map[RegionType]*RegionalContext
}

func NewDefaultRegionalProvider() *DefaultRegionalProvider {
	return &DefaultRegionalProvider{
		regionData: getDefaultRegionData(),
	}
}
func (p *DefaultRegionalProvider) GetRegionalContext(location Location) (*RegionalContext, error) {
	regionType := p.identifyRegionType(location)
	context, exists := p.regionData[regionType]
	if !exists {
		return nil, fmt.Errorf("지원하지 않는 지역 유형: %s", regionType)
	}
	adjustedContext := *context
	adjustedContext.CostOfLiving += p.getLocationVariance(location)
	adjustedContext.PopulationDensity += p.getDensityVariance(location)
	return &adjustedContext, nil
}
func (p *DefaultRegionalProvider) IsRegionSupported(regionType RegionType) bool {
	_, exists := p.regionData[regionType]
	return exists
}
func (p *DefaultRegionalProvider) identifyRegionType(location Location) RegionType {
	region := strings.ToLower(location.Region)
	if location.RegionType != "" {
		return location.RegionType
	}
	switch {
	case strings.Contains(region, "서울"):
		return RegionTypeSeoul
	case strings.Contains(region, "부산"):
		return RegionTypeBusan
	case strings.Contains(region, "대구"):
		return RegionTypeDaegu
	case strings.Contains(region, "인천"):
		return RegionTypeIncheon
	case strings.Contains(region, "광주"):
		return RegionTypeGwangju
	case strings.Contains(region, "대전"):
		return RegionTypeDaejeon
	case strings.Contains(region, "울산"):
		return RegionTypeUlsan
	case strings.Contains(region, "세종"):
		return RegionTypeSejong
	case strings.Contains(region, "경기"):
		return RegionTypeGyeonggi
	case strings.Contains(region, "강원"):
		return RegionTypeGangwon
	case strings.Contains(region, "충청북"):
		return RegionTypeChungbuk
	case strings.Contains(region, "충청남"):
		return RegionTypeChungnam
	case strings.Contains(region, "전라북"):
		return RegionTypeJeonbuk
	case strings.Contains(region, "전라남"):
		return RegionTypeJeonnam
	case strings.Contains(region, "경상북"):
		return RegionTypeGyeongbuk
	case strings.Contains(region, "경상남"):
		return RegionTypeGyeongnam
	case strings.Contains(region, "제주"):
		return RegionTypeJeju
	default:
		return RegionTypeOther
	}
}
func (p *DefaultRegionalProvider) getLocationVariance(location Location) float64 {
	hash := 0
	for _, char := range location.SubRegion {
		hash += int(char)
	}
	return float64(hash%20 - 10)
}
func (p *DefaultRegionalProvider) getDensityVariance(location Location) float64 {
	if strings.Contains(location.SubRegion, "강남") || strings.Contains(location.SubRegion, "서초") {
		return 500
	}
	if strings.Contains(location.SubRegion, "도심") || strings.Contains(location.SubRegion, "중심") {
		return 200
	}
	return 0
}
func getDefaultRegionData() map[RegionType]*RegionalContext {
	return map[RegionType]*RegionalContext{
		RegionTypeSeoul: {
			RegionType:         RegionTypeSeoul,
			PopulationDensity:  16154,
			TransportLevel:     TransportLevelHigh,
			CostOfLiving:       100.0,
			AvgApartmentPrice:  10.5,
			SchoolDistrictRank: 85,
			CrimeRate:          120.5,
			GreenSpaceRatio:    25.3,
		},
		RegionTypeBusan: {
			RegionType:         RegionTypeBusan,
			PopulationDensity:  4447,
			TransportLevel:     TransportLevelMedium,
			CostOfLiving:       85.0,
			AvgApartmentPrice:  4.2,
			SchoolDistrictRank: 65,
			CrimeRate:          145.2,
			GreenSpaceRatio:    35.1,
		},
		RegionTypeDaegu: {
			RegionType:         RegionTypeDaegu,
			PopulationDensity:  2758,
			TransportLevel:     TransportLevelMedium,
			CostOfLiving:       78.0,
			AvgApartmentPrice:  3.1,
			SchoolDistrictRank: 60,
			CrimeRate:          135.8,
			GreenSpaceRatio:    38.2,
		},
		RegionTypeIncheon: {
			RegionType:         RegionTypeIncheon,
			PopulationDensity:  2768,
			TransportLevel:     TransportLevelMedium,
			CostOfLiving:       82.0,
			AvgApartmentPrice:  3.8,
			SchoolDistrictRank: 55,
			CrimeRate:          138.9,
			GreenSpaceRatio:    32.4,
		},
		RegionTypeGwangju: {
			RegionType:         RegionTypeGwangju,
			PopulationDensity:  2972,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       75.0,
			AvgApartmentPrice:  2.9,
			SchoolDistrictRank: 58,
			CrimeRate:          142.3,
			GreenSpaceRatio:    40.5,
		},
		RegionTypeDaejeon: {
			RegionType:         RegionTypeDaejeon,
			PopulationDensity:  2860,
			TransportLevel:     TransportLevelMedium,
			CostOfLiving:       76.0,
			AvgApartmentPrice:  3.2,
			SchoolDistrictRank: 62,
			CrimeRate:          133.7,
			GreenSpaceRatio:    37.8,
		},
		RegionTypeUlsan: {
			RegionType:         RegionTypeUlsan,
			PopulationDensity:  1076,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       80.0,
			AvgApartmentPrice:  3.5,
			SchoolDistrictRank: 55,
			CrimeRate:          125.4,
			GreenSpaceRatio:    45.2,
		},
		RegionTypeSejong: {
			RegionType:         RegionTypeSejong,
			PopulationDensity:  211,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       88.0,
			AvgApartmentPrice:  3.8,
			SchoolDistrictRank: 70,
			CrimeRate:          95.3,
			GreenSpaceRatio:    55.8,
		},
		RegionTypeGyeonggi: {
			RegionType:         RegionTypeGyeonggi,
			PopulationDensity:  1338,
			TransportLevel:     TransportLevelHigh,
			CostOfLiving:       90.0,
			AvgApartmentPrice:  5.2,
			SchoolDistrictRank: 70,
			CrimeRate:          115.6,
			GreenSpaceRatio:    30.2,
		},
		RegionTypeGangwon: {
			RegionType:         RegionTypeGangwon,
			PopulationDensity:  90,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       82.0,
			AvgApartmentPrice:  2.8,
			SchoolDistrictRank: 50,
			CrimeRate:          108.9,
			GreenSpaceRatio:    75.3,
		},
		RegionTypeChungbuk: {
			RegionType:         RegionTypeChungbuk,
			PopulationDensity:  220,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       78.0,
			AvgApartmentPrice:  2.5,
			SchoolDistrictRank: 52,
			CrimeRate:          112.4,
			GreenSpaceRatio:    68.9,
		},
		RegionTypeChungnam: {
			RegionType:         RegionTypeChungnam,
			PopulationDensity:  256,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       79.0,
			AvgApartmentPrice:  2.6,
			SchoolDistrictRank: 53,
			CrimeRate:          114.7,
			GreenSpaceRatio:    65.4,
		},
		RegionTypeJeonbuk: {
			RegionType:         RegionTypeJeonbuk,
			PopulationDensity:  221,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       74.0,
			AvgApartmentPrice:  2.3,
			SchoolDistrictRank: 48,
			CrimeRate:          118.2,
			GreenSpaceRatio:    62.7,
		},
		RegionTypeJeonnam: {
			RegionType:         RegionTypeJeonnam,
			PopulationDensity:  145,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       73.0,
			AvgApartmentPrice:  2.1,
			SchoolDistrictRank: 46,
			CrimeRate:          105.8,
			GreenSpaceRatio:    70.1,
		},
		RegionTypeGyeongbuk: {
			RegionType:         RegionTypeGyeongbuk,
			PopulationDensity:  141,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       76.0,
			AvgApartmentPrice:  2.4,
			SchoolDistrictRank: 49,
			CrimeRate:          109.3,
			GreenSpaceRatio:    72.5,
		},
		RegionTypeGyeongnam: {
			RegionType:         RegionTypeGyeongnam,
			PopulationDensity:  315,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       77.0,
			AvgApartmentPrice:  2.7,
			SchoolDistrictRank: 51,
			CrimeRate:          111.6,
			GreenSpaceRatio:    58.9,
		},
		RegionTypeJeju: {
			RegionType:         RegionTypeJeju,
			PopulationDensity:  362,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       95.0,
			AvgApartmentPrice:  4.5,
			SchoolDistrictRank: 55,
			CrimeRate:          98.7,
			GreenSpaceRatio:    68.4,
		},
		RegionTypeOther: {
			RegionType:         RegionTypeOther,
			PopulationDensity:  200,
			TransportLevel:     TransportLevelLow,
			CostOfLiving:       80.0,
			AvgApartmentPrice:  2.8,
			SchoolDistrictRank: 50,
			CrimeRate:          120.0,
			GreenSpaceRatio:    40.0,
		},
	}
}
