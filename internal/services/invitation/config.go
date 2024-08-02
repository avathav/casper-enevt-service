package invitation

import (
	"event-service/internal/domain/event"
)

type Configuration func(i *Invitation) error

func WithInvitationRepository(inviter event.Inviter) Configuration {
	return func(i *Invitation) error {
		i.inviter = inviter

		return nil
	}
}

func WithEventFinderRepository(finder event.Finder) Configuration {
	return func(i *Invitation) error {
		i.eventFinder = finder

		return nil
	}
}

func WithInviteFinderRepository(finder event.InviteFinder) Configuration {
	return func(i *Invitation) error {
		i.inviteFinder = finder

		return nil
	}
}
