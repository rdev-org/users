package users

import (
	"fmt"

	"github.com/rdev-org/users/internal/models"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	if username == "testuser" {
		return &models.User{Username: username}, nil
	}

	return nil, fmt.Errorf("user not found")
}
