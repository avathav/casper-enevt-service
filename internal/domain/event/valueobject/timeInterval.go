package valueobject

import "time"

type TimeInterval struct {
	date     time.Time
	interval time.Duration
}

func NewTimeInterval(date time.Time, interval time.Duration) *TimeInterval {
	return &TimeInterval{date: date, interval: interval}
}

func (i TimeInterval) GetDate() time.Time {
	return i.date
}

func (i TimeInterval) GetInterval() time.Duration {
	return i.interval
}

func (i TimeInterval) IsSet() bool {
	return !i.date.IsZero() && i.interval.Nanoseconds() > 0
}

func (i TimeInterval) GetStartDate() time.Time {
	return i.date.Add(-i.interval)
}

func (i TimeInterval) GetEndDate() time.Time {
	return i.date.Add(i.interval)
}
