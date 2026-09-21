package models

import (
	"context"
	"time"
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

type RequestRegistration struct {
	bun.BaseModel `bun:"table:pending_registration,alias:pr"`

	Email     string    `json:"email"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
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
