package repository

import (
	"context"

	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type UserPreferenceRepository interface {
	GetPreference(ctx context.Context, user_id string) (*models.UserPreference, error)
	ListPreference(ctx context.Context, f *shared.Filter) ([]*models.UserPreference, int, error)
	CreatePreference(ctx context.Context, req *models.UserPreference) error
	UpdatePreference(ctx context.Context, req *models.UserPreference) error
	DeletePreference(ctx context.Context, uuid string) error
}

type preferenceRepo struct {
	db *bun.DB
}

func NewUserPreferenceRepository(db *bun.DB) *preferenceRepo {
	return &preferenceRepo{
		db: db,
	}
}

func (r *preferenceRepo) GetPreference(ctx context.Context, user_id string) (*models.UserPreference, error) {
	var preference = new(models.UserPreference)

	err := r.db.NewSelect().
		Model(preference).
		Where("user_id = ?", user_id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return preference, nil
}

func (r *preferenceRepo) ListPreference(ctx context.Context, f *shared.Filter) ([]*models.UserPreference, int, error) {
	var preference []*models.UserPreference

	q := r.db.NewSelect().
		Model(&preference).
		Limit(f.PageSize).
		Offset(f.Page).
		Order(f.Sort)

	if f.Search != "" {
		search := "%" + f.Search + "%"
		q.WhereGroup("OR", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.
				WhereOr("language ILIKE ?", search).
				WhereOr("currency ILIKE ?", search)
		})
	}

	count, err := q.ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return preference, count, nil
}

func (r *preferenceRepo) CreatePreference(ctx context.Context, req *models.UserPreference) error {
	_, err := r.db.NewInsert().Model(req).Exec(ctx)

	return err
}

func (r *preferenceRepo) UpdatePreference(ctx context.Context, req *models.UserPreference) error {
	cols := []string{
		"language",
		"currency",
		"email_opt_in",
		"sms_opt_in",
		"settings",
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

func (r *preferenceRepo) DeletePreference(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*models.UserPreference)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
