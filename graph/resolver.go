package graph

import (
	"event-service/internal/services/eventcreator"
	"event-service/internal/services/eventfinder"
	"event-service/internal/services/eventupdater"
	"event-service/internal/services/invitation"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AddEventHandler    eventcreator.Handler
	FindEventsHandler  eventfinder.ListHandler
	UpdateEventHandler eventupdater.Handler
	InvitationHandler  invitation.Handler
}
