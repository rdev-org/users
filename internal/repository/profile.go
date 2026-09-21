package repository

import (
	"context"
	"time"

	"github.com/rdev-org/users/internal/dto"
	"github.com/rdev-org/users/internal/jwt"
	"github.com/rdev-org/users/internal/models"
	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type ProfileRepository interface {
	MyProfile(ctx context.Context, uuid string) (*models.Profile, error)
	GetProfile(ctx context.Context, uuid string) (*models.Profile, error)
	ListProfile(ctx context.Context, f *shared.Filter) ([]*models.Profile, int, error)
	CreateProfile(ctx context.Context, req *dto.ProfileRequestDTO) (string, string, error)
	UpdateProfile(ctx context.Context, req *models.Profile) error
	DeleteProfile(ctx context.Context, uuid string) error
}

type profileRepo struct {
	db *bun.DB
}

func NewProfileRepository(db *bun.DB) *profileRepo {
	return &profileRepo{
		db: db,
	}
}

func (r *profileRepo) MyProfile(ctx context.Context, uuid string) (*models.Profile, error) {
	var profile = new(models.Profile)

	err := r.db.NewSelect().
		Model(profile).
		Relation("Addresses").
		Relation("Preference").
		Where("up.user_id = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepo) GetProfile(ctx context.Context, uuid string) (*models.Profile, error) {
	var profile = new(models.Profile)

	err := r.db.NewSelect().
		Model(profile).
		Relation("Addresses").
		Relation("Preference").
		Where("up.uuid = ?", uuid).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepo) ListProfile(ctx context.Context, f *shared.Filter) ([]*models.Profile, int, error) {
	var profile []*models.Profile

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

func (r *profileRepo) CreateProfile(ctx context.Context, req *dto.ProfileRequestDTO) (string, string, error) {
	var pendingRegistration = new(models.RequestRegistration)

	err := r.db.NewSelect().
		Model(pendingRegistration).
		Where("token_hash = ?", jwt.HashToken(req.Token)).
		Where("expires_at > ?", time.Now()).
		Scan(ctx)
	if err != nil {
		return "", "", err
	}

	credential := &models.Credential{
		Username:     req.Username,
		Email:        pendingRegistration.Email,
		PasswordHash: req.Password,
	}

	profile := &models.Profile{
		FirstName:    req.FirstName,
		MiddleName:   req.MiddleName,
		LastName:     req.LastName,
		MobileNumber: req.MobileNumber,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		Avatar:       req.Avatar,
	}

	err = r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		err := tx.NewInsert().
			Model(credential).
			Returning("uuid").
			Scan(ctx)
		if err != nil {
			return err
		}

		profile.UserID = credential.Uuid
		err = tx.NewInsert().
			Model(profile).
			Returning("uuid").
			Scan(ctx)

		if err != nil {
			return err
		}

		if req.Address != nil {
			address := &models.Address{
				UserID:        credential.Uuid,
				AddressLine1:  req.Address.AddressLine1,
				AddressLine2:  req.Address.AddressLine2,
				City:          req.Address.City,
				StateProvince: req.Address.StateProvince,
				PostalCode:    req.Address.PostalCode,
				Country:       req.Address.Country,
				IsDefault:     req.Address.IsDefault,
				AddressType:   req.Address.AddressType,
			}

			_, err := tx.NewInsert().
				Model(address).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", "", err
	}

	return profile.Meta.Uuid, pendingRegistration.Email, nil
}

func (r *profileRepo) UpdateProfile(ctx context.Context, req *models.Profile) error {
	cols := []string{
		"first_name",
		"middle_name",
		"last_name",
		"mobile_number",
		"date_of_birth",
		"gender",
		"avatar",
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
		Model((*models.Profile)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
