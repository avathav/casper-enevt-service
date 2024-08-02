package eventfinder

import (
	"event-service/internal/domain/event/valueobject"
	"time"
)

type Request struct {
	User     string
	Name     string
	Location *LocationRequest
	Upcoming *UpcomingEventRequest
	Public   bool
}

type LocationRequest struct {
	Latitude  float64
	Longitude float64
	Distance  int64
}

type UpcomingEventRequest struct {
	Date     time.Time
	Interval time.Duration
}

func convertRequestToListRequest(r Request) valueobject.ListRequest {
	cfg := []valueobject.ListRequestConfiguration{
		valueobject.WithName(r.Name),
		valueobject.WithUser(r.User),
		valueobject.WithPublicParam(r.Public),
	}

	if r.Location != nil {
		cfg = append(cfg, valueobject.WithDistance(r.Location.Distance, r.Location.Latitude, r.Location.Longitude))
	}

	if r.Upcoming != nil {
		cfg = append(cfg, valueobject.WithTimeInterval(r.Upcoming.Date, r.Upcoming.Interval))
	}

	return valueobject.NewListRequest(cfg...)
}
