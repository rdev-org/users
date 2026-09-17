package user_address

import (
	"context"
	"time"

	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type UserAddressRepository interface {
	GetAddress(ctx context.Context, uuid string) (*models.UserAddress, error)
	ListAddress(ctx context.Context, f *shared.Filter) ([]*models.UserAddress, int, error)
	CreateAddress(ctx context.Context, req *models.UserAddress) error
	UpdateAddress(ctx context.Context, req *models.UserAddress) error
	DeleteAddress(ctx context.Context, uuid string) error
	DefaultAddress(ctx context.Context, user_id, uuid string) error
}

type Repository struct {
	db *bun.DB
}

func NewUserAddressRepository(db *bun.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAddress(ctx context.Context, uuid string) (*models.UserAddress, error) {
	var userAddr = new(models.UserAddress)

	err := r.db.NewSelect().
		Model(userAddr).
		Where("uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return userAddr, nil
}

func (r *Repository) ListAddress(ctx context.Context, f *shared.Filter) ([]*models.UserAddress, int, error) {
	var addr []*models.UserAddress

	q := r.db.NewSelect().
		Model(&addr).
		Limit(f.PageSize).
		Offset(f.Page).
		Order(f.Sort)

	if f.Search != "" {
		search := "%" + f.Search + "%"
		q.WhereGroup("OR", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.
				WhereOr("address_line1 ILIKE ?", search).
				WhereOr("address_line2 ILIKE ?", search).
				WhereOr("city ILIKE ?", search).
				WhereOr("state_province ILIKE ?", search).
				WhereOr("postal_code ILIKE ?", search).
				WhereOr("country ILIKE ?", search).
				WhereOr("address_type ILIKE ?", search)
		})
	}

	count, err := q.ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return addr, count, nil
}

func (r *Repository) CreateAddress(ctx context.Context, req *models.UserAddress) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var err error

		if req.IsDefault {
			_, err = tx.NewUpdate().
				Model((*models.UserAddress)(nil)).
				Set("is_default = ?", false).
				Set("updated_at = ?", time.Now()).
				Where("user_id = ?", req.UserID).
				WhereAllWithDeleted().
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		_, err = tx.NewInsert().Model(req).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAddress(ctx context.Context, req *models.UserAddress) error {
	cols := []string{
		"address_line1",
		"address_line2",
		"city",
		"state_province",
		"postal_code",
		"country",
		"is_default",
		"address_type",
	}

	_, err := r.db.NewUpdate().
		Model(req).
		Column(cols...).
		Set("updated_at = ?", time.Now()).
		Where("uuid = ?", req.Uuid).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteAddress(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*models.UserAddress)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *Repository) DefaultAddress(ctx context.Context, user_id, uuid string) error {
	var updatedAddress models.UserAddress

	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewUpdate().
			Model((*models.UserAddress)(nil)).
			Set("is_default = ?", false).
			Set("updated_at = ?", time.Now()).
			Where("user_id = ?", user_id).
			WhereAllWithDeleted().
			Exec(ctx)
		if err != nil {
			return err
		}

		err = tx.NewUpdate().
			Model(&updatedAddress).
			Set("is_default = ?", true).
			Set("updated_at = ?", time.Now()).
			Where("user_id = ? AND uuid = ?", user_id, uuid).
			Returning("*").
			Scan(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
