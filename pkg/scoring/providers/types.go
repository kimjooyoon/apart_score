package providers

import (
	"time"
)

type DataSourceType string

const (
	DataSourceRealEstateAPI DataSourceType = "real_estate_api"
	DataSourceEnvironmental DataSourceType = "environmental"
	DataSourceTraffic       DataSourceType = "traffic"
	DataSourceLocalCache    DataSourceType = "local_cache"
	DataSourceMock          DataSourceType = "mock"
)

type DataFreshness struct {
	LastUpdated time.Time     `json:"last_updated"`
	Age         time.Duration `json:"age"`
	IsStale     bool          `json:"is_stale"`
	MaxAge      time.Duration `json:"max_age"`
}
type RealEstateData struct {
	PropertyID       string    `json:"property_id"`
	Address          string    `json:"address"`
	TradeType        string    `json:"trade_type"`
	Price            float64   `json:"price"`
	Area             float64   `json:"area"`
	Floor            int       `json:"floor"`
	TradeDate        time.Time `json:"trade_date"`
	BuildYear        int       `json:"build_year"`
	ContractType     string    `json:"contract_type"`
	DataSource       string    `json:"data_source"`
	ReliabilityScore float64   `json:"reliability_score"`
}
type EnvironmentalData struct {
	Location      string    `json:"location"`
	AirQuality    float64   `json:"air_quality"`
	NoiseLevel    float64   `json:"noise_level"`
	SunlightHours float64   `json:"sunlight_hours"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	GreenSpace    float64   `json:"green_space"`
	MeasuredAt    time.Time `json:"measured_at"`
	DataSource    string    `json:"data_source"`
}
type TrafficData struct {
	Location          string    `json:"location"`
	PublicTransport   float64   `json:"public_transport"`
	TrafficCongestion float64   `json:"traffic_congestion"`
	Walkability       float64   `json:"walkability"`
	BikeFriendly      float64   `json:"bike_friendly"`
	ParkingEase       float64   `json:"parking_ease"`
	MeasuredAt        time.Time `json:"measured_at"`
	DataSource        string    `json:"data_source"`
}
type InfrastructureData struct {
	Location       string   `json:"location"`
	Hospitals      []string `json:"hospitals"`
	Schools        []string `json:"schools"`
	ShoppingMalls  []string `json:"shopping_malls"`
	Parks          []string `json:"parks"`
	SubwayStations []string `json:"subway_stations"`
	BusStops       []string `json:"bus_stops"`
	DistanceToCBD  float64  `json:"distance_to_cbd"`
	DataSource     string   `json:"data_source"`
}
type ProviderConfig struct {
	SourceType    DataSourceType `json:"source_type"`
	APIKey        string         `json:"api_key,omitempty"`
	BaseURL       string         `json:"base_url,omitempty"`
	CacheDuration time.Duration  `json:"cache_duration"`
	Timeout       time.Duration  `json:"timeout"`
	MaxRetries    int            `json:"max_retries"`
	EnableCache   bool           `json:"enable_cache"`
}
type ProviderResult struct {
	Data         interface{}    `json:"data"`
	Freshness    DataFreshness  `json:"freshness"`
	Error        error          `json:"error,omitempty"`
	SourceType   DataSourceType `json:"source_type"`
	ResponseTime time.Duration  `json:"response_time"`
}
type CompositeData struct {
	RealEstate       *RealEstateData     `json:"real_estate,omitempty"`
	Environmental    *EnvironmentalData  `json:"environmental,omitempty"`
	Traffic          *TrafficData        `json:"traffic,omitempty"`
	Infrastructure   *InfrastructureData `json:"infrastructure,omitempty"`
	LastUpdated      time.Time           `json:"last_updated"`
	DataSources      []DataSourceType    `json:"data_sources"`
	ReliabilityScore float64             `json:"reliability_score"`
}
