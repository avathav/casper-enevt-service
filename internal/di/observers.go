package di

import (
	"event-service/internal/observers"
)

func NewEventUpdateObserver() {
	observers.NewEventUpdateObserver(NewEventUpdateProducer())
}
