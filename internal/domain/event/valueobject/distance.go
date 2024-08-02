package valueobject

type Distance struct {
	distance      int64
	initLatitude  float64
	initLongitude float64
}

func NewDistance(distance int64, lat float64, long float64) Distance {
	return Distance{distance: distance, initLatitude: lat, initLongitude: long}
}

func (d Distance) IsSet() bool {
	return d.initLatitude > 0 && d.initLongitude > 0 && d.distance > 0
}

func (d Distance) Distance() int64 {
	return d.distance
}

func (d Distance) InitLatitude() float64 {
	return d.initLatitude
}

func (d Distance) InitLongitude() float64 {
	return d.initLongitude
}
