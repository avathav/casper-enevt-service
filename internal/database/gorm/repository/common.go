package repository

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel TODO Test Before Update functionality
type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()

	return nil
}
