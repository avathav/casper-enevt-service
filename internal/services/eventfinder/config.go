package eventfinder

import "event-service/internal/domain/event"

type Configuration func(*EventFinder) error

func WithFinderRepository(finder event.Finder) Configuration {
	return func(ec *EventFinder) error {
		ec.finder = finder

		return nil
	}
}
