package graph

import (
	"time"

	"event-service/graph/model"
	"event-service/internal/domain/event/aggregate"
	"event-service/internal/services/eventcreator"
	"event-service/internal/services/eventupdater"
)

func ConvertNewEventToRequest(e model.NewEvent) eventcreator.Request {
	return eventcreator.Request{
		User:                e.User,
		Name:                e.Name,
		Description:         getValueIfNotNull(e.Description),
		Capacity:            e.Capacity,
		Duration:            time.Duration(e.Duration) * time.Hour * 24,
		Longitude:           e.Longitude,
		Latitude:            e.Latitude,
		DateStart:           e.StartDate,
		DateRegistrationEnd: getValueIfNotNull(e.RegistrationEndDate),
		Public:              e.Public,
	}
}

func ConvertEventEntryToModel(entry *aggregate.Event) *model.Event {
	e := &model.Event{
		ID:                    entry.Event.ExternalID.String(),
		User:                  entry.UserID.String(),
		Name:                  entry.Event.Name,
		Description:           &entry.Event.Description,
		Capacity:              entry.Event.Capacity,
		Duration:              int(entry.EventPeriod.Duration().Hours() / 24),
		StartDate:             entry.EventPeriod.Start(),
		EndDate:               entry.EventPeriod.End(),
		RegistrationStartDate: entry.RegistrationPeriod.Start(),
		RegistrationEndDate:   entry.RegistrationPeriod.End(),
		Public:                entry.Event.Public,
	}

	if l := entry.Location; l != nil {
		e.Latitude = l.Spot.Lat()
		e.Longitude = l.Spot.Long()
	}

	if entry.ParticipantsNumber() > 0 {
		for _, p := range entry.Participants {
			e.Participants = append(e.Participants, &model.Participant{User: p.String()})
		}
	}

	return e
}

func ConvertEventToUpdateRequest(e model.UpdateEvent) eventupdater.Request {
	r := eventupdater.Request{
		ID:          e.ID,
		Name:        getValueIfNotNull(e.Name),
		Description: getValueIfNotNull(e.Description),
		Capacity:    getValueIfNotNull(e.Capacity),
		Longitude:   getValueIfNotNull(e.Longitude),
		Latitude:    getValueIfNotNull(e.Latitude),
		Public:      e.Public,
	}

	if e.EventDate != nil {
		r.DateStart = &e.EventDate.Start
		r.DateEnd = &e.EventDate.End
	}

	if e.RegistrationDate != nil {
		r.DateRegistrationStart = &e.RegistrationDate.Start
		r.DateRegistrationEnd = &e.RegistrationDate.End
	}

	return r
}
