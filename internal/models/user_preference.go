package models

import (
	"context"
	"time"
	"uuid"

	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type UserPreference struct {
	bun.BaseModel `bun:"table:user_preference,alias:upr"`
	shared.Meta

	UserID     string            `json:"user_id"`
	Language   string            `json:"language"`
	Currency   string            `json:"currency"`
	EmailOptIn bool              `json:"email_opt_in"`
	SmsOptIn   bool              `json:"sms_opt_in"`
	Settings   map[string]string `json:"settings"`
}

func (u *UserPreference) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
