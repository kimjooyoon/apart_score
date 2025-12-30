package metadata

type MetadataType int

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
