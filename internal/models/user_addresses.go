package models

import (
	"context"
	"time"
	"uuid"

	"github.com/uptrace/bun"
)

type UserAddress struct {
	bun.BaseModel `bun:"table:user_addresses,alias:ua"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Status    string    `json:"status"`
	Uuid      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`

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
	}
	return nil
}
