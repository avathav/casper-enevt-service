package di

import (
	"event-service/internal/database/gorm/repository"
)

func EventsRepository() *repository.EventRepository {
	return repository.NewEventRepository()
}

func InvitationRepository() *repository.InvitationRepository {
	return repository.NewInvitationRepository()
}
