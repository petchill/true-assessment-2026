package interfaces

import (
	"context"
	"recommendation-system/src/internal/model/entity"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID int) (entity.User, error)
}
