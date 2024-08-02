package repository

import (
	"event-service/internal/domain/common/valueobject"
	"event-service/internal/domain/event/aggregate"
)

type Location struct {
	BaseModel
	ID        uint `gorm:"primaryKey"`
	Latitude  float64
	Longitude float64
	Events    []Event
}

func (l Location) toLocationAggregate() *aggregate.Location {
	return &aggregate.Location{
		ID:   l.ID,
		Spot: valueobject.NewLocation(l.Latitude, l.Longitude),
	}
}

func RecordFromLocationAggregate(l aggregate.Location) Location {
	return Location{
		ID:        l.ID,
		Latitude:  l.Spot.Lat(),
		Longitude: l.Spot.Long(),
	}
}
