package user_addresses

import (
	"context"

	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type IUserAddressRepository interface {
	GetAddress(ctx context.Context, uuid string) (*models.UserAddress, error)
	GetAllAddress(ctx context.Context, f *shared.Filter) ([]*models.UserAddress, int, error)
	CreateAddress(ctx context.Context, req *models.UserAddress) (*models.UserAddress, error)
	UpdateAddress(ctx context.Context, uuid string, req *models.UserAddress) (*models.UserAddress, error)
	DeleteAddress(ctx context.Context, uuid string) error
	SetDefaultAddress(ctx context.Context, user_id, uuid string) (*models.UserAddress, error)
}

type UserAddressRepository struct {
	db *bun.DB
}

func NewUserAddressRepository(db *bun.DB) *UserAddressRepository {
	return &UserAddressRepository{db: db}
}

func (r *UserAddressRepository) GetAddress(ctx context.Context, uuid string) (*models.UserAddress, error) {
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

func (r *UserAddressRepository) GetAllAddress(ctx context.Context, f *shared.Filter) ([]*models.UserAddress, int, error) {
	var userAddr []*models.UserAddress

	q := r.db.NewSelect().
		Model(&userAddr).
		Limit(f.PageSize).
		Offset(f.Page).
		Order(f.Sort)

	// @TODO: Add where here <-

	count, err := q.ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return userAddr, count, nil
}

func (r *UserAddressRepository) CreateAddress(ctx context.Context, req *models.UserAddress) (*models.UserAddress, error) {
	_, err := r.db.NewInsert().Model(req).Exec(ctx)

	return nil, err
}

func (r *UserAddressRepository) UpdateAddress(ctx context.Context, uuid string, req *models.UserAddress) (*models.UserAddress, error) {
	return nil, nil
}

func (r *UserAddressRepository) DeleteAddress(ctx context.Context, uuid string) error {
	return nil
}

func (r *UserAddressRepository) SetDefaultAddress(ctx context.Context, user_id, uuid string) (*models.UserAddress, error) {
	return nil, nil
}
