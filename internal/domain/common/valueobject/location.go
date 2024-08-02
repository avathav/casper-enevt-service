package valueobject

type Location struct {
	lat  float64
	long float64
}

func (l Location) Lat() float64 {
	return l.lat
}

func (l Location) Long() float64 {
	return l.long
}

func NewLocation(lat float64, long float64) Location {
	return Location{lat: lat, long: long}
}
