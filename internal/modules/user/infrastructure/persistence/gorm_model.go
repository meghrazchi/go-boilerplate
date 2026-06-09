package persistence

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (UserModel) TableName() string {
	return "users"
}
