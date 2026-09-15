package users

import (
	"context"

	"github.com/rdev-org/users/internal/models"
	"github.com/uptrace/bun"
)

type IUserRepository interface {
	GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error)
	RegisterUser(ctx context.Context, user *models.User) error
}

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	var user = new(models.User)

	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", identifier).
		WhereOr("email = ?", identifier).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) RegisterUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)

	return err
}
