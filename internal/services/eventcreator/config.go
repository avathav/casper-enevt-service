package eventcreator

import (
	"event-service/internal/domain/event"
)

type Configuration func(*EventCreator) error

func WithAdderRepository(adder event.Adder) Configuration {
	return func(ec *EventCreator) error {
		ec.adder = adder

		return nil
	}
}

func WithObservers(observers ...Observer) Configuration {
	return func(ec *EventCreator) error {
		ec.observersList = append(ec.observersList, observers...)

		return nil
	}
}
