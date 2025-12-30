package metadata

type MetadataType int
type FactorType string

const (
	FactorInternal FactorType = "internal"
	FactorExternal FactorType = "external"
)
const (
	FloorLevel MetadataType = iota
	DistanceToStation
	ElevatorPresence
	ConstructionYear
	ConstructionCompany
	ApartmentSize
	NearbyAmenities
	TransportationAccess
	SchoolDistrict
	CrimeRate
	GreenSpaceRatio
	Parking
	MaintenanceFee
	HeatingSystem
	MetadataTypeCount
)
