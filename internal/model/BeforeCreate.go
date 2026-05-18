package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) error {

	u.ID = uuid.New()

	return nil
}
