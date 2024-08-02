package valueobject

import (
	"errors"
	"time"
)

var ErrPeriodEmptyDates = errors.New("dates cannot be empty")
var ErrPeriodInvalidDates = errors.New("start date must be before end Date")

type Period struct {
	startDate, endDate time.Time
}

func (r Period) WithStartAndEndDate(startDate time.Time, endDate time.Time) (Period, error) {
	p := Period{startDate: startDate, endDate: endDate}

	if err := p.Validate(); err != nil {
		return Period{}, err
	}

	return p, nil
}

func (r Period) Validate() error {
	switch {
	case r.startDate.IsZero() || r.endDate.IsZero():
		return ErrPeriodEmptyDates
	case r.startDate.After(r.endDate):
		return ErrPeriodInvalidDates
	}

	return nil
}

func (r Period) Before(np Period) bool {
	return r.startDate.Before(np.startDate) && r.endDate.Before(np.endDate)
}

func (r Period) After(np Period) bool {
	return r.startDate.After(np.startDate) && r.endDate.After(np.endDate)
}

func (r Period) Start() time.Time {
	return r.startDate
}

func (r Period) End() time.Time {
	return r.endDate
}

func (r Period) Contains(d time.Time) bool {
	return r.startDate.Before(d) && r.endDate.After(d)
}
