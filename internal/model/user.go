package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Username string `gorm:"type:varchar(100);unique;not null"`
	Email    string `gorm:"type:varchar(255);unique;not null"`
	Password string `gorm:"type:text;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
