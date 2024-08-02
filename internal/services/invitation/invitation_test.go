package invitation

import (
	"context"
	"errors"
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

type invitationServiceFields struct {
	inviter       event.Inviter
	inviteFinder  event.InviteFinder
	eventFinder   event.Finder
	observersList map[EventType][]Observer
}

func newInvitationServiceFields(aggregates ...any) *invitationServiceFields {
	eventStorage := repository.NewEventsStorage()
	invitationStorage := repository.NewInvitationsStorage()

	if len(aggregates) > 0 {
		for _, a := range aggregates {
			if ea, ok := a.(aggregate.Event); ok {
				_ = eventStorage.Add(context.Background(), &ea)
			}

			if ia, ok := a.(aggregate.Invitation); ok {
				_ = invitationStorage.Add(&ia)
			}
		}
	}

	i := &invitationServiceFields{
		inviter:       invitationStorage,
		inviteFinder:  invitationStorage,
		eventFinder:   eventStorage,
		observersList: nil,
	}

	return i
}

func TestInvitation_Accept(t *testing.T) {
	mockInvitation := newInvitationAggregateMock()

	type args struct {
		eventExternalID string
		user            string
	}

	tests := []struct {
		name    string
		fields  *invitationServiceFields
		args    args
		wantErr bool
		errType error
	}{
		{
			name:   "valid invitation can be accepted",
			fields: newInvitationServiceFields(mockInvitation),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: false,
		},
		{
			name:   "invitation not in storage",
			fields: newInvitationServiceFields(*mockInvitation.Event),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: ErrInvitationNotFound,
		},
		{
			name: "invitation already accepted",
			fields: newInvitationServiceFields(func(mi aggregate.Invitation) aggregate.Invitation {
				now := time.Now().Add(-10 * time.Hour)
				mi.AcceptedAt = &now

				return mi
			}(mockInvitation)),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: aggregate.ErrInvitationAlreadyAccepted,
		},
		{
			name: "event capacity reached",
			fields: newInvitationServiceFields(func(mi aggregate.Invitation) aggregate.Invitation {
				e := *mi.Event
				e.Event.Capacity = 2
				e.Participants = []uuid.UUID{uuid.New(), uuid.New()}
				mi.Event = &e

				return mi
			}(mockInvitation)),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: aggregate.ErrCapacityFull,
		},
		{
			name: "registration closed",
			fields: newInvitationServiceFields(func(mi aggregate.Invitation) aggregate.Invitation {
				now := time.Now()
				period, _ := valueobject.Period{}.WithStartAndEndDate(now.Add(-10*time.Hour), now.Add(-5*time.Hour))
				e := *mi.Event
				e.RegistrationPeriod = period
				mi.Event = &e

				return mi
			}(mockInvitation)),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: aggregate.ErrRegistrationClosed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Invitation{
				inviter:       tt.fields.inviter,
				inviteFinder:  tt.fields.inviteFinder,
				eventFinder:   tt.fields.eventFinder,
				observersList: tt.fields.observersList,
			}
			err := i.Accept(context.Background(), tt.args.eventExternalID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil && tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.errType)
			}
		})
	}
}

func TestInvitation_Invite(t *testing.T) {
	mockInvitation := newInvitationAggregateMock()

	type args struct {
		eventExternalID string
		user            string
	}
	tests := []struct {
		name    string
		fields  *invitationServiceFields
		args    args
		wantErr bool
		errType error
	}{
		{
			name:   "new invitation added to clear storage",
			fields: newInvitationServiceFields(*mockInvitation.Event),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: false,
		},
		{
			name:   "invitation already exists",
			fields: newInvitationServiceFields(*mockInvitation.Event, mockInvitation),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: ErrInvitationAlreadyExists,
		},
		{
			name:   "event does not exists",
			fields: newInvitationServiceFields(),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Invitation{
				inviter:       tt.fields.inviter,
				inviteFinder:  tt.fields.inviteFinder,
				eventFinder:   tt.fields.eventFinder,
				observersList: tt.fields.observersList,
			}
			err := i.Invite(context.Background(), tt.args.eventExternalID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invite() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil && tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("Invite() error = %v, wantErr %v", err, tt.errType)
			}
		})
	}
}

func newInvitationAggregateMock() aggregate.Invitation {
	now := time.Now()
	period, _ := valueobject.Period{}.WithStartAndEndDate(now.Add(-10*time.Hour), now.Add(10*time.Hour))

	return aggregate.Invitation{
		Event: &aggregate.Event{
			ID:     1,
			UserID: uuid.New(),
			Event: &entity.Event{
				ExternalID: uuid.New(),
				Capacity:   10,
				Public:     true,
			},
			RegistrationPeriod: period,
		},
		InvitedUser: uuid.New(),
		AcceptedAt:  nil,
	}
}

func TestInvitation_Join(t *testing.T) {
	mockInvitation := newInvitationAggregateMock()

	type args struct {
		eventExternalID string
		user            string
	}

	tests := []struct {
		name    string
		fields  *invitationServiceFields
		args    args
		wantErr bool
		errType error
	}{
		{
			name:   "valid invitation can be accepted",
			fields: newInvitationServiceFields(*mockInvitation.Event),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: false,
		},
		{
			name:   "invitation not in storage",
			fields: newInvitationServiceFields(),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
		},
		{
			name: "invitation already accepted",
			fields: newInvitationServiceFields(*mockInvitation.Event, func(mi aggregate.Invitation) aggregate.Invitation {
				now := time.Now().Add(-10 * time.Hour)
				mi.AcceptedAt = &now

				return mi
			}(mockInvitation)),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: aggregate.ErrInvitationAlreadyAccepted,
		},
		{
			name: "event capacity reached",
			fields: newInvitationServiceFields(func(me aggregate.Event) aggregate.Event {
				me.Event.Capacity = 2
				me.Participants = []uuid.UUID{uuid.New(), uuid.New()}

				return me
			}(*mockInvitation.Event)),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: aggregate.ErrCapacityFull,
		},
		{
			name: "registration closed",
			fields: newInvitationServiceFields(func(me aggregate.Event) aggregate.Event {
				now := time.Now()
				period, _ := valueobject.Period{}.WithStartAndEndDate(now.Add(-10*time.Hour), now.Add(-5*time.Hour))
				me.RegistrationPeriod = period

				return me
			}(*mockInvitation.Event)),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: aggregate.ErrRegistrationClosed,
		},
		{
			name: "event is private",
			fields: newInvitationServiceFields(func(me aggregate.Event) aggregate.Event {
				e := *me.Event
				e.Public = false
				me.Event = &e

				return me
			}(*mockInvitation.Event)),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: ErrEventIsNotPublic,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Invitation{
				inviter:       tt.fields.inviter,
				inviteFinder:  tt.fields.inviteFinder,
				eventFinder:   tt.fields.eventFinder,
				observersList: tt.fields.observersList,
			}
			err := i.Join(context.Background(), tt.args.eventExternalID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil && tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.errType)
			}
		})
	}
}

func TestInvitation_Remove(t *testing.T) {
	hourAgo := time.Now().Add(-1 * time.Hour)
	mockInvitation := newInvitationAggregateMock()
	mockInvitation.AcceptedAt = &hourAgo

	type args struct {
		eventExternalID string
		user            string
	}

	tests := []struct {
		name    string
		fields  *invitationServiceFields
		args    args
		wantErr bool
		errType error
	}{
		{
			name:   "valid accepted user can be removed",
			fields: newInvitationServiceFields(mockInvitation),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: false,
		},
		{
			name:   "invitation not in storage",
			fields: newInvitationServiceFields(*mockInvitation.Event),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: ErrInvitationNotFound,
		},
		{
			name: "invitation is not accepted",
			fields: newInvitationServiceFields(func(mi aggregate.Invitation) aggregate.Invitation {
				mi.AcceptedAt = nil

				return mi
			}(mockInvitation)),
			args: args{
				eventExternalID: mockInvitation.Event.Event.ExternalID.String(),
				user:            mockInvitation.InvitedUser.String(),
			},
			wantErr: true,
			errType: ErrParticipantNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Invitation{
				inviter:       tt.fields.inviter,
				inviteFinder:  tt.fields.inviteFinder,
				eventFinder:   tt.fields.eventFinder,
				observersList: tt.fields.observersList,
			}
			err := i.Remove(context.Background(), tt.args.eventExternalID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil && tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.errType)
			}
		})
	}
}

func TestInvitation_validateRequiredResources(t *testing.T) {
	tests := []struct {
		name    string
		fields  *invitationServiceFields
		wantErr bool
	}{
		{
			name:    "valid resources",
			fields:  newInvitationServiceFields(nil, nil),
			wantErr: false,
		},
		{
			name: "no inviter repository",
			fields: &invitationServiceFields{
				inviteFinder: repository.InvitationsStorage{},
				eventFinder:  repository.EventsStorage{},
			},
			wantErr: true,
		},
		{
			name: "no invite finder  repository",
			fields: &invitationServiceFields{
				inviter:     repository.InvitationsStorage{},
				eventFinder: repository.EventsStorage{},
			},
			wantErr: true,
		},
		{
			name: "no event finder repository",
			fields: &invitationServiceFields{
				inviter:      repository.InvitationsStorage{},
				inviteFinder: repository.InvitationsStorage{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Invitation{
				inviter:       tt.fields.inviter,
				inviteFinder:  tt.fields.inviteFinder,
				eventFinder:   tt.fields.eventFinder,
				observersList: tt.fields.observersList,
			}
			if err := i.validateRequiredResources(); (err != nil) != tt.wantErr {
				t.Errorf("validateRequiredResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewInvitation(t *testing.T) {
	tests := []struct {
		name          string
		configuration []Configuration
		want          *Invitation
		wantErr       bool
	}{
		{
			name: "valid configuration",
			configuration: []Configuration{
				WithInvitationRepository(repository.InvitationsStorage{}),
				WithInviteFinderRepository(repository.InvitationsStorage{}),
				WithEventFinderRepository(repository.EventsStorage{}),
			},
			want: &Invitation{
				inviter:       repository.InvitationsStorage{},
				inviteFinder:  repository.InvitationsStorage{},
				eventFinder:   repository.EventsStorage{},
				observersList: map[EventType][]Observer{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInvitation(tt.configuration...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInvitation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInvitation() got = %v, want %v", got, tt.want)
			}
		})
	}
}
