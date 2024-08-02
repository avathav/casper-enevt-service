package valueobject

import (
	"errors"
	"time"
)

var ErrDurationRequired = errors.New("duration must be set")

type EventPeriod struct {
	period   Period
	duration time.Duration
}

func (ep EventPeriod) Start() time.Time {
	return ep.period.Start()
}

func (ep EventPeriod) End() time.Time {
	return ep.period.End()
}

func (ep EventPeriod) Duration() time.Duration {
	return ep.duration
}

func NewEventPeriod(period Period, duration time.Duration) *EventPeriod {
	return &EventPeriod{period: period, duration: duration}
}

func (ep EventPeriod) WithStartAndEndDate(startDate time.Time, endDate time.Time) (newEventPeriod EventPeriod, err error) {
	newEventPeriod.period, err = Period{}.WithStartAndEndDate(startDate, endDate)
	if err != nil {
		return EventPeriod{}, err
	}

	newEventPeriod.duration = endDate.Sub(startDate)

	return newEventPeriod, nil
}

func (ep EventPeriod) WithStartAndDuration(startDate time.Time, duration time.Duration) (newEventPeriod EventPeriod, err error) {
	if duration.Nanoseconds() == 0 {
		return EventPeriod{}, ErrDurationRequired
	}

	newEventPeriod.duration = duration.Abs()

	newEventPeriod.period, err = Period{}.WithStartAndEndDate(startDate, startDate.Add(newEventPeriod.duration))
	if err != nil {
		return EventPeriod{}, err
	}

	return newEventPeriod, nil
}
