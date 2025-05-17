package entity

import "time"

type (
	Base struct {
		ID        int64      `db:"id" json:"id"`
		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		CreatedBy string     `db:"created_by" json:"created_by"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		UpdatedBy *string    `db:"updated_by" json:"updated_by"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
		DeletedBy *string    `db:"deleted_by" json:"deleted_by"`
		IsDeleted int8       `db:"is_deleted" json:"is_deleted"`
	}
)
