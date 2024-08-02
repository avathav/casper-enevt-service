package aggregate

import (
	"errors"
	"time"

	"event-service/internal/domain/common/entity"
	"event-service/internal/domain/common/valueobject"

	"github.com/google/uuid"
)

var (
	ErrUserIDRequired    = errors.New("user id cannot be empty")
	ErrEventNameRequired = errors.New("event must be named")
)

type Event struct {
	ID                 uint
	UserID             uuid.UUID
	Event              *entity.Event
	Location           *Location
	EventPeriod        valueobject.EventPeriod
	RegistrationPeriod valueobject.Period
	Participants       []uuid.UUID
}

type EventPayload struct {
	UserID              uuid.UUID
	Name                string
	Description         string
	Lat                 float64
	Long                float64
	Capacity            int
	Duration            time.Duration
	StartDate           time.Time
	RegistrationEndDate time.Time
	Public              bool
}

func NewEvent(cfg EventPayload) (event *Event, err error) {
	if cfg.UserID == uuid.Nil {
		return nil, ErrUserIDRequired
	}

	if cfg.Name == "" {
		return nil, ErrEventNameRequired
	}

	event = &Event{
		UserID: cfg.UserID,
		Event: &entity.Event{
			ExternalID:  uuid.New(),
			Name:        cfg.Name,
			Description: cfg.Description,
			Capacity:    cfg.Capacity,
			Public:      cfg.Public,
		},
		Location: &Location{
			Spot: valueobject.NewLocation(cfg.Lat, cfg.Long),
		},
	}

	if event.EventPeriod, err = event.EventPeriod.WithStartAndDuration(cfg.StartDate, cfg.Duration); err != nil {
		return nil, err
	}

	if err = event.SetUpRegistrationPeriod(); err != nil {
		return nil, err
	}

	return event, nil
}

func (e *Event) SetUpRegistrationPeriod() (err error) {
	if e.RegistrationPeriod, err = e.RegistrationPeriod.WithStartAndEndDate(time.Now(), e.EventPeriod.Start()); err != nil {
		return err
	}

	return
}

func (e *Event) ParticipantsNumber() int {
	return len(e.Participants)
}

func (e *Event) OpenToJoin() bool {
	return e.Event.Capacity > e.ParticipantsNumber() && e.RegistrationPeriod.Contains(time.Now())
}
