package entity

import (
	"github.com/google/uuid"
)

type Event struct {
	ExternalID  uuid.UUID
	Name        string
	Description string
	Capacity    int
	Public      bool
}
