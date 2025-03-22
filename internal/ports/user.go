package ports

import (
	"context"
	"unit-test-mongo/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.UserCreate) (userID string, err error)
	Get(ctx context.Context, userID string) (*domain.User, error)
	Delete(ctx context.Context, userID string) error
	Update(ctx context.Context, userID string, user *domain.User) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}
