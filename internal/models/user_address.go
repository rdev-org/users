package models

import (
	"context"
	"time"
	"uuid"

	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type UserAddress struct {
	bun.BaseModel `bun:"table:user_address,alias:ua"`
	shared.Meta

	UserID        string `json:"user_id"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2,omitempty"`
	City          string `json:"city"`
	StateProvince string `json:"state_province"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	IsDefault     bool   `json:"is_default"`
	AddressType   string `json:"address_type"`
}

func (u *UserAddress) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if u.Status == "" {
			u.Status = "A"
		}
		if u.Uuid == "" {
			u.Uuid = uuid.New().String()
		}
	case *bun.UpdateQuery:
		if u.UpdatedAt.IsZero() {
			u.UpdatedAt = time.Now()
		}
	}
	return nil
}
