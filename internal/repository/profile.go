package repository

import (
	"context"

	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type UserProfileRepository interface {
	GetProfile(ctx context.Context, user_id string) (*models.UserProfile, error)
	ListProfile(ctx context.Context, f *shared.Filter) ([]*models.UserProfile, int, error)
	CreateProfile(ctx context.Context, req *models.UserProfile) error
	UpdateProfile(ctx context.Context, req *models.UserProfile) error
	DeleteProfile(ctx context.Context, uuid string) error
}

type profileRepo struct {
	db *bun.DB
}

func NewUserProfileRepository(db *bun.DB) *profileRepo {
	return &profileRepo{
		db: db,
	}
}

func (r *profileRepo) GetProfile(ctx context.Context, user_id string) (*models.UserProfile, error) {
	var profile = new(models.UserProfile)

	err := r.db.NewSelect().
		Model(profile).
		Where("user_id = ?", user_id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepo) ListProfile(ctx context.Context, f *shared.Filter) ([]*models.UserProfile, int, error) {
	var profile []*models.UserProfile

	q := r.db.NewSelect().
		Model(&profile).
		Limit(f.PageSize).
		Offset(f.Page).
		Order(f.Sort)

	if f.Search != "" {
		search := "%" + f.Search + "%"
		q.WhereGroup("OR", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.
				WhereOr("first_name ILIKE ?", search).
				WhereOr("middle_name ILIKE ?", search).
				WhereOr("last_name ILIKE ?", search).
				WhereOr("mobile_number ILIKE ?", search)
		})
	}

	count, err := q.ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return profile, count, nil
}

func (r *profileRepo) CreateProfile(ctx context.Context, req *models.UserProfile) error {
	_, err := r.db.NewInsert().Model(req).Exec(ctx)

	return err
}

func (r *profileRepo) UpdateProfile(ctx context.Context, req *models.UserProfile) error {
	cols := []string{
		"first_name",
		"middle_name",
		"last_name",
		"mobile_number",
		"date_of_birth",
		"gender",
		"avatar",
		"updated_at",
	}

	_, err := r.db.NewUpdate().
		Model(req).
		Column(cols...).
		Where("user_id = ?", req.UserID).
		Where("uuid = ?", req.Uuid).
		Exec(ctx)

	return err
}

func (r *profileRepo) DeleteProfile(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*models.UserProfile)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
