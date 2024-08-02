package di

import (
	"event-service/internal/observers"
	"event-service/internal/services/eventcreator"
	"event-service/internal/services/eventfinder"
	"event-service/internal/services/eventupdater"
	"event-service/internal/services/invitation"
)

func DefaultEventsAddHandler() (*eventcreator.EventCreator, error) {
	return eventcreator.NewEventCreator(
		eventcreator.WithAdderRepository(EventsRepository()),
	)
}

func DefaultEventsListHandler() (*eventfinder.EventFinder, error) {
	return eventfinder.NewEventFinder(
		eventfinder.WithFinderRepository(EventsRepository()),
	)
}

func DefaultInvitationHandler() (*invitation.Invitation, error) {
	i, err := invitation.NewInvitation(
		invitation.WithEventFinderRepository(EventsRepository()),
		invitation.WithInvitationRepository(InvitationRepository()),
		invitation.WithInviteFinderRepository(InvitationRepository()),
	)

	if err != nil {
		return nil, err
	}

	i.AddObserver(invitation.UserInvitedEvent, observers.NewInvitationNotificationObserver())

	return i, nil
}

func DefaultEventsUpdateHandler() (*eventupdater.EventUpdater, error) {
	return eventupdater.NewEventUpdater(
		eventupdater.WithUpdaterRepository(EventsRepository()),
		eventupdater.WithFinderRepository(EventsRepository()),
		eventupdater.WithObservers(observers.NewEventUpdateObserver(NewEventUpdateProducer())),
	)
}
