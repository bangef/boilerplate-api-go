package domain

import (
	"context"
	"time"
)

type Permission struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Nama       string    `gorm:"type:varchar(100);not null;unqiue" json:"nama"`
	Desskripsi string    `gorm:"type:text;null;" json:"deskripsi,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type PermissionRepository interface {
	Create(ctx context.Context, permission *Permission) error
	GetAll(ctx context.Context, req *PermissionQuery) ([]Permission, int64, error)
	GetByID(ctx context.Context, id int) (*Permission, error)
	Update(ctx context.Context, id int, req *PermissionRequest) (*Permission, error)
	Delete(ctx context.Context, id int) error
}

type PermissionRequest struct {
	Nama      string `json:"nama" validate:"required"`
	Deskripsi string `json:"deskripsi"`
}

type PermissionQuery struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Search   string `query:"search"`
	OrderBy  string `query:"order_by" validate:"oneof=id nama created_at"`
	DirOrder string `query:"direction_order" validate:"oneof=ASC DESC asc desc"`
}
