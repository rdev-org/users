package shared

import "time"

type Meta struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	Status    string    `json:"status"`
	Uuid      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
}
