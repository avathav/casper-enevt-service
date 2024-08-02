package aggregate

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewEvent(t *testing.T) {
	startDate := time.Now().Add(24 * 10 * time.Hour)
	registrationDate := startDate.Add(24 * 10 * time.Hour)

	tests := []struct {
		name    string
		args    EventPayload
		wantErr bool
	}{
		{
			name: "valid config",
			args: EventPayload{
				UserID:              uuid.New(),
				Name:                "Test Event",
				Description:         "Test Description",
				Lat:                 48.85853948884809,
				Long:                2.2944598405858394,
				Capacity:            10,
				Duration:            7 * 24 * time.Hour,
				StartDate:           startDate,
				RegistrationEndDate: registrationDate,
				Public:              true,
			},
			wantErr: false,
		},
		{
			name: "empty user ID",
			args: EventPayload{
				Name:                "Test Event",
				Description:         "Test Description",
				Lat:                 48.85853948884809,
				Long:                2.2944598405858394,
				Capacity:            10,
				Duration:            7 * 24 * time.Hour,
				StartDate:           startDate,
				RegistrationEndDate: registrationDate,
				Public:              true,
			},
			wantErr: true,
		},
		{
			name: "no name of the event",
			args: EventPayload{
				UserID:              uuid.New(),
				Description:         "Test Description",
				Lat:                 48.85853948884809,
				Long:                2.2944598405858394,
				Capacity:            10,
				Duration:            7 * 24 * time.Hour,
				StartDate:           startDate,
				RegistrationEndDate: registrationDate,
				Public:              true,
			},
			wantErr: true,
		},
		{
			name: "Empty Latitude And Lang",
			args: EventPayload{
				UserID:              uuid.New(),
				Name:                "Test Event",
				Description:         "Test Description",
				Capacity:            10,
				Duration:            7 * 24 * time.Hour,
				StartDate:           startDate,
				RegistrationEndDate: registrationDate,
				Public:              false,
			},
			wantErr: false,
		},
		{
			name: "no duration",
			args: EventPayload{
				UserID:              uuid.New(),
				Name:                "Test Event",
				Description:         "Test Description",
				Lat:                 48.85853948884809,
				Long:                2.2944598405858394,
				Capacity:            10,
				StartDate:           startDate,
				RegistrationEndDate: registrationDate,
				Public:              true,
			},
			wantErr: true,
		},
		{
			name: "empty capacity",
			args: EventPayload{
				UserID:              uuid.New(),
				Name:                "Test Event",
				Description:         "Test Description",
				Lat:                 48.85853948884809,
				Long:                2.2944598405858394,
				Duration:            7 * 24 * time.Hour,
				StartDate:           startDate,
				RegistrationEndDate: registrationDate,
				Public:              true,
			},
			wantErr: false,
		},
		{
			name: "no start date",
			args: EventPayload{
				UserID:              uuid.New(),
				Name:                "Test Event",
				Description:         "Test Description",
				Lat:                 48.85853948884809,
				Long:                2.2944598405858394,
				Capacity:            10,
				Duration:            7 * 24 * time.Hour,
				RegistrationEndDate: registrationDate,
				Public:              true,
			},
			wantErr: true,
		},
		{
			name: "no registration date",
			args: EventPayload{
				UserID:      uuid.New(),
				Name:        "Test Event",
				Description: "Test Description",
				Lat:         48.85853948884809,
				Long:        2.2944598405858394,
				Capacity:    10,
				Duration:    7 * 24 * time.Hour,
				StartDate:   startDate,
				Public:      true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEvent(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
