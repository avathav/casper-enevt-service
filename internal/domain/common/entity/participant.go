package entity

import "github.com/google/uuid"

type Participant struct {
	ID         uint
	ExternalID uuid.UUID
}
