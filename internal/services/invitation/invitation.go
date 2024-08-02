package invitation

import (
	"context"
	"errors"

	"event-service/internal/domain/event"
	"event-service/internal/domain/event/aggregate"
	"event-service/internal/services"

	"github.com/google/uuid"
)

var (
	ServiceName                = "invitation"
	ErrInvalidUserID           = errors.New("cannot parse user ID")
	ErrInvalidEventID          = errors.New("cannot parse event ID")
	ErrInvitationAlreadyExists = errors.New("invitation already exists")
	ErrParticipantNotFound     = errors.New("there is no accepted participant in given event")
	ErrInvitationNotFound      = errors.New("there is no invitation found")
	ErrEventIsNotPublic        = errors.New("event is private")
)

type EventType string

const (
	UserAcceptedEvent EventType = "UserAcceptedEvent"
	UserInvitedEvent  EventType = "UserInvitedEvent"
	UserJoinedEvent   EventType = "UserJoinedEvent"
)

type Handler interface {
	Invite(ctx context.Context, eventID, userID string) error
	Accept(ctx context.Context, eventID, userID string) error
	Remove(ctx context.Context, eventID, userID string) error
	Join(ctx context.Context, eventID, userID string) error
}

type Observer interface {
	Notify(context.Context, aggregate.Invitation) error
}

type Invitation struct {
	inviter       event.Inviter
	inviteFinder  event.InviteFinder
	eventFinder   event.Finder
	observersList map[EventType][]Observer
}

// NewInvitation creates a new invitation service
func NewInvitation(configuration ...Configuration) (*Invitation, error) {
	i := &Invitation{
		observersList: map[EventType][]Observer{},
	}

	for _, cfg := range configuration {
		if err := cfg(i); err != nil {
			return nil, err
		}
	}

	if err := i.validateRequiredResources(); err != nil {
		return nil, err
	}

	return i, nil
}

func (i Invitation) AddObserver(eventType EventType, observer Observer) {
	i.observersList[eventType] = append(i.observersList[eventType], observer)
}

func (i Invitation) validateRequiredResources() error {
	if i.eventFinder == nil {
		return services.NewErrResourceIsRequired(ServiceName, "event finder repository")
	}

	if i.inviter == nil {
		return services.NewErrResourceIsRequired(ServiceName, "inviter repository")
	}

	if i.inviteFinder == nil {
		return services.NewErrResourceIsRequired(ServiceName, "invite finder repository")
	}

	return nil
}

func (i Invitation) parseInvitationData(eventID, userID string) (e uuid.UUID, u uuid.UUID, err error) {
	if u, err = uuid.Parse(userID); err != nil {
		return e, u, errors.Join(err, ErrInvalidUserID)
	}

	if e, err = uuid.Parse(eventID); err != nil {
		return e, u, errors.Join(err, ErrInvalidEventID)
	}

	return e, u, nil
}

// Invite invites a user to an event
func (i Invitation) Invite(ctx context.Context, eventExternalID, user string) error {
	eventID, userID, parseErr := i.parseInvitationData(eventExternalID, user)
	if parseErr != nil {
		return parseErr
	}

	if ia, err := i.inviteFinder.FindBy(ctx, eventID, userID); err != nil {
		return err
	} else if ia != nil {
		return ErrInvitationAlreadyExists
	}

	eventAggregate, finderErr := i.eventFinder.FindByExternalID(ctx, eventID)
	if finderErr != nil {
		return finderErr
	}

	ia, iaErr := aggregate.NewInvitation(eventAggregate, userID)
	if iaErr != nil {
		return iaErr
	}

	if err := i.inviter.Invite(ctx, ia); err != nil {
		return err
	}

	return i.runObservers(ctx, UserInvitedEvent, *ia)

}

// Accept accepts an invitation
func (i Invitation) Accept(ctx context.Context, eventExternalID, user string) error {
	eventID, userID, parseErr := i.parseInvitationData(eventExternalID, user)
	if parseErr != nil {
		return parseErr
	}

	ia, iaErr := i.inviteFinder.FindBy(ctx, eventID, userID)
	if iaErr != nil {
		return iaErr
	}

	if ia == nil {
		return ErrInvitationNotFound
	}

	if err := ia.Accept(); err != nil {
		return err
	}

	if err := i.inviter.Accept(ctx, ia); err != nil {
		return err
	}

	return i.runObservers(ctx, UserAcceptedEvent, *ia)
}

// Remove removes an invitation
func (i Invitation) Remove(ctx context.Context, eventExternalID, user string) error {
	eventID, userID, parseErr := i.parseInvitationData(eventExternalID, user)
	if parseErr != nil {
		return parseErr
	}

	ia, iaErr := i.inviteFinder.FindBy(ctx, eventID, userID)
	if iaErr != nil {
		return iaErr
	}

	if ia == nil {
		return ErrInvitationNotFound
	}

	if !ia.IsAccepted() {
		return ErrParticipantNotFound
	}

	if err := i.inviter.Remove(ctx, ia); err != nil {
		return err
	}

	return nil
}

// Join joins a user to an event
func (i Invitation) Join(ctx context.Context, eventExternalID, user string) error {
	eventID, userID, parseErr := i.parseInvitationData(eventExternalID, user)
	if parseErr != nil {
		return parseErr
	}

	eventAggregate, finderErr := i.eventFinder.FindByExternalID(ctx, eventID)
	if finderErr != nil {
		return finderErr
	}

	if !eventAggregate.Event.Public {
		return ErrEventIsNotPublic
	}

	ia, iaErr := i.inviteFinder.FindBy(ctx, eventID, userID)
	if iaErr != nil {
		return iaErr
	}

	if ia == nil {
		ia, iaErr = aggregate.NewInvitation(eventAggregate, userID)
		if iaErr != nil {
			return iaErr
		}
	}

	if err := ia.Accept(); err != nil {
		return err
	}

	if err := i.inviter.Invite(ctx, ia); err != nil {
		return err
	}

	return i.runObservers(ctx, UserJoinedEvent, *ia)
}

func (i Invitation) runObservers(ctx context.Context, eventType EventType, ia aggregate.Invitation) error {
	for _, observer := range i.observersList[eventType] {
		if err := observer.Notify(ctx, ia); err != nil {
			return err
		}
	}

	return nil
}
