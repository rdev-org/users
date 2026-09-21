package auth

import (
	"context"

	"github.com/rdev-org/users/internal/models"
	"github.com/uptrace/bun"
)

type CredentialRepository interface {
	Login(ctx context.Context, identifier string) (*models.Credential, error)
	Register(ctx context.Context, req *RegisterRequest) (string, error)
	DeleteCredential(ctx context.Context, uuid string) error
	UpdatePassword(ctx context.Context, req *UpdatePasswordRequest) error
	CheckEmail(ctx context.Context, identifier string) (bool, error)
}

type Repository struct {
	db *bun.DB
}

func NewAuthenticationRepository(db *bun.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CheckEmail(ctx context.Context, identifier string) (bool, error) {
	return r.db.NewSelect().
		Model((*models.Credential)(nil)).
		Where("email = ?", identifier).
		Exists(ctx)
}

func (r *Repository) Login(ctx context.Context, identifier string) (*models.Credential, error) {
	var cred = new(models.Credential)

	err := r.db.NewSelect().
		Model(cred).
		Relation("Profile").
		Where("username = ?", identifier).
		WhereOr("email = ?", identifier).
		Where("is_verified = ?", true).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (r *Repository) Register(ctx context.Context, req *RegisterRequest) (string, error) {
	user := &models.Credential{
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	_, err := r.db.NewInsert().Model(user).Exec(ctx)

	return "", err
}

func (r *Repository) DeleteCredential(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*models.Credential)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *Repository) UpdatePassword(ctx context.Context, req *UpdatePasswordRequest) error {
	_, err := r.db.NewUpdate().
		Model((*models.Credential)(nil)).
		Set("password_hash = ?", req.Password).
		Where("uuid = ?", req.Uuid).
		Exec(ctx)

	return err
}
