package repository

import (
	"context"
	"strings"

	"event-service/internal/domain/event/aggregate"

	"github.com/google/uuid"
)

type InvitationsStorage struct {
	items map[string]*aggregate.Invitation
}

func NewInvitationsStorage() *InvitationsStorage {
	return &InvitationsStorage{items: make(map[string]*aggregate.Invitation)}
}

func (i InvitationsStorage) FindBy(_ context.Context, eventID, userID uuid.UUID) (*aggregate.Invitation, error) {
	invitationID := getInvitationID(eventID, userID)
	invitation, ok := i.items[invitationID]
	if !ok {
		return nil, nil
	}

	return invitation, nil
}

func (i InvitationsStorage) Invite(_ context.Context, invitation *aggregate.Invitation) error {
	return i.Add(invitation)
}

func (i InvitationsStorage) Accept(_ context.Context, invitation *aggregate.Invitation) error {
	return i.Add(invitation)
}

func (i InvitationsStorage) Remove(_ context.Context, invitation *aggregate.Invitation) error {
	invitationID := getInvitationID(invitation.Event.Event.ExternalID, invitation.InvitedUser)
	delete(i.items, invitationID)

	return nil
}

func (i InvitationsStorage) Add(invitation *aggregate.Invitation) error {
	invitationID := getInvitationID(invitation.Event.Event.ExternalID, invitation.InvitedUser)
	i.items[invitationID] = invitation

	return nil
}

// getInvitationID returns a string that is a combination of eventID and userID
func getInvitationID(eventID, userID uuid.UUID) string {
	return strings.Join([]string{eventID.String(), userID.String()}, "_")
}
