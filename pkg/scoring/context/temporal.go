package context

import (
	"time"
)

type DefaultTemporalProvider struct {
	marketData map[int]map[time.Month]*TemporalContext
}

func NewDefaultTemporalProvider() *DefaultTemporalProvider {
	return &DefaultTemporalProvider{
		marketData: getDefaultMarketData(),
	}
}
func (p *DefaultTemporalProvider) GetTemporalContext(timestamp time.Time) (*TemporalContext, error) {
	year := timestamp.Year()
	month := timestamp.Month()
	yearData, exists := p.marketData[year]
	if !exists {
		yearData = p.getLatestYearData()
	}
	context, exists := yearData[month]
	if !exists {
		context = p.getYearlyAverage(year)
	}
	result := *context
	result.Timestamp = timestamp
	result.Year = year
	result.Month = int(month)
	result.Season = p.getSeason(month)
	return &result, nil
}
func (p *DefaultTemporalProvider) IsTimeRangeSupported(start, end time.Time) bool {
	startYear := start.Year()
	endYear := end.Year()
	for year := startYear; year <= endYear; year++ {
		if _, exists := p.marketData[year]; !exists {
			return false
		}
	}
	return true
}
func (p *DefaultTemporalProvider) getSeason(month time.Month) Season {
	switch month {
	case time.December, time.January, time.February:
		return SeasonWinter
	case time.March, time.April, time.May:
		return SeasonSpring
	case time.June, time.July, time.August:
		return SeasonSummer
	case time.September, time.October, time.November:
		return SeasonFall
	default:
		return SeasonSpring
	}
}
func (p *DefaultTemporalProvider) getLatestYearData() map[time.Month]*TemporalContext {
	var latestYear int
	for year := range p.marketData {
		if year > latestYear {
			latestYear = year
		}
	}
	return p.marketData[latestYear]
}
func (p *DefaultTemporalProvider) getYearlyAverage(year int) *TemporalContext {
	yearData, exists := p.marketData[year]
	if !exists {
		return p.getDefaultTemporalContext()
	}
	var totalEconomicIndex, totalInflation, totalInterest float64
	count := 0
	for _, context := range yearData {
		totalEconomicIndex += context.EconomicIndex
		totalInflation += context.InflationRate
		totalInterest += context.InterestRate
		count++
	}
	if count == 0 {
		return p.getDefaultTemporalContext()
	}
	avgEconomicIndex := totalEconomicIndex / float64(count)
	avgInflation := totalInflation / float64(count)
	avgInterest := totalInterest / float64(count)
	marketCondition := p.inferMarketCondition(avgEconomicIndex, avgInflation)
	return &TemporalContext{
		MarketCondition: marketCondition,
		EconomicIndex:   avgEconomicIndex,
		InflationRate:   avgInflation,
		InterestRate:    avgInterest,
	}
}
func (p *DefaultTemporalProvider) inferMarketCondition(economicIndex, inflation float64) MarketCondition {
	if economicIndex > 5.0 && inflation < 2.0 {
		return MarketConditionBoom
	} else if economicIndex > 2.0 && economicIndex <= 5.0 {
		return MarketConditionNormal
	} else if economicIndex > 0 && economicIndex <= 2.0 {
		return MarketConditionSlump
	} else {
		return MarketConditionRecession
	}
}
func (p *DefaultTemporalProvider) getDefaultTemporalContext() *TemporalContext {
	return &TemporalContext{
		MarketCondition: MarketConditionNormal,
		EconomicIndex:   2.5,
		InflationRate:   2.5,
		InterestRate:    3.5,
	}
}
func getDefaultMarketData() map[int]map[time.Month]*TemporalContext {
	return map[int]map[time.Month]*TemporalContext{
		2023: {
			time.January: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   1.8,
				InflationRate:   5.2,
				InterestRate:    3.75,
			},
			time.February: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   2.1,
				InflationRate:   4.8,
				InterestRate:    3.75,
			},
			time.March: {
				MarketCondition: MarketConditionSlump,
				EconomicIndex:   0.8,
				InflationRate:   4.2,
				InterestRate:    3.75,
			},
			time.April: {
				MarketCondition: MarketConditionSlump,
				EconomicIndex:   1.2,
				InflationRate:   3.7,
				InterestRate:    3.75,
			},
			time.May: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   2.5,
				InflationRate:   3.3,
				InterestRate:    3.75,
			},
			time.June: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   2.8,
				InflationRate:   2.7,
				InterestRate:    3.5,
			},
			time.July: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   3.2,
				InflationRate:   2.4,
				InterestRate:    3.5,
			},
			time.August: {
				MarketCondition: MarketConditionBoom,
				EconomicIndex:   4.5,
				InflationRate:   2.3,
				InterestRate:    3.5,
			},
			time.September: {
				MarketCondition: MarketConditionBoom,
				EconomicIndex:   5.1,
				InflationRate:   3.0,
				InterestRate:    3.5,
			},
			time.October: {
				MarketCondition: MarketConditionBoom,
				EconomicIndex:   4.8,
				InflationRate:   3.2,
				InterestRate:    3.5,
			},
			time.November: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   3.5,
				InflationRate:   3.3,
				InterestRate:    3.5,
			},
			time.December: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   2.9,
				InflationRate:   3.6,
				InterestRate:    3.5,
			},
		},
		2024: {
			time.January: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   2.5,
				InflationRate:   2.8,
				InterestRate:    3.5,
			},
			time.February: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   2.3,
				InflationRate:   2.6,
				InterestRate:    3.5,
			},
			time.March: {
				MarketCondition: MarketConditionSlump,
				EconomicIndex:   1.5,
				InflationRate:   2.7,
				InterestRate:    3.5,
			},
			time.April: {
				MarketCondition: MarketConditionSlump,
				EconomicIndex:   1.2,
				InflationRate:   2.9,
				InterestRate:    3.5,
			},
			time.May: {
				MarketCondition: MarketConditionSlump,
				EconomicIndex:   0.8,
				InflationRate:   2.6,
				InterestRate:    3.5,
			},
			time.June: {
				MarketCondition: MarketConditionSlump,
				EconomicIndex:   0.5,
				InflationRate:   2.5,
				InterestRate:    3.5,
			},
			time.July: {
				MarketCondition: MarketConditionSlump,
				EconomicIndex:   0.3,
				InflationRate:   2.4,
				InterestRate:    3.0,
			},
			time.August: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   1.5,
				InflationRate:   2.3,
				InterestRate:    3.0,
			},
			time.September: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   2.0,
				InflationRate:   2.2,
				InterestRate:    3.0,
			},
			time.October: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   2.5,
				InflationRate:   2.1,
				InterestRate:    3.0,
			},
			time.November: {
				MarketCondition: MarketConditionNormal,
				EconomicIndex:   3.0,
				InflationRate:   2.0,
				InterestRate:    3.0,
			},
			time.December: {
				MarketCondition: MarketConditionBoom,
				EconomicIndex:   4.0,
				InflationRate:   2.2,
				InterestRate:    3.0,
			},
		},
	}
}
