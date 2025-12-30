package providers

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/scoring"
	"fmt"
	"math/rand"
	"time"
)

type MockDataProvider struct {
	config      ProviderConfig
	mockData    map[string]*CompositeData
	lastUpdated time.Time
}

func NewMockDataProvider(config ProviderConfig) *MockDataProvider {
	return &MockDataProvider{
		config:      config,
		mockData:    generateMockData(),
		lastUpdated: time.Now(),
	}
}
func (p *MockDataProvider) GetAdditionalScores(apartmentID string, metadataTypes []metadata.MetadataType) (map[metadata.MetadataType]scoring.ScoreValue, error) {
	scores := make(map[metadata.MetadataType]scoring.ScoreValue)
	for _, mt := range metadataTypes {
		switch mt {
		case metadata.ConstructionCompany:
			scores[mt] = scoring.ScoreValue(75.0 + rand.Float64()*25.0)
		case metadata.NearbyAmenities:
			scores[mt] = scoring.ScoreValue(60.0 + rand.Float64()*35.0)
		case metadata.ConstructionYear:
			scores[mt] = scoring.ScoreValue(70.0 + rand.Float64()*25.0)
		case metadata.ApartmentSize:
			scores[mt] = scoring.ScoreValue(50.0 + rand.Float64()*40.0)
		default:
			scores[mt] = scoring.ScoreValue(60.0 + rand.Float64()*30.0)
		}
	}
	return scores, nil
}
func (p *MockDataProvider) GetContextData(location string) (*CompositeData, error) {
	if data, exists := p.mockData[location]; exists {
		return data, nil
	}
	return generateRandomCompositeData(location), nil
}
func (p *MockDataProvider) IsDataAvailable(apartmentID string) bool {
	return true
}
func (p *MockDataProvider) GetDataFreshness(apartmentID string) (*DataFreshness, error) {
	return &DataFreshness{
		LastUpdated: p.lastUpdated,
		Age:         time.Since(p.lastUpdated),
		IsStale:     time.Since(p.lastUpdated) > time.Hour*24,
		MaxAge:      time.Hour * 24 * 7,
	}, nil
}
func (p *MockDataProvider) GetProviderInfo() ProviderInfo {
	return ProviderInfo{
		Name:        "MockDataProvider",
		Type:        DataSourceMock,
		Version:     "1.0.0",
		Description: "테스트용 목업 데이터 공급자",
		Capabilities: []string{
			"additional_scores",
			"context_data",
			"real_estate",
			"environmental",
			"traffic",
			"infrastructure",
		},
	}
}
func generateMockData() map[string]*CompositeData {
	data := make(map[string]*CompositeData)
	locations := []string{
		"서울시 강남구 역삼동",
		"서울시 서초구 반포동",
		"부산시 해운대구 센텀동",
		"대구시 중구 동성로",
		"인천시 남동구 구월동",
		"광주시 서구 치평동",
		"대전시 유성구 봉명동",
		"울산시 남구 삼산동",
		"경기도 성남시 분당구",
		"경기도 고양시 일산동구",
	}
	for _, location := range locations {
		data[location] = generateRandomCompositeData(location)
	}
	return data
}
func generateRandomCompositeData(location string) *CompositeData {
	return &CompositeData{
		RealEstate: &RealEstateData{
			PropertyID:       fmt.Sprintf("mock-%d", rand.Intn(10000)),
			Address:          location,
			TradeType:        "매매",
			Price:            30000 + rand.Float64()*50000,
			Area:             80 + rand.Float64()*50,
			Floor:            rand.Intn(20) + 1,
			TradeDate:        time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
			BuildYear:        2000 + rand.Intn(24),
			DataSource:       "mock_api",
			ReliabilityScore: 0.8 + rand.Float64()*0.2,
		},
		Environmental: &EnvironmentalData{
			Location:      location,
			AirQuality:    20 + rand.Float64()*30,
			NoiseLevel:    40 + rand.Float64()*20,
			SunlightHours: 4 + rand.Float64()*4,
			Temperature:   10 + rand.Float64()*15,
			Humidity:      50 + rand.Float64()*30,
			GreenSpace:    20 + rand.Float64()*40,
			MeasuredAt:    time.Now(),
			DataSource:    "mock_env_api",
		},
		Traffic: &TrafficData{
			Location:          location,
			PublicTransport:   60 + rand.Float64()*35,
			TrafficCongestion: 30 + rand.Float64()*40,
			Walkability:       50 + rand.Float64()*40,
			BikeFriendly:      40 + rand.Float64()*45,
			ParkingEase:       45 + rand.Float64()*40,
			MeasuredAt:        time.Now(),
			DataSource:        "mock_traffic_api",
		},
		Infrastructure: &InfrastructureData{
			Location:       location,
			Hospitals:      []string{"가까운 병원 A", "가까운 병원 B"},
			Schools:        []string{"가까운 학교 A", "가까운 학교 B"},
			ShoppingMalls:  []string{"가까운 쇼핑몰 A"},
			Parks:          []string{"가까운 공원 A", "가까운 공원 B"},
			SubwayStations: []string{"가까운 지하철역 A"},
			BusStops:       []string{"가까운 버스정류장 A", "가까운 버스정류장 B"},
			DistanceToCBD:  1 + rand.Float64()*10,
			DataSource:     "mock_infra_api",
		},
		LastUpdated:      time.Now(),
		DataSources:      []DataSourceType{DataSourceMock},
		ReliabilityScore: 0.7 + rand.Float64()*0.3,
	}
}
