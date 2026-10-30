package domain

import (
	"context"
	"time"
)

type Role struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	Nama       string       `gorm:"type:varchar(100);not null;unqiue" json:"nama"`
	Desskripsi string       `gorm:"type:text;null;" json:"deskripsi,omitempty"`
	Permission []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"` // Relasi M2M to Permission
	CreatedAt  time.Time    `json:"created_at"`
}

type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	GetAll(ctx context.Context, page, limit int) ([]Role, int64, error)
	GetByID(ctx context.Context, id int) (*Role, error)
	Update(ctx context.Context, id int, req *RoleRequest) (*Role, error)
	Delete(ctx context.Context, id int) error
}

type RoleRequest struct {
	Nama      string `json:"nama" validate:"alphanum,required"`
	Deskripsi string `json:"deskripsi" validate:"alphanum,required"`
}
