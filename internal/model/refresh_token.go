package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	Token string `gorm:"type:text;not null"`

	Revoked bool `gorm:"default:false"`

	ExpiresAt time.Time

	CreatedAt time.Time
}
