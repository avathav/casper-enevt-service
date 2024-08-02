package eventupdater

import (
	"context"
	"reflect"
	"testing"
	"time"

	"event-service/internal/database/inmemmory/repository"
	"event-service/internal/domain/common/entity"
	"event-service/internal/domain/common/valueobject"
	"event-service/internal/domain/event"
	"event-service/internal/domain/event/aggregate"

	"github.com/google/uuid"
)

type eventServiceFields struct {
	updater       event.Updater
	finder        event.Finder
	observersList []Observer
}

func newValidEventServiceFields(ea ...aggregate.Event) eventServiceFields {
	storage := repository.NewEventsStorage()
	if len(ea) > 0 {
		for _, e := range ea {
			_ = storage.Add(context.Background(), &e)
		}
	}

	return eventServiceFields{
		updater:       storage,
		finder:        storage,
		observersList: []Observer{mockObserver{}},
	}
}

type mockObserver struct{}

func (m mockObserver) Notify(_ context.Context, _ aggregate.Event) error {
	return nil
}

func TestEventUpdater_UpdateEvent(t *testing.T) {
	eventID := uuid.New()
	mockEventPeriod, _ := valueobject.EventPeriod{}.WithStartAndEndDate(time.Now().Add(time.Hour*-24), time.Now().Add(time.Hour*48))
	mockPeriod, _ := valueobject.Period{}.WithStartAndEndDate(time.Now().Add(time.Hour*-24), time.Now().Add(time.Hour*48))
	trueValue := true
	newDateStart := time.Now().Add(time.Hour * 240)
	newDateEnd := time.Now().Add(time.Hour * 480)

	initialEvent := aggregate.Event{
		ID:     1,
		UserID: uuid.New(),
		Event: &entity.Event{
			ExternalID:  eventID,
			Name:        "Initial Event Name",
			Description: "Initial Event Description",
			Capacity:    10,
			Public:      false,
		},
		Location: &aggregate.Location{
			ID:   1,
			Spot: valueobject.NewLocation(10, 10),
		},
		EventPeriod:        mockEventPeriod,
		RegistrationPeriod: mockPeriod,
		Participants:       []uuid.UUID{uuid.New(), uuid.New()},
	}

	tests := []struct {
		name    string
		fields  eventServiceFields
		request Request
		want    aggregate.Event
		wantErr bool
	}{
		{
			name:   "update event name",
			fields: newValidEventServiceFields(initialEvent),
			request: Request{
				ID:   eventID.String(),
				Name: "New Name",
			},
			want: func() aggregate.Event {
				ia := initialEvent
				ia.Event.Name = "New Name"

				return ia
			}(),
			wantErr: false,
		}, {
			name:   "update event description",
			fields: newValidEventServiceFields(initialEvent),
			request: Request{
				ID:          eventID.String(),
				Description: "New Description",
			},
			want: func() aggregate.Event {
				ia := initialEvent
				ia.Event.Description = "New Description"

				return ia
			}(),
			wantErr: false,
		}, {
			name:   "update event capacity",
			fields: newValidEventServiceFields(initialEvent),
			request: Request{
				ID:       eventID.String(),
				Capacity: 20,
			},
			want: func() aggregate.Event {
				ia := initialEvent
				ia.Event.Capacity = 20

				return ia
			}(),
			wantErr: false,
		}, {
			name:   "update event public",
			fields: newValidEventServiceFields(initialEvent),
			request: Request{
				ID:     eventID.String(),
				Public: &trueValue,
			},
			want: func() aggregate.Event {
				ia := initialEvent
				ia.Event.Public = trueValue

				return ia
			}(),
			wantErr: false,
		}, {
			name:   "update event date",
			fields: newValidEventServiceFields(initialEvent),
			request: Request{
				ID:        eventID.String(),
				DateStart: &newDateStart,
				DateEnd:   &newDateEnd,
			},
			want: func() aggregate.Event {
				ia := initialEvent
				ia.EventPeriod, _ = valueobject.EventPeriod{}.WithStartAndEndDate(newDateStart, newDateEnd)

				return ia
			}(),
			wantErr: false,
		}, {
			name:   "update event registration dates",
			fields: newValidEventServiceFields(initialEvent),
			request: Request{
				ID:                    eventID.String(),
				DateRegistrationStart: &newDateStart,
				DateRegistrationEnd:   &newDateEnd,
			},
			want: func() aggregate.Event {
				ia := initialEvent
				ia.RegistrationPeriod, _ = valueobject.Period{}.WithStartAndEndDate(newDateStart, newDateEnd)

				return ia
			}(),
			wantErr: false,
		}, {
			name:   "update event that dont exists",
			fields: newValidEventServiceFields(),
			request: Request{
				ID:   eventID.String(),
				Name: "New Name",
			},
			want: func() aggregate.Event {
				ia := initialEvent
				ia.Event.Name = "New Name"

				return ia
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := EventUpdater{
				updater:       tt.fields.updater,
				finder:        tt.fields.finder,
				observersList: tt.fields.observersList,
			}
			got, err := ec.UpdateEvent(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("UpdateEvent() got = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestEventUpdater_validateRequiredResources(t *testing.T) {
	tests := []struct {
		name    string
		fields  eventServiceFields
		wantErr bool
	}{
		{
			name: "required fields added to event service", fields: eventServiceFields{
				updater: repository.EventsStorage{},
				finder:  repository.EventsStorage{},
			},
			wantErr: false,
		},
		{
			name: "no updater repository", fields: eventServiceFields{
				finder: repository.EventsStorage{},
			},
			wantErr: true,
		},
		{
			name: "no finder repository", fields: eventServiceFields{
				updater: repository.EventsStorage{},
			},
			wantErr: true,
		},
		{
			name: "no observers", fields: eventServiceFields{
				updater: repository.EventsStorage{},
				finder:  repository.EventsStorage{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := EventUpdater{
				updater:       tt.fields.updater,
				finder:        tt.fields.finder,
				observersList: tt.fields.observersList,
			}
			if err := ec.validateRequiredResources(); (err != nil) != tt.wantErr {
				t.Errorf("validateRequiredResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEventUpdater(t *testing.T) {
	tests := []struct {
		name          string
		configuration []Configuration
		want          *EventUpdater
		wantErr       bool
	}{
		{
			name: "creates valid event updater",
			configuration: []Configuration{
				WithUpdaterRepository(repository.EventsStorage{}),
				WithFinderRepository(repository.EventsStorage{}),
				WithObservers(mockObserver{}),
			},
			want: &EventUpdater{
				updater:       repository.EventsStorage{},
				finder:        repository.EventsStorage{},
				observersList: []Observer{mockObserver{}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEventUpdater(tt.configuration...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEventUpdater() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventUpdater() got = %v, want %v", got, tt.want)
			}
		})
	}
}
