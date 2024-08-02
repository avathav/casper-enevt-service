package valueobject

import (
	"time"

	"github.com/markphelps/optional"
)

//type optionalBool

type ListRequest struct {
	user     string
	name     string
	public   optional.Bool
	interval TimeInterval
	distance Distance
}

func (l ListRequest) User() (string, bool) {
	return l.user, l.user != ""
}

func (l ListRequest) Name() (string, bool) {
	return l.name, l.name != ""
}

func (l ListRequest) Public() (bool, bool) {
	p, err := l.public.Get()

	return p, err == nil
}

func (l ListRequest) Interval() (TimeInterval, bool) {
	return l.interval, l.interval.IsSet()
}

func (l ListRequest) Distance() (Distance, bool) {
	return l.distance, l.distance.IsSet()
}

type ListRequestConfiguration func(*ListRequest)

func NewListRequest(configs ...ListRequestConfiguration) ListRequest {
	r := ListRequest{}

	for _, cfg := range configs {
		cfg(&r)
	}

	return r
}

func WithTimeInterval(date time.Time, interval time.Duration) ListRequestConfiguration {
	return func(r *ListRequest) {
		r.interval = TimeInterval{
			date:     date,
			interval: interval,
		}
	}
}

func WithDistance(distance int64, lat, long float64) ListRequestConfiguration {
	return func(r *ListRequest) {
		r.distance = Distance{
			distance:      distance,
			initLatitude:  lat,
			initLongitude: long,
		}
	}
}

func WithPublicParam(public bool) ListRequestConfiguration {
	return func(r *ListRequest) {
		r.public = optional.NewBool(public)
	}
}

func WithName(name string) ListRequestConfiguration {
	return func(r *ListRequest) {
		r.name = name
	}
}

func WithUser(user string) ListRequestConfiguration {
	return func(r *ListRequest) {
		r.user = user
	}
}
