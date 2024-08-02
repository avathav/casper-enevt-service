package eventfinder

import (
	"context"
	"event-service/internal/services"

	"event-service/internal/domain/event"
	"event-service/internal/domain/event/aggregate"

	"github.com/google/uuid"
)

var ServiceName = "event finder"

type ListHandler interface {
	List(context.Context, Request) ([]*aggregate.Event, error)
	GetByID(context.Context, string) (*aggregate.Event, error)
}

type EventFinder struct {
	finder event.Finder
}

func NewEventFinder(configuration ...Configuration) (*EventFinder, error) {
	ef := &EventFinder{}

	for _, cfg := range configuration {
		if err := cfg(ef); err != nil {
			return nil, err
		}
	}

	if err := ef.validateRequiredResources(); err != nil {
		return nil, err
	}

	return ef, nil
}

func (ef EventFinder) validateRequiredResources() error {
	if ef.finder == nil {
		return services.NewErrResourceIsRequired(ServiceName, "event finder repository")
	}

	return nil
}

func (ef EventFinder) List(ctx context.Context, r Request) ([]*aggregate.Event, error) {
	return ef.finder.FindBy(ctx, convertRequestToListRequest(r))
}

func (ef EventFinder) GetByID(ctx context.Context, id string) (*aggregate.Event, error) {
	externalID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return ef.finder.FindByExternalID(ctx, externalID)
}
