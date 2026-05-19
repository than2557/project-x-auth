package repository

import (
	"auth-service/internal/model"

	"gorm.io/gorm"
)

type RefreshRepository struct {
	db *gorm.DB // ✅ lowercase
}

func NewRefreshRepository(db *gorm.DB) *RefreshRepository {
	return &RefreshRepository{db: db}
}

func (r *RefreshRepository) Create(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshRepository) FindByToken(token string) (*model.RefreshToken, error) {
	var refresh model.RefreshToken

	err := r.db.
		Where("token = ?", token).
		First(&refresh).Error

	if err != nil {
		return nil, err
	}

	return &refresh, nil
}

func (r *RefreshRepository) Revoke(id string) error {
	return r.db.Model(&model.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}
