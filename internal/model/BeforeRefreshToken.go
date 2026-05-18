package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {

	r.ID = uuid.New()

	return nil
}
