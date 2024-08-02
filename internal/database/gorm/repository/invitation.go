package repository

import (
	"context"
	"gorm.io/gorm"
	"time"

	dbgrom "event-service/internal/database/gorm"
	"event-service/internal/domain/event/aggregate"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Invitation struct {
	BaseModel
	EventID    uint      `gorm:"primaryKey;autoIncrement:false"`
	UserID     uuid.UUID `gorm:"primaryKey;"`
	AcceptedAt *time.Time
	Event      Event `gorm:"foreignKey:EventID"`
}

func (i Invitation) toAggregate() (ia *aggregate.Invitation, err error) {
	ia = &aggregate.Invitation{
		InvitedUser: i.UserID,
		AcceptedAt:  i.AcceptedAt,
	}

	if ia.Event, err = i.Event.ToEventAggregate(); err != nil {
		return nil, err
	}

	return ia, nil
}

func RecordFromInvitationAggregate(i aggregate.Invitation) Invitation {
	return Invitation{
		EventID:    i.Event.ID,
		UserID:     i.InvitedUser,
		AcceptedAt: i.AcceptedAt,
	}
}

type InvitationRepository struct{}

func NewInvitationRepository() *InvitationRepository {
	return &InvitationRepository{}
}

func (r InvitationRepository) FindBy(ctx context.Context, eventID, userID uuid.UUID) (*aggregate.Invitation, error) {
	db, dbErr := dbgrom.ConnectionFromContext(ctx)
	if dbErr != nil {
		return nil, errors.Wrap(dbErr, "invitation repository")
	}

	item := Invitation{
		UserID: userID,
	}

	if findErr := db.Joins("JOIN events on events.id = invitations.event_id AND events.external_id = ?", eventID).
		Preload("Event").
		First(&item).Error; findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errors.Wrap(findErr, "invitation repository find by")
	}

	return item.toAggregate()
}

func (r InvitationRepository) Invite(ctx context.Context, invitation *aggregate.Invitation) error {
	db, dbErr := dbgrom.ConnectionFromContext(ctx)
	if dbErr != nil {
		return errors.Wrap(dbErr, "invitation repository")
	}

	record := RecordFromInvitationAggregate(*invitation)

	if err := db.Create(&record).Error; err != nil {
		return errors.Wrap(err, "invitation repository invite")
	}

	return nil
}

func (r InvitationRepository) Accept(ctx context.Context, invitation *aggregate.Invitation) error {
	db, dbErr := dbgrom.ConnectionFromContext(ctx)
	if dbErr != nil {
		return errors.Wrap(dbErr, "invitation repository")
	}

	record := RecordFromInvitationAggregate(*invitation)

	if err := db.Save(&record).Error; err != nil {
		return errors.Wrap(err, "invitation repository accept")
	}

	return nil
}

func (r InvitationRepository) Remove(ctx context.Context, invitation *aggregate.Invitation) error {
	db, dbErr := dbgrom.ConnectionFromContext(ctx)
	if dbErr != nil {
		return errors.Wrap(dbErr, "invitation repository")
	}

	record := RecordFromInvitationAggregate(*invitation)

	if err := db.Delete(&record).Error; err != nil {
		return errors.Wrap(err, "invitation repository accept")
	}

	return nil
}
