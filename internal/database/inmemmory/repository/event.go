package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"event-service/internal/domain/event/aggregate"
	"event-service/internal/domain/event/valueobject"
)

type EventsStorage struct {
	items map[uuid.UUID]*aggregate.Event
}

func NewEventsStorage() *EventsStorage {
	return &EventsStorage{items: make(map[uuid.UUID]*aggregate.Event)}
}

func (e EventsStorage) Update(ctx context.Context, event *aggregate.Event) error {
	eventToUpdate, err := e.FindByExternalID(ctx, event.Event.ExternalID)
	if err != nil {
		return err
	}

	event.ID = eventToUpdate.ID
	e.items[event.Event.ExternalID] = event

	return nil
}

func (e EventsStorage) FindBy(_ context.Context, request valueobject.ListRequest) ([]*aggregate.Event, error) {
	panic("implement me")
}

func (e EventsStorage) FindByExternalID(_ context.Context, id uuid.UUID) (*aggregate.Event, error) {
	event, ok := e.items[id]
	if !ok {
		return nil, errors.New("events with requested id not found")
	}

	return event, nil
}

func (e EventsStorage) Add(_ context.Context, event *aggregate.Event) error {
	event.ID = uint(len(e.items) + 1)
	e.items[event.Event.ExternalID] = event
	return nil
}
