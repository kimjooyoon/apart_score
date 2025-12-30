package context

import "time"

type ContextType string

const (
	ContextTypeRegional ContextType = "regional"
	ContextTypeTemporal ContextType = "temporal"
)

type RegionType string

const (
	RegionTypeSeoul     RegionType = "seoul"
	RegionTypeBusan     RegionType = "busan"
	RegionTypeDaegu     RegionType = "daegu"
	RegionTypeIncheon   RegionType = "incheon"
	RegionTypeGwangju   RegionType = "gwangju"
	RegionTypeDaejeon   RegionType = "daejeon"
	RegionTypeUlsan     RegionType = "ulsan"
	RegionTypeSejong    RegionType = "sejong"
	RegionTypeGyeonggi  RegionType = "gyeonggi"
	RegionTypeGangwon   RegionType = "gangwon"
	RegionTypeChungbuk  RegionType = "chungbuk"
	RegionTypeChungnam  RegionType = "chungnam"
	RegionTypeJeonbuk   RegionType = "jeonbuk"
	RegionTypeJeonnam   RegionType = "jeonnam"
	RegionTypeGyeongbuk RegionType = "gyeongbuk"
	RegionTypeGyeongnam RegionType = "gyeongnam"
	RegionTypeJeju      RegionType = "jeju"
	RegionTypeOther     RegionType = "other"
)

type TransportLevel string

const (
	TransportLevelHigh   TransportLevel = "high"
	TransportLevelMedium TransportLevel = "medium"
	TransportLevelLow    TransportLevel = "low"
)

type Season string

const (
	SeasonSpring Season = "spring"
	SeasonSummer Season = "summer"
	SeasonFall   Season = "fall"
	SeasonWinter Season = "winter"
)

type MarketCondition string

const (
	MarketConditionBoom      MarketCondition = "boom"
	MarketConditionNormal    MarketCondition = "normal"
	MarketConditionSlump     MarketCondition = "slump"
	MarketConditionRecession MarketCondition = "recession"
)

type Location struct {
	Address    string     `json:"address"`
	Region     string     `json:"region"`
	SubRegion  string     `json:"sub_region"`
	Latitude   float64    `json:"latitude"`
	Longitude  float64    `json:"longitude"`
	RegionType RegionType `json:"region_type"`
}
type RegionalContext struct {
	RegionType         RegionType     `json:"region_type"`
	PopulationDensity  float64        `json:"population_density"`
	TransportLevel     TransportLevel `json:"transport_level"`
	CostOfLiving       float64        `json:"cost_of_living"`
	AvgApartmentPrice  float64        `json:"avg_apartment_price"`
	SchoolDistrictRank int            `json:"school_district_rank"`
	CrimeRate          float64        `json:"crime_rate"`
	GreenSpaceRatio    float64        `json:"green_space_ratio"`
}
type TemporalContext struct {
	Timestamp       time.Time       `json:"timestamp"`
	Year            int             `json:"year"`
	Month           int             `json:"month"`
	Season          Season          `json:"season"`
	MarketCondition MarketCondition `json:"market_condition"`
	EconomicIndex   float64         `json:"economic_index"`
	InflationRate   float64         `json:"inflation_rate"`
	InterestRate    float64         `json:"interest_rate"`
}
type ContextData struct {
	Location Location         `json:"location"`
	Regional *RegionalContext `json:"regional,omitempty"`
	Temporal *TemporalContext `json:"temporal,omitempty"`
}
