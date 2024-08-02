package aggregate

import (
	"event-service/internal/domain/common/valueobject"
)

type Location struct {
	ID   uint
	Spot valueobject.Location
}
