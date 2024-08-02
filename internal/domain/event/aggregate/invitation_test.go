package aggregate

import (
	"reflect"
	"testing"
	"time"

	"event-service/internal/domain/common/entity"
	"event-service/internal/domain/common/valueobject"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func TestInvitation_Accept(t *testing.T) {
	tests := []struct {
		name        string
		Event       *Event
		InvitedUser uuid.UUID
		Accepted    time.Time
		wantErr     bool
	}{
		{
			name:        "event valid for acceptance",
			Event:       newEventMock(2, time.Now().Add(24*time.Hour)),
			InvitedUser: uuid.New(),
			wantErr:     false,
		},
		{
			name:        "invitation for the event already accepted",
			Event:       newEventMock(2, time.Now().Add(24*time.Hour)),
			InvitedUser: uuid.New(),
			Accepted:    time.Now().Add(-240 * time.Hour),
			wantErr:     true,
		},
		{
			name:        "capacity for the event already reached",
			Event:       newEventMock(2, time.Now().Add(24*time.Hour), uuid.New(), uuid.New()),
			InvitedUser: uuid.New(),
			wantErr:     true,
		},
		{
			name:        "registration for the event already ended",
			Event:       newEventMock(2, time.Now().Add(-24*time.Hour)),
			InvitedUser: uuid.New(),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Invitation{
				Event:       tt.Event,
				InvitedUser: tt.InvitedUser,
				AcceptedAt:  &tt.Accepted,
			}
			if err := i.Accept(); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewInvitation(t *testing.T) {
	validParticipant := uuid.New()
	validEvent := Event{
		ID:     1,
		UserID: uuid.New(),
		Event: &entity.Event{
			ExternalID: uuid.New(),
			Name:       "TestEvent",
			Capacity:   10,
			Public:     false,
		},
	}

	type args struct {
		event         *Event
		participantID uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		wantI   *Invitation
		wantErr bool
	}{
		{
			name: "no event for invitation",
			args: args{
				event:         nil,
				participantID: validParticipant,
			},
			wantI:   nil,
			wantErr: true,
		},
		{
			name: "invalid participant id",
			args: args{
				event:         &validEvent,
				participantID: uuid.UUID{},
			},
			wantI:   nil,
			wantErr: true,
		},
		{
			name: "valid participant and event",
			args: args{
				event:         &validEvent,
				participantID: validParticipant,
			},
			wantI: &Invitation{
				Event:       &validEvent,
				InvitedUser: validParticipant,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotI, err := NewInvitation(tt.args.event, tt.args.participantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInvitation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotI, tt.wantI) {
				t.Errorf("NewInvitation() gotI = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func newEventMock(cap int, endRegistrationDate time.Time, participants ...uuid.UUID) *Event {
	e := &Event{
		ID:     1,
		UserID: uuid.New(),
		Event: &entity.Event{
			ExternalID:  uuid.New(),
			Name:        faker.Name(),
			Description: faker.Sentence(),
			Capacity:    cap,
			Public:      false,
		},
		Location: &Location{
			Spot: valueobject.NewLocation(faker.Latitude(), faker.Longitude()),
		},
		RegistrationPeriod: valueobject.Period{},
	}

	if !endRegistrationDate.IsZero() {
		e.RegistrationPeriod, _ = valueobject.Period{}.WithStartAndEndDate(endRegistrationDate.Add(-240*time.Hour), endRegistrationDate)
	}

	if len(participants) > 0 {
		e.Participants = append(e.Participants, participants...)
	}

	return e
}
