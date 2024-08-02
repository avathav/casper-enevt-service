package observers

import (
	"context"
	"encoding/json"
	"event-service/internal/exchange/event"

	"event-service/internal/domain/event/aggregate"
	"event-service/internal/exchange"
)

type EventUpdateObserver struct {
	producer exchange.Producer
}

func NewEventUpdateObserver(producer exchange.Producer) *EventUpdateObserver {
	return &EventUpdateObserver{producer: producer}
}

func (e EventUpdateObserver) Notify(_ context.Context, event aggregate.Event) error {
	body, marshErr := json.Marshal(eventupdate.Message{Event: &event})
	if marshErr != nil {
		return marshErr
	}

	return e.producer.Publish(body)
}
