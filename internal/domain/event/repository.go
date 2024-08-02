package event

import (
	"context"

	"event-service/internal/domain/event/aggregate"
	"event-service/internal/domain/event/valueobject"

	"github.com/google/uuid"
)

type Adder interface {
	Add(context.Context, *aggregate.Event) error
}

type Updater interface {
	Update(context.Context, *aggregate.Event) error
}

type Finder interface {
	FindBy(context.Context, valueobject.ListRequest) ([]*aggregate.Event, error)
	FindByExternalID(ctx context.Context, id uuid.UUID) (*aggregate.Event, error)
}

type InviteFinder interface {
	FindBy(ctx context.Context, eventID, userID uuid.UUID) (*aggregate.Invitation, error)
}

type Inviter interface {
	Invite(context.Context, *aggregate.Invitation) error
	Accept(context.Context, *aggregate.Invitation) error
	Remove(context.Context, *aggregate.Invitation) error
}
