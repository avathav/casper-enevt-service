package valueobject

import (
	"reflect"
	"testing"
	"time"
)

func TestEventPeriod_WithStartAndDuration(t *testing.T) {
	type args struct {
		startDate time.Time
		duration  time.Duration
	}
	tests := []struct {
		name               string
		args               args
		wantNewEventPeriod EventPeriod
		wantErr            bool
	}{
		{
			name: "start date with duration",
			args: args{
				startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
				duration:  24 * time.Hour,
			},
			wantNewEventPeriod: EventPeriod{
				period: Period{
					startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
					endDate:   time.Date(2022, 10, 2, 5, 0, 0, 0, time.UTC),
				},
				duration: 24 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "start date without duration",
			args: args{
				startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
			},
			wantNewEventPeriod: EventPeriod{},
			wantErr:            true,
		},
		{
			name: "duration without start date",
			args: args{
				startDate: time.Time{},
				duration:  24 * time.Hour,
			},
			wantNewEventPeriod: EventPeriod{},
			wantErr:            true,
		},
		{
			name: "negative duration with start date",
			args: args{
				startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
				duration:  -24 * time.Hour,
			},
			wantNewEventPeriod: EventPeriod{
				period: Period{
					startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
					endDate:   time.Date(2022, 10, 2, 5, 0, 0, 0, time.UTC),
				},
				duration: 24 * time.Hour,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := EventPeriod{}
			gotNewEventPeriod, err := ep.WithStartAndDuration(tt.args.startDate, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithStartAndDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewEventPeriod, tt.wantNewEventPeriod) {
				t.Errorf("WithStartAndDuration() gotNewEventPeriod = %v, want %v", gotNewEventPeriod, tt.wantNewEventPeriod)
			}
		})
	}
}

func TestEventPeriod_WithStartAndEndDate(t *testing.T) {
	type args struct {
		startDate time.Time
		endDate   time.Time
	}
	tests := []struct {
		name               string
		args               args
		wantNewEventPeriod EventPeriod
		wantErr            bool
	}{
		{
			name: "valid start date and end date",
			args: args{
				startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 10, 2, 5, 0, 0, 0, time.UTC),
			},
			wantNewEventPeriod: EventPeriod{
				period: Period{
					startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
					endDate:   time.Date(2022, 10, 2, 5, 0, 0, 0, time.UTC),
				},
				duration: 24 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "start date without end date",
			args: args{
				startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
			},
			wantNewEventPeriod: EventPeriod{},
			wantErr:            true,
		},
		{
			name: "end date without start date",
			args: args{
				endDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
			},
			wantNewEventPeriod: EventPeriod{},
			wantErr:            true,
		},
		{
			name: "start date after end date",
			args: args{
				startDate: time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 10, 2, 5, 0, 0, 0, time.UTC),
			},
			wantNewEventPeriod: EventPeriod{},
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := EventPeriod{}

			gotNewEventPeriod, err := ep.WithStartAndEndDate(tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithStartAndEndDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewEventPeriod, tt.wantNewEventPeriod) {
				t.Errorf("WithStartAndEndDate() gotNewEventPeriod = %v, want %v", gotNewEventPeriod, tt.wantNewEventPeriod)
			}
		})
	}
}
