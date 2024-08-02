package repository

import (
	"context"
	"time"

	dbgrom "event-service/internal/database/gorm"
	"event-service/internal/domain/common/entity"
	"event-service/internal/domain/event/aggregate"
	"event-service/internal/domain/event/valueobject"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Event struct {
	BaseModel
	ID                    uint `gorm:"primaryKey"`
	ExternalID            uuid.UUID
	User                  uuid.UUID
	Name                  string
	LocationID            uint
	Description           string
	Capacity              int
	StartDate             time.Time
	EndDate               time.Time
	RegistrationStartDate time.Time
	RegistrationEndDate   time.Time
	Public                bool
	Location              Location `gorm:"foreignKey:LocationID"`
	Invitations           []Invitation
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.Location.ID > 0 {
		return nil
	}

	tx.Where(e.Location).Find(&e.Location).Limit(1)

	return nil
}

func (e *Event) ToEventAggregate() (entry *aggregate.Event, err error) {
	entry = &aggregate.Event{
		ID:     e.ID,
		UserID: e.User,
		Event: &entity.Event{
			ExternalID:  e.ExternalID,
			Name:        e.Name,
			Description: e.Description,
			Capacity:    e.Capacity,
			Public:      e.Public,
		},
		Location: e.Location.toLocationAggregate(),
	}

	entry.EventPeriod, err = entry.EventPeriod.WithStartAndEndDate(e.StartDate, e.EndDate)
	if err != nil {
		return
	}

	entry.RegistrationPeriod, err = entry.RegistrationPeriod.WithStartAndEndDate(e.RegistrationStartDate, e.RegistrationEndDate)
	if err != nil {
		return
	}

	if e.Invitations != nil {
		for _, invitation := range e.Invitations {
			entry.Participants = append(entry.Participants, invitation.UserID)
		}
	}

	return entry, nil
}

func RecordFromEventAggregate(e aggregate.Event) Event {
	event := Event{
		ID:                    e.ID,
		ExternalID:            e.Event.ExternalID,
		User:                  e.UserID,
		Name:                  e.Event.Name,
		Description:           e.Event.Description,
		Capacity:              e.Event.Capacity,
		StartDate:             e.EventPeriod.Start(),
		EndDate:               e.EventPeriod.End(),
		RegistrationStartDate: e.RegistrationPeriod.Start(),
		RegistrationEndDate:   e.RegistrationPeriod.End(),
		Public:                e.Event.Public,
	}

	if e.Location != nil {
		event.Location = Location{
			ID:        e.Location.ID,
			Latitude:  e.Location.Spot.Lat(),
			Longitude: e.Location.Spot.Long(),
		}
	}

	return event
}

type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r EventRepository) Update(ctx context.Context, entry *aggregate.Event) error {
	db, dbErr := dbgrom.ConnectionFromContext(ctx)
	if dbErr != nil {
		return errors.Wrap(dbErr, "events repository")
	}

	event := RecordFromEventAggregate(*entry)

	if err := db.Updates(&event).Error; err != nil {
		return errors.Wrap(err, "events repository update")
	}

	return nil
}

func (r EventRepository) FindBy(ctx context.Context, request valueobject.ListRequest) ([]*aggregate.Event, error) {
	db, dbErr := dbgrom.ConnectionFromContext(ctx)
	if dbErr != nil {
		return nil, errors.Wrap(dbErr, "events repository")
	}

	if name, ok := request.Name(); ok {
		db = db.Where("name like ?", "%"+name+"%")
	}

	if user, ok := request.User(); ok {
		db = db.Where("user = ?", user)
	}

	if public, ok := request.Public(); ok {
		db = db.Where("public = ?", public)
	}

	if interval, ok := request.Interval(); ok {
		db = db.Where("start_date BETWEEN ? AND ?", interval.GetStartDate(), interval.GetEndDate())
		interval.GetStartDate()
	}

	if distance, ok := request.Distance(); ok {
		db.Where(`location_id IN (SELECT id FROM (SELECT id, (6371 *
				acos(cos(radians(?)) *
				cos(radians(latitude)) *
				cos(radians(longitude) -
				radians(?)) +
				sin(radians(?)) *
				sin(radians(latitude)))
		) AS distance FROM locations HAVING distance < ? ) AS nearby)`,
			distance.InitLatitude(), distance.InitLongitude(), distance.InitLatitude(), distance.Distance())
	}

	var items []Event
	if findErr := db.Preload("Location").Find(&items).Error; findErr != nil {
		return nil, errors.Wrap(findErr, "events repository find by")
	}

	entries := make([]*aggregate.Event, len(items))

	for i := 0; i < len(items); i++ {
		entries[i], _ = items[i].ToEventAggregate()
	}

	return entries, nil

}

func (r EventRepository) FindByExternalID(ctx context.Context, id uuid.UUID) (*aggregate.Event, error) {
	db, dbErr := dbgrom.ConnectionFromContext(ctx)
	if dbErr != nil {
		return nil, errors.Wrap(dbErr, "events repository")
	}

	item := Event{}

	if findErr := db.Preload("Location").Preload("Invitations").First(&item, "external_id = ?", id).Error; findErr != nil {
		return nil, errors.Wrap(findErr, "events repository find by external ID")
	}

	return item.ToEventAggregate()
}

func (r EventRepository) Add(ctx context.Context, entry *aggregate.Event) error {
	db, dbErr := dbgrom.ConnectionFromContext(ctx)
	if dbErr != nil {
		return errors.Wrap(dbErr, "events repository")
	}
	record := RecordFromEventAggregate(*entry)

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return errors.Wrap(err, "events repository add")
		}

		entry.ID = record.ID

		return nil
	}); err != nil {
		return err
	}

	return nil
}
