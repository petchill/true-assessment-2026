package repository

import (
	"context"
	"recommendation-system/src/internal/model/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByID(ctx context.Context, userID int) (entity.User, error) {
	user := entity.User{}

	err := r.db.WithContext(ctx).
		Table("users").
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
