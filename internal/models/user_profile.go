package models

import (
	"context"
	"time"
	"uuid"

	"github.com/rdev-org/users/internal/shared"
	"github.com/uptrace/bun"
)

type UserProfile struct {
	bun.BaseModel `bun:"table:user_profile,alias:up"`
	shared.Meta

	UserID       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	DateOfBirth  string `json:"date_of_birth"`
	Gender       string `json:"gender"`
	Avatar       string `json:"avatar"`
}

func (u *UserProfile) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
