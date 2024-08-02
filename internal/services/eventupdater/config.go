package eventupdater

import (
	"event-service/internal/domain/event"
)

type Configuration func(*EventUpdater) error

func WithUpdaterRepository(updater event.Updater) Configuration {
	return func(ec *EventUpdater) error {
		ec.updater = updater

		return nil
	}
}

func WithFinderRepository(finder event.Finder) Configuration {
	return func(ec *EventUpdater) error {
		ec.finder = finder

		return nil
	}
}

func WithObservers(observers ...Observer) Configuration {
	return func(ec *EventUpdater) error {
		ec.observersList = append(ec.observersList, observers...)

		return nil
	}
}
