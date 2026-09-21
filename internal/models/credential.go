package models

import (
	"context"
	"uuid"

	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type Credential struct {
	bun.BaseModel `bun:"table:user_credential,alias:u"`
	shared.Meta

	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-" bun:"password_hash"`
	IsVerified   bool   `json:"is_verified"`

	Profile *Profile `bun:"rel:has-one,join:uuid=user_id" json:"profile"`
}

func (u *Credential) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
