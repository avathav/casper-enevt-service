package aggregate

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvitedUserIDRequired     = errors.New("for proper invitation user identifier must be send ")
	ErrNoEvent                   = errors.New("invitation must have a valid event attached ")
	ErrCapacityFull              = errors.New("event capacity has been reached")
	ErrRegistrationClosed        = errors.New("event registration is closed")
	ErrInvitationAlreadyAccepted = errors.New("invitation already accepted")
)

type Invitation struct {
	Event       *Event
	InvitedUser uuid.UUID
	AcceptedAt  *time.Time
}

func NewInvitation(event *Event, participantID uuid.UUID) (i *Invitation, err error) {
	if participantID == uuid.Nil {
		return nil, ErrInvitedUserIDRequired
	}

	if event == nil {
		return nil, ErrNoEvent
	}

	return &Invitation{
		Event:       event,
		InvitedUser: participantID,
	}, nil
}

func (i *Invitation) IsAccepted() bool {
	return i.AcceptedAt != nil && !i.AcceptedAt.IsZero()
}

func (i *Invitation) Accept() error {
	if i.IsAccepted() {
		return ErrInvitationAlreadyAccepted
	}

	if i.Event.Event.Capacity <= i.Event.ParticipantsNumber() {
		return ErrCapacityFull
	}

	if !i.Event.RegistrationPeriod.Contains(time.Now()) {
		return ErrRegistrationClosed
	}

	now := time.Now()
	i.AcceptedAt = &now

	return nil
}
