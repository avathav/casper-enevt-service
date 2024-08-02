package eventupdater

import (
	"context"
	"time"

	"event-service/internal/domain/common/valueobject"
	"event-service/internal/domain/event"
	"event-service/internal/domain/event/aggregate"
	"event-service/internal/services"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ServiceName = "event updater"

type Handler interface {
	UpdateEvent(context.Context, Request) (*aggregate.Event, error)
}

type Observer interface {
	Notify(context.Context, aggregate.Event) error
}

type EventUpdater struct {
	updater       event.Updater
	finder        event.Finder
	observersList []Observer
}

func NewEventUpdater(configuration ...Configuration) (*EventUpdater, error) {
	ec := &EventUpdater{}

	for _, cfg := range configuration {
		if err := cfg(ec); err != nil {
			return nil, err
		}
	}

	if err := ec.validateRequiredResources(); err != nil {
		return nil, err
	}

	return ec, nil
}

func (ec EventUpdater) validateRequiredResources() error {
	if ec.updater == nil {
		return services.NewErrResourceIsRequired(ServiceName, "event updater repository")
	}

	if ec.finder == nil {
		return services.NewErrResourceIsRequired(ServiceName, "event finder repository")
	}

	return nil
}

func (ec EventUpdater) UpdateEvent(ctx context.Context, r Request) (*aggregate.Event, error) {
	id, idErr := uuid.Parse(r.ID)
	if idErr != nil {
		return nil, idErr
	}

	ia, findErr := ec.finder.FindByExternalID(ctx, id)
	if findErr != nil {
		return nil, errors.Wrap(findErr, "event not found")
	}

	if err := updateAggregateWithRequest(ia, r); err != nil {
		return nil, errors.Wrap(err, "cannot convert request to entry")
	}

	if err := ec.updater.Update(ctx, ia); err != nil {
		return nil, errors.Wrap(err, "creating new event failed")
	}

	for _, observer := range ec.observersList {
		if err := observer.Notify(ctx, *ia); err != nil {
			return nil, err
		}
	}

	return ia, nil
}

type Request struct {
	ID                    string
	Name, Description     string
	Capacity              int
	Longitude             float64
	Latitude              float64
	DateStart             *time.Time
	DateEnd               *time.Time
	DateRegistrationStart *time.Time
	DateRegistrationEnd   *time.Time
	Public                *bool
}

func updateAggregateWithRequest(a *aggregate.Event, r Request) (err error) {
	if r.Name != "" {
		a.Event.Name = r.Name
	}

	if r.Description != "" {
		a.Event.Description = r.Description
	}

	if r.Capacity != 0 {
		a.Event.Capacity = r.Capacity
	}

	if r.Public != nil {
		a.Event.Public = *r.Public
	}

	if r.Latitude != 0 && r.Longitude != 0 {
		a.Location = &aggregate.Location{
			Spot: valueobject.NewLocation(r.Latitude, r.Longitude),
		}
	}

	if r.DateStart != nil && r.DateEnd != nil {
		a.EventPeriod, err = valueobject.EventPeriod{}.WithStartAndEndDate(*r.DateStart, *r.DateEnd)
		if err != nil {
			return err
		}
	}

	if r.DateRegistrationStart != nil && r.DateRegistrationEnd != nil {
		a.RegistrationPeriod, err = valueobject.Period{}.WithStartAndEndDate(*r.DateRegistrationStart, *r.DateRegistrationEnd)
		if err != nil {
			return err
		}
	}

	return nil
}
