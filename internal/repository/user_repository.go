package repository

import (
	"auth-service/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *model.User) error {

	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {

	var user model.User

	err := r.DB.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
