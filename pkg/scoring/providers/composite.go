package providers

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"fmt"
)

func NewCompositeDataProvider(primary DataProvider, fallbacks ...DataProvider) *CompositeDataProvider {
	return &CompositeDataProvider{
		providers: append([]DataProvider{primary}, fallbacks...),
		primary:   primary,
		fallbacks: fallbacks,
	}
}
func (cp *CompositeDataProvider) AddProvider(provider DataProvider) {
	cp.providers = append(cp.providers, provider)
	cp.fallbacks = append(cp.fallbacks, provider)
}
func (cp *CompositeDataProvider) GetAdditionalScores(apartmentID string, metadataTypes []metadata.MetadataType) (map[metadata.MetadataType]scoring.ScoreValue, error) {
	var lastErr error
	results := make(map[metadata.MetadataType][]scoring.ScoreValue)
	for _, provider := range cp.providers {
		if scores, err := provider.GetAdditionalScores(apartmentID, metadataTypes); err == nil {
			for mt, score := range scores {
				results[mt] = append(results[mt], score)
			}
		} else {
			lastErr = err
		}
	}
	if len(results) == 0 {
		if lastErr != nil {
			return nil, fmt.Errorf("추가 점수 조회 실패: %w", lastErr)
		}
		return nil, fmt.Errorf("사용 가능한 데이터 공급자가 없습니다")
	}
	finalScores := make(map[metadata.MetadataType]scoring.ScoreValue)
	for mt, scores := range results {
		if len(scores) > 0 {
			sum := 0.0
			for _, score := range scores {
				sum += float64(score)
			}
			finalScores[mt] = scoring.ScoreValue(sum / float64(len(scores)))
		}
	}
	return finalScores, nil
}
func (cp *CompositeDataProvider) GetContextData(location string) (*CompositeData, error) {
	var lastErr error
	var bestData *CompositeData
	var bestReliability float64
	for _, provider := range cp.providers {
		if data, err := provider.GetContextData(location); err == nil {
			if data.ReliabilityScore > bestReliability {
				bestData = data
				bestReliability = data.ReliabilityScore
			}
		} else {
			lastErr = err
		}
	}
	if bestData != nil {
		bestData.DataSources = cp.getDataSources()
		return bestData, nil
	}
	if lastErr != nil {
		return nil, fmt.Errorf("컨텍스트 데이터 조회 실패: %w", lastErr)
	}
	return nil, fmt.Errorf("사용 가능한 데이터 공급자가 없습니다")
}
func (cp *CompositeDataProvider) IsDataAvailable(apartmentID string) bool {
	for _, provider := range cp.providers {
		if provider.IsDataAvailable(apartmentID) {
			return true
		}
	}
	return false
}
func (cp *CompositeDataProvider) GetDataFreshness(apartmentID string) (*DataFreshness, error) {
	var bestFreshness *DataFreshness
	for _, provider := range cp.providers {
		if freshness, err := provider.GetDataFreshness(apartmentID); err == nil {
			if bestFreshness == nil || freshness.LastUpdated.After(bestFreshness.LastUpdated) {
				bestFreshness = freshness
			}
		}
	}
	if bestFreshness != nil {
		return bestFreshness, nil
	}
	return nil, fmt.Errorf("데이터 신선도 정보를 가져올 수 없습니다")
}
func (cp *CompositeDataProvider) GetProviderInfo() ProviderInfo {
	return ProviderInfo{
		Name:        "CompositeDataProvider",
		Type:        "composite",
		Version:     "1.0.0",
		Description: "여러 데이터 공급자를 결합한 복합 공급자",
		Capabilities: []string{
			"additional_scores",
			"context_data",
			"load_balancing",
			"failover",
		},
	}
}
func (cp *CompositeDataProvider) getDataSources() []DataSourceType {
	sources := make(map[DataSourceType]bool)
	for _, provider := range cp.providers {
		info := provider.GetProviderInfo()
		sources[info.Type] = true
	}
	var result []DataSourceType
	for source := range sources {
		result = append(result, source)
	}
	return result
}
func (cp *CompositeDataProvider) GetProviders() []DataProvider {
	return cp.providers
}
func (cp *CompositeDataProvider) GetPrimaryProvider() DataProvider {
	return cp.primary
}
func (cp *CompositeDataProvider) GetFallbackProviders() []DataProvider {
	return cp.fallbacks
}
