package echange

import log "github.com/sirupsen/logrus"

type EventQueryHandler struct{}

func NewEventQueryHandler() *EventQueryHandler {
	return &EventQueryHandler{}
}

func (h *EventQueryHandler) Handle(data []byte) {
	log.Printf("Event handled: %s", data)
}
