package providers

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"time"
)

type DataProvider interface {
	GetAdditionalScores(apartmentID string, metadataTypes []metadata.MetadataType) (map[metadata.MetadataType]scoring.ScoreValue, error)
	GetContextData(location string) (*CompositeData, error)
	IsDataAvailable(apartmentID string) bool
	GetDataFreshness(apartmentID string) (*DataFreshness, error)
	GetProviderInfo() ProviderInfo
}
type ProviderInfo struct {
	Name         string         `json:"name"`
	Type         DataSourceType `json:"type"`
	Version      string         `json:"version"`
	Description  string         `json:"description"`
	Capabilities []string       `json:"capabilities"`
}
type CompositeDataProvider struct {
	providers []DataProvider
	primary   DataProvider
	fallbacks []DataProvider
}
type RealEstateProvider interface {
	DataProvider
	GetRealEstateData(apartmentID string) (*RealEstateData, error)
	GetPriceHistory(apartmentID string, months int) ([]RealEstateData, error)
	GetMarketTrends(location string, period time.Duration) (*MarketTrends, error)
}
type EnvironmentalProvider interface {
	DataProvider
	GetEnvironmentalData(location string) (*EnvironmentalData, error)
	GetAirQualityHistory(location string, days int) ([]EnvironmentalData, error)
	GetEnvironmentalScore(location string) (float64, error)
}
type TrafficProvider interface {
	DataProvider
	GetTrafficData(location string) (*TrafficData, error)
	GetCommuteTime(origin, destination string) (time.Duration, error)
	GetPublicTransportScore(location string) (float64, error)
}
type InfrastructureProvider interface {
	DataProvider
	GetInfrastructureData(location string) (*InfrastructureData, error)
	FindNearestAmenity(location string, amenityType string) (*AmenityInfo, error)
	GetLocationScore(location string) (float64, error)
}
type AmenityInfo struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Distance float64 `json:"distance"`
	Rating   float64 `json:"rating"`
}
type MarketTrends struct {
	Location       string    `json:"location"`
	Period         string    `json:"period"`
	PriceChange    float64   `json:"price_change"`
	VolumeChange   float64   `json:"volume_change"`
	AveragePrice   float64   `json:"average_price"`
	TrendDirection string    `json:"trend_direction"`
	Confidence     float64   `json:"confidence"`
	LastUpdated    time.Time `json:"last_updated"`
}
