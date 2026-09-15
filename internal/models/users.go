package models

import (
	"context"
	"time"
	"uuid"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Status    string    `json:"status"`
	Uuid      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"deleted_at" bun:"deleted_at,nullzero"`

	Username   string `json:"username"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile,omitempty" bun:"default:null"`
	Password   string `json:"password,omitempty" bun:"default:null"`
	FirstName  string `json:"first_name,omitempty" bun:"default:null"`
	MiddleName string `json:"middle_name,omitempty" bun:"default:null"`
	LastName   string `json:"last_name,omitempty" bun:"default:null"`
	IsVerified bool   `json:"is_verified"`
}

func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
