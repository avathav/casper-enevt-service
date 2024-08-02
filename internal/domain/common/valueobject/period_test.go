package valueobject

import (
	"testing"
	"time"
)

func TestPeriod_After(t *testing.T) {
	r := Period{
		startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
		endDate:   time.Date(2022, 12, 1, 5, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name string
		args Period
		want bool
	}{
		{
			name: "start and end date before period starts",
			args: Period{
				startDate: time.Date(2022, 8, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 9, 1, 5, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "start and end date after period starts",
			args: Period{
				startDate: time.Date(2023, 1, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2023, 2, 1, 5, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "start date before period starts end date in period",
			args: Period{
				startDate: time.Date(2022, 9, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "start date in period starts end date after period ends",
			args: Period{
				startDate: time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2023, 1, 1, 5, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "start and end date in the period",
			args: Period{
				startDate: time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := r.After(tt.args); got != tt.want {
				t.Errorf("After() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriod_Before(t *testing.T) {
	r := Period{
		startDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
		endDate:   time.Date(2022, 12, 1, 5, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name string
		args Period
		want bool
	}{
		{
			name: "start and end date before period starts",
			args: Period{
				startDate: time.Date(2022, 8, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 9, 1, 5, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "start and end date after period starts",
			args: Period{
				startDate: time.Date(2023, 1, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2023, 2, 1, 5, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "start date before period starts end date in period",
			args: Period{
				startDate: time.Date(2022, 9, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "start date in period starts end date after period ends",
			args: Period{
				startDate: time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2023, 1, 1, 5, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "start and end date in the period",
			args: Period{
				startDate: time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Before(tt.args); got != tt.want {
				t.Errorf("Before() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriod_Validate(t *testing.T) {
	type fields struct {
		StartDate time.Time
		EndDate   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "start date after end date",
			fields: fields{
				StartDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "only start date",
			fields: fields{
				StartDate: time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
		{
			name: "only end date",
			fields: fields{
				EndDate: time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
		{
			name: "start date after end date",
			fields: fields{
				StartDate: time.Date(2022, 11, 1, 5, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2022, 10, 1, 5, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Period{
				startDate: tt.fields.StartDate,
				endDate:   tt.fields.EndDate,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
