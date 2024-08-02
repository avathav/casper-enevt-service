package eventcreator

import (
	"context"
	"time"

	"event-service/internal/domain/event"
	"event-service/internal/domain/event/aggregate"
	"event-service/internal/services"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ServiceName = "event creator"

type Handler interface {
	CreateEvent(context.Context, Request) (*aggregate.Event, error)
}

type Observer interface {
	Notify(context.Context, aggregate.Event) error
}

type EventCreator struct {
	adder         event.Adder
	observersList []Observer
}

func NewEventCreator(configuration ...Configuration) (*EventCreator, error) {
	ec := &EventCreator{}

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

func (ec EventCreator) validateRequiredResources() error {
	if ec.adder == nil {
		return services.NewErrResourceIsRequired(ServiceName, "event adder repository")
	}

	return nil
}

func (ec EventCreator) CreateEvent(ctx context.Context, r Request) (*aggregate.Event, error) {
	e, convErr := convertRequestToEvent(r)
	if convErr != nil {
		return nil, errors.Wrap(convErr, "cannot convert request to entry")
	}

	if err := ec.adder.Add(ctx, e); err != nil {
		return nil, errors.Wrap(err, "creating new event failed")
	}

	for _, observer := range ec.observersList {
		if err := observer.Notify(ctx, *e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

type Request struct {
	User                string
	Name, Description   string
	Capacity            int
	Duration            time.Duration
	Longitude           float64
	Latitude            float64
	DateStart           time.Time
	DateRegistrationEnd time.Time
	Public              bool
}

func convertRequestToEvent(r Request) (*aggregate.Event, error) {
	userID, uuidErr := uuid.Parse(r.User)
	if uuidErr != nil {
		return nil, uuidErr
	}

	return aggregate.NewEvent(aggregate.EventPayload{
		UserID:              userID,
		Name:                r.Name,
		Description:         r.Description,
		Lat:                 r.Latitude,
		Long:                r.Longitude,
		Capacity:            r.Capacity,
		Duration:            r.Duration,
		StartDate:           r.DateStart,
		RegistrationEndDate: r.DateRegistrationEnd,
		Public:              r.Public,
	})
}
