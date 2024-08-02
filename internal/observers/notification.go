package observers

import (
	"context"
	"fmt"

	"event-service/internal/domain/event/aggregate"
)

type InvitationNotificationObserver struct{}

func NewInvitationNotificationObserver() *InvitationNotificationObserver {
	return &InvitationNotificationObserver{}
}

func (e InvitationNotificationObserver) Notify(ctx context.Context, invitation aggregate.Invitation) error {
	fmt.Printf("Sending notification after invitation to user %s", invitation.InvitedUser)

	return nil
}
